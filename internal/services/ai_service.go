package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"hpf_wails_vpk_manager/internal/models"
)

type AiService struct {
	client *http.Client
}

func NewAiService() *AiService {
	return &AiService{
		client: &http.Client{
			Timeout: 20 * time.Second,
		},
	}
}

// formatEndpoint standardizes the OpenAI chat completions endpoint URL
func formatEndpoint(baseURL string) string {
	u := strings.TrimSpace(baseURL)
	if u == "" {
		u = "https://api.openai.com/v1"
	}
	u = strings.TrimRight(u, "/")
	if strings.HasSuffix(u, "/chat/completions") {
		return u
	}
	if strings.HasSuffix(u, "/v1") {
		return u + "/chat/completions"
	}
	return u + "/v1/chat/completions"
}

// ChatCompletion calls the OpenAI-compatible chat completions endpoint
func (s *AiService) ChatCompletion(ctx context.Context, req models.AiChatRequest, settings models.AiSettings) (string, error) {
	if settings.ApiKey == "" {
		return "", fmt.Errorf("API Key 未配置")
	}

	endpoint := formatEndpoint(settings.BaseURL)
	modelName := settings.Model
	if modelName == "" {
		modelName = "gpt-4o-mini"
	}
	req.Model = modelName

	if req.Temperature == 0 && settings.Temperature > 0 {
		req.Temperature = settings.Temperature
	}
	if req.MaxTokens == 0 && settings.MaxTokens > 0 {
		req.MaxTokens = settings.MaxTokens
	}

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("创建 HTTP 请求失败: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+settings.ApiKey)

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("请求大模型接口失败: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	var chatResp struct {
		Choices []struct {
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error *struct {
			Message string `json:"message"`
			Type    string `json:"type"`
			Code    any    `json:"code"`
		} `json:"error,omitempty"`
	}

	if err := json.Unmarshal(respBytes, &chatResp); err != nil {
		return "", fmt.Errorf("解析响应失败 (HTTP %d): %s", resp.StatusCode, string(respBytes))
	}

	if chatResp.Error != nil && chatResp.Error.Message != "" {
		return "", fmt.Errorf("API 错误: %s", chatResp.Error.Message)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBytes))
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("模型未返回任何有效回复")
	}

	content := strings.TrimSpace(chatResp.Choices[0].Message.Content)
	return content, nil
}

// TestConnection tests whether the configured OpenAI API is reachable and valid
func (s *AiService) TestConnection(settings models.AiSettings) models.ApiResult {
	if settings.ApiKey == "" {
		return models.ApiResult{Success: false, Error: "API Key 不能为空"}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	testReq := models.AiChatRequest{
		Model: settings.Model,
		Messages: []models.AiChatMessage{
			{Role: "user", Content: "Say OK"},
		},
		MaxTokens: 10,
	}

	reply, err := s.ChatCompletion(ctx, testReq, settings)
	if err != nil {
		return models.ApiResult{Success: false, Error: err.Error()}
	}

	return models.ApiResult{
		Success: true,
		Message: fmt.Sprintf("连接成功！模型 [%s] 响应: %s", settings.Model, strings.TrimSpace(reply)),
	}
}

// PredictGroupName predicts a concise group name using AI or falls back to smart heuristics
func (s *AiService) PredictGroupName(files []models.VpkFile, settings models.AiSettings) models.GroupPredictResult {
	if len(files) == 0 {
		return models.GroupPredictResult{
			ApiResult:     models.ApiResult{Success: true},
			PredictedName: "新模组分组",
			Source:        "heuristic",
		}
	}

	// 1. 如果配置了 API Key，优先尝试大模型预测
	if strings.TrimSpace(settings.ApiKey) != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()

		var sb strings.Builder
		for i, f := range files {
			if i >= 15 {
				sb.WriteString(fmt.Sprintf("- 还有更多 %d 个文件...\n", len(files)-15))
				break
			}
			title := f.DisplayName
			if f.AddonTitle != nil && *f.AddonTitle != "" {
				title = *f.AddonTitle
			}
			mapStr := ""
			if len(f.Maps) > 0 {
				mapStr = fmt.Sprintf(" (包含 %d 个战役地图)", len(f.Maps))
			}
			noteStr := ""
			if f.Note != "" {
				noteStr = fmt.Sprintf(" [备注:%s]", f.Note)
			}
			sb.WriteString(fmt.Sprintf("- 文件名: %s, 标题: %s%s%s\n", f.Name, title, noteStr, mapStr))
		}

		prompt := fmt.Sprintf(`用户正在求生之路2(Left 4 Dead 2)模组管理器中将以下 %d 个 VPK 模组组合为一个分组：
%s
请分析这组模组的共同特征（例如：生还者角色皮肤套装、武器枪械合集、战役地图包、音效包、特感重置等），为其生成一个简明、贴切、规范的中文分组名称。
规则要求：
1. 长度严格控制在 2 到 10 个汉字；
2. 仅输出分组名称本身；
3. 不要包含任何引号、标点、前缀序号或解释说明。`, len(files), sb.String())

		req := models.AiChatRequest{
			Model: settings.Model,
			Messages: []models.AiChatMessage{
				{
					Role:    "system",
					Content: "你是一个专精求生之路2 (Left 4 Dead 2) 模组管理的智能助理。根据用户给出的模组信息，给出精炼准确的中文模组分组名称。",
				},
				{Role: "user", Content: prompt},
			},
			Temperature: 0.3,
			MaxTokens:   30,
		}

		if aiName, err := s.ChatCompletion(ctx, req, settings); err == nil {
			clean := cleanPredictedName(aiName)
			if clean != "" && utf8.RuneCountInString(clean) <= 15 {
				return models.GroupPredictResult{
					ApiResult:     models.ApiResult{Success: true},
					PredictedName: clean,
					Reason:        "由大模型智能分析生成",
					Source:        "ai",
				}
			}
		}
	}

	// 2. 本地智能启发式规则分析保底 (即便没配置 API Key 或断网也能智能推测)
	heuristicName := PredictGroupNameHeuristic(files)
	return models.GroupPredictResult{
		ApiResult:     models.ApiResult{Success: true},
		PredictedName: heuristicName,
		Reason:        "本地规则启发式分析生成",
		Source:        "heuristic",
	}
}

