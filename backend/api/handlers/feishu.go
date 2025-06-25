package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hellodeveye/report/internal/config"
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

// Login 飞书登录处理
func (h *FeishuHandler) Login(w http.ResponseWriter, r *http.Request) {
	// TODO: 实现飞书OAuth登录流程
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprint(w, "TODO: 实现飞书登录")
}

// Callback 飞书回调处理
func (h *FeishuHandler) Callback(w http.ResponseWriter, r *http.Request) {
	// TODO: 实现飞书OAuth回调处理
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprint(w, "TODO: 实现飞书回调")
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
