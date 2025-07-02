package feishu

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hellodeveye/report/internal/models"
)

// AuthService 飞书认证服务
type AuthService struct {
	config *models.FeishuConfig
}

// NewAuthService 创建新的飞书认证服务
func NewAuthService(config *models.FeishuConfig) *AuthService {
	return &AuthService{
		config: config,
	}
}

// GenerateAuthURL 生成授权URL
func (s *AuthService) GenerateAuthURL() (string, string, error) {
	// 生成state参数防CSRF
	state := fmt.Sprintf("%d", time.Now().UnixNano())

	// 构建授权URL (使用官方推荐的v1 authorize endpoint)
	authURL := fmt.Sprintf("https://accounts.feishu.cn/open-apis/authen/v1/authorize?client_id=%s&redirect_uri=%s&response_type=code&state=%s",
		s.config.AppID,
		s.config.RedirectURI,
		state)

	return authURL, state, nil
}

// ExchangeCodeForUser 用授权码换取用户信息
func (s *AuthService) ExchangeCodeForUser(code string) (*models.User, error) {
	// 首先用授权码获取access token
	token, err := s.getAccessTokenByCode(code)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %v", err)
	}

	// 使用access token获取用户信息
	user, err := s.getUserInfo(token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}

	return user, nil
}

// getAccessTokenByCode 用授权码获取access token (使用OAuth2标准参数)
func (s *AuthService) getAccessTokenByCode(code string) (string, error) {
	// 使用官方推荐的v2 API和OAuth2标准参数
	requestData := map[string]string{
		"grant_type":    "authorization_code",
		"client_id":     s.config.AppID,
		"client_secret": s.config.AppSecret,
		"code":          code,
		"redirect_uri":  s.config.RedirectURI,
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return "", fmt.Errorf("failed to prepare request: %v", err)
	}

	// 使用v2 API endpoint (官方推荐)
	tokenURL := s.config.BaseURL + "/open-apis/authen/v2/oauth/token"

	resp, err := http.Post(
		tokenURL,
		"application/json",
		strings.NewReader(string(requestBody)),
	)
	if err != nil {
		return "", fmt.Errorf("failed to request access token: %v", err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("access token request failed with status: %d", resp.StatusCode)
	}

	var tokenResp models.FeishuOAuthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode token response: %v", err)
	}

	return tokenResp.AccessToken, nil
}

// getUserInfo 获取用户信息
func (s *AuthService) getUserInfo(accessToken string) (*models.User, error) {
	// 创建请求
	req, err := http.NewRequest("GET", s.config.BaseURL+"/open-apis/authen/v1/user_info", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// 设置Authorization头
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// 发送请求
	client := &http.Client{
		Timeout: 10 * time.Second, // 设置超时时间
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request user info: %v", err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user info request failed with status: %d", resp.StatusCode)
	}

	var userResp models.FeishuUserInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return nil, fmt.Errorf("failed to decode user info response: %v", err)
	}

	if userResp.Code != 0 {
		return nil, fmt.Errorf("feishu API error: %s (code: %d)", userResp.Msg, userResp.Code)
	}

	// 设置provider字段
	user := userResp.Data
	user.Provider = models.ProviderFeishu

	return &user, nil
}
