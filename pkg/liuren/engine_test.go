package liuren

import (
	"fmt"
	"testing"
	"time"

	"github.com/kaecer68/liuren-zenith/pkg/client"
)

type mockDataSource struct {
	data *client.CalendarData
	err  error
}

func (m *mockDataSource) GetCalendarData(t time.Time) (*client.CalendarData, error) {
	return m.data, m.err
}

func TestEngineGetShiShen(t *testing.T) {
	e := NewEngine(nil)

	tests := []struct {
		dayStem Stem
		target  Branch
		want    string
	}{
		{StemJia, Zi, "正印"},    // 甲木(陽) + 子水(陽) → 生我同陽 = 正印
		{StemJia, Chou, "偏財"},  // 甲木(陽) + 丑土(陰) → 我剋異陽 = 偏財
		{StemJia, Yin, "比肩"},   // 甲木(陽) + 寅木(陽) → 同我同陽 = 比肩
		{StemJia, Mao, "劫財"},   // 甲木(陽) + 卯木(陰) → 同我異陽 = 劫財
		{StemYi, Zi, "偏印"},     // 乙木(陰) + 子水(陽) → 生我異陽 = 偏印
		{StemYi, Mao, "比肩"},    // 乙木(陰) + 卯木(陰) → 同我同陰 = 比肩
		{StemBing, Shen, "正財"}, // 丙火(陽) + 申金(陽) → 我剋同陽 = 正財
		{StemBing, You, "偏財"},  // 丙火(陽) + 酉金(陰) → 我剋異陽 = 偏財
		{StemGeng, Yin, "正財"},  // 庚金(陽) + 寅木(陽) → 我剋同陽 = 正財
		{StemRen, Wu, "正財"},    // 壬水(陽) + 午火(陰) → 我剋異陽 = 正財
		{StemJia, Si, "傷官"},    // 甲木(陽) + 巳火(陰) → 我生異陽 = 傷官
		{StemJia, Wu, "食神"},    // 甲木(陽) + 午火(陽) → 我生同陽 = 食神
		{StemGui, Shen, "偏印"},  // 癸水(陰) + 申金(陽) → 生我異陽 = 偏印
		{StemGui, You, "正印"},   // 癸水(陰) + 酉金(陰) → 生我同陰 = 正印
		{StemWu, Yin, "七殺"},    // 戊土(陽) + 寅木(陽) → 剋我同陽 = 七殺
		{StemWu, Mao, "正官"},    // 戊土(陽) + 卯木(陰) → 剋我異陽 = 正官
	}

	for _, tt := range tests {
		t.Run(StemNames[tt.dayStem]+BranchNames[tt.target], func(t *testing.T) {
			got := e.getShiShen(tt.dayStem, tt.target)
			if got != tt.want {
				t.Errorf("getShiShen(%s, %s) = %s, want %s", StemNames[tt.dayStem], BranchNames[tt.target], got, tt.want)
			}
		})
	}
}

func TestEngineGetLiuQin(t *testing.T) {
	e := NewEngine(nil)

	tests := []struct {
		dayStem Stem
		target  Branch
		want    string
	}{
		{StemJia, Zi, "父母"},    // 木 + 水 = 生我
		{StemJia, Chou, "妻財"},  // 木 + 土 = 我剋
		{StemJia, Yin, "兄弟"},   // 木 + 木 = 比和
		{StemBing, Shen, "妻財"}, // 火 + 金 = 我剋
		{StemGeng, Yin, "妻財"},  // 金 + 木 = 我剋
		{StemRen, Wu, "妻財"},    // 水 + 火 = 我剋
		{StemJia, Si, "子孫"},    // 木 + 火 = 我生
		{StemGui, Shen, "父母"},  // 水 + 金 = 生我
		{StemWu, Yin, "官鬼"},    // 土 + 木 = 剋我
	}

	for _, tt := range tests {
		t.Run(StemNames[tt.dayStem]+BranchNames[tt.target], func(t *testing.T) {
			got := e.getLiuQin(tt.dayStem, tt.target)
			if got != tt.want {
				t.Errorf("getLiuQin(%s, %s) = %s, want %s", StemNames[tt.dayStem], BranchNames[tt.target], got, tt.want)
			}
		})
	}
}