// cleanPredictedName strips surrounding quotes, markdown, and formatting from model output
func cleanPredictedName(s string) string {
	s = strings.TrimSpace(s)
	lines := strings.Split(s, "\n")
	if len(lines) > 0 {
		s = strings.TrimSpace(lines[0])
	}
	s = strings.Trim(s, "`'\"“”‘’《》【】[]()（）*#-_:： ")
	s = strings.TrimPrefix(s, "分组名称：")
	s = strings.TrimPrefix(s, "分组名称:")
	s = strings.TrimPrefix(s, "分组：")
	s = strings.TrimPrefix(s, "名称：")
	s = strings.Trim(s, "`'\"“”‘’《》【】[]()（）*#-_:： ")
	return strings.TrimSpace(s)
}

func containsAnyKeyword(text string, keywords []string) bool {
	lowerText := strings.ToLower(text)
	tokens := strings.FieldsFunc(lowerText, func(r rune) bool {
		return r == ' ' || r == '_' || r == '-' || r == '.' || r == ',' || r == '/' || r == '\\' || r == '[' || r == ']' || r == '(' || r == ')'
	})
	tokenSet := make(map[string]bool, len(tokens))
	for _, t := range tokens {
		tokenSet[t] = true
	}

	for _, kw := range keywords {
		kwLower := strings.ToLower(kw)
		isShortAscii := true
		for _, r := range kwLower {
			if r > 127 {
				isShortAscii = false
				break
			}
		}
		if isShortAscii && len(kwLower) <= 4 {
			if tokenSet[kwLower] {
				return true
			}
		} else {
			if strings.Contains(lowerText, kwLower) {
				return true
			}
		}
	}
	return false
}

