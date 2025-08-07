package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hellodeveye/report/internal/config"
	"github.com/hellodeveye/report/internal/models"
	"github.com/hellodeveye/report/pkg/auth"
	"github.com/hellodeveye/report/pkg/feishu"
)

// HealthCheckHandler 健康检查处理器
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Backend is running!")
}

// FeishuHandler 飞书相关处理器
type FeishuHandler struct {
	authService   *feishu.AuthService
	reportService *feishu.ReportService
}

// NewFeishuHandler 创建新的飞书处理器
func NewFeishuHandler() *FeishuHandler {
	feishuConfig := config.GetFeishuConfig()
	feishuClient := feishu.NewClient(feishuConfig)
	authService := feishu.NewAuthService(feishuConfig)
	reportService := feishu.NewReportService(feishuClient)

	return &FeishuHandler{
		authService:   authService,
		reportService: reportService,
	}
}

// Login 飞书登录处理 - 返回授权URL给前端
func (h *FeishuHandler) Login(w http.ResponseWriter, r *http.Request) {
	authURL, state, err := h.authService.GenerateAuthURL()
	if err != nil {
		fmt.Printf("Failed to generate auth URL: %v\n", err)
		http.Error(w, "Failed to generate login URL", http.StatusInternalServerError)
		return
	}

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

	var requestData models.AuthRequest

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		fmt.Printf("Failed to decode request body: %v\n", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	codePreview := requestData.Code
	if len(codePreview) > 10 {
		codePreview = codePreview[:10] + "..."
	}
	fmt.Printf("Exchange code request received - Code: %s, State: %s\n", codePreview, requestData.State)

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
	user, err := h.authService.ExchangeCodeForUser(requestData.Code)
	if err != nil {
		fmt.Printf("Failed to exchange code for user: %v\n", err)
		http.Error(w, fmt.Sprintf("Failed to authenticate: %v", err), http.StatusUnauthorized)
		return
	}

	// 生成JWT token
	token, expiresAt, err := auth.GenerateToken(user.OpenID, user.Name)
	if err != nil {
		fmt.Printf("Failed to generate JWT token: %v\n", err)
		http.Error(w, "Failed to generate authentication token", http.StatusInternalServerError)
		return
	}

	// 返回token和用户信息给前端
	response := models.AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      *user,
		Provider:  models.ProviderFeishu,
	}

	fmt.Printf("Successfully authenticated user: %s (OpenID: %s)\n", user.Name, user.OpenID)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Printf("Failed to encode response: %v\n", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetCurrentUser 获取当前登录用户信息
func (h *FeishuHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// 从上下文中获取用户信息
	openID, ok := auth.GetUserOpenID(r.Context())
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	userName, _ := auth.GetUserName(r.Context())

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
