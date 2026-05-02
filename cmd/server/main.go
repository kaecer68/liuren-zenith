package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	liurenpb "github.com/kaecer68/liuren-zenith/gen/liurenpb"
	"github.com/kaecer68/liuren-zenith/internal/httpapi"
	"github.com/kaecer68/liuren-zenith/internal/runtimeconfig"
	"github.com/kaecer68/liuren-zenith/pkg/client"
	"github.com/kaecer68/liuren-zenith/pkg/liuren"
	"github.com/kaecer68/liuren-zenith/pkg/server"
	"google.golang.org/grpc"
)

var serviceVersion = "dev"

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

func main() {
	dataSource, err := client.NewLunarZenithClientWithFallback()
	if err != nil {
		log.Fatalf("創建 lunar-zenith 客戶端失敗: %v", err)
	}
	handler := &httpapi.Handler{
		Engine:  liuren.NewEngine(dataSource),
		Version: serviceVersion,
	}
	restPort := mustResolvePort("LIUREN_REST_PORT", "REST_PORT")
	grpcPort := mustResolvePort("LIUREN_GRPC_PORT", "GRPC_PORT")

	// 啟動 gRPC 服務（在後台運行）
	go startGRPCServer(grpcPort)

	// 註冊 REST 路由
	http.HandleFunc("/health", corsMiddleware(handler.HandleHealth))
	http.HandleFunc("/v1/liuren/calculate", corsMiddleware(handler.HandleCalculate))
	http.HandleFunc("/v1/liuren/analyze", corsMiddleware(handler.HandleAnalyze))
	http.HandleFunc("/api/v1/divination", corsMiddleware(handler.HandleDivination))

	webDir := "./web"
	if info, err := os.Stat(webDir); err != nil || !info.IsDir() {
		if err != nil && !os.IsNotExist(err) {
			log.Printf("Warning: 無法存取 web 目錄: %v", err)
		}
		execPath, err := os.Executable()
		if err != nil {
			log.Printf("Warning: 無法獲取執行路徑: %v", err)
			execPath = "."
		}
		execDir := filepath.Dir(execPath)
		webDir = filepath.Join(execDir, "web")
	}

	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	// 啟動 REST 服務
	log.Printf("Liuren-Zenith REST 服務啟動於 :%s", restPort)
	log.Printf("Liuren-Zenith gRPC 服務啟動於 :%s", grpcPort)
	log.Printf("網頁界面已掛載於 /")
	log.Printf("REST API: POST /api/v1/divination")
	log.Printf("Contract API: POST /v1/liuren/calculate, POST /v1/liuren/analyze")
	log.Printf("gRPC 服務: LiurenInfoService")
	log.Printf("健康檢查: GET /health")

	if err := http.ListenAndServe(":"+restPort, nil); err != nil {
		log.Fatal("REST 服務啟動失敗: ", err)
	}
}

func mustResolvePort(contractEnvKey, legacyEnvKey string) string {
	defer func() {
		if recovered := recover(); recovered != nil {
			log.Fatalf("runtime config error: %v", recovered)
		}
	}()
	return runtimeconfig.MustLookupPort(contractEnvKey, legacyEnvKey)
}

// startGRPCServer 啟動 gRPC 信息調用服務
func startGRPCServer(grpcPort string) {
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("gRPC 監聽失敗: %v", err)
	}

	s := grpc.NewServer()
	infoServer := server.NewInfoServer()
	liurenpb.RegisterLiurenInfoServiceServer(s, infoServer)

	log.Printf("gRPC 信息調用服務已啟動於 :%s", grpcPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("gRPC 服務啟動失敗: %v", err)
	}
}