func TestEngineElements(t *testing.T) {
	e := NewEngine(nil)

	// 天干五行
	if e.getStemElement(StemJia) != "木" {
		t.Errorf("甲應為木")
	}
	if e.getStemElement(StemBing) != "火" {
		t.Errorf("丙應為火")
	}
	if e.getStemElement(StemWu) != "土" {
		t.Errorf("戊應為土")
	}
	if e.getStemElement(StemGeng) != "金" {
		t.Errorf("庚應為金")
	}
	if e.getStemElement(StemRen) != "水" {
		t.Errorf("壬應為水")
	}

	// 地支五行
	if e.getBranchElement(Zi) != "水" {
		t.Errorf("子應為水")
	}
	if e.getBranchElement(Yin) != "木" {
		t.Errorf("寅應為木")
	}
	if e.getBranchElement(Si) != "火" {
		t.Errorf("巳應為火")
	}
	if e.getBranchElement(Shen) != "金" {
		t.Errorf("申應為金")
	}
	if e.getBranchElement(Chou) != "土" {
		t.Errorf("丑應為土")
	}
}

func TestEngineIsGenerating(t *testing.T) {
	e := NewEngine(nil)
	if !e.isGenerating("木", "火") {
		t.Error("木應生火")
	}
	if !e.isGenerating("火", "土") {
		t.Error("火應生土")
	}
	if e.isGenerating("木", "土") {
		t.Error("木不應生土")
	}
}

func TestEngineIsOvercoming(t *testing.T) {
	e := NewEngine(nil)
	if !e.isOvercoming("木", "土") {
		t.Error("木應剋土")
	}
	if !e.isOvercoming("土", "水") {
		t.Error("土應剋水")
	}
	if e.isOvercoming("木", "火") {
		t.Error("木不應剋火")
	}
}

func TestEngineCalculateXunKong(t *testing.T) {
	e := NewEngine(nil)

	tests := []struct {
		pillar Sexagenary
		want   []string
	}{
		{Sexagenary{StemJia, Zi}, []string{"戌", "亥"}},   // 甲子旬
		{Sexagenary{StemJia, Xu}, []string{"申", "酉"}},   // 甲戌旬
		{Sexagenary{StemJia, Shen}, []string{"午", "未"}}, // 甲申旬
		{Sexagenary{StemJia, Wu}, []string{"辰", "巳"}},   // 甲午旬
		{Sexagenary{StemJia, Chen}, []string{"寅", "卯"}}, // 甲辰旬
		{Sexagenary{StemJia, Yin}, []string{"子", "丑"}},  // 甲寅旬
		{Sexagenary{StemGui, You}, []string{"戌", "亥"}},  // 癸酉旬（甲子）
	}

	for _, tt := range tests {
		t.Run(PillarToString(tt.pillar), func(t *testing.T) {
			got := e.calculateXunKong(tt.pillar)
			if len(got) != 2 || got[0] != tt.want[0] || got[1] != tt.want[1] {
				t.Errorf("calculateXunKong(%+v) = %v, want %v", tt.pillar, got, tt.want)
			}
		})
	}
}

func PillarToString(p Sexagenary) string {
	return StemNames[p.Stem] + BranchNames[p.Branch]
}

func TestEngineCalculateYiMa(t *testing.T) {
	e := NewEngine(nil)

	tests := []struct {
		branch Branch
		want   Branch
	}{
		{Shen, Yin}, {Zi, Yin}, {Chen, Yin},
		{Yin, Shen}, {Wu, Shen}, {Xu, Shen},
		{Si, Hai}, {You, Hai}, {Chou, Hai},
		{Hai, Si}, {Mao, Si}, {Wei, Si},
	}

	for _, tt := range tests {
		t.Run(BranchNames[tt.branch], func(t *testing.T) {
			got := e.calculateYiMa(tt.branch)
			if got != tt.want {
				t.Errorf("calculateYiMa(%s) = %s, want %s", BranchNames[tt.branch], BranchNames[got], BranchNames[tt.want])
			}
		})
	}
}

