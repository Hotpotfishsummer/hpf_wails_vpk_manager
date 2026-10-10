package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"hpf_wails_vpk_manager/internal/models"
)

type RestarterService struct {
	client *http.Client
}

func NewRestarterService() *RestarterService {
	return &RestarterService{
		client: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}
}

type restartResponse struct {
	Ok        bool    `json:"ok"`
	DurationS float64 `json:"duration_s"`
	Stderr    string  `json:"stderr"`
	Stdout    string  `json:"stdout"`
}

func (s *RestarterService) RestartServer() models.RestartResult {
	url := "http://100.69.225.112:9370/restart"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return models.RestartResult{
			ApiResult: models.ApiResult{Success: false, Error: err.Error()},
		}
	}
	req.Header.Set("X-Auth-Token", "9f3c81a4d2b6e750c4a8f1d3")

	resp, err := s.client.Do(req)
	if err != nil {
		return models.RestartResult{
			ApiResult: models.ApiResult{Success: false, Error: fmt.Sprintf("调用重启服务失败: %v", err)},
		}
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == 401 {
		return models.RestartResult{
			ApiResult: models.ApiResult{Success: false, Error: "重启服务鉴权失败 (Token 无效)"},
		}
	}

	if resp.StatusCode != 200 {
		return models.RestartResult{
			ApiResult: models.ApiResult{Success: false, Error: fmt.Sprintf("服务器返回 HTTP %d: %s", resp.StatusCode, string(bodyBytes))},
		}
	}

	var data restartResponse
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return models.RestartResult{
			ApiResult: models.ApiResult{Success: false, Error: fmt.Sprintf("解析返回数据失败: %v", err)},
		}
	}

	if data.Ok {
		return models.RestartResult{
			ApiResult: models.ApiResult{Success: true},
			DurationS: data.DurationS,
		}
	}

	errMsg := data.Stderr
	if errMsg == "" {
		errMsg = data.Stdout
	}
	if errMsg == "" {
		errMsg = "未知错误"
	}
	return models.RestartResult{
		ApiResult: models.ApiResult{Success: false, Error: errMsg},
	}
}
