package config

import (
	"os"

	"github.com/hellodeveye/report/internal/models"
)

// GetJWTSecret 获取JWT密钥
func GetJWTSecret() string {
	return getEnv("JWT_SECRET", "default-jwt-secret-change-in-production")
}

// GetDingTalkConfig 获取钉钉配置
func GetDingTalkConfig() *models.DingTalkConfig {
	return &models.DingTalkConfig{
		CorpId:      getEnv("DINGTALK_CORP_ID", ""),
		AppKey:      getEnv("DINGTALK_APP_KEY", ""),
		AppSecret:   getEnv("DINGTALK_APP_SECRET", ""),
		RedirectURI: getEnv("DINGTALK_REDIRECT_URI", ""),
		BaseURL:     getEnv("DINGTALK_BASE_URL", "https://oapi.dingtalk.com"),
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