func TestEngineCalculateTaoHua(t *testing.T) {
	e := NewEngine(nil)

	tests := []struct {
		branch Branch
		want   Branch
	}{
		{Hai, Zi}, {Mao, Zi}, {Wei, Zi},
		{Si, Wu}, {You, Wu}, {Chou, Wu},
		{Yin, Mao}, {Wu, Mao}, {Xu, Mao},
		{Shen, You}, {Zi, You}, {Chen, You},
	}

	for _, tt := range tests {
		t.Run(BranchNames[tt.branch], func(t *testing.T) {
			got := e.calculateTaoHua(tt.branch)
			if got != tt.want {
				t.Errorf("calculateTaoHua(%s) = %s, want %s", BranchNames[tt.branch], BranchNames[got], BranchNames[tt.want])
			}
		})
	}

	// default 分支
	got := e.calculateTaoHua(Branch(99))
	if got != Zi {
		t.Errorf("calculateTaoHua(default) = %s, want 子", BranchNames[got])
	}
}

func TestEngineGetDayGeneralName(t *testing.T) {
	e := NewEngine(nil)
	if e.getDayGeneralName(true) != "晝貴（陽貴）" {
		t.Error("晝貴名稱錯誤")
	}
	if e.getDayGeneralName(false) != "夜貴（陰貴）" {
		t.Error("夜貴名稱錯誤")
	}
}

func TestEngineDetermineKeTi(t *testing.T) {
	e := NewEngine(nil)
	diPan := GetDiPan()

	// 伏吟課
	fuYinTianPan := diPan
	keTi, _ := e.determineKeTi(SanChuan{}, FourKe{}, fuYinTianPan, diPan)
	if keTi != "伏吟課" {
		t.Errorf("應為伏吟課，得到 %s", keTi)
	}

	// 返吟課
	var fanYinTianPan [12]Branch
	for i := 0; i < 12; i++ {
		fanYinTianPan[i] = Branch((int(diPan[i]) + 6) % 12)
	}
	keTi, _ = e.determineKeTi(SanChuan{}, FourKe{}, fanYinTianPan, diPan)
	if keTi != "返吟課" {
		t.Errorf("應為返吟課，得到 %s", keTi)
	}

	// 賊克法課（使用非伏吟天盤）
	tianPan := [12]Branch{Wu, Wei, Shen, You, Xu, Hai, Zi, Chou, Yin, Mao, Chen, Si}
	keTi, _ = e.determineKeTi(SanChuan{Method: "賊克法（第1課）"}, FourKe{}, tianPan, diPan)
	if keTi != "賊克法課" {
		t.Errorf("應為賊克法課，得到 %s", keTi)
	}

	// 普通課（使用非伏吟天盤）
	keTi, _ = e.determineKeTi(SanChuan{Method: "未知法"}, FourKe{}, tianPan, diPan)
	if keTi != "普通課" {
		t.Errorf("應為普通課，得到 %s", keTi)
	}
}

func TestEngineGetKeTiExplanation(t *testing.T) {
	e := NewEngine(nil)

	keTiNames := []string{
		"賊克法課", "比用法課", "涉害法課", "遙克法課",
		"昴星法課", "別責法課", "八專法課",
		"伏吟課", "返吟課", "普通課", "不存在的課",
	}

	for _, name := range keTiNames {
		t.Run(name, func(t *testing.T) {
			exp := e.getKeTiExplanation(name)
			if exp.Name == "" {
				t.Error("課體解說名稱不應為空")
			}
		})
	}
}

