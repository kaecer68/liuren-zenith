package liuren

// JiuZongMenCalculator 九宗門計算器
type JiuZongMenCalculator struct {
	DayStem   Stem
	DayBranch Branch
	TianPan   [12]Branch
	DiPan     [12]Branch
	FourKe    FourKe
}

// NewJiuZongMenCalculator 創建九宗門計算器
func NewJiuZongMenCalculator(dayStem Stem, dayBranch Branch, tianPan, diPan [12]Branch, fourKe FourKe) *JiuZongMenCalculator {
	return &JiuZongMenCalculator{
		DayStem:   dayStem,
		DayBranch: dayBranch,
		TianPan:   tianPan,
		DiPan:     diPan,
		FourKe:    fourKe,
	}
}

// CalculateSanChuan 計算三傳（自動選擇九宗門）
func (c *JiuZongMenCalculator) CalculateSanChuan() SanChuan {
	// 優先順序：賊克 → 比用 → 涉害 → 遙克 → 昴星 → 別責 → 八專 → 伏吟 → 返吟

	// 1. 嘗試賊克法
	if chuan, ok := c.ZeiKeMethod(); ok {
		return chuan
	}

	// 2. 嘗試比用法
	if chuan, ok := c.BiYongMethod(); ok {
		return chuan
	}

	// 3. 嘗試涉害法
	if chuan, ok := c.SheHaiMethod(); ok {
		return chuan
	}

	// 4. 嘗試遙克法
	if chuan, ok := c.YaoKeMethod(); ok {
		return chuan
	}

	// 5. 嘗試昴星法
	if chuan, ok := c.MaoXingMethod(); ok {
		return chuan
	}

	// 6. 嘗試別責法
	if chuan, ok := c.BieZeMethod(); ok {
		return chuan
	}

	// 7. 嘗試八專法
	if chuan, ok := c.BaZhuanMethod(); ok {
		return chuan
	}

	// 8. 嘗試伏吟法
	if chuan, ok := c.FuYinMethod(); ok {
		return chuan
	}

	// 9. 返吟法（最後的手段）
	return c.FanYinMethod()
}

// ZeiKeMethod 賊克法：以下克上為重，以上克下為輕
// 四課中若有下克上，取相克者為初傳
func (c *JiuZongMenCalculator) ZeiKeMethod() (SanChuan, bool) {
	// 優先找下克上（賊）
	keList := []Ke{c.FourKe.Ke1, c.FourKe.Ke2, c.FourKe.Ke3, c.FourKe.Ke4}

	// 下克上為賊（重）
	for i, ke := range keList {
		if c.isKe(ke.Down, ke.Up) { // 下克上
			return c.makeSanChuan(ke.Up, ZeiKe, i+1), true
		}
	}

	// 次找上克下（克）
	for i, ke := range keList {
		if c.isKe(ke.Up, ke.Down) { // 上克下
			return c.makeSanChuan(ke.Up, ZeiKe, i+1), true
		}
	}

	return SanChuan{}, false
}

// BiYongMethod 比用法（知一法）：兩賊或兩克，取與日干陰陽相比者
func (c *JiuZongMenCalculator) BiYongMethod() (SanChuan, bool) {
	keList := []Ke{c.FourKe.Ke1, c.FourKe.Ke2, c.FourKe.Ke3, c.FourKe.Ke4}

	var candidates []struct {
		ke    Ke
		idx   int
		isBi  bool // 是否相比（同陰陽）
	}

	// 收集所有賊/克
	for i, ke := range keList {
		isKe := c.isKe(ke.Down, ke.Up) || c.isKe(ke.Up, ke.Down)
		if isKe {
			// 比用：日干與上神陰陽相同
			dayStemYinYang := c.DayStem%2 == 0 // 甲丙戊庚壬為陽
			upYinYang := ke.Up%2 == 0          // 子寅辰午申戌為陽
			candidates = append(candidates, struct {
				ke   Ke
				idx  int
				isBi bool
			}{ke, i + 1, dayStemYinYang == upYinYang})
		}
	}

	// 只有一個相比的，取之
	biCount := 0
	var selected struct {
		ke  Ke
		idx int
	}
	for _, candi := range candidates {
		if candi.isBi {
			biCount++
			selected = struct {
				ke  Ke
				idx int
			}{candi.ke, candi.idx}
		}
	}

	if biCount == 1 {
		return c.makeSanChuan(selected.ke.Up, BiYong, selected.idx), true
	}

	return SanChuan{}, false
}

