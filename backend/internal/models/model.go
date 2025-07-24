package models

// FeishuRule 飞书报告规则
type FeishuRule struct {
	RuleID   string `json:"rule_id"`
	RuleName string `json:"rule_name"`
	// 可以根据实际API响应添加更多字段
}

// FeishuRulesResponse 飞书规则查询响应
type FeishuRulesResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Items []FeishuRule `json:"items"`
		Total int          `json:"total"`
	} `json:"data"`
}

// FeishuAuthResponse 飞书认证响应
type FeishuAuthResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	} `json:"data"`
}

// FeishuConfig 飞书配置
type FeishuConfig struct {
	AppID       string
	AppSecret   string
	RedirectURI string
	BaseURL     string
}

// User 用户信息
type User struct {
	OpenID   string   `json:"open_id"`
	UnionID  string   `json:"union_id"`
	UserID   string   `json:"userid"`
	Name     string   `json:"name"`
	Avatar   string   `json:"avatar_url"`
	Email    string   `json:"email"`
	Mobile   string   `json:"mobile"`
	Provider Provider `json:"provider"`
}

// AuthToken JWT认证token
type AuthToken struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
	User      User   `json:"user"`
}

// FeishuOAuthTokenResponse 飞书OAuth token响应
type FeishuOAuthTokenResponse struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshToken     string `json:"refresh_token"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	Scope            string `json:"scope"`
}

// FeishuUserInfoResponse 飞书用户信息响应
type FeishuUserInfoResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data User   `json:"data"`
}

// DingTalkConfig 钉钉配置
type DingTalkConfig struct {
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

// Provider 表示认证提供商
type Provider string

const (
	ProviderFeishu   Provider = "feishu"
	ProviderDingTalk Provider = "dingtalk"
)

// AuthRequest 通用认证请求
type AuthRequest struct {
	Provider Provider `json:"provider"`
	Code     string   `json:"code"`
	State    string   `json:"state"`
}

// AuthResponse 通用认证响应
type AuthResponse struct {
	Token     string   `json:"token"`
	ExpiresAt int64    `json:"expires_at"`
	User      User     `json:"user"`
	Provider  Provider `json:"provider"`
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
