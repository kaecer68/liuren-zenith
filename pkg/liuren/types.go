// Package liuren 提供大六壬核心資料結構與推演算法
package liuren

// Branch 地支索引 (0-11: 子丑寅卯辰巳午未申酉戌亥)
type Branch int

const (
	Zi   Branch = iota // 子
	Chou               // 丑
	Yin                // 寅
	Mao                // 卯
	Chen               // 辰
	Si                 // 巳
	Wu                 // 午
	Wei                // 未
	Shen               // 申
	You                // 酉
	Xu                 // 戌
	Hai                // 亥
)

// BranchNames 地支名稱
var BranchNames = []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}

// Stem 天干索引 (0-9: 甲乙丙丁戊己庚辛壬癸)
type Stem int

const (
	StemJia  Stem = iota // 甲
	StemYi               // 乙
	StemBing             // 丙
	StemDing             // 丁
	StemWu               // 戊
	StemJi               // 己
	StemGeng             // 庚
	StemXin              // 辛
	StemRen              // 壬
	StemGui              // 癸
)

// StemNames 天干名稱
var StemNames = []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}

// MonthGeneral 月將（十二月將）
type MonthGeneral int

const (
	DengMing   MonthGeneral = iota // 登明（亥）
	HeKui                          // 河魁（戌）
	CongKui                        // 從魁（酉）
	ChuanSong                      // 傳送（申）
	XiaoJi                         // 小吉（未）
	ShengGuang                     // 勝光（午）
	TaiYi                          // 太乙（巳）
	TianGang                       // 天罡（辰）
	TaiChong                       // 太沖（卯）
	GongCao                        // 功曹（寅）
	DaJi                           // 大吉（丑）
	ShenHou                        // 神後（子）
)

// MonthGeneralNames 月將名稱
var MonthGeneralNames = []string{
	"登明", "河魁", "從魁", "傳送", "小吉", "勝光",
	"太乙", "天罡", "太沖", "功曹", "大吉", "神後",
}

// MonthGeneralBranches 月將對應地支
var MonthGeneralBranches = []Branch{Hai, Xu, You, Shen, Wei, Wu, Si, Chen, Mao, Yin, Chou, Zi}

// HeavenlyGeneral 天將（十二天將）
type HeavenlyGeneral int

const (
	GuiRen   HeavenlyGeneral = iota // 貴人
	TengShe                         // 螣蛇
	ZhuQue                          // 朱雀
	LiuHe                           // 六合
	GouChen                         // 勾陳
	QingLong                        // 青龍
	TianKong                        // 天空
	BaiHu                           // 白虎
	TaiChang                        // 太常
	XuanWu                          // 玄武
	TaiYin                          // 太陰
	TianHou                         // 天后
)

// HeavenlyGeneralNames 天將名稱
var HeavenlyGeneralNames = []string{
	"貴人", "螣蛇", "朱雀", "六合", "勾陳", "青龍",
	"天空", "白虎", "太常", "玄武", "太陰", "天后",
}

// HeavenlyGeneralElements 天將五行屬性
var HeavenlyGeneralElements = []string{
	"土", "火", "火", "木", "土", "木",
	"土", "金", "土", "水", "金", "水",
}

// Shengxiao 生肖
var Shengxiao = []string{"鼠", "牛", "虎", "兔", "龍", "蛇", "馬", "羊", "猴", "雞", "狗", "豬"}

// Sexagenary 干支組合
type Sexagenary struct {
	Stem   Stem   // 天干
	Branch Branch // 地支
}

// String 返回干支字串
func (s Sexagenary) String() string {
	return StemNames[s.Stem] + BranchNames[s.Branch]
}

// Animal 返回對應生肖
func (s Sexagenary) Animal() string {
	return Shengxiao[s.Branch]
}

// StemBranch 日干支結構
type StemBranch struct {
	Day   Sexagenary // 日干（主體）
	Hour  Sexagenary // 時干支
	Month Sexagenary // 月干支
	Year  Sexagenary // 年干支
}

// Pan 六壬課盤結構
type Pan struct {
	YearPillar   Sexagenary          // 年柱
	MonthPillar  Sexagenary          // 月柱
	DayPillar    Sexagenary          // 日柱（日干主體）
	HourPillar   Sexagenary          // 時柱
	MonthGeneral MonthGeneral        // 月將
	HourBranch   Branch              // 占時（時支）
	DiPan        [12]Branch          // 地盤（固定）
	TianPan      [12]Branch          // 天盤（月將加時）
	TianJiang    [12]HeavenlyGeneral // 天將佈局
	FourKe       FourKe              // 四課
	SanChuan     SanChuan            // 三傳
	DayGeneral   HeavenlyGeneral     // 晝貴或夜貴
	IsDay        bool                // 是否晝占
	VoidBranches []Branch            // 空亡地支
	KeTi         string              // 課體名稱
	XiongJi      string              // 吉凶判斷
	DuanYu       []string            // 詳細斷語
}

// Ke 單課結構
type Ke struct {
	Down Branch // 下神（地盤/日支或日干）
	Up   Branch // 上神（天盤）
}

// FourKe 四課結構
type FourKe struct {
	Ke1 Ke // 第一課：日干之上
	Ke2 Ke // 第二課：第一課上神之上
	Ke3 Ke // 第三課：日支之上
	Ke4 Ke // 第四課：第三課上神之上
}

// SanChuan 三傳結構
type SanChuan struct {
	Chu    ChuanInfo // 初傳（發用/發端）
	Zhong  ChuanInfo // 中傳（移易）
	Mo     ChuanInfo // 末傳（歸計）
	Method string    // 九宗門取法
}

// ChuanInfo 傳之詳細資訊
type ChuanInfo struct {
	Branch   Branch          // 傳之地支
	General  HeavenlyGeneral // 天將
	Stem     Stem            // 遁干
	Relation string          // 六親關係（生我、我生、克我、我克、比和）
	IsEmpty  bool            // 是否落空
	IsHorse  bool            // 是否驛馬
	IsTaoHua bool            // 是否桃花
}

// JiuZongMen 九宗門類型
type JiuZongMen int

const (
	ZeiKe   JiuZongMen = iota // 賊克法
	BiYong                    // 比用法（知一法）
	SheHai                    // 涉害法
	YaoKe                     // 遙克法
	MaoXing                   // 昴星法
	BieZe                     // 別責法
	BaZhuan                   // 八專法
	FuYin                     // 伏吟法
	FanYin                    // 返吟法
)

// JiuZongMenNames 九宗門名稱
var JiuZongMenNames = []string{
	"賊克法", "比用法", "涉害法", "遙克法", "昴星法",
	"別責法", "八專法", "伏吟法", "返吟法",
}

// StemAttachment 日干寄宮表（甲寄寅、乙寄辰...）
var StemAttachment = map[Stem]Branch{
	StemJia:  Yin,  // 甲寄寅
	StemYi:   Chen, // 乙寄辰
	StemBing: Si,   // 丙寄巳
	StemDing: Wei,  // 丁寄未
	StemWu:   Si,   // 戊寄巳
	StemJi:   Wei,  // 己寄未
	StemGeng: Shen, // 庚寄申
	StemXin:  Xu,   // 辛寄戌
	StemRen:  Hai,  // 壬寄亥
	StemGui:  Chou, // 癸寄丑
}

// XunKong 六甲旬空
var XunKong = map[Stem][]Branch{
	StemJia: {Xu, Hai}, // 甲子旬：戌亥空
	// 其他旬需要根據日干支計算
}
