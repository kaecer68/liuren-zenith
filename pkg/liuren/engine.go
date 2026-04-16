package liuren

import (
	"fmt"
	"time"

	"github.com/kaecer68/liuren-zenith/pkg/client"
)

// Engine 大六壬排盤引擎
type Engine struct {
	DataSource client.CalendarDataSource
}

// NewEngine 創建大六壬引擎
func NewEngine(ds client.CalendarDataSource) *Engine {
	return &Engine{
		DataSource: ds,
	}
}

// DivinationRequest 占卜請求
type DivinationRequest struct {
	Time         time.Time // 占卜時間
	Question     string    // 占事（可選）
	QuestionType string    // 問題類型：財運/官運/婚姻/疾病/出行/訴訟/尋人/失物/考試/其他
	IsDay        *bool     // 指定晝夜（可選，nil 則自動判斷）
}

// DivinationResult 占卜結果
type DivinationResult struct {
	// 基本信息
	GregorianDate string    // 陽曆日期
	LunarDate     string    // 農曆日期（簡化）
	Time          time.Time // 原始時間

	// 四柱
	YearPillar  string // 年柱
	MonthPillar string // 月柱
	DayPillar   string // 日柱
	HourPillar  string // 時柱

	// 課盤資訊
	MonthGeneral string // 月將
	HourBranch   string // 占時
	IsDay        bool   // 是否晝占

	// 天地盤
	DiPan   [12]string // 地盤
	TianPan [12]string // 天盤

	// 天將
	TianJiang      [12]string // 十二天將佈局
	DayGeneralName string     // 日貴人名稱（晝貴/夜貴）
	GuiRenPosition string     // 貴人位置

	// 四課
	FourKe []KeInfo // 四課詳情

	// 三傳
	SanChuan SanChuanInfo // 三傳詳情

	// 課體與判斷
	KeTi       string          // 課體名稱
	KeTiDetail KeTiExplanation // 課體詳細解說
	XiongJi    string          // 吉凶判斷
	DuanYu     []string        // 詳細斷語

	// 用神分析（新增）
	QuestionType string   // 問題類型
	YongShen     YongShen // 用神資訊

	// 神煞
	VoidBranches []string // 空亡
	YiMa         string   // 驛馬
	TaoHua       string   // 桃花
	TianMa       string   // 天馬
}

// KeTiExplanation 課體詳細解說
type KeTiExplanation struct {
	Name        string   // 課體名稱
	Category    string   // 分類（九宗門/特殊課體）
	Description string   // 課體描述
	Meaning     string   // 占斷含義
	Features    []string // 主要特徵
	Examples    []string // 古籍示例
}

// YongShen 用神資訊
type YongShen struct {
	Type          string   // 用神類型：財運/官運/婚姻/疾病/出行/訴訟/尋人/失物/考試
	SixRelation   string   // 六親關係：父母/官鬼/妻財/子孫/兄弟
	TargetBranch  string   // 用神地支
	TargetGeneral string   // 用神天將
	Explanation   string   // 用神解說
	Advice        string   // 占斷建議
	Analysis      []string // 詳細分析要點
}

// SanChuanExplanation 三傳詳細解說
type SanChuanExplanation struct {
	Chu   ChuanExplanation // 初傳解說
	Zhong ChuanExplanation // 中傳解說
	Mo    ChuanExplanation // 末傳解說
}

// ChuanExplanation 單傳解說
type ChuanExplanation struct {
	Role        string   // 角色：起因/過程/結果
	Explanation string   // 解說文案
	Advice      string   // 建議
	KeyPoints   []string // 關鍵點
}

// KeInfo 課課詳細資訊
type KeInfo struct {
	Number      int    // 課序（1-4）
	Down        string // 下神
	Up          string // 上神
	DownGeneral string // 下神天將
	UpGeneral   string // 上神天將
	Relation    string // 與日干關係（十神）
	Meaning     string // 課之意義
}

// SanChuanInfo 三傳詳細資訊
type SanChuanInfo struct {
	Chu         ChuanDetail      // 初傳
	ChuDetail   ChuanExplanation // 初傳解說
	Zhong       ChuanDetail      // 中傳
	ZhongDetail ChuanExplanation // 中傳解說
	Mo          ChuanDetail      // 末傳
	MoDetail    ChuanExplanation // 末傳解說
	Method      string           // 九宗門取法
}

// ChuanDetail 傳之詳細資訊
type ChuanDetail struct {
	Name     string // 名稱（初傳/中傳/末傳）
	Branch   string // 地支
	General  string // 天將
	Stem     string // 遁干
	Relation string // 六親關係
	IsEmpty  bool   // 是否落空
	IsHorse  bool   // 是否驛馬
	IsTaoHua bool   // 是否桃花
}

