package models

// User 用户信息
type User struct {
	OpenID  string `json:"open_id"`
	UnionID string `json:"union_id"`
	UserID  string `json:"userid"`
	Name    string `json:"name"`
	Avatar  string `json:"avatar_url"`
	Email   string `json:"email"`
	Mobile  string `json:"mobile"`
}

// AuthToken JWT认证token
type AuthToken struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
	User      User   `json:"user"`
}

// DingTalkConfig 钉钉配置
type DingTalkConfig struct {
	CorpId      string
	AppKey      string
	AppSecret   string
	RedirectURI string
	BaseURL     string
}

// DingTalkOAuthTokenResponse 钉钉OAuth token响应
type DingTalkOAuthTokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpireIn     int64  `json:"expireIn"`
	CorpId       string `json:"corpId"`
}

// DingTalkUserInfoResponse 钉钉用户信息响应
type DingTalkUserInfoResponse struct {
	UserID    string `json:"userid"`
	Nick      string `json:"nick"`
	AvatarUrl string `json:"avatarUrl"`
	Mobile    string `json:"mobile"`
	OpenId    string `json:"openId"`
	UnionId   string `json:"unionId"`
	Email     string `json:"email"`
	StateCode string `json:"stateCode"`
}

// AuthRequest 通用认证请求
type AuthRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

// AuthResponse 通用认证响应
type AuthResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
	User      User   `json:"user"`
}

type DingTalkAccessTokenResponse struct {
	ErrCode     int    `json:"errcode"`
	AccessToken string `json:"access_token"`
	ErrMsg      string `json:"errmsg"`
	ExpiresIn   int    `json:"expires_in"`
}

// TemplateInfo 模板信息
type TemplateInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// DingTalkUserByUnionIdResponse 钉钉通过UnionID获取用户响应
type DingTalkUserByUnionIdResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	Result  struct {
		ContactType int    `json:"contact_type"`
		UserID      string `json:"userid"`
	} `json:"result"`
	RequestID string `json:"request_id"`
}

type DingTalkTemplate struct {
	Name       string `json:"name"`
	ReportCode string `json:"report_code"`
}

type DingTalkTemplateResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	Result  struct {
		TemplateList []DingTalkTemplate `json:"template_list"`
	} `json:"result"`
	RequestID string `json:"request_id"`
}
