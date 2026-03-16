package server

import (
	"context"

	"github.com/kaecer68/liuren-zenith/api/proto"
)

// InfoServer gRPC 信息調用服務實現
type InfoServer struct {
	proto.UnimplementedLiurenInfoServiceServer
}

// NewInfoServer 創建信息服務器
func NewInfoServer() *InfoServer {
	return &InfoServer{}
}

// GetKeTiInfo 查詢課體信息
func (s *InfoServer) GetKeTiInfo(ctx context.Context, req *proto.GetKeTiInfoRequest) (*proto.GetKeTiInfoResponse, error) {
	// 課體信息庫
	keTiDB := map[string]*proto.KeTiInfo{
		"伏吟課": {
			Name:        "伏吟課",
			Category:    "特殊課體",
			Description: "天盤地盤相同，如人伏地呻吟，故名伏吟",
			Meaning:     "主事多滯礙，遲疑不決，靜守為宜，動則有悔",
			Features:    []string{"天盤地盤完全重合", "月將加時後不動", "事多反覆，難以速決"},
			Examples:    []string{"《畢法賦》：伏吟課，動不如靜", "《大六壬課經》：伏吟者，呻吟不動也"},
		},
		"返吟課": {
			Name:        "返吟課",
			Category:    "特殊課體",
			Description: "天盤地盤相對，如人返覆呻吟，故名返吟",
			Meaning:     "主事多反覆，來去不定，速則有成，遲則生變",
			Features:    []string{"天盤地盤完全相對", "月將加時後全動", "事多變化，速決為宜"},
			Examples:    []string{"《畢法賦》：返吟課，速決為宜", "《大六壬課經》：返吟者，反覆不定也"},
		},
		"進茹課": {
			Name:        "進茹課",
			Category:    "茹格",
			Description: "三傳依次進行，如人進步之狀",
			Meaning:     "主事漸進，宜進不宜退，循序渐进則有成",
			Features:    []string{"三傳順行", "事情漸進", "宜進取"},
			Examples:    []string{"《畢法賦》：進茹課，循序渐进"},
		},
		"退茹課": {
			Name:        "退茹課",
			Category:    "茹格",
			Description: "三傳依次退行，如人退縮之狀",
			Meaning:     "主事漸退，宜退不宜進，退守為安",
			Features:    []string{"三傳逆行", "事情漸退", "宜退守"},
			Examples:    []string{"《畢法賦》：退茹課，退守為安"},
		},
		"三奇課": {
			Name:        "三奇課",
			Category:    "奇格",
			Description: "三傳見乙丙丁三奇",
			Meaning:     "主有奇緣奇遇，凡事逢凶化吉，遇難呈祥",
			Features:    []string{"三傳見乙丙丁", "奇遇奇緣", "逢凶化吉"},
			Examples:    []string{"《畢法賦》：三奇臨課，百禍消散"},
		},
		"別責課": {
			Name:        "別責課",
			Category:    "特殊課體",
			Description: "四課無克，別責取用神",
			Meaning:     "主事無頭緒，須借他人之力，依附而行",
			Features:    []string{"四課無克", "別責取用神", "依附他人"},
			Examples:    []string{"《畢法賦》：別責課，依附而行"},
		},
		"八專課": {
			Name:        "八專課",
			Category:    "特殊課體",
			Description: "日干同位，八專之日",
			Meaning:     "主事專一，但易生執念，宜變通",
			Features:    []string{"日干同位", "八專之日", "專一執著"},
			Examples:    []string{"《畢法賦》：八專課，專一而行"},
		},
	}

	resp := &proto.GetKeTiInfoResponse{}

	if req.KeTiName != "" {
		// 查詢特定課體
		if keTi, ok := keTiDB[req.KeTiName]; ok {
			resp.KeTiList = append(resp.KeTiList, keTi)
		}
	} else {
		// 返回所有課體
		for _, keTi := range keTiDB {
			resp.KeTiList = append(resp.KeTiList, keTi)
		}
	}

	return resp, nil
}

