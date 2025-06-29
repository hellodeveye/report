package config

import (
	"os"

	"github.com/hellodeveye/report/internal/models"
)

// GetFeishuConfig 获取飞书配置
func GetFeishuConfig() *models.FeishuConfig {
	return &models.FeishuConfig{
		AppID:       getEnv("FEISHU_APP_ID", ""),
		AppSecret:   getEnv("FEISHU_APP_SECRET", ""),
		RedirectURI: getEnv("FEISHU_REDIRECT_URI", ""),
		BaseURL:     getEnv("FEISHU_BASE_URL", "https://open.feishu.cn"),
	}
}

// GetFeishuAccessToken 获取飞书访问令牌
func GetFeishuAccessToken() string {
	return getEnv("FEISHU_ACCESS_TOKEN", "")
}

// GetJWTSecret 获取JWT密钥
func GetJWTSecret() string {
	return getEnv("JWT_SECRET", "default-jwt-secret-change-in-production")
}

// GetDingTalkConfig 获取钉钉配置
func GetDingTalkConfig() *models.DingTalkConfig {
	return &models.DingTalkConfig{
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
