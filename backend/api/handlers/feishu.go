package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hellodeveye/report/internal/config"
	"github.com/hellodeveye/report/internal/models"
	"github.com/hellodeveye/report/pkg/auth"
	"github.com/hellodeveye/report/pkg/feishu"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkreport "github.com/larksuite/oapi-sdk-go/v3/service/report/v1"
)

// HealthCheckHandler 健康检查处理器
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Backend is running!")
}

// FeishuHandler 飞书相关处理器
type FeishuHandler struct {
	client *lark.Client
}

// NewFeishuHandler 创建新的飞书处理器
func NewFeishuHandler() *FeishuHandler {
	feishuConfig := config.GetFeishuConfig()
	feishuClient := feishu.NewClient(feishuConfig)

	return &FeishuHandler{
		client: feishuClient,
	}
}

// GetRules 获取报告规则
func (h *FeishuHandler) GetRuleDetail(w http.ResponseWriter, r *http.Request) {
	// 创建请求对象
	req := larkreport.NewQueryRuleReqBuilder().
		RuleName(`test-日报模版`).
		IncludeDeleted(0).
		UserIdType(`open_id`).
		Build()

	// 发起请求
	resp, err := h.client.Report.V1.Rule.Query(context.Background(), req)

	// 处理错误
	if err != nil {
		http.Error(w, fmt.Sprintf("查询规则失败: %v", err), http.StatusInternalServerError)
		return
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// 返回响应
	if err := json.NewEncoder(w).Encode(resp.Data.Rules); err != nil {
		http.Error(w, fmt.Sprintf("编码响应失败: %v", err), http.StatusInternalServerError)
		return
	}
}

// GetRules 获取报告规则
func (h *FeishuHandler) GetRules(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	// 创建请求对象
	req := larkreport.NewQueryRuleReqBuilder().
		RuleName(name).
		IncludeDeleted(0).
		UserIdType(`open_id`).
		Build()

	// 发起请求
	resp, err := h.client.Report.V1.Rule.Query(context.Background(), req)

	// 处理错误
	if err != nil {
		http.Error(w, fmt.Sprintf("查询规则失败: %v", err), http.StatusInternalServerError)
		return
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// 返回响应
	if err := json.NewEncoder(w).Encode(resp.Data.Rules); err != nil {
		http.Error(w, fmt.Sprintf("编码响应失败: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *FeishuHandler) GetReports(w http.ResponseWriter, r *http.Request) {
	ruleId := r.URL.Query().Get("rule_id")
	startTime := r.URL.Query().Get("start_time")
	endTime := r.URL.Query().Get("end_time")

	// 解析时间戳
	var startTimeInt, endTimeInt int
	if startTime != "" {
		if parsed, err := strconv.ParseInt(startTime, 10, 64); err == nil {
			startTimeInt = int(parsed)
		}
	}
	if endTime != "" {
		if parsed, err := strconv.ParseInt(endTime, 10, 64); err == nil {
			endTimeInt = int(parsed)
		}
	}

	// 创建请求对象
	req := larkreport.NewQueryTaskReqBuilder().
		UserIdType(`open_id`).
		Body(larkreport.NewQueryTaskReqBodyBuilder().
			CommitStartTime(startTimeInt).
			CommitEndTime(endTimeInt).
			RuleId(ruleId).
			PageSize(10).
			PageToken("").
			Build()).
		Build()

	// 发起请求
	resp, err := h.client.Report.V1.Task.Query(context.Background(), req)

	if err != nil {
		http.Error(w, fmt.Sprintf("查询报告失败: %v", err), http.StatusInternalServerError)
		return
	}
	if !resp.Success() {
		http.Error(w, fmt.Sprintf("查询报告失败: %s", resp.Msg), http.StatusInternalServerError)
		return
	}
	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// 返回响应
	if err := json.NewEncoder(w).Encode(resp.Data); err != nil {
		http.Error(w, fmt.Sprintf("编码响应失败: %v", err), http.StatusInternalServerError)
		return
	}

}

// Login 飞书登录处理 - 返回授权URL给前端
func (h *FeishuHandler) Login(w http.ResponseWriter, r *http.Request) {
	feishuConfig := config.GetFeishuConfig()

	// 生成state参数防CSRF
	state := fmt.Sprintf("%d", time.Now().UnixNano())

	// 构建授权URL (使用官方推荐的v1 authorize endpoint)
	// 注意：authorize使用v1，但token exchange使用v2
	authURL := fmt.Sprintf("https://accounts.feishu.cn/open-apis/authen/v1/authorize?client_id=%s&redirect_uri=%s&response_type=code&state=%s",
		feishuConfig.AppID,
		feishuConfig.RedirectURI,
		state)

	fmt.Printf("Generated authURL: %s\n", authURL)
	fmt.Printf("State: %s\n", state)

	// 返回授权URL和state给前端
	response := map[string]string{
		"auth_url": authURL,
		"state":    state,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Printf("Failed to encode login response: %v\n", err)
		http.Error(w, "Failed to generate login URL", http.StatusInternalServerError)
		return
	}
}

// ExchangeCode 处理前端发送的授权码，返回JWT token
func (h *FeishuHandler) ExchangeCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Printf("Invalid method: %s, expected POST\n", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		Code  string `json:"code"`
		State string `json:"state"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		fmt.Printf("Failed to decode request body: %v\n", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Printf("Exchange code request received - Code: %s, State: %s\n",
		requestData.Code[:min(len(requestData.Code), 10)]+"...", requestData.State)

	// 验证授权码
	if requestData.Code == "" {
		fmt.Printf("Authorization code is empty\n")
		http.Error(w, "Authorization code is required", http.StatusBadRequest)
		return
	}

	// TODO: 在生产环境中，应该验证state参数以防止CSRF攻击
	// 这里需要实现session或其他状态存储机制来验证state
	if requestData.State == "" {
		fmt.Printf("Warning: State parameter is empty, CSRF protection is disabled\n")
	}

	// 用授权码换取用户信息
	user, err := h.exchangeCodeForUser(requestData.Code)
	if err != nil {
		fmt.Printf("Failed to exchange code for user: %v\n", err)
		http.Error(w, fmt.Sprintf("Failed to authenticate: %v", err), http.StatusUnauthorized)
		return
	}

	// 生成JWT token
	token, err := auth.GenerateToken(*user)
	if err != nil {
		fmt.Printf("Failed to generate JWT token: %v\n", err)
		http.Error(w, "Failed to generate authentication token", http.StatusInternalServerError)
		return
	}

	// 返回token和用户信息给前端
	response := models.AuthToken{
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		User:      *user,
	}

	fmt.Printf("Successfully authenticated user: %s (OpenID: %s)\n", user.Name, user.OpenID)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Printf("Failed to encode response: %v\n", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// min helper function for Go versions that don't have it built-in
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GenerateDraftHandler 生成报告草稿
func GenerateDraftHandler(w http.ResponseWriter, r *http.Request) {
	// For now, we'll just simulate the draft generation
	// This mimics the logic that was previously on the frontend

	// In the future, we will parse the request to get the template ID
	// and selected source reports.

	// Simulate creating a response structure.
	// This structure should match what the frontend expects.
	// Based on the frontend code, it expects a map of field IDs to values.
	draft := make(map[string]interface{})

	// This is a simplified simulation. A real implementation would inspect
	// the requested template and generate appropriate content.
	draft["summary"] = fmt.Sprintf("<p>这是由Go后端为<b>工作总结</b>在 %s 生成的动态内容。</p>", time.Now().Format("15:04:05"))
	draft["plan"] = fmt.Sprintf("<p>这是由Go后端为<b>下周计划</b>在 %s 生成的动态内容。</p>", time.Now().Format("15:04:05"))
	draft["risk"] = "<p>后端生成：目前没有可见风险。</p>"
	draft["kpi"] = "后端生成：完成了80%的后端开发任务。"
	draft["achievements"] = "<p>后端生成：成功搭建了API框架。</p>"
	draft["learnings"] = "<p>后端生成：学习了如何在Go中处理JSON。</p>"
	draft["next_month_goals"] = "<p>后端生成：完成与飞书API的对接。</p>"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(draft)
}

// exchangeCodeForUser 用授权码换取用户信息
func (h *FeishuHandler) exchangeCodeForUser(code string) (*models.User, error) {
	// 首先用授权码获取access token
	token, err := h.getAccessTokenByCode(code)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %v", err)
	}

	// 使用access token获取用户信息
	user, err := h.getUserInfo(token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}

	return user, nil
}

// getAccessTokenByCode 用授权码获取access token (使用OAuth2标准参数)
func (h *FeishuHandler) getAccessTokenByCode(code string) (string, error) {
	feishuConfig := config.GetFeishuConfig()

	// 使用官方推荐的v2 API和OAuth2标准参数
	requestData := map[string]string{
		"grant_type":    "authorization_code",
		"client_id":     feishuConfig.AppID,
		"client_secret": feishuConfig.AppSecret,
		"code":          code,
		"redirect_uri":  feishuConfig.RedirectURI,
	}

	fmt.Printf("OAuth2 config: AppID=%s, RedirectURI=%s\n", feishuConfig.AppID, feishuConfig.RedirectURI)

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		fmt.Printf("Failed to marshal request data: %v\n", err)
		return "", fmt.Errorf("failed to prepare request: %v", err)
	}

	fmt.Printf("Requesting access token with code: %s\n", code[:min(len(code), 20)]+"...")

	// 使用v2 API endpoint (官方推荐)
	tokenURL := feishuConfig.BaseURL + "/open-apis/authen/v2/oauth/token"
	fmt.Printf("Token endpoint: %s\n", tokenURL)

	resp, err := http.Post(
		tokenURL,
		"application/json",
		strings.NewReader(string(requestBody)),
	)
	if err != nil {
		fmt.Printf("Failed to request access token: %v\n", err)
		return "", fmt.Errorf("failed to request access token: %v", err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Access token request failed with status: %d\n", resp.StatusCode)
		return "", fmt.Errorf("access token request failed with status: %d", resp.StatusCode)
	}

	// 读取响应体用于调试
	respBody := make([]byte, 0)
	var tokenResp models.FeishuOAuthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		fmt.Printf("Failed to decode token response: %v\n", err)
		fmt.Printf("Response body: %s\n", string(respBody))
		return "", fmt.Errorf("failed to decode token response: %v", err)
	}

	fmt.Printf("Successfully obtained access token\n")
	return tokenResp.AccessToken, nil
}

// getUserInfo 获取用户信息
func (h *FeishuHandler) getUserInfo(accessToken string) (*models.User, error) {
	feishuConfig := config.GetFeishuConfig()

	fmt.Printf("Requesting user info with access token\n")

	// 创建请求
	req, err := http.NewRequest("GET", feishuConfig.BaseURL+"/open-apis/authen/v1/user_info", nil)
	if err != nil {
		fmt.Printf("Failed to create user info request: %v\n", err)
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
		fmt.Printf("Failed to request user info: %v\n", err)
		return nil, fmt.Errorf("failed to request user info: %v", err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("User info request failed with status: %d\n", resp.StatusCode)
		return nil, fmt.Errorf("user info request failed with status: %d", resp.StatusCode)
	}

	var userResp models.FeishuUserInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		fmt.Printf("Failed to decode user info response: %v\n", err)
		return nil, fmt.Errorf("failed to decode user info response: %v", err)
	}

	if userResp.Code != 0 {
		fmt.Printf("Feishu API error when getting user info: %s (code: %d)\n", userResp.Msg, userResp.Code)
		return nil, fmt.Errorf("feishu API error: %s (code: %d)", userResp.Msg, userResp.Code)
	}

	fmt.Printf("Successfully obtained user info for: %s\n", userResp.Data.Name)
	return &userResp.Data, nil
}

// GetCurrentUser 获取当前登录用户信息
func (h *FeishuHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// 从上下文中获取用户信息
	openID, ok := r.Context().Value("user_open_id").(string)
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	userName, _ := r.Context().Value("user_name").(string)

	user := models.User{
		OpenID: openID,
		Name:   userName,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Logout 退出登录
func (h *FeishuHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// 由于JWT是无状态的，退出登录主要在前端处理（删除本地存储的token）
	// 这里可以添加token黑名单逻辑，或者其他服务端清理逻辑

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logged out successfully",
	})
}
