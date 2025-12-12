# Test Results

| Suite | ROM       | Status |
|-------|-----------|--------|
| **Ungrouped** |  |  |
|  | `Test__lucca` | ✅ |
| **blargg** |  |  |
|  | `cpu_instrs.gb` | ⏩ |
|  | `instr_timing.gb` | ✅ |
|  | `mem_timing.gb` | ✅ |
|  | `halt_bug.gb` | ❌ |
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
| **mooneye / oam_dma** |  |  |
|  | `basic.gb` | ✅ |
|  | `reg_read.gb` | ✅ |
|  | `sources-GS.gb` | ❌ |
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
