package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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

func (h *DingTalkHandler) GetTemplates(w http.ResponseWriter, r *http.Request) {
	openID, ok := auth.GetUserOpenID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	templates, err := h.reportService.GetTemplates(openID)
	if err != nil {
		http.Error(w, "Failed to get templates", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(templates); err != nil {
		http.Error(w, "Failed to encode templates", http.StatusInternalServerError)
		return
	}
}

func (h *DingTalkHandler) GetTemplateDetail(w http.ResponseWriter, r *http.Request) {
	openID, ok := auth.GetUserOpenID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	templateName := r.URL.Query().Get("template_name")
	if templateName == "" {
		http.Error(w, "Missing template name", http.StatusBadRequest)
		return
	}

	detail, err := h.reportService.GetTemplateDetail(openID, templateName)
	if err != nil {
		http.Error(w, "Failed to get template detail", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(detail); err != nil {
		http.Error(w, "Failed to encode template detail", http.StatusInternalServerError)
		return
	}
}

type CreateReportPayload struct {
	TemplateID   string `json:"template_id"`
	TemplateName string `json:"template_name"`
	Contents     []struct {
		Key   string      `json:"key"`
		Value interface{} `json:"value"`
	} `json:"contents"`
}

func (h *DingTalkHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserOpenID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var payload CreateReportPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 1. Get template detail to map fields
	templateDetail, err := h.reportService.GetTemplateDetail(userID, payload.TemplateName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get template details: %v", err), http.StatusInternalServerError)
		return
	}
	if templateDetail.ErrCode != 0 {
		// Marshal the result to a string for better error logging
		resultBytes, _ := json.Marshal(templateDetail.Result)
		http.Error(w, fmt.Sprintf("Failed to get template details from DingTalk, errorCode: %d, result: %s", templateDetail.ErrCode, string(resultBytes)), http.StatusInternalServerError)
		return
	}

	fieldMap := make(map[string]dingtalk.Field)
	for _, field := range templateDetail.Result.Fields {
		fieldMap[field.FieldName] = field
	}

	// 2. Build the final request for DingTalk API
	var reportContents []dingtalk.ContentItem
	for _, content := range payload.Contents {
		field, exists := fieldMap[content.Key]
		if !exists {
			// Skip content if key doesn't match any field in the template
			continue
		}

		reportContents = append(reportContents, dingtalk.ContentItem{
			Key:         field.FieldName,
			Sort:        field.Sort,
			Type:        field.Type,
			Content:     content.Value.(string),
			ContentType: "markdown", // 默认markdown
		})
	}

	createReq := dingtalk.CreateReportRequest{
		CreateReportParam: struct {
			Contents   []dingtalk.ContentItem `json:"contents"`
			DDFrom     string                 `json:"dd_from"`
			TemplateID string                 `json:"template_id"`
			UserID     string                 `json:"userid"`
			ToChat     bool                   `json:"to_chat"`
			ToCIDs     []string               `json:"to_cids"`
			ToUserIDs  []string               `json:"to_userids"`
		}{
			TemplateID: payload.TemplateID,
			UserID:     userID,
			Contents:   reportContents,
			// Other fields can be set if needed
		},
	}

	response, err := h.reportService.Create(userID, &createReq)
	if err != nil {
		http.Error(w, "Failed to create report", http.StatusInternalServerError)
		return
	}

	if response.ErrCode != 0 {
		http.Error(w, fmt.Sprintf("Failed to create report on DingTalk: %s", response.ErrMsg), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *DingTalkHandler) GetReports(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserOpenID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	templateName := r.URL.Query().Get("template_name")
	if templateName == "" {
		http.Error(w, "template_name is required", http.StatusBadRequest)
		return
	}
	startTimeStr := r.URL.Query().Get("start_time")
	endTimeStr := r.URL.Query().Get("end_time")
	cursorStr := r.URL.Query().Get("cursor")
	sizeStr := r.URL.Query().Get("size")

	var startTime, endTime int64
	var cursor, size int
	var err error

	if startTimeStr != "" {
		startTime, err = strconv.ParseInt(startTimeStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid start_time", http.StatusBadRequest)
			return
		}
	}
	if endTimeStr != "" {
		endTime, err = strconv.ParseInt(endTimeStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid end_time", http.StatusBadRequest)
			return
		}
	}
	if cursorStr != "" {
		cursor, err = strconv.Atoi(cursorStr)
		if err != nil {
			http.Error(w, "Invalid cursor", http.StatusBadRequest)
			return
		}
	} else {
		cursor = 0
	}
	if sizeStr != "" {
		size, err = strconv.Atoi(sizeStr)
		if err != nil {
			http.Error(w, "Invalid size", http.StatusBadRequest)
			return
		}
	} else {
		size = 20 // Default size
	}

	reportsResponse, err := h.reportService.GetReports(userID, templateName, startTime*1000, endTime*1000, cursor, size)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get reports: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(reportsResponse); err != nil {
		http.Error(w, "Failed to encode reports", http.StatusInternalServerError)
		return
	}
}
