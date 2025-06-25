package config

import (
	"os"

	"github.com/hellodeveye/report/internal/models"
)

// GetFeishuConfig 获取飞书配置
func GetFeishuConfig() *models.FeishuConfig {
	return &models.FeishuConfig{
		AppID:     getEnv("FEISHU_APP_ID", ""),
		AppSecret: getEnv("FEISHU_APP_SECRET", ""),
	}
}

// GetFeishuAccessToken 获取飞书访问令牌
func GetFeishuAccessToken() string {
	return getEnv("FEISHU_ACCESS_TOKEN", "")
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
