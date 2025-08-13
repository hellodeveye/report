package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hellodeveye/report/internal/config"
	"github.com/hellodeveye/report/internal/models"
	"github.com/hellodeveye/report/pkg/auth"
	"github.com/hellodeveye/report/pkg/dingtalk"
)

// DingTalkHandler 钉钉相关处理器
type DingTalkHandler struct {
	authService   *dingtalk.AuthService
	reportService *dingtalk.ReportService
}

// NewDingTalkHandler 创建新的钉钉处理器
func NewDingTalkHandler() *DingTalkHandler {
	dingTalkConfig := config.GetDingTalkConfig()
	authService := dingtalk.NewAuthService(dingTalkConfig)
	dingtalkClient := dingtalk.NewClient(dingTalkConfig)
	reportService := dingtalk.NewReportService(dingtalkClient)

	return &DingTalkHandler{
		authService:   authService,
		reportService: reportService,
	}
}

// Login 钉钉登录处理 - 返回授权URL给前端
func (h *DingTalkHandler) Login(w http.ResponseWriter, r *http.Request) {
	authURL, state, err := h.authService.GenerateAuthURL()
	if err != nil {
		fmt.Printf("Failed to generate DingTalk auth URL: %v\n", err)
		http.Error(w, "Failed to generate login URL", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Generated DingTalk authURL: %s\n", authURL)
	fmt.Printf("State: %s\n", state)

	// 返回授权URL和state给前端
	response := map[string]string{
		"auth_url": authURL,
		"state":    state,
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
	user, err := h.authService.ExchangeCodeForUser(requestData.Code)
	if err != nil {
		fmt.Printf("Failed to exchange code for user: %v\n", err)
		http.Error(w, fmt.Sprintf("Authentication failed: %v", err), http.StatusUnauthorized)
		return
	}

	// 生成JWT token
	token, expiresAt, err := auth.GenerateToken(user.UserID, user.Name)
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
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(authResponse); err != nil {
		fmt.Printf("Failed to encode auth response: %v\n", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	fmt.Printf("DingTalk authentication successful for user: %s\n", user.Name)
}

func (h *DingTalkHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// 由于JWT是无状态的，退出登录主要在前端处理（删除本地存储的token）
	// 这里可以添加token黑名单逻辑，或者其他服务端清理逻辑

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logged out successfully",
	})
}