// GetShenShaInfo 查詢神煞信息
func (s *InfoServer) GetShenShaInfo(ctx context.Context, req *proto.GetShenShaInfoRequest) (*proto.GetShenShaInfoResponse, error) {
	shenShaDB := map[string]*proto.ShenShaInfo{
		"天乙貴人": {
			Name:        "天乙貴人",
			Description: "諸神之首，主貴氣、權威、扶助",
			Meaning:     "凡事宜求助貴人，逢凶化吉，遇難呈祥",
			Calculation: "甲戊庚牛羊，乙己鼠猴鄉，丙丁豬雞位，壬癸蛇兔藏，六辛逢虎馬，此是貴人方",
		},
		"驛馬": {
			Name:        "驛馬",
			Description: "主出行、遷移、變動",
			Meaning:     "宜出行、遷移、變動，主快速、奔波",
			Calculation: "申子辰馬在寅，巳酉丑馬在亥，寅午戌馬在申，亥卯未馬在巳",
		},
		"桃花": {
			Name:        "桃花",
			Description: "主人緣、異性緣、感情",
			Meaning:     "主人緣好，異性緣佳，但防桃花劫",
			Calculation: "申子辰桃花在酉，巳酉丑桃花在午，寅午戌桃花在卯，亥卯未桃花在子",
		},
		"空亡": {
			Name:        "空亡",
			Description: "主虛無、不實、落空",
			Meaning:     "事情易落空，虛而不實，但宜出家、修道",
			Calculation: "按旬查空亡，如甲子旬空亡戌亥",
		},
		"文昌": {
			Name:        "文昌",
			Description: "主學業、文書、智慧",
			Meaning:     "宜學習、考試、文書之事，主聰明智慧",
			Calculation: "甲蛇乙馬丙戊猴，丁己雞宮庚豬遊，辛犬癸虎是文昌",
		},
		"華蓋": {
			Name:        "華蓋",
			Description: "主孤獨、清高、出世",
			Meaning:     "性格孤高，宜出家、修行，防孤獨",
			Calculation: "申子辰華蓋在辰，巳酉丑華蓋在丑，寅午戌華蓋在戌，亥卯未華蓋在未",
		},
	}

	resp := &proto.GetShenShaInfoResponse{}

	if req.ShenShaName != "" {
		if shenSha, ok := shenShaDB[req.ShenShaName]; ok {
			resp.ShenShaList = append(resp.ShenShaList, shenSha)
		}
	} else {
		for _, shenSha := range shenShaDB {
			resp.ShenShaList = append(resp.ShenShaList, shenSha)
		}
	}

	return resp, nil
}

// GetLiuQinInfo 查詢六親關係
func (s *InfoServer) GetLiuQinInfo(ctx context.Context, req *proto.GetLiuQinInfoRequest) (*proto.GetLiuQinInfoResponse, error) {
	// 六親關係說明
	liuQinDB := map[string]*proto.GetLiuQinInfoResponse{
		"比肩": {
			Relation:    "比肩",
			Description: "與日干同五行同陰陽",
			Meaning:     "主兄弟姐妹、同事、競爭者，宜合作但防爭鬥",
		},
		"劫財": {
			Relation:    "劫財",
			Description: "與日干同五行異陰陽",
			Meaning:     "主破財、損耗、爭奪，宜謹慎理財",
		},
		"食神": {
			Relation:    "食神",
			Description: "日干所生同陰陽",
			Meaning:     "主福氣、享受、口才、子孫，宜享受生活",
		},
		"傷官": {
			Relation:    "傷官",
			Description: "日干所生異陰陽",
			Meaning:     "主才華、創意、叛逆，宜發揮才華但防口舌",
		},
		"偏財": {
			Relation:    "偏財",
			Description: "日干所克同陰陽",
			Meaning:     "主意外之財、流動資產、父親，宜投資但防風險",
		},
		"正財": {
			Relation:    "正財",
			Description: "日干所克異陰陽",
			Meaning:     "主正當收入、穩定財源、妻子，宜勤奮工作",
		},
		"偏官": {
			Relation:    "偏官",
			Description: "克日干同陰陽",
			Meaning:     "主壓力、挑戰、小人，宜謹慎應對",
		},
		"正官": {
			Relation:    "正官",
			Description: "克日干異陰陽",
			Meaning:     "主權威、地位、上司、丈夫，宜順從規矩",
		},
		"偏印": {
			Relation:    "偏印",
			Description: "生日干同陰陽",
			Meaning:     "主學習、偏門學問、繼母，宜學習新知",
		},
		"正印": {
			Relation:    "正印",
			Description: "生日干異陰陽",
			Meaning:     "主學業、文書、母親、貴人，宜進修學習",
		},
	}

	// 如果提供了干支，計算具體關係
	if req.Stem != "" && req.Branch != "" {
		// 這裡簡化處理，實際應根據干支計算
		relation := "待計算"
		return &proto.GetLiuQinInfoResponse{
			Relation:    relation,
			Description: "日干 " + req.Stem + " 與地支 " + req.Branch + " 的關係",
			Meaning:     "請參考具體排盤結果",
		}, nil
	}

	// 返回通用說明（取第一個）
	for _, info := range liuQinDB {
		return info, nil
	}

	return &proto.GetLiuQinInfoResponse{}, nil
}

