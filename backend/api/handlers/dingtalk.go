package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/hellodeveye/report/internal/config"
	"github.com/hellodeveye/report/internal/models"
	"github.com/hellodeveye/report/pkg/auth"
	"github.com/hellodeveye/report/pkg/dingtalk"
)

// DingTalkHandler 钉钉相关处理器
type DingTalkHandler struct {
	client *dingtalk.Client
}

// NewDingTalkHandler 创建新的钉钉处理器
func NewDingTalkHandler() *DingTalkHandler {
	dingTalkConfig := config.GetDingTalkConfig()
	dingTalkClient := dingtalk.NewClient(dingTalkConfig)

	return &DingTalkHandler{
		client: dingTalkClient,
	}
}

// Login 钉钉登录处理 - 返回授权URL给前端
func (h *DingTalkHandler) Login(w http.ResponseWriter, r *http.Request) {
	dingTalkConfig := config.GetDingTalkConfig()

	// 生成state参数防CSRF
	state := fmt.Sprintf("%d", time.Now().UnixNano())

	// 构建授权URL
	authURL := fmt.Sprintf("https://login.dingtalk.com/oauth2/auth?redirect_uri=%s&response_type=code&client_id=%s&scope=openid corpid&state=%s&prompt=consent",
		url.QueryEscape(dingTalkConfig.RedirectURI),
		dingTalkConfig.AppKey,
		state)

	fmt.Printf("Generated DingTalk authURL: %s\n", authURL)
	fmt.Printf("State: %s\n", state)

	// 返回授权URL和state给前端
	response := map[string]string{
		"auth_url": authURL,
		"state":    state,
		"provider": "dingtalk",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Printf("Failed to encode DingTalk login response: %v\n", err)
		http.Error(w, "Failed to generate login URL", http.StatusInternalServerError)
		return
	}
}

// ExchangeCode 处理前端发送的授权码，返回JWT token
func (h *DingTalkHandler) ExchangeCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Printf("Invalid method: %s, expected POST\n", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestData models.AuthRequest

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		fmt.Printf("Failed to decode request body: %v\n", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received DingTalk auth request - Code: %s, State: %s\n", requestData.Code, requestData.State)

	if requestData.Code == "" {
		fmt.Printf("Missing authorization code\n")
		http.Error(w, "Missing authorization code", http.StatusBadRequest)
		return
	}

	// 用授权码换取用户信息
	user, err := h.exchangeCodeForUser(requestData.Code)
	if err != nil {
		fmt.Printf("Failed to exchange code for user: %v\n", err)
		http.Error(w, fmt.Sprintf("Authentication failed: %v", err), http.StatusUnauthorized)
		return
	}

	// 生成JWT token
	token, expiresAt, err := auth.GenerateToken(user.OpenID, user.Name)
	if err != nil {
		fmt.Printf("Failed to generate JWT token: %v\n", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// 返回JWT token和用户信息
	authResponse := models.AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      *user,
		Provider:  models.ProviderDingTalk,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(authResponse); err != nil {
		fmt.Printf("Failed to encode auth response: %v\n", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	fmt.Printf("DingTalk authentication successful for user: %s\n", user.Name)
}

// exchangeCodeForUser 用授权码换取用户信息
func (h *DingTalkHandler) exchangeCodeForUser(code string) (*models.User, error) {
	// 1. 用授权码获取用户访问令牌
	tokenResp, err := h.getUserAccessToken(code)
	if err != nil {
		return nil, fmt.Errorf("failed to get user access token: %v", err)
	}

	// 2. 用用户访问令牌获取用户信息
	user, err := h.getUserInfo(tokenResp.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}

	return user, nil
}

// getUserAccessToken 获取用户访问令牌
func (h *DingTalkHandler) getUserAccessToken(code string) (*models.DingTalkOAuthTokenResponse, error) {
	return h.client.GetUserAccessToken(code)
}

// getUserInfo 用用户访问令牌获取用户信息
func (h *DingTalkHandler) getUserInfo(accessToken string) (*models.User, error) {
	userResp, err := h.client.GetUserInfo(accessToken)
	if err != nil {
		return nil, err
	}

	// 转换为统一的用户模型
	user := &models.User{
		OpenID:   userResp.OpenId,         // 钉钉的openId
		UnionID:  userResp.UnionId,        // 钉钉的unionId
		UserID:   userResp.OpenId,         // 使用openId作为UserID
		Name:     userResp.Nick,           // 用户昵称
		Avatar:   userResp.AvatarUrl,      // 头像URL
		Email:    userResp.Email,          // 邮箱地址
		Mobile:   userResp.Mobile,         // 手机号码
		Provider: models.ProviderDingTalk, // 标识为钉钉用户
	}

	return user, nil
}