// SheHaiMethod 涉害法：俱比或俱不比，以涉害深者為用
// 涉害深淺 = 該地支在地盤上受克的次數
func (c *JiuZongMenCalculator) SheHaiMethod() (SanChuan, bool) {
	keList := []Ke{c.FourKe.Ke1, c.FourKe.Ke2, c.FourKe.Ke3, c.FourKe.Ke4}

	var candidates []struct {
		ke       Ke
		idx      int
		sheHai   int // 涉害深淺
		isSiMeng bool // 是否四孟
		isSiZhong bool // 是否四仲
	}

	// 收集所有賊/克並計算涉害
	for i, ke := range keList {
		isKe := c.isKe(ke.Down, ke.Up) || c.isKe(ke.Up, ke.Down)
		if isKe {
			sheHai := c.calculateSheHai(ke.Up)
			candidates = append(candidates, struct {
				ke       Ke
				idx      int
				sheHai   int
				isSiMeng bool
				isSiZhong bool
			}{
				ke:       ke,
				idx:      i + 1,
				sheHai:   sheHai,
				isSiMeng: c.isSiMeng(ke.Up),
				isSiZhong: c.isSiZhong(ke.Up),
			})
		}
	}

	if len(candidates) == 0 {
		return SanChuan{}, false
	}

	// 找涉害最深者
	maxSheHai := -1
	var deepest []struct {
		ke  Ke
		idx int
		isSiMeng bool
		isSiZhong bool
	}

	for _, candi := range candidates {
		if candi.sheHai > maxSheHai {
			maxSheHai = candi.sheHai
			deepest = []struct {
				ke  Ke
				idx int
				isSiMeng bool
				isSiZhong bool
			}{{candi.ke, candi.idx, candi.isSiMeng, candi.isSiZhong}}
		} else if candi.sheHai == maxSheHai {
			deepest = append(deepest, struct {
				ke  Ke
				idx int
				isSiMeng bool
				isSiZhong bool
			}{candi.ke, candi.idx, candi.isSiMeng, candi.isSiZhong})
		}
	}

	// 只有一個最深，取之
	if len(deepest) == 1 {
		return c.makeSanChuan(deepest[0].ke.Up, SheHai, deepest[0].idx), true
	}

	// 涉害相等，取四孟上神
	for _, d := range deepest {
		if d.isSiMeng {
			return c.makeSanChuan(d.ke.Up, SheHai, d.idx), true
		}
	}

	// 無四孟，取四仲上神
	for _, d := range deepest {
		if d.isSiZhong {
			return c.makeSanChuan(d.ke.Up, SheHai, d.idx), true
		}
	}

	// 四仲也無（都是四季），取第一個（實務上較少見）
	return c.makeSanChuan(deepest[0].ke.Up, SheHai, deepest[0].idx), true
}

// YaoKeMethod 遙克法：四課俱無賊克，取神遙克日或日遙克神
func (c *JiuZongMenCalculator) YaoKeMethod() (SanChuan, bool) {
	keList := []Ke{c.FourKe.Ke1, c.FourKe.Ke2, c.FourKe.Ke3, c.FourKe.Ke4}
	dayAttachment := StemAttachment[c.DayStem]

	// 優先：神遙克日（上神克日干寄宮）
	for i, ke := range keList {
		if c.isKe(ke.Up, dayAttachment) {
			return c.makeSanChuan(ke.Up, YaoKe, i+1), true
		}
	}

	// 其次：日遙克神（日干寄宮克上神）
	for i, ke := range keList {
		if c.isKe(dayAttachment, ke.Up) {
			return c.makeSanChuan(ke.Up, YaoKe, i+1), true
		}
	}

	return SanChuan{}, false
}

// MaoXingMethod 昴星法：剛日取酉上神，柔日取卯上神
func (c *JiuZongMenCalculator) MaoXingMethod() (SanChuan, bool) {
	// 檢查是否有四課齊全（無缺一）
	// 簡化：直接判斷剛柔
	dayStemYang := c.DayStem%2 == 0 // 甲丙戊庚壬為剛（陽）

	var chu Branch
	if dayStemYang {
		// 剛日：從天盤酉(10)位置取上神
		chuanPos := c.TianPan[You]
		chu = chuanPos
	} else {
		// 柔日：從天盤卯(4)位置取上神
		chuanPos := c.TianPan[Mao]
		chu = chuanPos
	}

	// 中傳：剛日取支上神，柔日取干上神
	var zhong Branch
	if dayStemYang {
		zhong = c.FourKe.Ke3.Up // 支上神
	} else {
		zhong = c.FourKe.Ke1.Up // 干上神
	}

	// 末傳：根據中傳再取
	mo := c.TianPan[zhong]

	return SanChuan{
		Chu:  ChuanInfo{Branch: chu},
		Zhong: ChuanInfo{Branch: zhong},
		Mo:   ChuanInfo{Branch: mo},
		Method: "昴星法",
	}, true
}

