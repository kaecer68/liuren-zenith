package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/kaecer68/liuren-zenith/pkg/liuren"
)

type Handler struct {
	Engine *liuren.Engine
}

type LegacyDivinationRequest struct {
	Date         string `json:"date"`
	Time         string `json:"time"`
	Question     string `json:"question"`
	QuestionType string `json:"question_type"`
}

type LegacyDivinationResponse struct {
	Success bool                     `json:"success"`
	Message string                   `json:"message,omitempty"`
	Data    *liuren.DivinationResult `json:"data,omitempty"`
}

type CalculateRequest struct {
	Datetime     string `json:"datetime"`
	QuestionType string `json:"question_type"`
}

type AnalyzeRequest struct {
	Chart    CalculateResponse `json:"chart"`
	Question string            `json:"question"`
}

type CalculateResponse struct {
	KeTi     string              `json:"ke_ti"`
	SiKe     SiKeResponse        `json:"si_ke"`
	SanChuan SanChuanResponse    `json:"san_chuan"`
	TianPan  []HeavenlyPlateItem `json:"tian_pan,omitempty"`
	DiPan    []EarthlyPlateItem  `json:"di_pan,omitempty"`
	ShenSha  []ShenShaItem       `json:"shen_sha,omitempty"`
}

type AnalyzeResponse struct {
	Overall  string         `json:"overall"`
	Analysis map[string]int `json:"analysis"`
	Advice   string         `json:"advice"`
}

type SiKeResponse struct {
	FirstKe  string `json:"first_ke"`
	SecondKe string `json:"second_ke"`
	ThirdKe  string `json:"third_ke"`
	FourthKe string `json:"fourth_ke"`
}

type SanChuanResponse struct {
	InitialChuan string `json:"initial_chuan"`
	MiddleChuan  string `json:"middle_chuan"`
	FinalChuan   string `json:"final_chuan"`
}

type HeavenlyPlateItem struct {
	Position        int    `json:"position"`
	HeavenlyGeneral string `json:"heavenly_general"`
}

type EarthlyPlateItem struct {
	Position      int    `json:"position"`
	EarthlyBranch string `json:"earthly_branch"`
}

type ShenShaItem struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Position string `json:"position"`
}

func (h *Handler) HandleDivination(w http.ResponseWriter, r *http.Request) {
	writeJSONHeader(w)
	if r.Method != http.MethodPost {
		writeLegacyError(w, http.StatusMethodNotAllowed, "只支援 POST 方法")
		return
	}

	var req LegacyDivinationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeLegacyError(w, http.StatusBadRequest, "請求格式錯誤: "+err.Error())
		return
	}

	divinationTime, err := parseLegacyTime(req.Date, req.Time)
	if err != nil {
		writeLegacyError(w, http.StatusBadRequest, "時間格式錯誤: "+err.Error())
		return
	}

	result, err := h.Engine.Calculate(liuren.DivinationRequest{
		Time:         divinationTime,
		Question:     req.Question,
		QuestionType: req.QuestionType,
	})
	if err != nil {
		writeLegacyError(w, http.StatusInternalServerError, "排盤失敗: "+err.Error())
		return
	}

	_ = json.NewEncoder(w).Encode(LegacyDivinationResponse{
		Success: true,
		Data:    result,
	})
}

