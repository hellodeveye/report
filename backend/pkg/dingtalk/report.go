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
