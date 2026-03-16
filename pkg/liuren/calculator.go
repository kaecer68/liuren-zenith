package liuren

import (
	"fmt"
	"strings"
	"time"
)

// Calculator 大六壬計算器
type Calculator struct{}

// NewCalculator 創建計算器
func NewCalculator() *Calculator {
	return &Calculator{}
}

// ParsePillar 解析干支字串（如 "甲子"）為 Sexagenary
func ParsePillar(s string) (Sexagenary, error) {
	// 清理字串：移除空白
	s = strings.TrimSpace(s)

	// 使用 rune 處理 UTF-8 中文字符
	runes := []rune(s)
	if len(runes) < 2 {
		return Sexagenary{}, fmt.Errorf("干支格式錯誤（需要至少2個字符）: %s (長度:%d)", s, len(runes))
	}

	stemChar := string(runes[0])
	branchChar := string(runes[1])

	var stem Stem = -1
	var branch Branch = -1

	for i, name := range StemNames {
		if name == stemChar {
			stem = Stem(i)
			break
		}
	}

	for i, name := range BranchNames {
		if name == branchChar {
			branch = Branch(i)
			break
		}
	}

	if stem == -1 {
		return Sexagenary{}, fmt.Errorf("無法解析天干: '%s' (輸入:%s)", stemChar, s)
	}
	if branch == -1 {
		return Sexagenary{}, fmt.Errorf("無法解析地支: '%s' (輸入:%s)", branchChar, s)
	}

	return Sexagenary{Stem: stem, Branch: branch}, nil
}

// CalculateMonthGeneral 計算月將（根據節氣）
// 月將以中氣為換將基準，逆佈十二支
// 正月寅月：亥將（登明）| 二月卯月：戌將（河魁）| ... | 十二月丑月：子將（神後）
func CalculateMonthGeneral(solarTermIdx int) MonthGeneral {
	// 節氣索引 0-23 對應：立春、雨水、驚蟄、春分...小寒、大寒
	// 月將換將在中氣（odd indices: 1, 3, 5, ...）
	// 立春(0)在寅月，月將為亥(登明)
	// 簡化：根據節氣索引推算月將
	// 立春(0)/驚蟄(2)/清明(4)/立夏(6)/芒種(8)/小暑(10)/立秋(12)/白露(14)/寒露(16)/立冬(18)/大雪(20)/小寒(22)
	// 對應月將：亥(0)、戌(1)、酉(2)、申(3)、未(4)、午(5)、巳(6)、辰(7)、卯(8)、寅(9)、丑(10)、子(11)

	// 簡化計算：節氣索引 // 2 得到月份索引 (0-11)，月將 = (11 - 月份索引) % 12
	monthIdx := solarTermIdx / 2
	generalIdx := (11 - monthIdx + 12) % 12
	return MonthGeneral(generalIdx)
}

// CalculateTianPan 計算天盤（月將加時）
// 將月將放置於占時地支上，順佈十二支
func CalculateTianPan(monthGeneral MonthGeneral, hourBranch Branch) [12]Branch {
	var tianPan [12]Branch
	generalBranch := MonthGeneralBranches[monthGeneral]

	// 月將加時：將月將放在占時位置上，其餘順佈
	// 計算偏移量：月將需要移動到占時位置
	// 天盤[占時] = 月將
	// 然後順時針（地支順序）佈置其他

	// 找到月將在標準地支中的位置
	generalPos := int(generalBranch)
	// 找到占時位置
	hourPos := int(hourBranch)

	// 計算偏移：天盤位置 p 對應的地支 = (p - hourPos + generalPos + 12) % 12
	for i := 0; i < 12; i++ {
		// 天盤位置 i 應該放置的地支
		branchIdx := (i - hourPos + generalPos + 12) % 12
		tianPan[i] = Branch(branchIdx)
	}

	return tianPan
}

// GetDiPan 獲取地盤（固定不變）
// 地盤右下為亥，順佈十二支
func GetDiPan() [12]Branch {
	var diPan [12]Branch
	// 地盤固定：子丑寅卯辰巳午未申酉戌亥
	// 但位置安排：右下為亥（索引 0），順時針
	// 標準六壬盤地盤：
	// 巳午未申
	// 辰    酉
	// 卯    戌
	// 寅丑子亥
	// 對應索引：
	// 0: 亥, 1: 子, 2: 丑, 3: 寅
	// 11: 戌          4: 卯
	// 10: 酉          5: 辰
	// 9: 申, 8: 未, 7: 午, 6: 巳
	// 簡化：使用順序排列，位置 0-11 對應地支 亥子丑寅卯辰巳午未申酉戌
	diPan = [12]Branch{Hai, Zi, Chou, Yin, Mao, Chen, Si, Wu, Wei, Shen, You, Xu}
	return diPan
}