// GetTianJiangInfo 查詢十二天將信息
func (s *InfoServer) GetTianJiangInfo(ctx context.Context, req *proto.GetTianJiangInfoRequest) (*proto.GetTianJiangInfoResponse, error) {
	tianJiangDB := map[string]*proto.TianJiangInfo{
		"貴人": {
			Name:        "天乙貴人",
			Description: "諸神之首，主貴氣、權威",
			Nature:      "大吉",
			Meaning:     "凡事順利，逢凶化吉，有貴人扶助",
		},
		"螣蛇": {
			Name:        "螣蛇",
			Description: "主虛驚、怪異、纏繞",
			Nature:      "凶",
			Meaning:     "主虛驚不安，事情纏繞不清，但亦主靈異",
		},
		"朱雀": {
			Name:        "朱雀",
			Description: "主口舌、文書、訴訟",
			Nature:      "平",
			Meaning:     "主口舌是非，文書之事，訴訟之憂",
		},
		"六合": {
			Name:        "六合",
			Description: "主和合、婚姻、合作",
			Nature:      "吉",
			Meaning:     "主和合美滿，婚姻順利，合作愉快",
		},
		"勾陳": {
			Name:        "勾陳",
			Description: "主牽連、遲滯、田土",
			Nature:      "凶",
			Meaning:     "事情牽連遲滯，田土之憂，防小人暗算",
		},
		"青龍": {
			Name:        "青龍",
			Description: "主喜慶、財帛、婚姻",
			Nature:      "大吉",
			Meaning:     "主喜事臨門，財源廣進，婚姻美滿",
		},
		"天空": {
			Name:        "天空",
			Description: "主虛空、詐偽、僧道",
			Nature:      "平",
			Meaning:     "主虛而不實，防詐騙，宜僧道之事",
		},
		"白虎": {
			Name:        "白虎",
			Description: "主凶喪、疾病、官非",
			Nature:      "大凶",
			Meaning:     "主凶喪之事，疾病之憂，官非之災",
		},
		"太常": {
			Name:        "太常",
			Description: "主衣食、田宅、婦女",
			Nature:      "吉",
			Meaning:     "主衣食豐足，田宅安穩，婦女賢良",
		},
		"玄武": {
			Name:        "玄武",
			Description: "主盜賊、陰私、暗昧",
			Nature:      "凶",
			Meaning:     "主盜賊之憂，陰私之事，暗昧不明",
		},
		"太陰": {
			Name:        "太陰",
			Description: "主陰私、密謀、婦女",
			Nature:      "平",
			Meaning:     "主陰私之事，密謀之計，婦女相關",
		},
		"天后": {
			Name:        "天后",
			Description: "主恩澤、婦女、隱蔽",
			Nature:      "吉",
			Meaning:     "主恩澤降臨，婦女之佑，隱蔽之事",
		},
	}

	resp := &proto.GetTianJiangInfoResponse{}

	if req.TianJiangName != "" {
		if tj, ok := tianJiangDB[req.TianJiangName]; ok {
			resp.TianJiangList = append(resp.TianJiangList, tj)
		}
	} else {
		// 按順序返回
		order := []string{"貴人", "螣蛇", "朱雀", "六合", "勾陳", "青龍",
			"天空", "白虎", "太常", "玄武", "太陰", "天后"}
		for _, name := range order {
			if tj, ok := tianJiangDB[name]; ok {
				resp.TianJiangList = append(resp.TianJiangList, tj)
			}
		}
	}

	return resp, nil
}

