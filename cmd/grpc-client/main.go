package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kaecer68/liuren-zenith/api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 連接 gRPC 服務
	conn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("連接失敗: %v", err)
	}
	defer conn.Close()

	client := proto.NewLiurenInfoServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	fmt.Println("=== 大六壬 gRPC 信息調用服務測試 ===\n")

	// 1. 測試查詢課體信息
	fmt.Println("1. 查詢所有課體信息:")
	keTiResp, err := client.GetKeTiInfo(ctx, &proto.GetKeTiInfoRequest{})
	if err != nil {
		log.Printf("查詢課體失敗: %v", err)
	} else {
		for _, keTi := range keTiResp.KeTiList {
			fmt.Printf("   - %s (%s): %s\n", keTi.Name, keTi.Category, keTi.Description)
		}
	}

	// 2. 測試查詢特定課體
	fmt.Println("\n2. 查詢特定課體（伏吟課）:")
	keTiSingle, err := client.GetKeTiInfo(ctx, &proto.GetKeTiInfoRequest{KeTiName: "伏吟課"})
	if err != nil {
		log.Printf("查詢失敗: %v", err)
	} else if len(keTiSingle.KeTiList) > 0 {
		k := keTiSingle.KeTiList[0]
		fmt.Printf("   名稱: %s\n", k.Name)
		fmt.Printf("   類別: %s\n", k.Category)
		fmt.Printf("   描述: %s\n", k.Description)
		fmt.Printf("   占斷: %s\n", k.Meaning)
	}

	// 3. 測試查詢神煞信息
	fmt.Println("\n3. 查詢所有神煞:")
	shenShaResp, err := client.GetShenShaInfo(ctx, &proto.GetShenShaInfoRequest{})
	if err != nil {
		log.Printf("查詢神煞失敗: %v", err)
	} else {
		for _, ss := range shenShaResp.ShenShaList {
			fmt.Printf("   - %s: %s\n", ss.Name, ss.Description)
		}
	}

	// 4. 測試查詢天將信息
	fmt.Println("\n4. 查詢十二天將:")
	tjResp, err := client.GetTianJiangInfo(ctx, &proto.GetTianJiangInfoRequest{})
	if err != nil {
		log.Printf("查詢天將失敗: %v", err)
	} else {
		for _, tj := range tjResp.TianJiangList {
			fmt.Printf("   - %s (%s): %s\n", tj.Name, tj.Nature, tj.Description)
		}
	}

	// 5. 測試查詢月將信息
	fmt.Println("\n5. 查詢十二月將:")
	mgResp, err := client.GetMonthGeneralInfo(ctx, &proto.GetMonthGeneralInfoRequest{})
	if err != nil {
		log.Printf("查詢月將失敗: %v", err)
	} else {
		for _, mg := range mgResp.MonthGeneralList {
			fmt.Printf("   - %s (%s): %s\n", mg.Month, mg.Name, mg.Description)
		}
	}

	// 6. 測試查詢歷史記錄
	fmt.Println("\n6. 查詢歷史排盤記錄:")
	historyResp, err := client.QueryDivinationHistory(ctx, &proto.QueryDivinationHistoryRequest{Limit: 5})
	if err != nil {
		log.Printf("查詢歷史失敗: %v", err)
	} else {
		fmt.Printf("   總數: %d\n", historyResp.Total)
		for _, record := range historyResp.Records {
			fmt.Printf("   - [%s] %s %s | 課體: %s | 問題: %s\n",
				record.Id, record.Date, record.Time, record.KeTi, record.Question)
		}
	}

	fmt.Println("\n=== 測試完成 ===")
}
