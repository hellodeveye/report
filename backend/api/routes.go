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

	// 飞书相关路由
	api.HandleFunc("/auth/feishu/login", feishuHandler.Login).Methods("GET")
	api.HandleFunc("/auth/feishu/callback", feishuHandler.Callback).Methods("GET")
	api.HandleFunc("/rules", feishuHandler.GetRules).Methods("GET")

	// 报告生成路由
	api.HandleFunc("/generate-draft", handlers.GenerateDraftHandler).Methods("POST")

	return r
}
