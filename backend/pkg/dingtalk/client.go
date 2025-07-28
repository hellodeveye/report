package dingtalk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hellodeveye/report/internal/models"
	"resty.dev/v3"
)

// Client 钉钉API客户端
type Client struct {
	config     *models.DingTalkConfig
	httpClient *resty.Client
}

// NewClient 创建新的钉钉客户端
func NewClient(config *models.DingTalkConfig) *Client {
	return &Client{
		config:     config,
		httpClient: resty.New().SetTimeout(30 * time.Second).EnableDebug(),
	}
}

func (c *Client) GetAccessToken() (*models.DingTalkAccessTokenResponse, error) {
	url := "https://oapi.dingtalk.com/gettoken?appkey=" + c.config.AppKey + "&appsecret=" + c.config.AppSecret
	resp, err := c.httpClient.R().Get(url)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %v", err)
	}

	var response models.DingTalkAccessTokenResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %v", err)
	}

	return &response, nil
}

// GetUserAccessToken 通过授权码获取用户访问令牌
func (c *Client) GetUserAccessToken(code string) (*models.DingTalkOAuthTokenResponse, error) {
	url := "https://api.dingtalk.com/v1.0/oauth2/userAccessToken"

	requestBody := map[string]string{
		"clientId":     c.config.AppKey,
		"clientSecret": c.config.AppSecret,
		"code":         code,
		"grantType":    "authorization_code",
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request body failed: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.R().SetBody(jsonData).Post(url)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode(), string(body))
	}

	var tokenResp models.DingTalkOAuthTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %v", err)
	}

	return &tokenResp, nil
}

// GetUserInfo 通过用户访问令牌获取用户信息
func (c *Client) GetUserInfo(accessToken string) (*models.DingTalkUserInfoResponse, error) {
	url := "https://api.dingtalk.com/v1.0/contact/users/me"

	resp, err := c.httpClient.R().SetHeader("x-acs-dingtalk-access-token", accessToken).Get(url)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode(), string(body))
	}

	var userResp models.DingTalkUserInfoResponse
	if err := json.Unmarshal(body, &userResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %v", err)
	}

	return &userResp, nil
}

func (c *Client) GetUserByUnionId(accessToken string, unionId string) (*models.DingTalkUserByUnionIdResponse, error) {
	url := "https://oapi.dingtalk.com/topapi/user/getbyunionid?access_token=" + accessToken

	requestBody := map[string]string{
		"unionid": unionId,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request body failed: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.R().SetBody(jsonData).Post(url)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode(), string(body))
	}

	var userResp models.DingTalkUserByUnionIdResponse
	if err := json.Unmarshal(body, &userResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %v", err)
	}

	return &userResp, nil
}