func TestEngineJudgeXiongJi(t *testing.T) {
	e := NewEngine(nil)

	// 空亡
	result := &DivinationResult{
		KeTi: "普通課",
		SanChuan: SanChuanInfo{
			Chu: ChuanDetail{IsEmpty: true},
		},
	}
	xj, dy := e.judgeXiongJi(result)
	if xj != "平" {
		t.Errorf("空亡應為平，得到 %s", xj)
	}
	if len(dy) == 0 {
		t.Error("空亡應有斷語")
	}

	// 貴人
	result = &DivinationResult{
		KeTi: "普通課",
		SanChuan: SanChuanInfo{
			Chu: ChuanDetail{General: "貴人"},
		},
	}
	xj, dy = e.judgeXiongJi(result)
	if xj != "吉" {
		t.Errorf("貴人應為吉，得到 %s", xj)
	}

	// 白虎
	result = &DivinationResult{
		KeTi: "普通課",
		SanChuan: SanChuanInfo{
			Chu: ChuanDetail{General: "白虎"},
		},
	}
	xj, dy = e.judgeXiongJi(result)
	if xj != "凶" {
		t.Errorf("白虎應為凶，得到 %s", xj)
	}

	// 伏吟
	result = &DivinationResult{
		KeTi: "伏吟課",
		SanChuan: SanChuanInfo{
			Chu: ChuanDetail{},
		},
	}
	xj, dy = e.judgeXiongJi(result)
	if len(dy) == 0 {
		t.Error("伏吟應有斷語")
	}

	// 返吟
	result = &DivinationResult{
		KeTi: "返吟課",
		SanChuan: SanChuanInfo{
			Chu: ChuanDetail{},
		},
	}
	xj, dy = e.judgeXiongJi(result)
	if len(dy) == 0 {
		t.Error("返吟應有斷語")
	}

	// 平穩
	result = &DivinationResult{
		KeTi: "普通課",
		SanChuan: SanChuanInfo{
			Chu: ChuanDetail{},
		},
	}
	xj, dy = e.judgeXiongJi(result)
	if xj != "平" {
		t.Errorf("普通課應為平，得到 %s", xj)
	}
}

func TestEngineConvertFourKe(t *testing.T) {
	e := NewEngine(nil)
	tianJiang := [12]HeavenlyGeneral{GuiRen, TengShe, ZhuQue, LiuHe, GouChen, QingLong, TianKong, BaiHu, TaiChang, XuanWu, TaiYin, TianHou}
	dunGan := [12]Stem{StemJia, StemYi, StemBing, StemDing, StemWu, StemJi, StemGeng, StemXin, StemRen, StemGui, StemJia, StemYi}

	fourKe := FourKe{
		Ke1: Ke{Down: Zi, Up: Yin},
		Ke2: Ke{Down: Yin, Up: Si},
		Ke3: Ke{Down: Chou, Up: Mao},
		Ke4: Ke{Down: Mao, Up: Wu},
	}

	result := e.convertFourKe(fourKe, tianJiang, StemJia, dunGan)
	if len(result) != 4 {
		t.Fatalf("四課應有4條，得到 %d", len(result))
	}

	if result[0].Number != 1 {
		t.Error("第一課編號錯誤")
	}
	if result[0].Down != "子" || result[0].Up != "寅" {
		t.Errorf("第一課錯誤：%s/%s", result[0].Down, result[0].Up)
	}
	if result[0].UpGeneral != "朱雀" {
		t.Errorf("第一課上神天將錯誤：%s", result[0].UpGeneral)
	}
}

func TestEngineConvertSanChuan(t *testing.T) {
	e := NewEngine(nil)
	tianJiang := [12]HeavenlyGeneral{GuiRen, TengShe, ZhuQue, LiuHe, GouChen, QingLong, TianKong, BaiHu, TaiChang, XuanWu, TaiYin, TianHou}
	dunGan := [12]Stem{StemJia, StemYi, StemBing, StemDing, StemWu, StemJi, StemGeng, StemXin, StemRen, StemGui, StemJia, StemYi}

	sanChuan := SanChuan{
		Chu:    ChuanInfo{Branch: Zi},
		Zhong:  ChuanInfo{Branch: Yin},
		Mo:     ChuanInfo{Branch: Si},
		Method: "賊克法（第1課）",
	}

	result := e.convertSanChuan(sanChuan, tianJiang, StemJia, []string{"戌", "亥"}, Xu, Zi, dunGan)

	if result.Chu.Branch != "子" {
		t.Errorf("初傳錯誤：%s", result.Chu.Branch)
	}
	if result.Chu.General != "貴人" {
		t.Errorf("初傳天將錯誤：%s", result.Chu.General)
	}
	if result.Chu.IsEmpty {
		t.Error("子不應為空亡")
	}
	if result.Chu.IsHorse {
		t.Error("子不應為驛馬")
	}
	if !result.Chu.IsTaoHua {
		t.Error("子應為桃花")
	}
	if result.Zhong.Branch != "寅" {
		t.Errorf("中傳錯誤：%s", result.Zhong.Branch)
	}
	if result.Mo.Branch != "巳" {
		t.Errorf("末傳錯誤：%s", result.Mo.Branch)
	}
	if result.Method != "賊克法（第1課）" {
		t.Errorf("方法錯誤：%s", result.Method)
	}
}