// Calculate 執行大六壬排盤
func (e *Engine) Calculate(req DivinationRequest) (*DivinationResult, error) {
	// 1. 獲取曆法基礎數據
	calendarData, err := e.DataSource.GetCalendarData(req.Time)
	if err != nil {
		return nil, fmt.Errorf("獲取曆法數據失敗: %w", err)
	}

	// 2. 解析干支
	dayPillar, err := ParsePillar(calendarData.DayPillar)
	if err != nil {
		return nil, fmt.Errorf("解析日柱失敗: %w", err)
	}

	hourPillar, err := ParsePillar(calendarData.HourPillar)
	if err != nil {
		return nil, fmt.Errorf("解析時柱失敗: %w", err)
	}

	// 3. 判斷晝夜
	isDay := req.IsDay
	if isDay == nil {
		day := IsDayTime(req.Time)
		isDay = &day
	}

	// 4. 計算月將（根據節氣索引）
	monthGeneral := CalculateMonthGeneral(calendarData.SolarTermIdx)

	// 5. 計算天地盤
	diPan := GetDiPan()
	tianPan := CalculateTianPan(monthGeneral, hourPillar.Branch)

	// 6. 尋找貴人並佈置天將
	guiRenPos, ok := FindGuiRen(dayPillar.Stem, *isDay)
	if !ok {
		return nil, fmt.Errorf("無法確定貴人位置")
	}
	tianJiang := CalculateTianJiang(guiRenPos, *isDay, tianPan)

	// 7. 計算四課
	fourKe := CalculateFourKe(dayPillar.Stem, dayPillar.Branch, tianPan, diPan)

	// 8. 計算三傳（九宗門）
	jzCalc := NewJiuZongMenCalculator(dayPillar.Stem, dayPillar.Branch, tianPan, diPan, fourKe)
	sanChuan := jzCalc.CalculateSanChuan()

	// 9. 計算空亡
	voidBranches := e.calculateXunKong(dayPillar)

	// 10. 計算神煞
	yearPillar, _ := ParsePillar(calendarData.YearPillar)
	yiMa := e.calculateYiMa(yearPillar.Branch)
	taoHua := e.calculateTaoHua(yearPillar.Branch)
	tianMa := e.calculateTianMa(yearPillar.Branch)

	// 11. 組裝結果
	result := &DivinationResult{
		GregorianDate:  calendarData.GregorianDate,
		Time:           req.Time,
		YearPillar:     calendarData.YearPillar,
		MonthPillar:    calendarData.MonthPillar,
		DayPillar:      calendarData.DayPillar,
		HourPillar:     calendarData.HourPillar,
		MonthGeneral:   MonthGeneralNames[monthGeneral],
		HourBranch:     BranchNames[hourPillar.Branch],
		IsDay:          *isDay,
		DayGeneralName: e.getDayGeneralName(*isDay),
		GuiRenPosition: BranchNames[guiRenPos],
		VoidBranches:   voidBranches,
		YiMa:           BranchNames[yiMa],
		TaoHua:         BranchNames[taoHua],
		TianMa:         BranchNames[tianMa],
		QuestionType:   req.QuestionType, // 問題類型
	}

	// 轉換天地盤為字串
	for i := 0; i < 12; i++ {
		result.DiPan[i] = BranchNames[diPan[i]]
		result.TianPan[i] = BranchNames[tianPan[i]]
		result.TianJiang[i] = HeavenlyGeneralNames[tianJiang[i]]
	}

	// 計算遁干
	dunGan := CalculateDunGan(dayPillar, tianPan)

	// 轉換四課
	result.FourKe = e.convertFourKe(fourKe, tianJiang, dayPillar.Stem, dunGan)

	// 轉換三傳
	result.SanChuan = e.convertSanChuan(sanChuan, tianJiang, dayPillar.Stem, voidBranches, yiMa, taoHua, dunGan)

	// 12. 課體判斷（現在返回詳細解說）
	result.KeTi, result.KeTiDetail = e.determineKeTi(sanChuan, fourKe, tianPan, diPan)

	// 13. 吉凶判斷（簡化版，完整版需畢法賦）
	result.XiongJi, result.DuanYu = e.judgeXiongJi(result)

	// 14. 計算用神（根據問題類型）
	result.YongShen = e.calculateYongShen(req.QuestionType, result, dayPillar.Stem)

	return result, nil
}

// convertFourKe 轉換四課為顯示格式
func (e *Engine) convertFourKe(fourKe FourKe, tianJiang [12]HeavenlyGeneral, dayStem Stem, dunGan [12]Stem) []KeInfo {
	keList := []struct {
		ke      Ke
		num     int
		meaning string // 四課意義
	}{
		{fourKe.Ke1, 1, "占者自身狀態、表象（日干之寄宮）"},
		{fourKe.Ke2, 2, "第一課之延伸、自身深層狀態（日上神之寄宮）"},
		{fourKe.Ke3, 3, "所問之事、客方狀態（日支之寄宮）"},
		{fourKe.Ke4, 4, "第三課之延伸、事之深層變化（支上神之寄宮）"},
	}

	var result []KeInfo
	for _, item := range keList {
		ki := KeInfo{
			Number:      item.num,
			Down:        BranchNames[item.ke.Down],
			Up:          BranchNames[item.ke.Up],
			DownGeneral: HeavenlyGeneralNames[tianJiang[item.ke.Down]],
			UpGeneral:   HeavenlyGeneralNames[tianJiang[item.ke.Up]],
			Relation:    e.getShiShen(dayStem, item.ke.Up), // 使用十神
			Meaning:     item.meaning,
		}
		if int(item.ke.Down) >= 0 && int(item.ke.Down) < 12 {
			_ = dunGan[item.ke.Down]
		}
		result = append(result, ki)
	}
	return result
}

