package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kaecer68/liuren-zenith/api/proto"
	"github.com/kaecer68/liuren-zenith/pkg/client"
	"github.com/kaecer68/liuren-zenith/pkg/liuren"
	"github.com/kaecer68/liuren-zenith/pkg/server"
	"google.golang.org/grpc"
)

// LiurenHandler 處理大六壬請求
type LiurenHandler struct {
	Engine *liuren.Engine
}

// DivinationRequest HTTP 請求結構
type DivinationRequest struct {
	Date         string `json:"date"`          // 格式: YYYY-MM-DD
	Time         string `json:"time"`          // 格式: HH:MM (可選，默認 12:00)
	Question     string `json:"question"`      // 占事（可選）
	QuestionType string `json:"question_type"` // 問題類型（可選）
}

// DivinationResponse HTTP 響應結構
type DivinationResponse struct {
	Success bool                     `json:"success"`
	Message string                   `json:"message,omitempty"`
	Data    *liuren.DivinationResult `json:"data,omitempty"`
}

// NewLiurenHandler 創建處理器
func NewLiurenHandler() *LiurenHandler {
	// 使用 lunar-zenith 服務獲取精確曆法數據
	dataSource := client.NewLunarZenithClient("")
	engine := liuren.NewEngine(dataSource)

	return &LiurenHandler{
		Engine: engine,
	}
}

// HandleDivination 處理排盤請求
func (h *LiurenHandler) HandleDivination(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "只支援 POST 方法")
		return
	}

	var req DivinationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "請求格式錯誤: "+err.Error())
		return
	}

	// 解析時間
	divinationTime, err := parseTime(req.Date, req.Time)
	if err != nil {
		writeError(w, http.StatusBadRequest, "時間格式錯誤: "+err.Error())
		return
	}

	// 執行排盤
	result, err := h.Engine.Calculate(liuren.DivinationRequest{
		Time:         divinationTime,
		Question:     req.Question,
		QuestionType: req.QuestionType,
	})

	if err != nil {
		writeError(w, http.StatusInternalServerError, "排盤失敗: "+err.Error())
		return
	}

	resp := DivinationResponse{
		Success: true,
		Data:    result,
	}

	json.NewEncoder(w).Encode(resp)
}

// HandleHealth 健康檢查
func (h *LiurenHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	resp := map[string]interface{}{
		"status":    "ok",
		"service":   "liuren-zenith",
		"version":   "v1.0.0",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	json.NewEncoder(w).Encode(resp)
}

// corsMiddleware 添加跨域支援
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}
func parseTime(dateStr, timeStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, fmt.Errorf("日期不能為空")
	}

	if timeStr == "" {
		timeStr = "12:00"
	}

	layout := "2006-01-02 15:04"
	return time.Parse(layout, dateStr+" "+timeStr)
}

// writeError 寫入錯誤響應
func writeError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(DivinationResponse{
		Success: false,
		Message: message,
	})
}

func resolvePort(contractEnvKey, legacyEnvKey string) string {
	if value := os.Getenv(contractEnvKey); value != "" {
		return value
	}
	if value := os.Getenv(legacyEnvKey); value != "" {
		return value
	}

	envFile, err := os.Open(filepath.Clean(".env.ports"))
	if err == nil {
		defer envFile.Close()
		scanner := bufio.NewScanner(envFile)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if key == contractEnvKey || key == legacyEnvKey {
				return value
			}
		}
		if scanErr := scanner.Err(); scanErr == nil {
			log.Fatalf("missing runtime port configuration for %s/%s", contractEnvKey, legacyEnvKey)
		} else {
			log.Fatalf("failed to read .env.ports: %v", scanErr)
		}
	}
	if !errors.Is(err, os.ErrNotExist) {
		log.Fatalf("failed to read .env.ports: %v", err)
	}

	log.Fatalf("missing runtime port configuration for %s/%s", contractEnvKey, legacyEnvKey)
	return ""
}

func main() {
	handler := NewLiurenHandler()
	restPort := resolvePort("LIUREN_REST_PORT", "REST_PORT")
	grpcPort := resolvePort("LIUREN_GRPC_PORT", "GRPC_PORT")

	// 啟動 gRPC 服務（在後台運行）
	go startGRPCServer(grpcPort)

	// 註冊 REST 路由
	http.HandleFunc("/health", corsMiddleware(handler.HandleHealth))
	http.HandleFunc("/api/v1/divination", corsMiddleware(handler.HandleDivination))

	// 提供靜態網頁 - 使用絕對路徑
	execPath, err := os.Executable()
	if err != nil {
		log.Printf("Warning: 無法獲取執行路徑: %v", err)
		execPath = "."
	}
	execDir := filepath.Dir(execPath)
	webDir := filepath.Join(execDir, "web")

	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	// 啟動 REST 服務
	log.Printf("Liuren-Zenith REST 服務啟動於 http://localhost:%s", restPort)
	log.Printf("Liuren-Zenith gRPC 服務啟動於 localhost:%s", grpcPort)
	log.Printf("網頁界面: http://localhost:%s/", restPort)
	log.Printf("REST API: POST /api/v1/divination")
	log.Printf("gRPC 服務: LiurenInfoService")
	log.Printf("健康檢查: GET /health")

	if err := http.ListenAndServe(":"+restPort, nil); err != nil {
		log.Fatal("REST 服務啟動失敗: ", err)
	}
}

// startGRPCServer 啟動 gRPC 信息調用服務
func startGRPCServer(grpcPort string) {
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("gRPC 監聽失敗: %v", err)
	}

	s := grpc.NewServer()
	infoServer := server.NewInfoServer()
	proto.RegisterLiurenInfoServiceServer(s, infoServer)

	log.Printf("gRPC 信息調用服務已啟動於 :%s", grpcPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("gRPC 服務啟動失敗: %v", err)
	}
}
