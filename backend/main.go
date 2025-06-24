package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// A simple health check endpoint
	r.HandleFunc("/", HealthCheckHandler).Methods("GET")

	api := r.PathPrefix("/api").Subrouter()

	// TODO: Task 4 - Implement Feishu OAuth 2.0 login flow
	api.HandleFunc("/auth/feishu/login", FeishuLoginHandler).Methods("GET")
	api.HandleFunc("/auth/feishu/callback", FeishuCallbackHandler).Methods("GET")

	// TODO: Task 7 - Create a protected API endpoint to get report rules
	api.HandleFunc("/rules", GetRulesHandler).Methods("GET")

	// TODO: Task 10 - Create the core draft generation endpoint
	api.HandleFunc("/generate-draft", GenerateDraftHandler).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server starting on port " + port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Backend is running!")
}

// FeishuLoginHandler redirects users to the Feishu authorization page.
func FeishuLoginHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Construct the Feishu OAuth URL and redirect
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprint(w, "TODO: Implement Feishu Login")
}

// FeishuCallbackHandler handles the callback from Feishu after user authorization.
func FeishuCallbackHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Exchange the authorization code for an access_token
	// TODO: Store the access_token (e.g., in Redis)
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprint(w, "TODO: Implement Feishu Callback")
}

// GetRulesHandler fetches the report rules (templates) from Feishu.
func GetRulesHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Get access_token from storage
	// TODO: Call Feishu API to get rules
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprint(w, "TODO: Implement Get Rules")
}

// GenerateDraftHandler generates a report draft.
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