// convertSanChuan 轉換三傳為顯示格式
func (e *Engine) convertSanChuan(sanChuan SanChuan, tianJiang [12]HeavenlyGeneral, dayStem Stem, voidBranches []string, yiMa, taoHua Branch, dunGan [12]Stem) SanChuanInfo {
	chuanList := []struct {
		info  ChuanInfo
		name  string
		role  string
		isChu bool
	}{
		{sanChuan.Chu, "初傳", "起因/發端", true},
		{sanChuan.Zhong, "中傳", "過程/變化", false},
		{sanChuan.Mo, "末傳", "結果/歸宿", false},
	}

	var details []ChuanDetail
	var explanations []ChuanExplanation
	for _, item := range chuanList {
		isEmpty := false
		branchName := BranchNames[item.info.Branch]
		for _, vb := range voidBranches {
			if vb == branchName {
				isEmpty = true
				break
			}
		}

		relation := e.getShiShen(dayStem, item.info.Branch) // 使用十神
		general := HeavenlyGeneralNames[tianJiang[item.info.Branch]]

		stemName := ""
		if int(item.info.Branch) >= 0 && int(item.info.Branch) < 12 {
			stemName = StemNames[dunGan[item.info.Branch]]
		}
		cd := ChuanDetail{
			Name:     item.name,
			Branch:   branchName,
			General:  general,
			Stem:     stemName,
			Relation: relation,
			IsEmpty:  isEmpty,
			IsHorse:  item.info.Branch == yiMa,
			IsTaoHua: item.info.Branch == taoHua,
		}
		details = append(details, cd)

		// 生成解說
		explanation := e.generateChuanExplanation(item.name, item.role, branchName, general, relation, isEmpty, item.info.Branch == yiMa)
		explanations = append(explanations, explanation)
	}

	return SanChuanInfo{
		Chu:         details[0],
		ChuDetail:   explanations[0],
		Zhong:       details[1],
		ZhongDetail: explanations[1],
		Mo:          details[2],
		MoDetail:    explanations[2],
		Method:      sanChuan.Method,
	}
}

// generateChuanExplanation 生成單傳解說（結合天將與十神）
func (e *Engine) generateChuanExplanation(name, role, branch, general, shiShen string, isEmpty, isHorse bool) ChuanExplanation {
	ce := ChuanExplanation{
		Role:      role,
		KeyPoints: []string{},
	}

	// 基礎解說
	ce.Explanation = fmt.Sprintf("%s臨%s（%s），十神為%s。", general, branch, name, shiShen)

	// 空亡判斷
	if isEmpty {
		ce.Explanation += "此傳落空亡，主事多虛詐，難以成就。"
		ce.KeyPoints = append(ce.KeyPoints, "落空亡，事多虛詐")
	}

	// 驛馬判斷
	if isHorse {
		ce.Explanation += "逢驛馬，主變動、遠行。"
		ce.KeyPoints = append(ce.KeyPoints, "逢驛馬，主變動")
	}

	// 結合天將與十神生成具體建議
	ce.Advice = e.generateSmartAdvice(name, general, shiShen, isEmpty)

	// 根據角色補充說明
	if name == "初傳" {
		ce.Explanation = "【發用·起因】" + ce.Explanation
		if !isEmpty {
			ce.Explanation += e.getChuChuanTiming(shiShen)
		}
	} else if name == "中傳" {
		ce.Explanation = "【移易·過程】" + ce.Explanation
		if !isEmpty {
			ce.Explanation += e.getZhongChuanProcess(shiShen)
		}
	} else if name == "末傳" {
		ce.Explanation = "【歸計·結果】" + ce.Explanation
		if !isEmpty {
			ce.Explanation += e.getMoChuanOutcome(general, shiShen)
		}
	}

	// 添加關鍵點
	ce.KeyPoints = append(ce.KeyPoints, e.getKeyPointByCombination(general, shiShen))

	return ce
}