// BieZeMethod 別責法：四課缺一（無克），剛日取干合，柔日取支前三合
func (c *JiuZongMenCalculator) BieZeMethod() (SanChuan, bool) {
	// 簡化實現：直接判斷剛柔
	dayStemYang := c.DayStem%2 == 0

	var chu Branch
	if dayStemYang {
		// 剛日（陽日）：取日干合（干合：甲己、乙庚、丙辛、丁壬、戊癸）
		chu = c.getStemHe(c.DayStem)
	} else {
		// 柔日（陰日）：取日支前三合
		chu = c.getSanHe(c.DayBranch)
	}

	// 中末傳：俱取干上神
	zhong := c.FourKe.Ke1.Up
	mo := c.FourKe.Ke1.Up

	return SanChuan{
		Chu:  ChuanInfo{Branch: chu},
		Zhong: ChuanInfo{Branch: zhong},
		Mo:   ChuanInfo{Branch: mo},
		Method: "別責法",
	}, true
}

// BaZhuanMethod 八專法：八專日（干支同位），取剛柔比用
func (c *JiuZongMenCalculator) BaZhuanMethod() (SanChuan, bool) {
	// 八專日：干支同位（如甲寅、乙卯、丁未、己未、庚申、辛酉...）
	// 簡化：直接給出結果

	dayStemYang := c.DayStem%2 == 0

	var chu Branch
	if dayStemYang {
		// 剛日：取天盤魁（戌）或罡（辰）
		chu = Xu // 簡化取戌
	} else {
		// 柔日：取天盤從魁（酉）或大吉（丑）
		chu = You // 簡化取酉
	}

	// 中末傳：根據初傳順推
	zhong := c.TianPan[chu]
	mo := c.TianPan[zhong]

	return SanChuan{
		Chu:  ChuanInfo{Branch: chu},
		Zhong: ChuanInfo{Branch: zhong},
		Mo:   ChuanInfo{Branch: mo},
		Method: "八專法",
	}, true
}

// FuYinMethod 伏吟法：天盤地盤相同
func (c *JiuZongMenCalculator) FuYinMethod() (SanChuan, bool) {
	// 檢查是否伏吟
	isFuYin := true
	for i := 0; i < 12; i++ {
		if c.TianPan[i] != c.DiPan[i] {
			isFuYin = false
			break
		}
	}

	if !isFuYin {
		return SanChuan{}, false
	}

	// 伏吟課：初傳取刑神，末傳取沖神
	// 簡化：取日干寄宮刑神
	attachment := StemAttachment[c.DayStem]
	chu := c.getXingShen(attachment)

	// 中傳：根據初傳順行
	zhong := c.TianPan[chu]
	// 末傳：根據中傳順行
	mo := c.TianPan[zhong]

	return SanChuan{
		Chu:  ChuanInfo{Branch: chu},
		Zhong: ChuanInfo{Branch: zhong},
		Mo:   ChuanInfo{Branch: mo},
		Method: "伏吟法",
	}, true
}

// FanYinMethod 返吟法：天盤地盤對沖
func (c *JiuZongMenCalculator) FanYinMethod() SanChuan {
	// 返吟：天盤地盤相沖（位置相對）
	// 簡化處理：取馬星為初傳
	chu := c.getMaXing()

	// 中末傳：根據初傳順行
	zhong := c.TianPan[chu]
	mo := c.TianPan[zhong]

	return SanChuan{
		Chu:  ChuanInfo{Branch: chu},
		Zhong: ChuanInfo{Branch: zhong},
		Mo:   ChuanInfo{Branch: mo},
		Method: "返吟法",
	}
}

// 輔助函數

