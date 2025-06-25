package main

import (
	"log"
	"net/http"
	"os"

	"github.com/hellodeveye/report/api"
)

func main() {
	// 设置路由
	router := api.SetupRoutes()

	// 获取端口
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server starting on port " + port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}
