package liuren

import (
	"testing"
)

func TestParsePillar(t *testing.T) {
	tests := []struct {
		input    string
		expected Sexagenary
		wantErr  bool
	}{
		{"甲子", Sexagenary{Stem: StemJia, Branch: Zi}, false},
		{"乙丑", Sexagenary{Stem: StemYi, Branch: Chou}, false},
		{"丙寅", Sexagenary{Stem: StemBing, Branch: Yin}, false},
		{"丁卯", Sexagenary{Stem: StemDing, Branch: Mao}, false},
		{"戊辰", Sexagenary{Stem: StemWu, Branch: Chen}, false},
		{"己巳", Sexagenary{Stem: StemJi, Branch: Si}, false},
		{"庚午", Sexagenary{Stem: StemGeng, Branch: Wu}, false},
		{"辛未", Sexagenary{Stem: StemXin, Branch: Wei}, false},
		{"壬申", Sexagenary{Stem: StemRen, Branch: Shen}, false},
		{"癸酉", Sexagenary{Stem: StemGui, Branch: You}, false},
		{"甲戌", Sexagenary{Stem: StemJia, Branch: Xu}, false},
		{"乙亥", Sexagenary{Stem: StemYi, Branch: Hai}, false},
		{"己丑", Sexagenary{Stem: StemJi, Branch: Chou}, false}, // 這個是問題案例
		{"", Sexagenary{}, true},
		{"甲", Sexagenary{}, true},
		{"invalid", Sexagenary{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParsePillar(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePillar(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && (got.Stem != tt.expected.Stem || got.Branch != tt.expected.Branch) {
				t.Errorf("ParsePillar(%q) = {Stem:%v, Branch:%v}, want {Stem:%v, Branch:%v}",
					tt.input, got.Stem, got.Branch, tt.expected.Stem, tt.expected.Branch)
			}
		})
	}
}

func TestStemNames(t *testing.T) {
	t.Logf("StemNames: %v", StemNames)
	for i, name := range StemNames {
		t.Logf("Stem %d: %s (rune: %v)", i, name, []rune(name))
	}
}

func TestBranchNames(t *testing.T) {
	t.Logf("BranchNames: %v", BranchNames)
	for i, name := range BranchNames {
		t.Logf("Branch %d: %s (rune: %v)", i, name, []rune(name))
	}
}

func TestCalculateMonthGeneral(t *testing.T) {
	tests := []struct {
		solarTermIdx int
		want         MonthGeneral
	}{
		{0, ShenHou},     // 立春，沿用大寒中氣 → 子
		{1, DengMing},    // 雨水 → 亥
		{2, DengMing},    // 驚蟄，沿用雨水 → 亥
		{3, HeKui},       // 春分 → 戌
		{4, HeKui},       // 清明，沿用春分 → 戌
		{5, CongKui},     // 谷雨 → 酉
		{6, CongKui},     // 立夏，沿用谷雨 → 酉
		{7, ChuanSong},   // 小滿 → 申
		{8, ChuanSong},   // 芒種，沿用小滿 → 申
		{9, XiaoJi},      // 夏至 → 未
		{10, XiaoJi},     // 小暑，沿用夏至 → 未
		{11, ShengGuang}, // 大暑 → 午
		{12, ShengGuang}, // 立秋，沿用大暑 → 午
		{13, TaiYi},      // 處暑 → 巳
		{14, TaiYi},      // 白露，沿用處暑 → 巳
		{15, TianGang},   // 秋分 → 辰
		{16, TianGang},   // 寒露，沿用秋分 → 辰
		{17, TaiChong},   // 霜降 → 卯
		{18, TaiChong},   // 立冬，沿用霜降 → 卯
		{19, GongCao},    // 小雪 → 寅
		{20, GongCao},    // 大雪，沿用小雪 → 寅
		{21, DaJi},       // 冬至 → 丑
		{22, DaJi},       // 小寒，沿用冬至 → 丑
		{23, ShenHou},    // 大寒 → 子
	}

	for _, tt := range tests {
		t.Run(MonthGeneralNames[tt.want], func(t *testing.T) {
			got := CalculateMonthGeneral(tt.solarTermIdx)
			if got != tt.want {
				t.Errorf("CalculateMonthGeneral(%d) = %v (%s), want %v (%s)",
					tt.solarTermIdx, got, MonthGeneralNames[got], tt.want, MonthGeneralNames[tt.want])
			}
		})
	}
}

func TestGetDiPan(t *testing.T) {
	diPan := GetDiPan()
	// 驗證地盤索引約定：0=亥(右下)，順時針
	want := [12]Branch{Hai, Zi, Chou, Yin, Mao, Chen, Si, Wu, Wei, Shen, You, Xu}
	if diPan != want {
		t.Errorf("GetDiPan() = %v, want %v", diPan, want)
	}
}

func TestCalculateTianPan(t *testing.T) {
	// 範例：亥將加子時，天盤應與地盤錯位對應
	// 地盤子位是索引1，應放亥；丑位索引2放子...
	tianPan := CalculateTianPan(DengMing, Zi)
	// diPan[1]=子，應放亥(11)
	if tianPan[1] != Hai {
		t.Errorf("CalculateTianPan(亥將, 子時)[1] = %v, want 亥", BranchNames[tianPan[1]])
	}
	// diPan[2]=丑，應放子(0)
	if tianPan[2] != Zi {
		t.Errorf("CalculateTianPan(亥將, 子時)[2] = %v, want 子", BranchNames[tianPan[2]])
	}
	// diPan[0]=亥，應放戌(10)
	if tianPan[0] != Xu {
		t.Errorf("CalculateTianPan(亥將, 子時)[0] = %v, want 戌", BranchNames[tianPan[0]])
	}

	// 範例：戌將加卯時
	// diPan[4]=卯，應放戌(10)
	// diPan[5]=辰，應放亥(11)
	// diPan[3]=寅，應放酉(9)
	tianPan2 := CalculateTianPan(HeKui, Mao)
	if tianPan2[4] != Xu {
		t.Errorf("CalculateTianPan(戌將, 卯時)[4] = %v, want 戌", BranchNames[tianPan2[4]])
	}
	if tianPan2[5] != Hai {
		t.Errorf("CalculateTianPan(戌將, 卯時)[5] = %v, want 亥", BranchNames[tianPan2[5]])
	}
	if tianPan2[3] != You {
		t.Errorf("CalculateTianPan(戌將, 卯時)[3] = %v, want 酉", BranchNames[tianPan2[3]])
	}
}

func TestCalculateDunGan(t *testing.T) {
	// 甲子日，旬首為子
	dayPillar := Sexagenary{StemJia, Zi}
	// 構造一個天盤：亥將加子時
	// 實際天盤：[戌, 亥, 子, 丑, 寅, 卯, 辰, 巳, 午, 未, 申, 酉]
	// 旬首地支「子」在天盤中的位置是索引 2
	tianPan := CalculateTianPan(DengMing, Zi)
	dunGan := CalculateDunGan(dayPillar, tianPan)

	// 索引2應為甲，索引3應為乙，索引4應為丙...
	if dunGan[2] != StemJia {
		t.Errorf("CalculateDunGan(甲子)[2] = %v, want 甲", StemNames[dunGan[2]])
	}
	if dunGan[3] != StemYi {
		t.Errorf("CalculateDunGan(甲子)[3] = %v, want 乙", StemNames[dunGan[3]])
	}
	if dunGan[4] != StemBing {
		t.Errorf("CalculateDunGan(甲子)[4] = %v, want 丙", StemNames[dunGan[4]])
	}

	// 驗證循環：索引11應為癸(9)
	if dunGan[11] != StemGui {
		t.Errorf("CalculateDunGan(甲子)[11] = %v, want 癸", StemNames[dunGan[11]])
	}
	// 索引0應為甲（10%10=0），索引1應為乙（11%10=1）
	if dunGan[0] != StemJia {
		t.Errorf("CalculateDunGan(甲子)[0] = %v, want 甲", StemNames[dunGan[0]])
	}
	if dunGan[1] != StemYi {
		t.Errorf("CalculateDunGan(甲子)[1] = %v, want 乙", StemNames[dunGan[1]])
	}
}

func TestCalculateFourKe(t *testing.T) {
	// 範例：癸酉日、子時、亥將（來源：《六壬指南》）
	// 癸寄丑，丑上天盤應為子
	// 子上天盤應為亥
	// 酉上天盤應為申
	// 申上天盤應為未
	// 天盤（亥將加子時）：
	// 0:戌 1:亥 2:子 3:丑 4:寅 5:卯 6:辰 7:巳 8:午 9:未 10:申 11:酉
	tianPan := [12]Branch{Xu, Hai, Zi, Chou, Yin, Mao, Chen, Si, Wu, Wei, Shen, You}
	diPan := GetDiPan()
	fourKe := CalculateFourKe(StemGui, You, tianPan, diPan)

	if fourKe.Ke1.Down != Chou || fourKe.Ke1.Up != Zi {
		t.Errorf("Ke1 want 丑上子, got %s上%s", BranchNames[fourKe.Ke1.Down], BranchNames[fourKe.Ke1.Up])
	}
	if fourKe.Ke2.Down != Zi || fourKe.Ke2.Up != Hai {
		t.Errorf("Ke2 want 子上亥, got %s上%s", BranchNames[fourKe.Ke2.Down], BranchNames[fourKe.Ke2.Up])
	}
	if fourKe.Ke3.Down != You || fourKe.Ke3.Up != Shen {
		t.Errorf("Ke3 want 酉上申, got %s上%s", BranchNames[fourKe.Ke3.Down], BranchNames[fourKe.Ke3.Up])
	}
	if fourKe.Ke4.Down != Shen || fourKe.Ke4.Up != Wei {
		t.Errorf("Ke4 want 申上未, got %s上%s", BranchNames[fourKe.Ke4.Down], BranchNames[fourKe.Ke4.Up])
	}
}

func TestCalculateTianJiang(t *testing.T) {
	tianPan := GetDiPan()
	tj1 := CalculateTianJiang(Chou, true, tianPan)
	if tj1[Chou] != GuiRen || tj1[Mao] != ZhuQue || tj1[Chen] != LiuHe {
		t.Errorf("貴人在丑順行錯誤")
	}

	tj2 := CalculateTianJiang(Wu, false, tianPan)
	if tj2[Wu] != GuiRen || tj2[Si] != TengShe || tj2[Mao] != LiuHe {
		t.Errorf("貴人在午逆行錯誤")
	}
}

func TestIsKe(t *testing.T) {
	calc := &JiuZongMenCalculator{}

	// 正確相剋
	if !calc.isKe(Zi, Wu) {
		t.Error("子應克午")
	}
	if !calc.isKe(Chou, Zi) {
		t.Error("丑應克子")
	}
	if !calc.isKe(Yin, Chou) {
		t.Error("寅應克丑")
	}
	if !calc.isKe(Mao, Chen) {
		t.Error("卯應克辰")
	}
	if !calc.isKe(Si, Shen) {
		t.Error("巳應克申")
	}
	if !calc.isKe(Wu, You) {
		t.Error("午應克酉")
	}
	if !calc.isKe(Shen, Yin) {
		t.Error("申應克寅")
	}
	if !calc.isKe(You, Mao) {
		t.Error("酉應克卯")
	}
	if !calc.isKe(Hai, Si) {
		t.Error("亥應克巳")
	}

	// 不應相剋
	if calc.isKe(Zi, Chou) {
		t.Error("子不應克丑")
	}
	if calc.isKe(Yin, Zi) {
		t.Error("寅不應克子")
	}
	if calc.isKe(Si, Zi) {
		t.Error("巳不應克子")
	}
}

func TestCalculateXunKong(t *testing.T) {
	engine := &Engine{}

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
		{Sexagenary{StemGui, You}, []string{"戌", "亥"}},  // 癸酉，甲子旬
		{Sexagenary{StemGui, Wei}, []string{"申", "酉"}},  // 癸未，甲戌旬
		{Sexagenary{StemBing, Wu}, []string{"寅", "卯"}},  // 丙午，甲辰旬
	}

	for _, tt := range tests {
		t.Run(tt.pillar.String(), func(t *testing.T) {
			got := engine.calculateXunKong(tt.pillar)
			if len(got) != 2 || got[0] != tt.want[0] || got[1] != tt.want[1] {
				t.Errorf("calculateXunKong(%s) = %v, want %v", tt.pillar.String(), got, tt.want)
			}
		})
	}
}