// generateSmartAdvice 結合天將與十神生成智能建議
func (e *Engine) generateSmartAdvice(chuanName, general, shiShen string, isEmpty bool) string {
	if isEmpty {
		return "此傳落空，所求之事難以成就，宜靜待時機。"
	}

	// 天將分類
	guiShen := map[string]bool{"貴人": true, "青龍": true, "六合": true, "太常": true}
	xiongShen := map[string]bool{"螣蛇": true, "朱雀": true, "白虎": true, "玄武": true, "天空": true}
	zhongShen := map[string]bool{"勾陳": true, "天后": true}

	// 十神分類
	bangFu := map[string]bool{"比肩": true, "劫財": true, "正印": true, "偏印": true} // 幫扶類
	xieHao := map[string]bool{"食神": true, "傷官": true}                         // 洩耗類
	keZhi := map[string]bool{"正官": true, "七殺": true}                          // 克制類
	caiLi := map[string]bool{"正財": true, "偏財": true}                          // 財利類

	// 組合判斷
	isGui := guiShen[general]
	isXiong := xiongShen[general]
	isBangFu := bangFu[shiShen]
	isXieHao := xieHao[shiShen]
	isKeZhi := keZhi[shiShen]
	isCaiLi := caiLi[shiShen]

	// 初傳/中傳/末傳的階段性建議
	stageAdvice := ""
	switch chuanName {
	case "初傳":
		stageAdvice = "【發端】"
	case "中傳":
		stageAdvice = "【過程】"
	case "末傳":
		stageAdvice = "【結果】"
	}

	// 組合邏輯
	switch {
	// 吉神 + 幫扶 = 大吉，順勢而為
	case isGui && isBangFu:
		return stageAdvice + "得吉神護佑，且有助力，宜積極進取，順勢而為。此時行動易得支持，事半功倍。"

	// 吉神 + 財利 = 求財有利
	case isGui && isCaiLi:
		if general == "青龍" {
			return stageAdvice + "青龍主財祿，又逢財星，求財大吉。宜把握機會，投資或交易可獲利。"
		}
		if general == "六合" {
			return stageAdvice + "六合主和合，財星臨之，利於合作求財。宜與人合作，共分利益。"
		}
		return stageAdvice + "吉神臨財星，求財有利，但需防範風險，不可貪心。"

	// 吉神 + 克制 = 有貴人解厄
	case isGui && isKeZhi:
		if shiShen == "正官" {
			return stageAdvice + "官星得吉神扶持，利於求官或處理公務。得貴人提攜，前途可期。"
		}
		return stageAdvice + "雖有壓力，但有貴人相助可解。宜謙卑求助，不可獨斷。"

	// 吉神 + 洩耗 = 順勢發揮
	case isGui && isXieHao:
		return stageAdvice + "才華得以發揮，又有吉神護佑。宜展現能力，求名求利皆順。"

	// 凶神 + 克制 = 壓力重重，需防範
	case isXiong && isKeZhi:
		if shiShen == "七殺" && general == "白虎" {
			return stageAdvice + "七殺臨白虎，凶性倍增。此階段壓力極大，或有血光之災。宜靜守避禍，萬勿冒進。"
		}
		if shiShen == "七殺" {
			return stageAdvice + "七殺臨凶神，主有強敵或重大壓力。需謹慎應對，尋求化解之道。"
		}
		if general == "朱雀" {
			return stageAdvice + "官星臨朱雀，主有口舌是非或官司糾紛。宜謹言慎行，防範文書之災。"
		}
		return stageAdvice + "此階段阻力較大，需步步為營。宜守不宜攻，待機而動。"

	// 凶神 + 洩耗 = 防範暗算
	case isXiong && isXieHao:
		if general == "玄武" {
			return stageAdvice + "玄武臨洩神，主有暗耗或暗中損失。宜謹守財物，防範小人暗算。"
		}
		if general == "螣蛇" {
			return stageAdvice + "螣蛇主虛驚，又逢洩神，主有莫名煩憂。宜靜心應對，勿被表象迷惑。"
		}
		return stageAdvice + "此階段需防範暗中損耗，不宜輕信他人。宜低調行事，減少支出。"

	// 凶神 + 財利 = 求財有風險
	case isXiong && isCaiLi:
		if general == "玄武" {
			return stageAdvice + "玄武臨財星，主財來路不正或有被盜之虞。求財需謹慎，防範風險。"
		}
		return stageAdvice + "凶神臨財，求財有風險。宜穩健經營，不可貪求暴利。"

	// 凶神 + 幫扶 = 有助力可解凶
	case isXiong && isBangFu:
		if shiShen == "正印" {
			return stageAdvice + "雖有凶神，但印星護身，可得長輩或貴人庇佑。危機可解，不必過憂。"
		}
		return stageAdvice + "雖遇困境，但有同輩或朋友相助。宜團結合作，共渡難關。"

	// 中神組合
	case zhongShen[general]:
		if shiShen == "正官" || shiShen == "正印" {
			return stageAdvice + "事態平穩，宜按常規處理。保持現狀，不宜大動。"
		}
		return stageAdvice + "此階段宜審慎觀察，待時而動。不急不躁，靜待轉機。"

	// 默認
	default:
		return stageAdvice + "依常規審慎行事，觀察形勢再作決定。"
	}
}

// getChuChuanTiming 初傳時機判斷
func (e *Engine) getChuChuanTiming(shiShen string) string {
	switch shiShen {
	case "正官", "正印":
		return "此為最佳時機，宜立即行動。"
	case "七殺", "偏印":
		return "時機雖有，但需謹慎，充分準備後再動。"
	case "正財", "偏財":
		return "有利可圖，宜把握機會。"
	case "食神", "傷官":
		return "可展現才華，但需量力而為。"
	case "比肩", "劫財":
		return "有同伴相助，宜合作而行。"
	default:
		return ""
	}
}

// getZhongChuanProcess 中傳過程判斷
func (e *Engine) getZhongChuanProcess(shiShen string) string {
	switch shiShen {
	case "正官", "正印":
		return "過程順遂，按部就班即可。"
	case "七殺", "偏印":
		return "過程有阻，需調整策略應對。"
	case "正財", "偏財":
		return "過程中有收益，但需防範風險。"
	case "食神", "傷官":
		return "可發揮創意，但勿過度消耗。"
	case "比肩", "劫財":
		return "過程中需與人協調，防範競爭。"
	default:
		return ""
	}
}

// getMoChuanOutcome 末傳結果判斷
func (e *Engine) getMoChuanOutcome(general, shiShen string) string {
	// 吉神 + 吉神組合
	guiShen := map[string]bool{"貴人": true, "青龍": true, "六合": true}
	if guiShen[general] {
		switch shiShen {
		case "正官", "正印", "正財":
			return "結果圓滿，所求皆得。"
		case "食神":
			return "結果稱心，名利雙收。"
		case "比肩":
			return "結果平穩，與人共享。"
		}
	}

	// 凶神組合
	xiongShen := map[string]bool{"白虎": true, "螣蛇": true, "玄武": true}
	if xiongShen[general] {
		switch shiShen {
		case "七殺":
			return "結果不佳，需防損失。"
		case "傷官":
			return "結果有損，宜早收斂。"
		case "劫財":
			return "結果有競爭，防範損失。"
		}
	}

	return ""
}

