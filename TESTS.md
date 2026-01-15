# Test Results

| Suite | ROM       | Status |
|-------|-----------|--------|
| **blargg** |  |  |
|  | `cpu_instrs.gb` | ✅ |
|  | `dmg_sound.gb` | ❌ |
|  | `halt_bug.gb` | ❌ |
|  | `instr_timing.gb` | ✅ |
|  | `mem_timing.gb` | ✅ |
|  | `oam_bug.gb` | ❌ |
| **mooneye** |  |  |
|  | `add_sp_e_timing.gb` | ✅ |
|  | `call_cc_timing.gb` | ✅ |
|  | `call_cc_timing2.gb` | ✅ |
|  | `call_timing.gb` | ✅ |
|  | `call_timing2.gb` | ✅ |
|  | `di_timing-GS.gb` | ✅ |
|  | `div_timing.gb` | ✅ |
|  | `ei_sequence.gb` | ✅ |
|  | `ei_timing.gb` | ✅ |
|  | `halt_ime0_ei.gb` | ✅ |
|  | `halt_ime0_nointr_timing.gb` | ✅ |
|  | `halt_ime1_timing.gb` | ✅ |
|  | `halt_ime1_timing2-GS.gb` | ✅ |
|  | `if_ie_registers.gb` | ✅ |
|  | `intr_timing.gb` | ✅ |
|  | `jp_cc_timing.gb` | ✅ |
|  | `jp_timing.gb` | ✅ |
|  | `ld_hl_sp_e_timing.gb` | ✅ |
|  | `oam_dma_restart.gb` | ✅ |
|  | `oam_dma_start.gb` | ✅ |
|  | `oam_dma_timing.gb` | ✅ |
|  | `pop_timing.gb` | ✅ |
|  | `push_timing.gb` | ✅ |
|  | `rapid_di_ei.gb` | ✅ |
|  | `ret_cc_timing.gb` | ✅ |
|  | `ret_timing.gb` | ✅ |
|  | `reti_intr_timing.gb` | ✅ |
|  | `reti_timing.gb` | ✅ |
|  | `rst_timing.gb` | ✅ |
| **mooneye / bits** |  |  |
|  | `mem_oam.gb` | ✅ |
|  | `reg_f.gb` | ✅ |
|  | `unused_hwio-GS.gb` | ✅ |
| **mooneye / instr** |  |  |
|  | `daa.gb` | ✅ |
| **mooneye / interrupts** |  |  |
|  | `ie_push.gb` | ✅ |
| **mooneye / mbc1** |  |  |
|  | `bits_bank1.gb` | ✅ |
|  | `bits_bank2.gb` | ✅ |
|  | `bits_mode.gb` | ✅ |
|  | `bits_ramg.gb` | ✅ |
|  | `multicart_rom_8Mb.gb` | ❌ |
|  | `ram_256kb.gb` | ✅ |
|  | `ram_64kb.gb` | ✅ |
|  | `rom_16Mb.gb` | ✅ |
|  | `rom_1Mb.gb` | ✅ |
|  | `rom_2Mb.gb` | ✅ |
|  | `rom_4Mb.gb` | ✅ |
|  | `rom_512kb.gb` | ✅ |
|  | `rom_8Mb.gb` | ✅ |
| **mooneye / mbc2** |  |  |
|  | `bits_ramg.gb` | ✅ |
|  | `bits_romb.gb` | ✅ |
|  | `bits_unused.gb` | ✅ |
|  | `ram.gb` | ✅ |
|  | `rom_1Mb.gb` | ✅ |
|  | `rom_2Mb.gb` | ✅ |
|  | `rom_512kb.gb` | ✅ |
| **mooneye / mbc5** |  |  |
|  | `rom_16Mb.gb` | ✅ |
|  | `rom_1Mb.gb` | ✅ |
|  | `rom_2Mb.gb` | ✅ |
|  | `rom_32Mb.gb` | ✅ |
|  | `rom_4Mb.gb` | ✅ |
|  | `rom_512kb.gb` | ✅ |
|  | `rom_64Mb.gb` | ✅ |
|  | `rom_8Mb.gb` | ✅ |
| **mooneye / oam_dma** |  |  |
|  | `basic.gb` | ✅ |
|  | `reg_read.gb` | ✅ |
|  | `sources-GS.gb` | ✅ |
| **mooneye / ppu** |  |  |
|  | `hblank_ly_scx_timing-GS.gb` | ✅ |
|  | `intr_1_2_timing-GS.gb` | ✅ |
|  | `intr_2_0_timing.gb` | ✅ |
|  | `intr_2_mode0_timing_sprites.gb` | ❌ |
|  | `intr_2_mode0_timing.gb` | ✅ |
|  | `intr_2_mode3_timing.gb` | ✅ |
|  | `intr_2_oam_ok_timing.gb` | ✅ |
|  | `stat_irq_blocking.gb` | ✅ |
|  | `vblank_stat_intr-GS.gb` | ✅ |
| **mooneye / timer** |  |  |
|  | `div_write.gb` | ✅ |
|  | `rapid_toggle.gb` | ✅ |
|  | `tim00_div_trigger.gb` | ✅ |
|  | `tim00.gb` | ✅ |
|  | `tim01_div_trigger.gb` | ✅ |
|  | `tim01.gb` | ✅ |
|  | `tim10_div_trigger.gb` | ✅ |
|  | `tim10.gb` | ✅ |
|  | `tim11_div_trigger.gb` | ✅ |
|  | `tim11.gb` | ✅ |
|  | `tima_reload.gb` | ✅ |
|  | `tima_write_reloading.gb` | ✅ |
|  | `tma_write_reloading.gb` | ✅ |

