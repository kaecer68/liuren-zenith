# 大六壬 Zenith 🔮

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![API](https://img.shields.io/badge/API-REST%20%7C%20gRPC-orange.svg)](docs/api.md)

> 三式之首 · 百占百驗 · 來人不用問

大六壬 Zenith 是一個現代化的大六壬排盤系統，提供精準的占卜推演與詳盡的解說分析。

## ✨ 核心功能

### 🔮 占卜推演
- **四柱八字** - 年月日時四柱干支計算
- **天地盤** - 月將加時起天盤，十二天將佈局
- **四課** - 日干支之寄宮，含十神與課之意義
- **三傳** - 九宗門取法，初傳/中傳/末傳詳解
- **課體判斷** - 伏吟、返吟、九宗門課體識別
- **用神分析** - 根據問題類型自動推斷用神

### 🎯 智能解說
- **天將 + 十神組合算法** - 提供具體化建議
- **階段性分析** - 初傳（發端）、中傳（過程）、末傳（結果）
- **神煞標註** - 空亡、驛馬、桃花自動標識
- **吉凶判斷** - 綜合課體與傳位評估

### 🌐 服務接口
- **REST API** - HTTP JSON 接口（端口 8081）
- **gRPC 服務** - 高效 protobuf 接口（端口 50052）
- **Web 前端** - 現代化響面設計

## 🚀 快速開始

### 環境要求
- Go 1.21+
- 可選：lunar-zenith 曆法服務

### 安裝運行

```bash
# 克隆專案
git clone https://github.com/yourusername/liuren-zenith.git
cd liuren-zenith

# 編譯服務器
go build -o liuren-server ./cmd/server/main.go

# 啟動服務（REST + gRPC）
./liuren-server
```

### API 使用範例

**REST API:**
```bash
curl -X POST http://localhost:8081/api/v1/divination \
  -H "Content-Type: application/json" \
  -d '{
    "date": "2026-03-16",
    "time": "14:00",
    "question": "問財運",
    "question_type": "財運"
  }'
```

**gRPC:**
```bash
go run ./cmd/grpc-client/main.go
```

## 📁 項目結構

```
liuren-zenith/
├── api/                    # API 定義
│   └── proto/             # Protocol Buffers
├── cmd/                   # 應用入口
│   ├── server/           # REST + gRPC 服務器
│   └── grpc-client/      # gRPC 測試客戶端
├── pkg/                   # 核心套件
│   ├── liuren/           # 大六壬引擎
│   │   ├── engine.go     # 排盤邏輯
│   │   ├── types.go      # 類型定義
│   │   ├── jiuzongmen.go # 九宗門取法
│   │   └── calculator.go # 曆法計算
│   ├── client/           # 外部服務客戶端
│   └── server/           # gRPC 服務器實現
├── web/                   # 前端界面
│   └── index.html        # 排盤頁面
├── configs/              # 配置檔案
└── docs/                 # 文檔
    ├── PRD.md            # 產品需求文檔
    └── api.md            # API 文檔
```

## 🛠️ 技術棧

- **後端**: Go 1.21+
- **通信**: REST (gin) / gRPC (protobuf)
- **曆法**: lunar-zenith 服務（可選）
- **前端**: HTML5 + CSS3 + JavaScript

## 📝 算法特性

### 十神系統
取代傳統六親，更細緻區分陰陽：
- 比肩/劫財（同我）
- 食神/傷官（我生）
- 正財/偏財（我剋）
- 正官/七殺（剋我）
- 正印/偏印（生我）

### 智能建議算法
結合天將（貴人、青龍、白虎等）與十神，生成階段性建議：
- 初傳：時機判斷（宜立即行動/需謹慎準備）
- 中傳：過程評估（順遂/有阻/需調整）
- 末傳：結果預測（圓滿/防損/平穩）

## 🤝 貢獻指南

歡迎提交 Issue 和 Pull Request！請參考 [CONTRIBUTING.md](CONTRIBUTING.md)。

## 📜 開源許可

本專案採用 [MIT License](LICENSE) 開源許可證。

## 🙏 致謝

- 大六壬傳統典籍《課經》《畢法賦》
- 曆法計算參考 lunar-zenith 專案

[德凱/KAECER](https://github.com/kaecer68) 
- 對傳統文化數位化有興趣的前端工程師
- Blog: https://goluck.im/
- Twitter: [@kaecer](https://twitter.com/kaecer)

相關專案：
- [lunar-zenith](https://github.com/kaecer68/lunar-zenith) - 高精度農曆節氣 API
- [ziwei-zenith](https://github.com/kaecer68/ziwei-zenith) - 紫微斗數排盤 API
- [liuren-zenith](https://github.com/kaecer68/liuren-zenith) - 六壬排盤 API
- [bazi-zenith](https://github.com/kaecer68/bazi-zenith) - 八字排盤 API


---

<div align="center">
  <sub>Built with ❤️ by Zenith Team</sub>
</div>