// getKeyPointByCombination 根據天將與十神組合返回關鍵點
func (e *Engine) getKeyPointByCombination(general, shiShen string) string {
	// 定義關鍵點映射
	keyPoints := map[string]map[string]string{
		"貴人": {
			"正官": "貴人提攜，利於仕途",
			"正印": "貴人相助，學業有成",
			"正財": "貴人引財，財源穩定",
			"食神": "貴人賞識，才華得展",
		},
		"青龍": {
			"正財": "青龍送財，財運亨通",
			"偏財": "偏財得利，意外之喜",
			"正官": "龍德護官，升遷有望",
			"食神": "龍德生輝，名利雙收",
		},
		"六合": {
			"正財": "合作生財，共分利益",
			"正官": "聯合發展，官運順遂",
			"比肩": "伙伴同心，共創佳績",
		},
		"朱雀": {
			"正官": "文書謹慎，防範官司",
			"七殺": "口舌是非，避免爭執",
			"傷官": "言多必失，謹言慎行",
		},
		"白虎": {
			"七殺": "殺氣太重，防血光災",
			"正官": "官非纏身，需防訴訟",
			"偏印": "長輩有憂，注意安全",
		},
		"玄武": {
			"偏財": "偏財有險，防被盜騙",
			"正財": "正財被耗，謹守錢財",
			"劫財": "防範小人，謹守財物",
		},
		"螣蛇": {
			"偏印": "虛驚一場，勿憂心過度",
			"七殺": "驚憂交加，防範詐騙",
			"傷官": "憂慮過度，宜放寬心",
		},
	}

	if points, ok := keyPoints[general]; ok {
		if point, ok := points[shiShen]; ok {
			return point
		}
	}

	// 默認關鍵點
	return general + "臨" + shiShen + "，宜審慎"
}

// calculateYongShen 計算用神
func (e *Engine) calculateYongShen(questionType string, result *DivinationResult, dayStem Stem) YongShen {
	if questionType == "" {
		questionType = "其他"
	}

	yongShen := YongShen{
		Type: questionType,
	}

	// 根據問題類型確定用神六親
	switch questionType {
	case "財運":
		yongShen.SixRelation = "妻財"
		yongShen.Explanation = "財運以妻財爻為用神。妻財者，我克之神也，主錢財、利潤。"
		yongShen.Advice = "看妻財爻臨何天將，逢吉神則財源順利，逢凶神則財多破耗。"
	case "官運":
		yongShen.SixRelation = "官鬼"
		yongShen.Explanation = "官運以官鬼爻為用神。官鬼者，克我之神也，主官職、權勢。"
		yongShen.Advice = "看官鬼爻臨何天將，旺相則官運亨通，休囚則仕途受阻。"
	case "婚姻":
		yongShen.SixRelation = "妻財" // 男看妻財，女看官鬼（簡化處理）
		yongShen.Explanation = "婚姻以妻財爻（男占）或官鬼爻（女占）為用神，六合為媒。"
		yongShen.Advice = "看用神的上下神及天將，逢六合、青龍則婚姻可成。"
	case "疾病":
		yongShen.SixRelation = "官鬼"
		yongShen.Explanation = "疾病以官鬼爻為用神。官鬼為克我之神，主疾厄、災病。"
		yongShen.Advice = "看官鬼是否落空，空則病易癒；看白虎是否加臨，臨則病重。"
	case "出行":
		yongShen.SixRelation = ""
		yongShen.Explanation = "出行以驛馬為用神，看初傳及日干支。"
		yongShen.Advice = "看驛馬是否落空，空則不行；看初傳吉凶，吉則路順。"
	case "訴訟":
		yongShen.SixRelation = "官鬼"
		yongShen.Explanation = "訴訟以官鬼爻為用神，朱雀為訴狀。"
		yongShen.Advice = "看官鬼旺衰，旺則官司難解；看朱雀動靜，動則口舌紛爭。"
	case "尋人":
		yongShen.SixRelation = ""
		yongShen.Explanation = "尋人以日干為人，看天盤日干上神及驛馬。"
		yongShen.Advice = "看驛馬方位，可知人所往；看空亡，空則難尋。"
	case "失物":
		yongShen.SixRelation = "妻財"
		yongShen.Explanation = "失物以妻財爻為用神（財為物）。"
		yongShen.Advice = "看妻財是否落空，空則難尋；看玄武，臨則被盜。"
	case "考試":
		yongShen.SixRelation = "父母"
		yongShen.Explanation = "考試以父母爻為用神。父母者，生我之神，主文書、學業。"
		yongShen.Advice = "看父母旺相，旺則考運佳；看朱雀，臨則文書順利。"
	default:
		yongShen.SixRelation = ""
		yongShen.Explanation = "一般占斷以日干為主，看三傳及四課吉凶。"
		yongShen.Advice = "審視課體，觀察三傳，依常法斷之。"
	}

	// 在三傳中尋找用神
	if yongShen.SixRelation != "" {
		yongShen = e.findYongShenInSanChuan(yongShen, result, dayStem)
	}

	return yongShen
}