func TestEngineCalculate(t *testing.T) {
	day := true
	ds := &mockDataSource{
		data: &client.CalendarData{
			GregorianDate: "2024-02-15",
			YearPillar:    "甲辰",
			MonthPillar:   "丙寅",
			DayPillar:     "戊午",
			HourPillar:    "己未",
			SolarTermIdx:  3, // 春分
		},
	}

	e := NewEngine(ds)
	req := DivinationRequest{
		Time:         time.Date(2024, 2, 15, 14, 0, 0, 0, time.UTC),
		QuestionType: "求財",
		IsDay:        &day,
	}

	result, err := e.Calculate(req)
	if err != nil {
		t.Fatalf("Calculate 失敗: %v", err)
	}

	if result.DayPillar != "戊午" {
		t.Errorf("日柱錯誤: %s", result.DayPillar)
	}
	if result.HourPillar != "己未" {
		t.Errorf("時柱錯誤: %s", result.HourPillar)
	}
	if result.MonthGeneral != "河魁" {
		t.Errorf("月將錯誤: %s", result.MonthGeneral)
	}
	if result.DayGeneralName != "晝貴（陽貴）" {
		t.Errorf("貴人名稱錯誤: %s", result.DayGeneralName)
	}
	if len(result.FourKe) != 4 {
		t.Errorf("四課數量錯誤: %d", len(result.FourKe))
	}
	if result.SanChuan.Chu.Branch == "" {
		t.Error("初傳不應為空")
	}
	if len(result.VoidBranches) != 2 {
		t.Errorf("空亡數量錯誤: %d", len(result.VoidBranches))
	}
	if result.KeTi == "" {
		t.Error("課體不應為空")
	}
	if result.YongShen.Type == "" {
		t.Error("用神不應為空")
	}
}

func TestEngineCalculateDataSourceError(t *testing.T) {
	ds := &mockDataSource{
		err: fmt.Errorf("mock error"),
	}
	e := NewEngine(ds)
	_, err := e.Calculate(DivinationRequest{Time: time.Now()})
	if err == nil {
		t.Error("應返回錯誤")
	}
}

func TestEngineCalculateGuiRenError(t *testing.T) {
	ds := &mockDataSource{
		data: &client.CalendarData{
			GregorianDate: "2024-02-15",
			YearPillar:    "甲辰",
			MonthPillar:   "丙寅",
			DayPillar:     "戊午",
			HourPillar:    "己未",
			SolarTermIdx:  3,
		},
	}
	e := NewEngine(ds)
	// 使用無效的日干（ Stem(-1) 無法通過 ParsePillar，但這裡我們直接構造一個無法找到貴人的情況）
	// 實際上 ParsePillar 會失敗，讓我們測試 ParsePillar 失敗路徑
	ds.data.DayPillar = "無效"
	_, err := e.Calculate(DivinationRequest{Time: time.Now()})
	if err == nil {
		t.Error("應返回錯誤")
	}
}

func TestEngineGenerateChuanExplanation(t *testing.T) {
	e := NewEngine(nil)

	// 初傳 + 空亡
	ce := e.generateChuanExplanation("初傳", "起因/發端", "子", "貴人", "正印", true, false)
	if ce.Role != "起因/發端" {
		t.Error("Role 錯誤")
	}
	if !contains(ce.KeyPoints, "落空亡，事多虛詐") {
		t.Error("應包含空亡關鍵點")
	}

	// 中傳 + 驛馬
	ce = e.generateChuanExplanation("中傳", "過程/變化", "寅", "青龍", "正財", false, true)
	if !contains(ce.KeyPoints, "逢驛馬，主變動") {
		t.Error("應包含驛馬關鍵點")
	}

	// 末傳 + 吉神
	ce = e.generateChuanExplanation("末傳", "結果/歸宿", "午", "六合", "正官", false, false)
	if !containsSubstring(ce.Explanation, "【歸計·結果】") {
		t.Error("應包含末傳標記")
	}
}

func contains(arr []string, s string) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}
	return false
}

func containsSubstring(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) > 0 && findSub(s, sub))
}

