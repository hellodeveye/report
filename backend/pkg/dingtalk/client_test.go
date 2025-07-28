package dingtalk

import (
	"os"
	"testing"

	"github.com/hellodeveye/report/internal/models"
)

func TestGetAccessToken(t *testing.T) {
	client := NewClient(&models.DingTalkConfig{
		AppKey:    os.Getenv("DINGTALK_APP_KEY"),
		AppSecret: os.Getenv("DINGTALK_APP_SECRET"),
	})
	accessToken, err := client.GetAccessToken()
	if err != nil {
		t.Fatalf("GetAccessToken failed: %v", err)
	}

	if accessToken.ErrCode != 0 {
		t.Fatalf("GetAccessToken failed: %v", accessToken)
	}
	t.Logf("accessToken: %v", accessToken.AccessToken)
}
