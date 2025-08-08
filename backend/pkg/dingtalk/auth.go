package dingtalk

import (
	"fmt"
	"net/url"
	"time"

	"github.com/hellodeveye/report/internal/models"
)

// AuthService 钉钉认证服务
type AuthService struct {
	client *Client
	config *models.DingTalkConfig
}

// NewAuthService 创建新的钉钉认证服务
func NewAuthService(config *models.DingTalkConfig) *AuthService {
	return &AuthService{
		client: NewClient(config),
		config: config,
	}
}

// GenerateAuthURL 生成授权URL
func (s *AuthService) GenerateAuthURL() (string, string, error) {
	// 生成state参数防CSRF
	state := fmt.Sprintf("%d", time.Now().UnixNano())

	// 构建授权URL
	authURL := fmt.Sprintf("https://login.dingtalk.com/oauth2/auth?redirect_uri=%s&response_type=code&client_id=%s&scope=openid corpid&state=%s&prompt=consent&corpId=%s",
		url.QueryEscape(s.config.RedirectURI),
		s.config.AppKey,
		state,
		s.config.CorpId,
	)

	return authURL, state, nil
}

// ExchangeCodeForUser 用授权码换取用户信息
func (s *AuthService) ExchangeCodeForUser(code string) (*models.User, error) {
	// 1. 用授权码获取用户访问令牌
	tokenResp, err := s.client.GetUserAccessToken(code)
	if err != nil {
		return nil, fmt.Errorf("failed to get user access token: %v", err)
	}

	// 2. 用用户访问令牌获取用户信息
	userResp, err := s.client.GetUserInfo(tokenResp.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}

	accessToken, err := s.client.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %v", err)
	}

	userByUnionIdResp, err := s.client.GetUserByUnionId(accessToken.AccessToken, userResp.UnionId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by union id: %v", err)
	}
	userResp.UserID = userByUnionIdResp.Result.UserID

	// 3. 转换为统一的用户模型
	user := s.convertToUser(userResp)
	return user, nil
}

// convertToUser 将钉钉用户信息转换为统一的用户模型
func (s *AuthService) convertToUser(userResp *models.DingTalkUserInfoResponse) *models.User {
	return &models.User{
		OpenID:  userResp.OpenId,    // 钉钉的openId
		UnionID: userResp.UnionId,   // 钉钉的unionId
		UserID:  userResp.UserID,    // 使用openId作为UserID
		Name:    userResp.Nick,      // 用户昵称
		Avatar:  userResp.AvatarUrl, // 头像URL
		Email:   userResp.Email,     // 邮箱地址
		Mobile:  userResp.Mobile,    // 手机号码
	}
}
