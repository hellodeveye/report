package api

import (
	"github.com/gorilla/mux"
	"github.com/graphql-go/handler"
	"github.com/hellodeveye/report/api/handlers"
	"github.com/hellodeveye/report/api/middleware"
	"github.com/hellodeveye/report/graphql"
)

// SetupRoutes 设置所有API路由
func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// 添加CORS中间件
	r.Use(middleware.CORS)

	// API路由组
	api := r.PathPrefix("/api").Subrouter()

	// 创建钉钉处理器
	dingTalkHandler := handlers.NewDingTalkHandler()

	// 认证相关路由（无需登录）
	api.HandleFunc("/auth/dingtalk/login", dingTalkHandler.Login).Methods("GET")
	api.HandleFunc("/auth/dingtalk/exchange", dingTalkHandler.ExchangeCode).Methods("POST")
	api.HandleFunc("/auth/logout", dingTalkHandler.Logout).Methods("POST")

	// 创建 GraphQL HTTP 处理器
	schema := graphql.SetupGraphQLSchema()
	h := handler.New(&handler.Config{
		Schema:   schema,
		Pretty:   true,
		GraphiQL: true,
	})
	r.HandleFunc("/graphql", h.ServeHTTP)

	// 需要认证的路由
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	protected.HandleFunc("/graphql", h.ServeHTTP)

	return r
}