func (h *Handler) HandleCalculate(w http.ResponseWriter, r *http.Request) {
	writeJSONHeader(w)
	if r.Method != http.MethodPost {
		http.Error(w, "只支援 POST 方法", http.StatusMethodNotAllowed)
		return
	}

	var req CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "請求格式錯誤: "+err.Error(), http.StatusBadRequest)
		return
	}

	dt, err := time.Parse(time.RFC3339, req.Datetime)
	if err != nil {
		http.Error(w, "datetime 必須為 RFC3339 格式", http.StatusBadRequest)
		return
	}

	result, err := h.Engine.Calculate(liuren.DivinationRequest{
		Time:         dt,
		QuestionType: req.QuestionType,
	})
	if err != nil {
		http.Error(w, "排盤失敗: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(buildCalculateResponse(result))
}

func (h *Handler) HandleAnalyze(w http.ResponseWriter, r *http.Request) {
	writeJSONHeader(w)
	if r.Method != http.MethodPost {
		http.Error(w, "只支援 POST 方法", http.StatusMethodNotAllowed)
		return
	}

	var req AnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "請求格式錯誤: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.Question == "" {
		http.Error(w, "question 不能為空", http.StatusBadRequest)
		return
	}

	initial := scoreChuan(req.Chart.SanChuan.InitialChuan)
	middle := scoreChuan(req.Chart.SanChuan.MiddleChuan)
	final := scoreChuan(req.Chart.SanChuan.FinalChuan)

	resp := AnalyzeResponse{
		Overall: overallFromScore((initial+middle+final)/3, req.Chart.KeTi),
		Analysis: map[string]int{
			"initial_strength": initial,
			"middle_strength":  middle,
			"final_strength":   final,
		},
		Advice: adviceFromChart(req.Chart.KeTi, req.Question),
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) HandleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSONHeader(w)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status":    "ok",
		"service":   "liuren-zenith",
		"version":   "v1.0.0",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func buildCalculateResponse(result *liuren.DivinationResult) CalculateResponse {
	resp := CalculateResponse{
		KeTi: result.KeTi,
		SiKe: SiKeResponse{},
		SanChuan: SanChuanResponse{
			InitialChuan: "初傳：" + result.SanChuan.Chu.Branch,
			MiddleChuan:  "中傳：" + result.SanChuan.Zhong.Branch,
			FinalChuan:   "末傳：" + result.SanChuan.Mo.Branch,
		},
		TianPan: make([]HeavenlyPlateItem, 0, len(result.TianJiang)),
		DiPan:   make([]EarthlyPlateItem, 0, len(result.DiPan)),
		ShenSha: buildShenSha(result),
	}

	if len(result.FourKe) > 0 {
		resp.SiKe.FirstKe = formatKe(result.FourKe, 0)
		resp.SiKe.SecondKe = formatKe(result.FourKe, 1)
		resp.SiKe.ThirdKe = formatKe(result.FourKe, 2)
		resp.SiKe.FourthKe = formatKe(result.FourKe, 3)
	}

	for idx, general := range result.TianJiang {
		resp.TianPan = append(resp.TianPan, HeavenlyPlateItem{
			Position:        idx,
			HeavenlyGeneral: general,
		})
	}
	for idx, branch := range result.DiPan {
		resp.DiPan = append(resp.DiPan, EarthlyPlateItem{
			Position:      idx,
			EarthlyBranch: branch,
		})
	}

	return resp
}

func buildShenSha(result *liuren.DivinationResult) []ShenShaItem {
	items := make([]ShenShaItem, 0, 4)
	for _, branch := range result.VoidBranches {
		items = append(items, ShenShaItem{
			Name:     "空亡",
			Type:     "中性",
			Position: branch,
		})
	}
	if result.YiMa != "" {
		items = append(items, ShenShaItem{Name: "驛馬", Type: "吉神", Position: result.YiMa})
	}
	if result.TaoHua != "" {
		items = append(items, ShenShaItem{Name: "桃花", Type: "中性", Position: result.TaoHua})
	}
	if result.TianMa != "" {
		items = append(items, ShenShaItem{Name: "天馬", Type: "吉神", Position: result.TianMa})
	}
	return items
}

func formatKe(fourKe []liuren.KeInfo, index int) string {
	if index >= len(fourKe) {
		return ""
	}
	return fourKe[index].Down + "上神：" + fourKe[index].Up
}

func parseLegacyTime(dateStr, timeStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, errInvalidDate
	}
	if timeStr == "" {
		timeStr = "12:00"
	}
	return time.Parse("2006-01-02 15:04", dateStr+" "+timeStr)
}

func writeLegacyError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(LegacyDivinationResponse{
		Success: false,
		Message: message,
	})
}

func writeJSONHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func scoreChuan(value string) int {
	switch {
	case value == "":
		return 40
	case containsAny(value, "貴", "龍", "合"):
		return 80
	case containsAny(value, "空", "虎"):
		return 35
	default:
		return 60
	}
}

func overallFromScore(score int, keTi string) string {
	switch {
	case score >= 75:
		return keTi + "偏吉，可順勢推進"
	case score <= 45:
		return keTi + "偏守，宜先避險再行動"
	default:
		return keTi + "中平，宜審勢而動"
	}
}

func adviceFromChart(keTi, question string) string {
	switch keTi {
	case "伏吟課":
		return "課勢反覆，針對「" + question + "」建議保守觀察，避免躁進。"
	case "返吟課":
		return "變動訊號較強，針對「" + question + "」宜預留調整空間。"
	default:
		return "可結合具體占事「" + question + "」與三傳強弱做分段決策。"
	}
}

func containsAny(value string, needles ...string) bool {
	for _, needle := range needles {
		if needle != "" && strings.Contains(value, needle) {
			return true
		}
	}
	return false
}

var errInvalidDate = errors.New("日期不能為空")
