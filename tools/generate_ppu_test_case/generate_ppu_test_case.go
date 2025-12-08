package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/png" // Register PNG decoder
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/davidyorr/LuccaGB/gameboy"
	"github.com/davidyorr/LuccaGB/hasher"
)

func main() {
	// 1. CONFIGURATION
	romPath := flag.String("rom", "", "Path to the ROM file")
	imgPath := flag.String("image", "", "Path to the reference PNG screenshot")
	maxFrames := flag.Int("max", 1000, "Stop after this many frames if not found")

	// Custom Color Flags (Defaults are empty strings)
	c0 := flag.String("c0", "", "Color 0 (Darkest/Black) - e.g. #384828")
	c1 := flag.String("c1", "", "Color 1 (Dark Gray)     - e.g. #607028")
	c2 := flag.String("c2", "", "Color 2 (Light Gray)    - e.g. #a0a830")
	c3 := flag.String("c3", "", "Color 3 (Lightest/White)- e.g. #d0e040")

	flag.Parse()

	if *romPath == "" || *imgPath == "" {
		fmt.Println("Usage: go run tools/find_frame/find_frame.go -rom <rom.gb> -image <screenshot.png> -max <max_frames> [colors...]")
		os.Exit(1)
	}

	// 2. CALCULATE TARGET HASH FROM IMAGE
	fmt.Printf("Loading reference image: %s...\n", *imgPath)

	// Build the palette map from flags (if provided)
	var palette map[color.RGBA]byte
	if *c0 != "" {
		palette = make(map[color.RGBA]byte)
		// Map the input colors to the canonical hash bytes (00=Darkest, FF=Lightest)
		addColor(palette, *c0, 0x00)
		addColor(palette, *c1, 0x55)
		addColor(palette, *c2, 0xAA)
		addColor(palette, *c3, 0xFF)
	}

	// Open and Hash the PNG
	file, err := os.Open(*imgPath)
	if err != nil {
		die(fmt.Errorf("failed to open image: %w", err))
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		die(fmt.Errorf("failed to decode image: %w", err))
	}

	targetHash, err := hasher.HashImage(img, palette)
	if err != nil {
		die(fmt.Errorf("failed to hash image: %w", err))
	}

	fmt.Printf("Target Hash calculated: %s\n", targetHash)

	// 3. SEARCH FOR MATCHING FRAME
	fmt.Printf("Scanning %s for match (max %d frames)...\n", *romPath, *maxFrames)

	romData, err := os.ReadFile(*romPath)
	if err != nil {
		die(fmt.Errorf("failed to read ROM: %w", err))
	}

	gb := gameboy.New()
	gb.LoadRom(romData)

	for f := 1; f <= *maxFrames; f++ {
		// Run until one frame is ready
		for {
			_, frameReady, _ := gb.Step()
			if frameReady {
				break
			}
		}

		// Hash current emulator state
		buffer := gb.FrameBuffer()
		currentHash := hasher.HashFrameBuffer(buffer)

		// Compare
		if currentHash == targetHash {
			// 1. Calculate path relative to "roms/test" (the helper's base dir)
			cwd, _ := os.Getwd()
			// We assume the tool is run from project root, so we look for "roms/test"
			expectedBaseDir := filepath.Join(cwd, "roms", "test")
			absoluteRomPath, _ := filepath.Abs(*romPath)

			// Get relative path
			relPath, err := filepath.Rel(expectedBaseDir, absoluteRomPath)
			if err != nil {
				// Fallback: just use filename if path math fails
				relPath = filepath.Base(*romPath)
			}

			// 2. Format the argument string (remove extension, force forward slashes)
			argPath := strings.TrimSuffix(relPath, filepath.Ext(relPath))
			argPath = filepath.ToSlash(argPath)

			// 3. Generate Function Name
			funcNameSuffix := strings.ReplaceAll(argPath, "/", "__")
			// Clean up ".." (parent dir) to make it valid Go syntax
			funcNameSuffix = strings.ReplaceAll(funcNameSuffix, "..", "")
			// Clean up "." (current dir) just in case
			funcNameSuffix = strings.ReplaceAll(funcNameSuffix, ".", "_")

			funcName := "Test" + funcNameSuffix

			fmt.Printf("\n✅ MATCH FOUND AT FRAME: %d\n", f)
			fmt.Println("---------------------------------------------------")
			fmt.Println("Copy this function into your test file:")

			// 4. Print the Code Block
			fmt.Printf("\nfunc %s(t *testing.T) {\n", funcName)
			fmt.Printf("\trunPpuTest(t, \"%s\", %d, \"%s\")\n", argPath, f, targetHash)
			fmt.Printf("}\n\n")

			os.Exit(0)
		}
	}

	fmt.Printf("\n❌ No match found after %d frames.\n", *maxFrames)
	os.Exit(1)
}

func die(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}

func addColor(m map[color.RGBA]byte, hexStr string, val byte) {
	s := strings.TrimPrefix(hexStr, "#")
	v, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		die(fmt.Errorf("invalid hex color: %s", hexStr))
	}
	// We force Alpha to 0xFF (255) for consistency
	m[color.RGBA{R: uint8(v >> 16), G: uint8(v >> 8), B: uint8(v), A: 0xFF}] = val
}