func findSub(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func TestEngineGenerateSmartAdvice(t *testing.T) {
	e := NewEngine(nil)

	tests := []struct {
		name     string
		general  string
		shiShen  string
		isEmpty  bool
		contains string
	}{
		{"初傳", "貴人", "正印", false, "助力"},
		{"初傳", "青龍", "正財", false, "財祿"},
		{"初傳", "六合", "正財", false, "合作"},
		{"中傳", "貴人", "正官", false, "提攜"},
		{"中傳", "白虎", "七殺", false, "血光"},
		{"末傳", "玄武", "偏財", false, "被盜"},
		{"末傳", "螣蛇", "傷官", false, "煩憂"},
		{"初傳", "朱雀", "正官", false, "官司"},
		{"中傳", "勾陳", "正印", false, "平穩"},
		{"末傳", "天空", "比肩", false, "同輩"},
		{"初傳", "貴人", "正印", true, "難以成就"},
		{"中傳", "太常", "食神", false, "才華"},
		{"末傳", "白虎", "傷官", false, "暗中損耗"},
		{"初傳", "玄武", "劫財", false, "困境"},
		{"中傳", "天后", "偏財", false, "審慎"},
	}

	for _, tt := range tests {
		t.Run(tt.general+tt.shiShen, func(t *testing.T) {
			got := e.generateSmartAdvice(tt.name, tt.general, tt.shiShen, tt.isEmpty)
			if !containsSubstring(got, tt.contains) {
				t.Errorf("generateSmartAdvice 應包含 %q，得到: %s", tt.contains, got)
			}
		})
	}
}

func TestEngineChuanTimingAndProcess(t *testing.T) {
	e := NewEngine(nil)

	// getChuChuanTiming
	if e.getChuChuanTiming("正官") == "" {
		t.Error("正官初傳應有時機判斷")
	}
	if e.getChuChuanTiming("七殺") == "" {
		t.Error("七殺初傳應有時機判斷")
	}

	// getZhongChuanProcess
	if e.getZhongChuanProcess("正財") == "" {
		t.Error("正財中傳應有過程判斷")
	}
	if e.getZhongChuanProcess("劫財") == "" {
		t.Error("劫財中傳應有過程判斷")
	}

	// getMoChuanOutcome
	if e.getMoChuanOutcome("貴人", "正官") == "" {
		t.Error("貴人正官末傳應有結果判斷")
	}
	if e.getMoChuanOutcome("白虎", "七殺") == "" {
		t.Error("白虎七殺末傳應有結果判斷")
	}
	if e.getMoChuanOutcome("青龍", "食神") == "" {
		t.Error("青龍食神末傳應有結果判斷")
	}
	// 未匹配組合應返回空
	if e.getMoChuanOutcome("天后", "正印") != "" {
		t.Error("天后正印末傳應無特定結果判斷")
	}
}

func TestEngineGetKeyPointByCombination(t *testing.T) {
	e := NewEngine(nil)

	tests := []struct {
		general string
		shiShen string
		want    string
	}{
		{"貴人", "正官", "貴人提攜"},
		{"青龍", "正財", "青龍送財"},
		{"六合", "比肩", "伙伴同心"},
		{"朱雀", "七殺", "口舌是非"},
		{"白虎", "七殺", "殺氣太重"},
		{"玄武", "偏財", "偏財有險"},
		{"螣蛇", "偏印", "虛驚一場"},
		{"天后", "正印", "天后臨正印"},
	}

	for _, tt := range tests {
		t.Run(tt.general+tt.shiShen, func(t *testing.T) {
			got := e.getKeyPointByCombination(tt.general, tt.shiShen)
			if !containsSubstring(got, tt.want) {
				t.Errorf("getKeyPointByCombination(%s, %s) = %s, 應包含 %s", tt.general, tt.shiShen, got, tt.want)
			}
		})
	}
}

func TestEngineCalculateYongShen(t *testing.T) {
	e := NewEngine(nil)

	result := &DivinationResult{
		SanChuan: SanChuanInfo{
			Chu:   ChuanDetail{Branch: "寅", General: "青龍", Relation: "偏財"},
			Zhong: ChuanDetail{Branch: "午", General: "朱雀", Relation: "正官"},
			Mo:    ChuanDetail{Branch: "戌", General: "白虎", Relation: "七殺"},
		},
		FourKe: []KeInfo{
			{Number: 1, Up: "子", UpGeneral: "貴人", Relation: "正印"},
			{Number: 2, Up: "卯", UpGeneral: "六合", Relation: "劫財"},
			{Number: 3, Up: "巳", UpGeneral: "螣蛇", Relation: "食神"},
			{Number: 4, Up: "酉", UpGeneral: "天空", Relation: "正官"},
		},
	}

	// 財運 - 用神在三傳中找到
	ys := e.calculateYongShen("財運", result, StemJia)
	if ys.SixRelation != "妻財" {
		t.Errorf("財運用神應為妻財，得到 %s", ys.SixRelation)
	}
	if ys.TargetBranch != "寅" {
		t.Errorf("財運用神應在寅，得到 %s", ys.TargetBranch)
	}

	// 官運 - 用神在三傳中找到
	ys = e.calculateYongShen("官運", result, StemJia)
	if ys.SixRelation != "官鬼" {
		t.Errorf("官運用神應為官鬼，得到 %s", ys.SixRelation)
	}
	if ys.TargetBranch != "午" {
		t.Errorf("官運用神應在午，得到 %s", ys.TargetBranch)
	}

	// 考試 - 用神在四課中找到
	ys = e.calculateYongShen("考試", result, StemJia)
	if ys.SixRelation != "父母" {
		t.Errorf("考試用神應為父母，得到 %s", ys.SixRelation)
	}
	if ys.TargetBranch != "子" {
		t.Errorf("考試用神應在子（四課第一課），得到 %s", ys.TargetBranch)
	}

	// 出行 - 無固定六親用神
	ys = e.calculateYongShen("出行", result, StemJia)
	if ys.SixRelation != "" {
		t.Errorf("出行用神應為空，得到 %s", ys.SixRelation)
	}

	// 未知問題類型
	ys = e.calculateYongShen("", result, StemJia)
	if ys.Type != "其他" {
		t.Errorf("空問題類型應為其他，得到 %s", ys.Type)
	}
}

func TestEngineFindYongShenInSanChuan(t *testing.T) {
	e := NewEngine(nil)

	// 三傳中找到用神
	result := &DivinationResult{
		SanChuan: SanChuanInfo{
			Chu:   ChuanDetail{Branch: "寅", General: "青龍", Relation: "偏財", IsEmpty: true, IsHorse: true, IsTaoHua: true},
			Zhong: ChuanDetail{Branch: "午", General: "朱雀", Relation: "正官"},
			Mo:    ChuanDetail{Branch: "戌", General: "白虎", Relation: "七殺"},
		},
		FourKe: []KeInfo{},
	}
	ys := e.findYongShenInSanChuan(YongShen{SixRelation: "偏財"}, result, StemJia)
	if ys.TargetBranch != "寅" {
		t.Errorf("應找到寅，得到 %s", ys.TargetBranch)
	}
	if len(ys.Analysis) < 4 {
		t.Errorf("應包含空亡、驛馬、桃花分析，得到 %d 條", len(ys.Analysis))
	}

	// 三傳中未找到，在四課中找到
	result = &DivinationResult{
		SanChuan: SanChuanInfo{
			Chu:   ChuanDetail{Branch: "寅", General: "青龍", Relation: "偏財"},
			Zhong: ChuanDetail{Branch: "午", General: "朱雀", Relation: "正官"},
			Mo:    ChuanDetail{Branch: "戌", General: "白虎", Relation: "七殺"},
		},
		FourKe: []KeInfo{
			{Number: 1, Up: "子", UpGeneral: "貴人", Relation: "正印"},
		},
	}
	ys = e.findYongShenInSanChuan(YongShen{SixRelation: "正印"}, result, StemJia)
	if ys.TargetBranch != "子" {
		t.Errorf("應在四課找到子，得到 %s", ys.TargetBranch)
	}

	// 四課中也未找到
	ys = e.findYongShenInSanChuan(YongShen{SixRelation: "劫財"}, result, StemJia)
	if ys.TargetBranch != "" {
		t.Error("應未找到用神")
	}
	if len(ys.Analysis) == 0 {
		t.Error("應有提示分析")
	}
}
