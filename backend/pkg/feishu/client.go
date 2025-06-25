package feishu

import (
	"github.com/hellodeveye/report/internal/models"
	lark "github.com/larksuite/oapi-sdk-go/v3"
)

// NewClient 创建新的飞书API客户端
func NewClient(config *models.FeishuConfig) *lark.Client {
	client := lark.NewClient(config.AppID, config.AppSecret)
	return client
}
