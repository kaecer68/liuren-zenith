package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// CalendarDataSource 曆法數據源接口
type CalendarDataSource interface {
	GetCalendarData(t time.Time) (*CalendarData, error)
}

// CalendarData 曆法數據結構
type CalendarData struct {
	GregorianDate string  // 陽曆日期 YYYY-MM-DD
	JulianDay     float64 // 儒略日
	YearPillar    string  // 年柱 (如：甲辰)
	MonthPillar   string  // 月柱 (如：丙寅)
	DayPillar     string  // 日柱 (如：戊午)
	HourPillar    string  // 時柱 (如：己未)
	SolarTerm     string  // 節氣名稱
	SolarTermIdx  int     // 節氣索引 (0-23)
}

// LunarZenithClient 調用 lunar-zenith 服務的客戶端
type LunarZenithClient struct {
	BaseURL string
	Client  *http.Client
}

// NewLunarZenithClient 創建 lunar-zenith 客戶端
func NewLunarZenithClient(baseURL string) *LunarZenithClient {
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	return &LunarZenithClient{
		BaseURL: baseURL,
		Client:  &http.Client{Timeout: 10 * time.Second},
	}
}

// GetCalendarData 從 lunar-zenith 獲取曆法數據，失敗時返回錯誤
func (c *LunarZenithClient) GetCalendarData(t time.Time) (*CalendarData, error) {
	dateStr := t.Format("2006-01-02")
	url := fmt.Sprintf("%s/v1/calendar?date=%s", c.BaseURL, dateStr)

	resp, err := c.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("連接 lunar-zenith 失敗: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("lunar-zenith 返回錯誤狀態: %d", resp.StatusCode)
	}

	var result struct {
		Pillars struct {
			Year struct {
				StemIndex   int `json:"StemIndex"`
				BranchIndex int `json:"BranchIndex"`
			} `json:"Year"`
			Month struct {
				StemIndex   int `json:"StemIndex"`
				BranchIndex int `json:"BranchIndex"`
			} `json:"Month"`
			Day struct {
				StemIndex   int `json:"StemIndex"`
				BranchIndex int `json:"BranchIndex"`
			} `json:"Day"`
			Hour struct {
				StemIndex   int `json:"StemIndex"`
				BranchIndex int `json:"BranchIndex"`
			} `json:"Hour"`
		} `json:"pillars"`
		SolarTerm struct {
			Index int    `json:"Index"`
			Name  string `json:"Name"`
		} `json:"solar_term"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析 lunar-zenith 響應失敗: %w", err)
	}

	stems := []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}
	branches := []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}

	// 使用傳入的時間 t 計算時柱（lunar-zenith 只根據日期返回默認時柱）
	hour := t.Hour()
	// 時支計算：子(23-1), 丑(1-3), 寅(3-5), 卯(5-7), 辰(7-9), 巳(9-11), 午(11-13), 未(13-15), 申(15-17), 酉(17-19), 戌(19-21), 亥(21-23)
	// 公式：(hour + 1) / 2 % 12，但 23點和0點都屬於子時
	var hourBranch int
	if hour == 23 || hour == 0 {
		hourBranch = 0 // 子時
	} else {
		hourBranch = ((hour + 1) / 2) % 12
	}

	// 時干計算：五鼠遁
	// 甲己還加甲(甲日子時起甲子)，乙庚丙作初，丙辛從戊起，丁壬庚子居，戊癸何方發，壬子是真途
	dayStem := result.Pillars.Day.StemIndex
	var startStem int
	switch dayStem % 5 {
	case 0: // 甲己
		startStem = 0 // 甲
	case 1: // 乙庚
		startStem = 2 // 丙
	case 2: // 丙辛
		startStem = 4 // 戊
	case 3: // 丁壬
		startStem = 6 // 庚
	case 4: // 戊癸
		startStem = 8 // 壬
	}
	hourStem := (startStem + hourBranch) % 10

	return &CalendarData{
		GregorianDate: dateStr,
		YearPillar:    stems[result.Pillars.Year.StemIndex] + branches[result.Pillars.Year.BranchIndex],
		MonthPillar:   stems[result.Pillars.Month.StemIndex] + branches[result.Pillars.Month.BranchIndex],
		DayPillar:     stems[result.Pillars.Day.StemIndex] + branches[result.Pillars.Day.BranchIndex],
		HourPillar:    stems[hourStem] + branches[hourBranch],
		SolarTerm:     result.SolarTerm.Name,
		SolarTermIdx:  result.SolarTerm.Index,
	}, nil
}

// LocalDataSource 本地數據源實現（備用）
type LocalDataSource struct{}

// NewLocalDataSource 創建本地數據源
func NewLocalDataSource() *LocalDataSource {
	return &LocalDataSource{}
}

// GetCalendarData 獲取曆法數據（簡易實現，後續可連接 lunar-zenith）
func (l *LocalDataSource) GetCalendarData(t time.Time) (*CalendarData, error) {
	// 簡易實現：使用固定規則計算干支
	// 實際應用中應調用 lunar-zenith 服務
	year, month, day := t.Date()
	hour := t.Hour()

	stems := []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}
	branches := []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}

	// 計算年干支（以 1984 甲子年為基準）
	// 注意：年柱以立春為界，這裡簡化處理
	yearStem := (year - 1984) % 10
	if yearStem < 0 {
		yearStem += 10
	}
	yearBranch := (year - 1984) % 12
	if yearBranch < 0 {
		yearBranch += 12
	}

	// 計算日干支（簡易算法，以 2000-01-01 為基準）
	baseDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	daysDiff := int(t.Sub(baseDate).Hours() / 24)
	dayStem := (4 + daysDiff) % 10 // 2000-01-01 是 戊午 (4, 6)
	if dayStem < 0 {
		dayStem += 10
	}
	dayBranch := (6 + daysDiff) % 12
	if dayBranch < 0 {
		dayBranch += 12
	}

	// 計算月柱（根據節氣和五虎遁）
	// 月支：正月寅起，根據節氣確定
	// 簡化：使用陽曆月份估算節氣（精確計算需天文算法）
	monthBranch, solarTermIdx := calculateMonthBranchByDate(month, day)

	// 月干：五虎遁
	// 甲己之年丙作首，乙庚之歲戊為頭，丙辛之年尋庚起，丁壬壬位順行流，戊癸甲寅之上求
	monthStem := calculateMonthStem(yearStem, monthBranch)

	// 計算時干支
	hourBranch := ((hour + 1) / 2) % 12
	// 五鼠遁：甲己還加甲...
	startStem := (dayStem % 5 * 2) % 10
	hourStem := (startStem + hourBranch) % 10

	return &CalendarData{
		GregorianDate: t.Format("2006-01-02"),
		YearPillar:    stems[yearStem] + branches[yearBranch],
		MonthPillar:   stems[monthStem] + branches[monthBranch],
		DayPillar:     stems[dayStem] + branches[dayBranch],
		HourPillar:    stems[hourStem] + branches[hourBranch],
		SolarTermIdx:  solarTermIdx,
	}, nil
}

// calculateMonthBranchByDate 根據陽曆日期估算月支和節氣索引
// 節氣對應：立春(0-1月)、驚蟄(2-3月)、清明(4-5月)、立夏(6-7月)、芒種(8-9月)、小暑(10-11月)
// 立秋(12-13月)、白露(14-15月)、寒露(16-17月)、立冬(18-19月)、大雪(20-21月)、小寒(22-23月)
func calculateMonthBranchByDate(month time.Month, day int) (int, int) {
	// 簡化節氣判斷（約略值，實際節氣每年略有不同）
	// 寅月(立春2/4): 2月4日左右
	// 卯月(驚蟄3/6): 3月6日左右
	// 辰月(清明4/5): 4月5日左右
	// 巳月(立夏5/6): 5月6日左右
	// 午月(芒種6/6): 6月6日左右
	// 未月(小暑7/7): 7月7日左右
	// 申月(立秋8/8): 8月8日左右
	// 酉月(白露9/8): 9月8日左右
	// 戌月(寒露10/8): 10月8日左右
	// 亥月(立冬11/7): 11月7日左右
	// 子月(大雪12/7): 12月7日左右
	// 丑月(小寒1/6): 1月6日左右

	// 節氣分界日（簡化）
	jieQiDay := []int{6, 4, 6, 5, 6, 6, 7, 8, 8, 8, 7, 7} // 各月節氣約略日期

	m := int(month)
	isAfterJieQi := day >= jieQiDay[m-1]

	// 計算節氣索引（0-23）
	// 每個月有兩個節氣：節和中氣
	// 簡化：每月第一個節氣為 (month-1)*2，中氣為 (month-1)*2+1
	solarTermIdx := (m - 1) * 2
	if !isAfterJieQi && m > 1 {
		solarTermIdx = (m-2)*2 + 1 // 上個月的中氣
	}
	if solarTermIdx < 0 {
		solarTermIdx = 23 // 大寒
	}
	if solarTermIdx > 23 {
		solarTermIdx = 0
	}

	// 月支：正月寅(2)，順行
	// 節氣後才換月
	monthIndex := m - 1 // 0-11
	if !isAfterJieQi {
		monthIndex = m - 2
		if monthIndex < 0 {
			monthIndex = 11
		}
	}
	// 轉換為地支索引：正月寅=2
	monthBranch := (monthIndex + 2) % 12

	return monthBranch, solarTermIdx
}

// calculateMonthStem 五虎遁計算月干
// 甲己之年丙作首(丙寅)，乙庚之歲戊為頭(戊寅)，丙辛之年尋庚起(庚寅)，
// 丁壬壬位順行流(壬寅)，戊癸甲寅之上求(甲寅)
func calculateMonthStem(yearStem, monthBranch int) int {
	// 年干對應的月干起始
	// 甲(0)、己(5) → 丙(2) 起
	// 乙(1)、庚(6) → 戊(4) 起
	// 丙(2)、辛(7) → 庚(6) 起
	// 丁(3)、壬(8) → 壬(8) 起
	// 戊(4)、癸(9) → 甲(0) 起

	var startStem int
	switch yearStem % 5 {
	case 0: // 甲己
		startStem = 2 // 丙
	case 1: // 乙庚
		startStem = 4 // 戊
	case 2: // 丙辛
		startStem = 6 // 庚
	case 3: // 丁壬
		startStem = 8 // 壬
	case 4: // 戊癸
		startStem = 0 // 甲
	}

	// 月干 = 起始干 + 月支偏移（從寅開始算）
	// 寅=2, 卯=3, 辰=4...
	offset := (monthBranch - 2 + 12) % 12
	return (startStem + offset) % 10
}
