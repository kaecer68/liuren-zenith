# 貢獻指南

感謝您對大六壬 Zenith 專案的關注！我們歡迎各種形式的貢獻。

## 📋 目錄

- [行為準則](#行為準則)
- [如何貢獻](#如何貢獻)
- [開發流程](#開發流程)
- [代碼規範](#代碼規範)
- [提交規範](#提交規範)

## 行為準則

請尊重所有參與者，保持友善和專業的態度。我們致力於建立一個開放、包容的社區。

## 如何貢獻

### 報告問題

如果您發現了 bug 或有功能建議，請通過 GitHub Issues 提交：

1. 檢查是否已有相關 issue
2. 使用 issue 模板提供詳細信息
3. 包括重現步驟（對於 bug）或詳細描述（對於功能建議）

### 提交代碼

1. Fork 本專案
2. 創建您的功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交您的更改 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 開啟 Pull Request

## 開發流程

### 環境設置

```bash
# 克隆專案
git clone https://github.com/yourusername/liuren-zenith.git
cd liuren-zenith

# 安裝依賴
go mod download

# 編譯
go build -o liuren-server ./cmd/server/main.go

# 運行測試
go test ./...
```

### 項目結構

```
liuren-zenith/
├── api/           # API 定義
├── cmd/           # 應用入口
├── pkg/           # 核心套件
├── web/           # 前端
└── docs/          # 文檔
```

### 主要模塊

- `pkg/liuren/` - 大六壬排盤核心算法
- `pkg/server/` - gRPC 服務實現
- `pkg/client/` - 外部服務客戶端
- `cmd/server/` - 服務器入口
- `web/` - 前端界面

## 代碼規範

### Go 代碼規範

- 遵循 [Effective Go](https://golang.org/doc/effective_go.html)
- 使用 `gofmt` 格式化代碼
- 編寫清晰的註釋（特別是導出函數）
- 保持函數簡短（建議 < 50 行）
- 錯誤處理要明確

### 命名規範

- 使用 CamelCase 命名導出類型/函數
- 使用 camelCase 命名未導出類型/函數
- 常量使用 PascalCase
- 接口名稱以 `er` 結尾（如 `Reader`, `Writer`）

### 註釋規範

```go
// Calculate 執行大六壬排盤計算
// 接收 DivinationRequest 參數，返回完整的排盤結果
func (e *Engine) Calculate(req DivinationRequest) (*DivinationResult, error) {
    // ...
}
```

## 提交規範

### Commit Message 格式

```
<類型>: <簡短描述>

<詳細描述（可選）>

<頁腳（可選）>
```

### 類型說明

- `feat` - 新功能
- `fix` - Bug 修復
- `docs` - 文檔更新
- `style` - 代碼格式調整（不影響功能）
- `refactor` - 代碼重構
- `perf` - 性能優化
- `test` - 測試相關
- `chore` - 構建/工具相關

### 示例

```
feat: 添加十神計算函數

實現完整的十神系統：
- 比肩、劫財（同我）
- 食神、傷官（我生）
- 正財、偏財（我剋）
- 正官、七殺（剋我）
- 正印、偏印（生我）

包含陰陽五行雙重判斷邏輯。
```

## 測試要求

- 新功能必須包含測試
- 保持測試覆蓋率 > 80%
- 使用表格驅動測試

```go
func TestCalculateShiShen(t *testing.T) {
    tests := []struct {
        dayStem Stem
        target  Branch
        want    string
    }{
        {Jia, Zi, "正印"},
        {Yi, Chou, "偏印"},
        // ...
    }
    
    for _, tt := range tests {
        t.Run(tt.want, func(t *testing.T) {
            got := engine.getShiShen(tt.dayStem, tt.target)
            if got != tt.want {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

## 文檔要求

- 更新 README.md（如影響使用方式）
- 更新 PRD.md（如添加新功能）
- 添加/更新代碼註釋

## 審查流程

1. Pull Request 必須通過 CI 檢查
2. 至少一名維護者審查批准
3. 保持提交歷史整潔（建議 rebase）

## 聯繫方式

- GitHub Issues: [問題報告](https://github.com/yourusername/liuren-zenith/issues)
- Discussions: [功能討論](https://github.com/yourusername/liuren-zenith/discussions)

## 許可證

通過提交代碼，您同意將您的貢獻置於本專案的 MIT 許可證下。

---

再次感謝您的貢獻！🙏
