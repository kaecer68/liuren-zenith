# Changelog

所有對本專案的重要更改都將記錄在此文件中。

格式基於 [Keep a Changelog](https://keepachangelog.com/zh-TW/1.0.0/)，並且本專案遵循 [語義化版本](https://semver.org/lang/zh-TW/)。

## [Unreleased]

### Added
- 初始專案架構建立
- 大六壬排盤核心引擎
- 四柱八字計算（含時柱五鼠遁本地計算）
- 天地盤計算（月將加時）
- 十二天將佈局（晝夜貴人順逆佈局）
- 四課計算與顯示
- 三傳計算（九宗門取法）
- 十神系統完整實現（比肩、劫財、食神、傷官、正財、偏財、正官、七殺、正印、偏印）
- 智能建議算法（結合天將與十神）
- 課體判斷（伏吟、返吟、九宗門）
- 用神分析（按問題類型）
- 神煞計算（空亡、驛馬、桃花）
- REST API 服務
- gRPC 服務與 Protocol Buffers
- Web 前端界面
- 開源文檔（README、LICENSE、PRD、CONTRIBUTING）

### Changed
- 時柱計算改為本地實現（五鼠遁算法）

### Fixed
- 修復時柱始終為子時的問題
- 改進四課意義描述

## [1.0.0] - 2026-03-16

### Added
- 專案初始發布
- 完整的大六壬排盤功能
- 雙 API 接口支持（REST + gRPC）
- 現代化 Web 界面

---

## 版本說明

- `Added` - 新增功能
- `Changed` - 功能變更
- `Deprecated` - 即將移除的功能
- `Removed` - 已移除功能
- `Fixed` - 錯誤修復
- `Security` - 安全相關
