package api

import (
	"github.com/gorilla/mux"
	"github.com/hellodeveye/report/api/handlers"
	"github.com/hellodeveye/report/api/middleware"
)

// SetupRoutes 设置所有API路由
func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// 添加CORS中间件
	r.Use(middleware.CORS)

	// 健康检查端点
	r.HandleFunc("/", handlers.HealthCheckHandler).Methods("GET")

	// API路由组
	api := r.PathPrefix("/api").Subrouter()

	// 创建飞书处理器
	feishuHandler := handlers.NewFeishuHandler()
	// 创建钉钉处理器
	dingTalkHandler := handlers.NewDingTalkHandler()

	// 认证相关路由（无需登录）
	api.HandleFunc("/auth/feishu/login", feishuHandler.Login).Methods("GET")
	api.HandleFunc("/auth/feishu/exchange", feishuHandler.ExchangeCode).Methods("POST")
	api.HandleFunc("/auth/dingtalk/login", dingTalkHandler.Login).Methods("GET")
	api.HandleFunc("/auth/dingtalk/exchange", dingTalkHandler.ExchangeCode).Methods("POST")
	api.HandleFunc("/auth/logout", feishuHandler.Logout).Methods("POST")

	// 需要认证的路由
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	protected.HandleFunc("/auth/user", feishuHandler.GetCurrentUser).Methods("GET")

	protected.HandleFunc("/dingtalk/templates", dingTalkHandler.GetTemplates).Methods("GET")
	protected.HandleFunc("/dingtalk/templates/detail", dingTalkHandler.GetTemplateDetail).Methods("GET")
	protected.HandleFunc("/dingtalk/reports", dingTalkHandler.GetReports).Methods("GET")
	protected.HandleFunc("/dingtalk/reports", dingTalkHandler.Create).Methods("POST")

	protected.HandleFunc("/feishu/templates", feishuHandler.GetTemplates).Methods("GET")
	protected.HandleFunc("/feishu/templates/detail", feishuHandler.GetTemplateDetail).Methods("GET")
	protected.HandleFunc("/feishu/reports", feishuHandler.GetReports).Methods("GET")

	return r
}