// findYongShenInSanChuan 在三傳中尋找用神
func (e *Engine) findYongShenInSanChuan(yongShen YongShen, result *DivinationResult, dayStem Stem) YongShen {
	chuanList := []struct {
		chuan ChuanDetail
		name  string
	}{
		{result.SanChuan.Chu, "初傳"},
		{result.SanChuan.Zhong, "中傳"},
		{result.SanChuan.Mo, "末傳"},
	}

	// 六親與十神對應表
	liuQinToShiShen := map[string][]string{
		"父母": {"正印", "偏印"},
		"官鬼": {"正官", "七殺"},
		"妻財": {"正財", "偏財"},
		"子孫": {"食神", "傷官"},
		"兄弟": {"比肩", "劫財"},
	}

	// 檢查 relation 是否匹配用神六親
	isMatch := func(relation string) bool {
		if relation == yongShen.SixRelation {
			return true
		}
		for _, shiShen := range liuQinToShiShen[yongShen.SixRelation] {
			if relation == shiShen {
				return true
			}
		}
		return false
	}

	// 優先在初傳尋找用神
	for _, item := range chuanList {
		if isMatch(item.chuan.Relation) {
			yongShen.TargetBranch = item.chuan.Branch
			yongShen.TargetGeneral = item.chuan.General
			yongShen.Analysis = append(yongShen.Analysis,
				fmt.Sprintf("用神%s臨%s（%s），為%s。", yongShen.SixRelation, item.name, item.chuan.Branch, item.chuan.General))

			if item.chuan.IsEmpty {
				yongShen.Analysis = append(yongShen.Analysis, "用神落空亡，所求難成。")
			}
			if item.chuan.IsHorse {
				yongShen.Analysis = append(yongShen.Analysis, "用神逢驛馬，主有變動。")
			}
			if item.chuan.IsTaoHua {
				yongShen.Analysis = append(yongShen.Analysis, "用神逢桃花，主人緣情事。")
			}

			break
		}
	}

	// 若三傳無用神，在四課中尋找
	if yongShen.TargetBranch == "" {
		for _, ke := range result.FourKe {
			if isMatch(ke.Relation) {
				yongShen.TargetBranch = ke.Up
				yongShen.TargetGeneral = ke.UpGeneral
				yongShen.Analysis = append(yongShen.Analysis,
					fmt.Sprintf("用神%s臨第%d課上神（%s），為%s。", yongShen.SixRelation, ke.Number, ke.Up, ke.UpGeneral))
				break
			}
		}
	}

	// 仍未找到，提示在四課中查看
	if yongShen.TargetBranch == "" {
		yongShen.Analysis = append(yongShen.Analysis,
			fmt.Sprintf("%s爻未在三傳出現，請查看四課中%s爻的位置。", yongShen.SixRelation, yongShen.SixRelation))
	}

	return yongShen
}

// getLiuQin 獲取六親關係（父母、官鬼、妻財、子孫、兄弟/比和）
func (e *Engine) getLiuQin(dayStem Stem, target Branch) string {
	dayStemElement := e.getStemElement(dayStem)
	targetElement := e.getBranchElement(target)

	// 五行相生：木生火、火生土、土生金、金生水、水生木
	// 五行相剋：木剋土、土剋水、水剋火、火剋金、金剋木

	if dayStemElement == targetElement {
		return "兄弟" // 比和
	}

	// 生我者為父母
	if e.isGenerating(targetElement, dayStemElement) {
		return "父母"
	}

	// 我生者為子孫
	if e.isGenerating(dayStemElement, targetElement) {
		return "子孫"
	}

	// 剋我者為官鬼
	if e.isOvercoming(targetElement, dayStemElement) {
		return "官鬼"
	}

	// 我剋者為妻財
	if e.isOvercoming(dayStemElement, targetElement) {
		return "妻財"
	}

	return "-"
}

// getShiShen 獲取十神關係
// 十神：比肩、劫財、食神、傷官、偏財、正財、七殺、正官、偏印、正印
func (e *Engine) getShiShen(dayStem Stem, target Branch) string {
	dayStemElement := e.getStemElement(dayStem)
	targetElement := e.getBranchElement(target)

	// 判斷陰陽
	// 天干陰陽：甲(0)丙(2)戊(4)庚(6)壬(8)為陽，乙(1)丁(3)己(5)辛(7)癸(9)為陰
	dayStemYang := dayStem%2 == 0
	// 地支陰陽：子(0)寅(2)辰(4)午(6)申(8)戌(10)為陽，丑(1)卯(3)巳(5)未(7)酉(9)亥(11)為陰
	targetYang := target%2 == 0

	// 同陰陽
	isSameYinYang := dayStemYang == targetYang

	// 五行關係判斷
	if dayStemElement == targetElement {
		// 同我者：比肩（同陰陽）/劫財（異陰陽）
		if isSameYinYang {
			return "比肩"
		}
		return "劫財"
	}

	// 我生者：食神（同陰陽）/傷官（異陰陽）
	if e.isGenerating(dayStemElement, targetElement) {
		if isSameYinYang {
			return "食神"
		}
		return "傷官"
	}

	// 我剋者：偏財（異陰陽）/正財（同陰陽）
	if e.isOvercoming(dayStemElement, targetElement) {
		if isSameYinYang {
			return "正財"
		}
		return "偏財"
	}

	// 剋我者：七殺（同陰陽）/正官（異陰陽）
	if e.isOvercoming(targetElement, dayStemElement) {
		if isSameYinYang {
			return "七殺"
		}
		return "正官"
	}

	// 生我者：偏印（異陰陽）/正印（同陰陽）
	if e.isGenerating(targetElement, dayStemElement) {
		if isSameYinYang {
			return "正印"
		}
		return "偏印"
	}

	return "-"
}

// isGenerating 判斷是否相生（from 生 to）
func (e *Engine) isGenerating(from, to string) bool {
	// 木生火、火生土、土生金、金生水、水生木
	generating := map[string]string{
		"木": "火",
		"火": "土",
		"土": "金",
		"金": "水",
		"水": "木",
	}
	return generating[from] == to
}

// isOvercoming 判斷是否相剋（from 剋 to）
func (e *Engine) isOvercoming(from, to string) bool {
	// 木剋土、土剋水、水剋火、火剋金、金剋木
	overcoming := map[string]string{
		"木": "土",
		"土": "水",
		"水": "火",
		"火": "金",
		"金": "木",
	}
	return overcoming[from] == to
}

// getStemElement 獲取天干五行
func (e *Engine) getStemElement(stem Stem) string {
	elements := []string{"木", "木", "火", "火", "土", "土", "金", "金", "水", "水"}
	return elements[stem]
}