# Screenshot Test Results

| Suite | ROM       | Status |
|-------|-----------|--------|
| **boop** |  |  |
|  | `solid-color-0-background.gb` | ✅ |
|  | `solid-color-0-window.gb` | ✅ |
|  | `solid-color-1-background.gb` | ✅ |
|  | `solid-color-1-window.gb` | ✅ |
|  | `solid-color-2-background.gb` | ✅ |
|  | `solid-color-2-window.gb` | ✅ |
|  | `solid-color-3-background.gb` | ✅ |
|  | `solid-color-3-window.gb` | ✅ |
|  | `sprite-8x8.gb` | ✅ |
|  | `sprite-8x16.gb` | ✅ |
| **dmg_acid2** |  |  |
|  | `dmg-acid2.gb` | ✅ |
| **mealybug_tearoom** |  |  |
|  | `m2_win_en_toggle.gb` | ✅ |
|  | `m3_bgp_change.gb` | ❌ |
|  | `m3_bgp_change_sprites.gb` | ❌ |
|  | `m3_lcdc_bg_en_change.gb` | ❌ |
|  | `m3_lcdc_bg_map_change.gb` | ❌ |
|  | `m3_lcdc_obj_en_change.gb` | ❌ |
|  | `m3_lcdc_obj_en_change_variant.gb` | ❌ |
|  | `m3_lcdc_obj_size_change.gb` | ❌ |
|  | `m3_lcdc_obj_size_change_scx.gb` | ❌ |
|  | `m3_lcdc_tile_sel_change.gb` | ❌ |
|  | `m3_lcdc_tile_sel_win_change.gb` | ❌ |
|  | `m3_lcdc_win_en_change_multiple.gb` | ❌ |
|  | `m3_lcdc_win_en_change_multiple_wx.gb` | ❌ |
|  | `m3_lcdc_win_map_change.gb` | ❌ |
|  | `m3_obp0_change.gb` | ❌ |
|  | `m3_scx_high_5_bits.gb` | ❌ |
|  | `m3_scx_low_3_bits.gb` | ❌ |
|  | `m3_scy_change.gb` | ❌ |
|  | `m3_window_timing.gb` | ❌ |
|  | `m3_window_timing_wx_0.gb` | ❌ |
|  | `m3_wx_4_change.gb` | ❌ |
|  | `m3_wx_4_change_sprites.gb` | ❌ |
|  | `m3_wx_5_change.gb` | ❌ |
|  | `m3_wx_6_change.gb` | ❌ |
| **mooneye** |  |  |
|  | `sprite_priority.gb` | ✅ |
| **other** |  |  |
|  | `lucca.gb` | ✅ |