// PredictGroupNameHeuristic analyzes file names, addon titles and notes using domain heuristics
func PredictGroupNameHeuristic(files []models.VpkFile) string {
	if len(files) == 0 {
		return "新模组分组"
	}

	// 1. 如果有任何文件识别出了战役地图，直接归为战役地图包
	for _, f := range files {
		if len(f.Maps) > 0 {
			return "战役地图模组包"
		}
	}

	// 收集所有文本源
	allTexts := make([]string, 0, len(files)*2)
	for _, f := range files {
		stem := strings.TrimSuffix(f.Name, filepath.Ext(f.Name))
		stem = strings.TrimSuffix(stem, ".disabled")
		allTexts = append(allTexts, stem)
		if f.AddonTitle != nil && *f.AddonTitle != "" {
			allTexts = append(allTexts, *f.AddonTitle)
		}
		if f.DisplayName != "" {
			allTexts = append(allTexts, f.DisplayName)
		}
		if f.Note != "" {
			allTexts = append(allTexts, f.Note)
		}
	}
	joined := strings.Join(allTexts, " ")

	// 角色匹配 (有序列表确保匹配确定性)
	type rule struct {
		key  string
		name string
	}
	survivorRules := []rule{
		{"zoey", "佐伊皮肤套装"},
		{"佐伊", "佐伊皮肤套装"},
		{"nick", "尼克皮肤套装"},
		{"尼克", "尼克皮肤套装"},
		{"ellis", "艾利斯皮肤套装"},
		{"艾利斯", "艾利斯皮肤套装"},
		{"rochelle", "罗雪儿皮肤套装"},
		{"罗雪儿", "罗雪儿皮肤套装"},
		{"coach", "教练皮肤套装"},
		{"教练", "教练皮肤套装"},
		{"bill", "比尔皮肤套装"},
		{"比尔", "比尔皮肤套装"},
		{"francis", "弗朗西斯皮肤套装"},
		{"弗朗西斯", "弗朗西斯皮肤套装"},
		{"louis", "路易斯皮肤套装"},
		{"路易斯", "路易斯皮肤套装"},
	}
	for _, r := range survivorRules {
		if containsAnyKeyword(joined, []string{r.key}) {
			return r.name
		}
	}

	// 特感匹配 (有序列表)
	infectedRules := []rule{
		{"tank", "Tank坦克特感模组"},
		{"坦克", "Tank坦克特感模组"},
		{"witch", "Witch女巫特感模组"},
		{"女巫", "Witch女巫特感模组"},
		{"hunter", "Hunter特感模组"},
		{"smoker", "Smoker特感模组"},
		{"boomer", "Boomer特感模组"},
		{"charger", "Charger特感模组"},
		{"spitter", "Spitter特感模组"},
		{"jockey", "Jockey特感模组"},
	}
	for _, r := range infectedRules {
		if containsAnyKeyword(joined, []string{r.key}) {
			return r.name
		}
	}

	// 武器匹配
	weaponKeywords := []string{"weapon", "weapons", "gun", "guns", "ak47", "m16", "deagle", "shotgun", "sniper", "rifle", "pistol", "melee", "武器", "枪械", "近战", "手枪", "步枪", "散弹", "连喷"}
	if containsAnyKeyword(joined, weaponKeywords) {
		return "武器枪械模组包"
	}

	// 战役地图匹配
	mapKeywords := []string{"campaign", "map", "maps", "mission", "战役", "地图", "章节"}
	if containsAnyKeyword(joined, mapKeywords) {
		return "战役地图模组包"
	}

	// 音效匹配
	soundKeywords := []string{"sound", "voice", "music", "bgm", "audio", "音效", "语音", "音乐", "配乐"}
	if containsAnyKeyword(joined, soundKeywords) {
		return "语音音效模组包"
	}

	// 界面 HUD 匹配
	uiKeywords := []string{"hud", "ui", "crosshair", "准星", "界面", "菜单", "血条"}
	if containsAnyKeyword(joined, uiKeywords) {
		return "HUD与界面模组"
	}

	// 若只有一个文件，采用其主标题或去字符后的名称
	if len(files) == 1 {
		title := files[0].DisplayName
		if files[0].AddonTitle != nil && *files[0].AddonTitle != "" {
			title = *files[0].AddonTitle
		}
		title = cleanTitle(title)
		if title != "" && utf8.RuneCountInString(title) <= 12 {
			return title
		}
	}

	// 寻找公共词/前缀 (例如 miku_zoey 与 miku_voice -> miku)
	if len(files) >= 2 {
		common := findCommonPrefixOrKeyword(files)
		if common != "" {
			return common + "模组包"
		}
	}

	// 兜底：使用第一个文件的名称简写
	firstTitle := files[0].DisplayName
	if files[0].AddonTitle != nil && *files[0].AddonTitle != "" {
		firstTitle = *files[0].AddonTitle
	}
	firstClean := cleanTitle(firstTitle)
	if utf8.RuneCountInString(firstClean) > 6 {
		runes := []rune(firstClean)
		firstClean = string(runes[:6])
	}
	return firstClean + "等模组"
}

var nonAlphanumericChinese = regexp.MustCompile(`[^\p{Han}a-zA-Z0-9]+`)

func cleanTitle(title string) string {
	title = strings.TrimSpace(title)
	title = regexp.MustCompile(`\[.*?\]`).ReplaceAllString(title, "")
	title = regexp.MustCompile(`\(.*?\)`).ReplaceAllString(title, "")
	title = strings.TrimSuffix(title, ".vpk")
	title = strings.TrimSuffix(title, ".disabled")
	title = strings.TrimSpace(title)
	return title
}

func findCommonPrefixOrKeyword(files []models.VpkFile) string {
	stems := make([]string, len(files))
	for i, f := range files {
		s := strings.TrimSuffix(f.Name, filepath.Ext(f.Name))
		s = strings.TrimSuffix(s, ".disabled")
		stems[i] = strings.ToLower(s)
	}

	// 尝试以下划线或连字符分词后的第一个词
	firstTokens := strings.Split(stems[0], "_")
	if len(firstTokens) == 1 {
		firstTokens = strings.Split(stems[0], "-")
	}
	if len(firstTokens) > 1 && len(firstTokens[0]) >= 3 {
		prefix := firstTokens[0]
		allMatch := true
		for _, stem := range stems[1:] {
			if !strings.HasPrefix(stem, prefix) {
				allMatch = false
				break
			}
		}
		if allMatch {
			return prefix
		}
	}

	return ""
}
