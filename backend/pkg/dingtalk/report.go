package dingtalk

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
)

type ReportService struct {
	client *Client
}

func NewReportService(client *Client) *ReportService {
	return &ReportService{client: client}
}

type TemplateItem struct {
	IconURL    string `json:"icon_url"`
	Name       string `json:"name"`
	ReportCode string `json:"report_code"`
	URL        string `json:"url"`
}

type TemplateListResult struct {
	TemplateList []TemplateItem `json:"template_list"`
}

type TemplateListResponse struct {
	ErrCode   int                `json:"errcode"`
	ErrMsg    string             `json:"errmsg"`
	Result    TemplateListResult `json:"result"`
	RequestID string             `json:"request_id"`
}

func (s *ReportService) GetTemplates(userId string) (*TemplateListResponse, error) {

	requestBody := map[string]interface{}{
		"userid": userId,
		"offset": "0",
		"size":   "100",
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.client.GetAccessToken()
	if err != nil {
		return nil, err
	}

	resp, err := s.client.httpClient.Post("https://oapi.dingtalk.com/topapi/report/template/listbyuserid?access_token="+accessToken.AccessToken, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("request failed:", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response TemplateListResponse

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// 获取模板详情
func (s *ReportService) GetTemplateDetail(userId, template_name string) (*TemplateDetailResponse, error) {
	accessToken, err := s.client.GetAccessToken()
	if err != nil {
		return nil, err
	}

	url := "https://oapi.dingtalk.com/topapi/report/template/getbyname?access_token=" + accessToken.AccessToken

	requestBody := map[string]interface{}{
		"userid":        userId,
		"template_name": template_name,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response TemplateDetailResponse

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

type TemplateDetailResponse struct {
	ErrCode   int                  `json:"errcode"`
	Result    TemplateDetailResult `json:"result"`
	RequestID string               `json:"request_id"`
}

type TemplateDetailResult struct {
	DefaultReceivedConvs []Conversation `json:"default_received_convs"`
	DefaultReceivers     []Receiver     `json:"default_receivers"`
	Fields               []Field        `json:"fields"`
	ID                   string         `json:"id"`
	Name                 string         `json:"name"`
	UserName             string         `json:"user_name"`
	UserID               string         `json:"userid"`
}

type Conversation struct {
	ConversationID string `json:"conversation_id"`
	Title          string `json:"title"`
}

type Receiver struct {
	UserName string `json:"user_name"`
	UserID   string `json:"userid"`
}

type Field struct {
	FieldName string `json:"field_name"`
	Sort      int    `json:"sort"`
	Type      int    `json:"type"`
}
