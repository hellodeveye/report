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
	AppID     string
	AppSecret string
}