// getBranchElement 獲取地支五行
func (e *Engine) getBranchElement(branch Branch) string {
	elements := []string{"水", "土", "木", "木", "土", "火", "火", "土", "金", "金", "土", "水"}
	return elements[branch]
}

// getDayGeneralName 獲取貴人名稱
func (e *Engine) getDayGeneralName(isDay bool) string {
	if isDay {
		return "晝貴（陽貴）"
	}
	return "夜貴（陰貴）"
}

// calculateXunKong 計算旬空
func (e *Engine) calculateXunKong(dayPillar Sexagenary) []string {
	// 旬首地支 = (日支 - 日干 + 12) % 12
	// 空亡 = (旬首 + 10) % 12, (旬首 + 11) % 12
	// 例：甲子旬（旬首子=0）→ 戌(10)亥(11)空
	//     甲戌旬（旬首戌=10）→ 申(8)酉(9)空
	//     甲寅旬（旬首寅=2）→ 子(0)丑(1)空
	xunShouBranch := (int(dayPillar.Branch) - int(dayPillar.Stem) + 12) % 12
	void1 := (xunShouBranch + 10) % 12
	void2 := (xunShouBranch + 11) % 12
	return []string{BranchNames[void1], BranchNames[void2]}
}

// calculateYiMa 計算驛馬
func (e *Engine) calculateYiMa(yearBranch Branch) Branch {
	// 申子辰馬在寅、寅午戌馬在申、巳酉丑馬在亥、亥卯未馬在巳
	maMap := map[Branch]Branch{
		Shen: Yin, Zi: Yin, Chen: Yin,
		Yin: Shen, Wu: Shen, Xu: Shen,
		Si: Hai, You: Hai, Chou: Hai,
		Hai: Si, Mao: Si, Wei: Si,
	}
	if ma, ok := maMap[yearBranch]; ok {
		return ma
	}
	return Yin
}

// calculateTaoHua 計算桃花
func (e *Engine) calculateTaoHua(yearBranch Branch) Branch {
	// 亥卯未見子、巳酉丑見午、寅午戌見卯、申子辰見酉
	taoMap := map[Branch]Branch{
		Hai: Zi, Mao: Zi, Wei: Zi,
		Si: Wu, You: Wu, Chou: Wu,
		Yin: Mao, Wu: Mao, Xu: Mao,
		Shen: You, Zi: You, Chen: You,
	}
	if tao, ok := taoMap[yearBranch]; ok {
		return tao
	}
	return Zi
}

// calculateTianMa 計算天馬
// 天馬：申子辰馬在午、寅午戌馬在申、巳酉丑馬在亥、亥卯未馬在巳
func (e *Engine) calculateTianMa(yearBranch Branch) Branch {
	tianMaMap := map[Branch]Branch{
		Shen: Wu, Zi: Wu, Chen: Wu, // 申子辰 → 午
		Yin: Shen, Wu: Shen, Xu: Shen, // 寅午戌 → 申
		Si: Hai, You: Hai, Chou: Hai, // 巳酉丑 → 亥
		Hai: Si, Mao: Si, Wei: Si, // 亥卯未 → 巳
	}
	if tianMa, ok := tianMaMap[yearBranch]; ok {
		return tianMa
	}
	return Wu // 預設午
}

// determineKeTi 判斷課體並返回詳細解說
func (e *Engine) determineKeTi(sanChuan SanChuan, fourKe FourKe, tianPan, diPan [12]Branch) (string, KeTiExplanation) {
	// 檢查伏吟
	isFuYin := true
	for i := 0; i < 12; i++ {
		if tianPan[i] != diPan[i] {
			isFuYin = false
			break
		}
	}
	if isFuYin {
		return "伏吟課", e.getKeTiExplanation("伏吟課")
	}

	// 檢查返吟
	isFanYin := true
	for i := 0; i < 12; i++ {
		if int(tianPan[i]) != (int(diPan[i])+6)%12 {
			isFanYin = false
			break
		}
	}
	if isFanYin {
		return "返吟課", e.getKeTiExplanation("返吟課")
	}

	// 根據九宗門判斷
	for _, name := range JiuZongMenNames {
		if sanChuan.Method == name || (len(sanChuan.Method) > len(name) && sanChuan.Method[:len(name)] == name) {
			keTiName := name + "課"
			return keTiName, e.getKeTiExplanation(keTiName)
		}
	}

	return "普通課", e.getKeTiExplanation("普通課")
}

