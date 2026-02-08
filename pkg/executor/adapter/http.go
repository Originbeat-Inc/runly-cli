package adapter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/originbeat-inc/runly-cli/internal/config"
	"github.com/originbeat-inc/runly-cli/internal/i18n"
	"github.com/schollz/progressbar/v3"
)

type RunlyRequest struct {
	Header  map[string]string      `json:"header"`
	Payload map[string]interface{} `json:"payload"`
}

type RunlyResponse struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type RequestClient struct {
	BaseURL string
	Token   string
	Timeout time.Duration
	// 存储 Profile 引用以便内部切换地址
	activeProfile config.Profile
}

// NewClient 初始化，默认使用 HubServer 满足大部分资产操作需求
func NewClient() *RequestClient {
	cfg, _ := config.LoadConfig()
	active := cfg.GetActive()

	return &RequestClient{
		BaseURL:       active.HubServer, // 修正：默认设为 HubServer，解决 init/search/publish 报错
		Token:         active.AccessToken,
		Timeout:       60 * time.Second,
		activeProfile: active,
	}
}

// SetToMeServer 切换到身份服务器地址 (用于 keys 指令)
func (c *RequestClient) SetToMeServer() *RequestClient {
	c.BaseURL = c.activeProfile.MeServer
	return c
}

// SetToHubServer 切换到资产服务器地址 (用于 init, search, pull, publish 指令)
func (c *RequestClient) SetToHubServer() *RequestClient {
	c.BaseURL = c.activeProfile.HubServer
	return c
}

// Post 通用请求（用于 Search, Keys 等不带进度条的操作）
func (c *RequestClient) Post(path string, payload map[string]interface{}) (map[string]interface{}, error) {
	fullURL := fmt.Sprintf("%s%s", c.BaseURL, path)
	reqData := RunlyRequest{
		Header:  map[string]string{"Authorization": "Bearer " + c.Token},
		Payload: payload,
	}
	jsonBytes, _ := json.Marshal(reqData)

	req, _ := http.NewRequest("POST", fullURL, bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: c.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("errors.network_err"), err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	var runlyResp RunlyResponse
	if err := json.Unmarshal(bodyBytes, &runlyResp); err != nil {
		return nil, fmt.Errorf(i18n.T("errors.yaml_unmarshal_fail"), err)
	}

	if runlyResp.Status != "success" {
		return nil, fmt.Errorf(runlyResp.Message)
	}
	return runlyResp.Data, nil
}

// PostWithProgress 带进度的上传/下载操作
func (c *RequestClient) PostWithProgress(path string, payload map[string]interface{}, progressKey string) (map[string]interface{}, error) {
	fullURL := fmt.Sprintf("%s%s", c.BaseURL, path)
	progressText := i18n.T(progressKey)

	reqData := RunlyRequest{
		Header: map[string]string{
			"Authorization": "Bearer " + c.Token,
			"Content-Type":  "application/json",
		},
		Payload: payload,
	}

	jsonBytes, _ := json.Marshal(reqData)
	totalSize := int64(len(jsonBytes))

	bar := progressbar.NewOptions64(
		totalSize,
		progressbar.OptionSetDescription(progressText),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]█[reset]",
			SaucerHead:    "[green]█[reset]",
			SaucerPadding: " ",
			BarStart:      "|",
			BarEnd:        "|",
		}),
	)

	bodyReader := io.TeeReader(bytes.NewReader(jsonBytes), bar)
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "POST", fullURL, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.ContentLength = totalSize

	client := &http.Client{}
	fmt.Print("\n")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("errors.network_err"), err)
	}
	defer resp.Body.Close()

	bar.Finish()
	fmt.Print("\n")

	bodyBytes, _ := io.ReadAll(resp.Body)
	var runlyResp RunlyResponse
	if err := json.Unmarshal(bodyBytes, &runlyResp); err != nil {
		return nil, fmt.Errorf(i18n.T("errors.yaml_unmarshal_fail"), err)
	}

	if runlyResp.Status != "success" {
		return nil, fmt.Errorf("%s: %s", i18n.T("errors.server_err"), runlyResp.Message)
	}

	return runlyResp.Data, nil
}