// GetMonthGeneralInfo 查詢月將信息
func (s *InfoServer) GetMonthGeneralInfo(ctx context.Context, req *proto.GetMonthGeneralInfoRequest) (*proto.GetMonthGeneralInfoResponse, error) {
	monthGeneralDB := map[string]*proto.MonthGeneralInfo{
		"正月": {Name: "登明", Month: "正月（寅月）", Description: "水神，主智慧、流動", Meaning: "事情多變，宜靈活應對"},
		"二月": {Name: "河魁", Month: "二月（卯月）", Description: "土神，主收藏、穩固", Meaning: "事情漸穩，宜守成"},
		"三月": {Name: "從魁", Month: "三月（辰月）", Description: "金神，主肅殺、果斷", Meaning: "宜果斷決策，不宜拖延"},
		"四月": {Name: "傳送", Month: "四月（巳月）", Description: "金神，主傳遞、變動", Meaning: "事情多變動，宜靈活"},
		"五月": {Name: "小吉", Month: "五月（午月）", Description: "火神，主文明、喜慶", Meaning: "主喜慶之事，但防過急"},
		"六月": {Name: "勝光", Month: "六月（未月）", Description: "火神，主光明、顯達", Meaning: "事情顯明，宜積極進取"},
		"七月": {Name: "太卜", Month: "七月（申月）", Description: "土神，主占卜、決斷", Meaning: "宜占卜決策，果斷行事"},
		"八月": {Name: "天罡", Month: "八月（酉月）", Description: "土神，主剛強、威嚴", Meaning: "事情剛強，宜堅持原則"},
		"九月": {Name: "太沖", Month: "九月（戌月）", Description: "木神，主生發、開始", Meaning: "宜開始新事，生發向上"},
		"十月": {Name: "功曹", Month: "十月（亥月）", Description: "木神，主功勳、成就", Meaning: "宜追求成就，建立功勳"},
		"十一月": {Name: "大吉", Month: "十一月（子月）", Description: "土神，主吉祥、安穩", Meaning: "主吉祥如意，平安順利"},
		"十二月": {Name: "神后", Month: "十二月（丑月）", Description: "土神，主尊貴、庇佑", Meaning: "得神佑護，尊貴吉祥"},
	}

	resp := &proto.GetMonthGeneralInfoResponse{}

	if req.Month != "" {
		if mg, ok := monthGeneralDB[req.Month]; ok {
			resp.MonthGeneralList = append(resp.MonthGeneralList, mg)
		}
	} else {
		order := []string{"正月", "二月", "三月", "四月", "五月", "六月",
			"七月", "八月", "九月", "十月", "十一月", "十二月"}
		for _, month := range order {
			if mg, ok := monthGeneralDB[month]; ok {
				resp.MonthGeneralList = append(resp.MonthGeneralList, mg)
			}
		}
	}

	return resp, nil
}

// QueryDivinationHistory 查詢歷史排盤記錄（示例實現）
func (s *InfoServer) QueryDivinationHistory(ctx context.Context, req *proto.QueryDivinationHistoryRequest) (*proto.QueryDivinationHistoryResponse, error) {
	// 這裡可以連接數據庫查詢歷史記錄
	// 目前返回示例數據
	resp := &proto.QueryDivinationHistoryResponse{
		Total: 0,
	}

	// 示例記錄
	if req.Limit > 0 {
		resp.Records = []*proto.DivinationRecord{
			{
				Id:           "1",
				Date:         "2026-03-16",
				Time:         "12:00",
				YearPillar:   "丙午",
				MonthPillar:  "辛卯",
				DayPillar:    "己丑",
				HourPillar:   "甲子",
				KeTi:         "伏吟課",
				XiongJi:      "吉",
				Question:     "財運如何",
				QuestionType: "財運",
				CreatedAt:    "2026-03-16T12:00:00Z",
			},
		}
		resp.Total = 1
	}

	return resp, nil
}