// isKe 判斷是否相克（a 克 b）
func (c *JiuZongMenCalculator) isKe(a, b Branch) bool {
	// 地支五行：子丑水、寅卯木、辰巳土、午未火、申酉金、戌亥水（簡化）
	// 實際應根據五行生克表
	keMap := map[Branch][]Branch{
		Zi:  {Si, Wu, Wei}, // 子水克巳午未火
		Chou: {Si, Wu, Wei}, // 丑水克巳午未火
		Yin:  {Chen, Si, Chou}, // 寅木克辰巳丑土
		Mao:  {Chen, Si, Chou}, // 卯木克辰巳丑土
		Chen: {Zi, Chou, Hai}, // 辰土克子丑亥水
		Si:   {Shen, You},    // 巳火克申酉金
		Wu:   {Shen, You},    // 午火克申酉金
		Wei:  {Shen, You},    // 未火克申酉金
		Shen: {Yin, Mao},     // 申金克寅卯木
		You:  {Yin, Mao},     // 酉金克寅卯木
		Xu:   {Zi, Chou, Hai}, // 戌土克子丑亥水
		Hai:  {Si, Wu, Wei},  // 亥水克巳午未火
	}

	targets, ok := keMap[a]
	if !ok {
		return false
	}

	for _, target := range targets {
		if target == b {
			return true
		}
	}
	return false
}

// calculateSheHai 計算涉害深淺（地盤上受克的次數）
func (c *JiuZongMenCalculator) calculateSheHai(branch Branch) int {
	count := 0
	for i := 0; i < 12; i++ {
		if c.isKe(Branch(i), branch) {
			count++
		}
	}
	return count
}

// isSiMeng 是否四孟（寅申巳亥）
func (c *JiuZongMenCalculator) isSiMeng(b Branch) bool {
	return b == Yin || b == Shen || b == Si || b == Hai
}

// isSiZhong 是否四仲（子午卯酉）
func (c *JiuZongMenCalculator) isSiZhong(b Branch) bool {
	return b == Zi || b == Wu || b == Mao || b == You
}

// makeSanChuan 創建三傳（根據初傳推算中末傳）
func (c *JiuZongMenCalculator) makeSanChuan(chu Branch, method JiuZongMen, keIdx int) SanChuan {
	// 中傳：初傳天盤上的地支
	zhong := c.TianPan[chu]
	// 末傳：中傳天盤上的地支
	mo := c.TianPan[zhong]

	return SanChuan{
		Chu:  ChuanInfo{Branch: chu},
		Zhong: ChuanInfo{Branch: zhong},
		Mo:   ChuanInfo{Branch: mo},
		Method: JiuZongMenNames[method] + "（第" + string(rune('0'+keIdx)) + "課）",
	}
}

// getStemHe 獲取日干合（干合）
func (c *JiuZongMenCalculator) getStemHe(stem Stem) Branch {
	// 干合：甲己、乙庚、丙辛、丁壬、戊癸
	heMap := map[Stem]Branch{
		StemJia: Chou, // 甲己合土，取丑
		StemYi:  Zi,   // 乙庚合金，取子
		StemBing: Hai, // 丙辛合水，取亥
		StemDing: Xu,  // 丁壬合木，取戌
		StemWu:   You, // 戊癸合火，取酉
	}

	if branch, ok := heMap[stem]; ok {
		return branch
	}
	return Zi
}

// getSanHe 獲取日支前三合
func (c *JiuZongMenCalculator) getSanHe(branch Branch) Branch {
	// 三合：申子辰、寅午戌、巳酉丑、亥卯未
	// 前三合：取該三合局的下一個地支
	sanHeNext := map[Branch]Branch{
		Shen: Zi, // 申子辰，申的下一個是子
		Zi:   Chen, // 子辰...
		Chen: Shen,
		Yin:  Wu,
		Wu:   Xu,
		Xu:   Yin,
		Si:   You,
		You:  Chou,
		Chou: Si,
		Hai:  Mao,
		Mao:  Wei,
		Wei:  Hai,
	}

	if next, ok := sanHeNext[branch]; ok {
		return next
	}
	return branch
}

// getXingShen 獲取刑神
func (c *JiuZongMenCalculator) getXingShen(branch Branch) Branch {
	// 地支三刑：寅巳申、丑戌未、子卯、辰午酉亥自刑
	// 簡化：取沖位
	return Branch((int(branch) + 6) % 12)
}

// getMaXing 獲取馬星（驛馬）
func (c *JiuZongMenCalculator) getMaXing() Branch {
	// 驛馬：申子辰馬在寅、寅午戌馬在申、巳酉丑馬在亥、亥卯未馬在巳
	yearBranch := c.DayBranch // 簡化使用日支

	maMap := map[Branch]Branch{
		Shen: Yin,
		Zi:   Yin,
		Chen: Yin,
		Yin:  Shen,
		Wu:   Shen,
		Xu:   Shen,
		Si:   Hai,
		You:  Hai,
		Chou: Hai,
		Hai:  Si,
		Mao:  Si,
		Wei:  Si,
	}

	if ma, ok := maMap[yearBranch]; ok {
		return ma
	}
	return Yin
}
