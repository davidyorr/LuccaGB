# LuccaGB

A Game Boy emulator written in Go and compiled to WebAssembly.

![LuccaGB Demo](https://github.com/user-attachments/assets/6d1e2af6-d65e-469e-86b2-9ac34beba13e)

## Input

| Game Boy | Keyboard             | Gamepad            |
| -------- | -------------------- | ------------------ |
| `B`      | <kbd>Z</kbd>         | <kbd>×</kbd>       |
| `A`      | <kbd>X</kbd>         | <kbd>○</kbd>       |
| `START`  | <kbd>Enter</kbd>     | <kbd>Options</kbd> |
| `SELECT` | <kbd>Backspace</kbd> | <kbd>Share</kbd>   |
| `UP`     | <kbd>↑</kbd>         | <kbd>↑</kbd>       |
| `DOWN`   | <kbd>↓</kbd>         | <kbd>↓</kbd>       |
| `LEFT`   | <kbd>←</kbd>         | <kbd>←</kbd>       |
| `RIGHT`  | <kbd>→</kbd>         | <kbd>→</kbd>       |

_Gamepad mappings use the standard [Gamepad API](https://w3c.github.io/gamepad/) button layout. Symbols above are shown using PlayStation terminology._

## Tests

Tests are run after each push to `main`. Auto-generated results are written to [TESTS.md](TESTS.md).