// getKeTiExplanation 獲取課體詳細解說
func (e *Engine) getKeTiExplanation(keTiName string) KeTiExplanation {
	explanations := map[string]KeTiExplanation{
		"賊克法課": {
			Name:        "賊克法課",
			Category:    "九宗門",
			Description: "四課中有下克上（賊）或上克下（克）者。下克上為重，上克下為輕，故名賊克。",
			Meaning:     "主事有侵犯、爭鬥之象。賊者陰私暗昧，克者明爭明鬥。",
			Features: []string{
				"以下犯上為賊，以上制下為克",
				"四課中必須有相克關係",
				"取相克之神為初傳",
			},
			Examples: []string{
				"《大六壬課經》：賊克者，以下犯上也。",
				"《畢法賦》：賊克課，利於先發制人。",
			},
		},
		"比用法課": {
			Name:        "比用法課",
			Category:    "九宗門",
			Description: "兩賊或兩克，取與日干陰陽相比者為用，又名知一法。",
			Meaning:     "主事有同類相助、比鄰相親之象。",
			Features: []string{
				"兩賊或兩克同時存在",
				"日干與上神陰陽相同者取之",
				"甲丙戊庚壬為陽，乙丁己辛癸為陰",
			},
			Examples: []string{
				"《大六壬課經》：比用者，同類相比也。",
			},
		},
		"涉害法課": {
			Name:        "涉害法課",
			Category:    "九宗門",
			Description: "俱比或俱不比，以涉害深淺為用。涉害深者受克多，淺者受克少。",
			Meaning:     "主事有深淺輕重之分，涉害深者事大，淺者事小。",
			Features: []string{
				"計算各地支在地盤上受克的次數",
				"涉害相等時取四孟上神",
				"無四孟取四仲上神",
			},
			Examples: []string{
				"《大六壬課經》：涉害者，計較深淺也。",
			},
		},
		"遙克法課": {
			Name:        "遙克法課",
			Category:    "九宗門",
			Description: "四課俱無賊克，取神遙克日或日遙克神者。",
			Meaning:     "主事有遠方來克或克遠方之象，遙遠相制。",
			Features: []string{
				"四課無上下相克",
				"神遙克日為重",
				"日遙克神為輕",
			},
			Examples: []string{
				"《大六壬課經》：遙克者，遠方相克也。",
			},
		},
		"昴星法課": {
			Name:        "昴星法課",
			Category:    "九宗門",
			Description: "剛日取酉上神，柔日取卯上神，如昴星之特立。",
			Meaning:     "主事有獨特、孤立之象，如昴星之在二十八宿中特立獨行。",
			Features: []string{
				"剛日從酉上起課",
				"柔日從卯上起課",
				"四課俱無賊克時用之",
			},
			Examples: []string{
				"《大六壬課經》：昴星者，特立獨行也。",
			},
		},
		"別責法課": {
			Name:        "別責法課",
			Category:    "九宗門",
			Description: "四課缺一（無克），剛日取干合，柔日取支前三合。",
			Meaning:     "主事有別責他屬、責任轉移之象。",
			Features: []string{
				"四課缺一，無上下克",
				"剛日取日干合處",
				"柔日取日支三合處",
			},
			Examples: []string{
				"《大六壬課經》：別責者，責任別屬也。",
			},
		},
		"八專法課": {
			Name:        "八專法課",
			Category:    "九宗門",
			Description: "八專日（干支同位），取剛柔比用。",
			Meaning:     "主事有專注、集中之象，但亦有偏頗之虞。",
			Features: []string{
				"甲寅、乙卯、丁未、己未、庚申、辛酉等日",
				"剛日取戌，柔日取酉",
				"干支同宮，陰陽不分",
			},
			Examples: []string{
				"《大六壬課經》：八專者，干支同宮也。",
			},
		},
		"伏吟課": {
			Name:        "伏吟課",
			Category:    "特殊課體",
			Description: "天盤地盤相同，如人伏地呻吟，故名伏吟。",
			Meaning:     "主事多滯礙，遲疑不決，靜守為宜，動則有悔。",
			Features: []string{
				"天盤地盤完全重合",
				"月將加时後不動",
				"事多反覆，難以速決",
			},
			Examples: []string{
				"《畢法賦》：伏吟課，動不如靜。",
				"《大六壬課經》：伏吟者，呻吟不動也。",
			},
		},
		"返吟課": {
			Name:        "返吟課",
			Category:    "特殊課體",
			Description: "天盤地盤對沖，如人反覆來回，故名返吟。",
			Meaning:     "主事多反覆，來回不定，聚散無常，往來反覆。",
			Features: []string{
				"天盤地盤相對沖（相差六位）",
				"事多反覆，難以定局",
				"利於遠行，不利安居",
			},
			Examples: []string{
				"《畢法賦》：返吟課，反覆不定。",
				"《大六壬課經》：返吟者，反覆來回也。",
			},
		},
		"普通課": {
			Name:        "普通課",
			Category:    "一般課體",
			Description: "無特殊課象，按一般法則占斷。",
			Meaning:     "課象平穩，宜審慎行事，依常規推斷。",
			Features: []string{
				"無伏吟返吟等特殊現象",
				"按九宗門常法取三傳",
				"吉凶依常規判斷",
			},
			Examples: []string{
				"《大六壬課經》：普通課，依常法斷之。",
			},
		},
	}

	if exp, ok := explanations[keTiName]; ok {
		return exp
	}
	return explanations["普通課"]
}

// judgeXiongJi 判斷吉凶（簡化版）
func (e *Engine) judgeXiongJi(result *DivinationResult) (string, []string) {
	var duanYu []string
	xiongJi := "平"

	// 基本判斷邏輯
	if result.SanChuan.Chu.IsEmpty {
		duanYu = append(duanYu, "初傳落空，事多虛詐，所求難成。")
	}

	if result.SanChuan.Chu.General == "貴人" {
		duanYu = append(duanYu, "初傳見貴人，主有貴人扶持，吉利。")
		xiongJi = "吉"
	}

	if result.SanChuan.Chu.General == "螣蛇" || result.SanChuan.Chu.General == "白虎" {
		duanYu = append(duanYu, "初傳見凶神，主有驚憂或血光之災。")
		xiongJi = "凶"
	}

	if result.KeTi == "伏吟課" {
		duanYu = append(duanYu, "伏吟課主事多滯礙，遲疑不決。")
	}

	if result.KeTi == "返吟課" {
		duanYu = append(duanYu, "返吟課主事多反覆，來回不定。")
	}

	if len(duanYu) == 0 {
		duanYu = append(duanYu, "課象平穩，宜審慎行事。")
	}

	return xiongJi, duanYu
}
