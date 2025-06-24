package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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
	// TODO: Get parameters from request body
	// TODO: Call Feishu API to get source reports
	// TODO: Aggregate content (and optionally call an LLM)
	// TODO: Return the draft
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprint(w, "TODO: Implement Draft Generation")
}