// FindGuiRen 尋找貴人位置（起貴人訣）
// 甲戊庚牛羊，乙己鼠猴鄉，丙丁豬雞位，壬癸蛇兔藏，六辛逢馬虎
// 陽貴（晝貴）：從丑順行至未
// 陰貴（夜貴）：從未逆行至丑
func FindGuiRen(dayStem Stem, isDay bool) (Branch, bool) {
	// 起貴人訣
	guiRenMap := map[Stem]struct {
		yang Branch // 陽貴（牛/丑）
		yin  Branch // 陰貴（羊/未）
	}{
		StemJia:  {Chou, Wei}, // 甲：牛羊
		StemWu:   {Chou, Wei}, // 戊：牛羊
		StemGeng: {Chou, Wei}, // 庚：牛羊
		StemYi:   {Zi, Shen},  // 乙：鼠猴
		StemJi:   {Zi, Shen},  // 己：鼠猴
		StemBing: {Hai, You},  // 丙：豬雞
		StemDing: {Hai, You},  // 丁：豬雞
		StemRen:  {Si, Mao},   // 壬：蛇兔
		StemGui:  {Si, Mao},   // 癸：蛇兔
		StemXin:  {Wu, Yin},   // 辛：馬虎
	}

	gr, ok := guiRenMap[dayStem]
	if !ok {
		return 0, false
	}

	if isDay {
		return gr.yang, true // 陽貴（晝貴）
	}
	return gr.yin, true // 陰貴（夜貴）
}

// CalculateTianJiang 計算天將佈局
// 貴人定位後，按順序佈置其他天將
// 晝貴（陽貴）：亥至辰順行（貴螣朱六勾青空白常玄陰后）
// 夜貴（陰貴）：巳至戌逆行（貴螣朱六勾青空白常玄陰后）
func CalculateTianJiang(guiRenPos Branch, isDay bool, tianPan [12]Branch) [12]HeavenlyGeneral {
	var tianJiang [12]HeavenlyGeneral

	// 天將順序：貴人、螣蛇、朱雀、六合、勾陳、青龍、天空、白虎、太常、玄武、太陰、天后
	generals := []HeavenlyGeneral{GuiRen, TengShe, ZhuQue, LiuHe, GouChen, QingLong, TianKong, BaiHu, TaiChang, XuanWu, TaiYin, TianHou}

	if isDay {
		// 晝貴順行：從貴人位置開始，順時針佈置
		startIdx := int(guiRenPos)
		for i, general := range generals {
			pos := (startIdx + i) % 12
			tianJiang[pos] = general
		}
	} else {
		// 夜貴逆行：從貴人位置開始，逆時針佈置
		startIdx := int(guiRenPos)
		for i, general := range generals {
			pos := (startIdx - i + 12) % 12
			tianJiang[pos] = general
		}
	}

	return tianJiang
}

// CalculateFourKe 計算四課
func CalculateFourKe(dayStem Stem, dayBranch Branch, tianPan [12]Branch) FourKe {
	// 日干寄宮
	dayStemAttachment := StemAttachment[dayStem]

	// 第一課：日干寄宮之上神
	ke1Down := dayStemAttachment
	ke1Up := tianPan[ke1Down]

	// 第二課：第一課上神之上神
	ke2Down := ke1Up
	ke2Up := tianPan[ke2Down]

	// 第三課：日支之上神
	ke3Down := dayBranch
	ke3Up := tianPan[ke3Down]

	// 第四課：第三課上神之上神
	ke4Down := ke3Up
	ke4Up := tianPan[ke4Down]

	return FourKe{
		Ke1: Ke{Down: ke1Down, Up: ke1Up},
		Ke2: Ke{Down: ke2Down, Up: ke2Up},
		Ke3: Ke{Down: ke3Down, Up: ke3Up},
		Ke4: Ke{Down: ke4Down, Up: ke4Up},
	}
}

// IsDayTime 判斷是否為白天（用於晝夜貴人）
// 簡化：6:00-18:00 為晝
func IsDayTime(t time.Time) bool {
	hour := t.Hour()
	return hour >= 6 && hour < 18
}
