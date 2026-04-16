package liuren

import (
	"testing"
)

func TestZeiKeMethod(t *testing.T) {
	// 範例：第四課巳火克酉金（上克下）→ 元首課，初傳為巳
	// 確保 Ke1-Ke3 無克，只有 Ke4 有上克下
	tianPan := [12]Branch{Wu, Wei, Shen, You, Xu, Hai, Zi, Chou, Yin, Mao, Chen, Si}
	diPan := GetDiPan()
	fourKe := FourKe{
		Ke1: Ke{Down: Wei, Up: Zi},    // 無克
		Ke2: Ke{Down: Zi, Up: Shen},   // 無克
		Ke3: Ke{Down: Chou, Up: Shen}, // 無克
		Ke4: Ke{Down: You, Up: Si},    // 上克下：巳克酉
	}
	calc := NewJiuZongMenCalculator(StemDing, Chou, tianPan, diPan, fourKe)

	chuan, ok := calc.ZeiKeMethod()
	if !ok {
		t.Fatal("賊克法應找到初傳")
	}
	if chuan.Chu.Branch != Si {
		t.Errorf("初傳應為巳，得到 %s", BranchNames[chuan.Chu.Branch])
	}
}

func TestBiYongMethod(t *testing.T) {
	// 比用法測試：多個克中，只有一個與日干陰陽相比
	// 戊為陽日，Ke2 上克下（卯克辰，up=卯為陰，不匹配），Ke4 下克上（子克午，up=午為陽，匹配）
	tianPan := [12]Branch{Si, Wu, Wei, Shen, You, Xu, Hai, Zi, Chou, Yin, Mao, Chen}
	diPan := GetDiPan()
	fourKe := FourKe{
		Ke1: Ke{Down: Si, Up: Chou},  // 無克
		Ke2: Ke{Down: Chen, Up: Mao}, // 上克下：卯克辰，up=卯(陰)，不匹配
		Ke3: Ke{Down: Yin, Up: Wu},   // 無克
		Ke4: Ke{Down: Zi, Up: Wu},    // 下克上：子克午，up=午(陽)，匹配
	}
	calc := NewJiuZongMenCalculator(StemWu, Yin, tianPan, diPan, fourKe)

	chuan, ok := calc.BiYongMethod()
	if !ok {
		t.Fatal("比用法應找到初傳")
	}
	if chuan.Chu.Branch != Wu {
		t.Errorf("初傳應為午，得到 %s", BranchNames[chuan.Chu.Branch])
	}
}

func TestMaoXingMethod(t *testing.T) {
	// 剛日取酉上神，柔日取卯上神
	// 構造一個 tianPan，使得地盤酉位(10)天盤為申，地盤卯位(4)天盤為寅
	tianPan := [12]Branch{Hai, Zi, Chou, Yin, Yin, Chen, Si, Wu, Wei, Shen, Shen, Xu}
	diPan := GetDiPan()
	fourKe := CalculateFourKe(StemJia, Zi, tianPan, diPan)

	// 陽日（甲日）
	calcYang := NewJiuZongMenCalculator(StemJia, Zi, tianPan, diPan, fourKe)
	chuanYang, ok := calcYang.MaoXingMethod()
	if !ok {
		t.Fatal("昴星法應找到初傳")
	}
	// 天盤酉(10)位上是申
	if chuanYang.Chu.Branch != Shen {
		t.Errorf("陽日昴星初傳應為申，得到 %s", BranchNames[chuanYang.Chu.Branch])
	}

	// 陰日（乙日）
	calcYin := NewJiuZongMenCalculator(StemYi, Zi, tianPan, diPan, fourKe)
	chuanYin, ok := calcYin.MaoXingMethod()
	if !ok {
		t.Fatal("昴星法應找到初傳")
	}
	// 天盤卯(4)位上是寅
	if chuanYin.Chu.Branch != Yin {
		t.Errorf("陰日昴星初傳應為寅，得到 %s", BranchNames[chuanYin.Chu.Branch])
	}
}

func TestBieZeMethod(t *testing.T) {
	// 構造 tianPan：地盤丑位(2)天盤為子（干合丑，丑上天盤為子）
	tianPan := [12]Branch{Hai, Zi, Zi, Yin, Mao, Chen, Si, Wu, Wei, Shen, You, Xu}
	diPan := GetDiPan()
	fourKe := CalculateFourKe(StemJia, Zi, tianPan, diPan)

	// 陽日：甲日，干合為己→丑，丑上天盤為子
	calc := NewJiuZongMenCalculator(StemJia, Zi, tianPan, diPan, fourKe)
	chuan, ok := calc.BieZeMethod()
	if !ok {
		t.Fatal("別責法應找到初傳")
	}
	if chuan.Chu.Branch != Zi {
		t.Errorf("陽日別責初傳應為子，得到 %s", BranchNames[chuan.Chu.Branch])
	}
}

func TestBaZhuanMethod(t *testing.T) {
	// 構造 tianPan：地盤辰位(4)天盤為辰，地盤戌位(10)天盤為戌
	tianPan := [12]Branch{Hai, Zi, Chou, Yin, Chen, Si, Wu, Wei, Shen, You, Xu, Mao}
	diPan := GetDiPan()
	fourKe := CalculateFourKe(StemJia, Zi, tianPan, diPan)

	// 甲寅日（陽八專日）
	calc := NewJiuZongMenCalculator(StemJia, Yin, tianPan, diPan, fourKe)
	chuan, ok := calc.BaZhuanMethod()
	if !ok {
		t.Fatal("八專法應找到初傳")
	}
	// 陽日取辰上神 = 辰位(4)天盤 = 辰
	if chuan.Chu.Branch != Chen {
		t.Errorf("陽日八專初傳應為辰，得到 %s", BranchNames[chuan.Chu.Branch])
	}

	// 乙卯日（陰八專日）
	calc2 := NewJiuZongMenCalculator(StemYi, Mao, tianPan, diPan, fourKe)
	chuan2, ok := calc2.BaZhuanMethod()
	if !ok {
		t.Fatal("八專法應找到初傳")
	}
	// 陰日取戌上神 = 戌位(10)天盤 = 戌
	if chuan2.Chu.Branch != Xu {
		t.Errorf("陰日八專初傳應為戌，得到 %s", BranchNames[chuan2.Chu.Branch])
	}
}

func TestFuYinMethod(t *testing.T) {
	// 伏吟課：天盤地盤相同
	diPan := GetDiPan()
	fourKe := CalculateFourKe(StemJia, Zi, diPan, diPan)
	calc := NewJiuZongMenCalculator(StemJia, Zi, diPan, diPan, fourKe)

	chuan, ok := calc.FuYinMethod()
	if !ok {
		t.Fatal("伏吟法應找到初傳")
	}
	// 甲日寄寅，寅刑巳
	if chuan.Chu.Branch != Si {
		t.Errorf("伏吟初傳應為巳（寅刑巳），得到 %s", BranchNames[chuan.Chu.Branch])
	}

	// 非伏吟課應返回 false
	tianPan := [12]Branch{Wu, Wei, Shen, You, Xu, Hai, Zi, Chou, Yin, Mao, Chen, Si}
	calc2 := NewJiuZongMenCalculator(StemJia, Zi, tianPan, diPan, fourKe)
	_, ok2 := calc2.FuYinMethod()
	if ok2 {
		t.Error("非伏吟課應返回 false")
	}
}

func TestFanYinMethod(t *testing.T) {
	// 返吟課：天盤地盤對沖
	diPan := GetDiPan()
	// 構造對沖天盤
	var tianPan [12]Branch
	for i := 0; i < 12; i++ {
		tianPan[i] = Branch((int(diPan[i]) + 6) % 12)
	}
	fourKe := CalculateFourKe(StemJia, Zi, tianPan, diPan)
	calc := NewJiuZongMenCalculator(StemJia, Zi, tianPan, diPan, fourKe)

	chuan := calc.FanYinMethod()
	// 甲日支為子，驛馬在寅
	if chuan.Chu.Branch != Yin {
		t.Errorf("返吟初傳應為寅（驛馬），得到 %s", BranchNames[chuan.Chu.Branch])
	}
}

func TestGetXingShen(t *testing.T) {
	calc := &JiuZongMenCalculator{}

	tests := []struct {
		branch Branch
		want   Branch
	}{
		{Yin, Si},    // 寅刑巳
		{Si, Shen},   // 巳刑申
		{Shen, Yin},  // 申刑寅
		{Chou, Xu},   // 丑刑戌
		{Xu, Wei},    // 戌刑未
		{Wei, Chou},  // 未刑丑
		{Zi, Mao},    // 子刑卯
		{Mao, Zi},    // 卯刑子
		{Chen, Chen}, // 辰自刑
		{Wu, Wu},     // 午自刑
		{You, You},   // 酉自刑
		{Hai, Hai},   // 亥自刑
	}

	for _, tt := range tests {
		t.Run(BranchNames[tt.branch], func(t *testing.T) {
			got := calc.getXingShen(tt.branch)
			if got != tt.want {
				t.Errorf("getXingShen(%s) = %s, want %s", BranchNames[tt.branch], BranchNames[got], BranchNames[tt.want])
			}
		})
	}
}

func TestCalculateSanChuanPriority(t *testing.T) {
	// 驗證九宗門優先順序：賊克 > 比用 > 涉害 > ...
	diPan := GetDiPan()
	tianPan := [12]Branch{Hai, Zi, Chou, Yin, Mao, Chen, Si, Wu, Wei, Shen, You, Xu}

	// 構造一個只有第4課有克的四課
	fourKe := FourKe{
		Ke1: Ke{Down: Zi, Up: Shen},  // 無克
		Ke2: Ke{Down: Shen, Up: Zi},  // 無克
		Ke3: Ke{Down: Chou, Up: Wu},  // 無克
		Ke4: Ke{Down: Mao, Up: Chen}, // 下克上：卯克辰
	}
	calc := NewJiuZongMenCalculator(StemJia, Zi, tianPan, diPan, fourKe)
	chuan := calc.CalculateSanChuan()
	if chuan.Method != "賊克法（第4課）" {
		t.Errorf("應為賊克法（第4課），得到 %s", chuan.Method)
	}
}
