package dingtalk

import (
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

	resp, err := s.client.httpClient.R().SetBody(jsonData).Post("https://oapi.dingtalk.com/topapi/report/template/listbyuserid?access_token=" + accessToken.AccessToken)
	if err != nil {
		log.Println("request failed:", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("read response body failed:", err)
		return nil, err
	}

	var response TemplateListResponse

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

type ReportContent struct {
	Key   string `json:"key"`
	Sort  string `json:"sort"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type ReportData struct {
	Contents     []ReportContent `json:"contents"`
	CreateTime   int64           `json:"create_time"`
	CreatorID    string          `json:"creator_id"`
	CreatorName  string          `json:"creator_name"`
	DeptName     string          `json:"dept_name"`
	ModifiedTime int64           `json:"modified_time"`
	ReportID     string          `json:"report_id"`
	TemplateName string          `json:"template_name"`
}

type ReportListResult struct {
	DataList   []ReportData `json:"data_list"`
	HasMore    bool         `json:"has_more"`
	NextCursor int64        `json:"next_cursor"`
	Size       int          `json:"size"`
}

type ReportListResponse struct {
	ErrCode   int              `json:"errcode"`
	ErrMsg    string           `json:"errmsg"`
	Result    ReportListResult `json:"result"`
	RequestID string           `json:"request_id"`
}

func (s *ReportService) GetReports(userID string, templateName string, startTime, endTime int64, cursor, size int) (*ReportListResponse, error) {
	accessToken, err := s.client.GetAccessToken()
	if err != nil {
		return nil, err
	}

	url := "https://oapi.dingtalk.com/topapi/report/list?access_token=" + accessToken.AccessToken

	requestBody := map[string]interface{}{
		"userid":        userID,
		"template_name": templateName,
		"start_time":    startTime,
		"end_time":      endTime,
		"cursor":        cursor,
		"size":          size,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.httpClient.R().SetBody(jsonData).Post(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response ReportListResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("Failed to unmarshal DingTalk report list response. Body: %s", string(body))
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

	resp, err := s.client.httpClient.R().SetBody(jsonData).Post(url)
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

type CreateReportRequest struct {
	CreateReportParam struct {
		Contents   []ContentItem `json:"contents"`
		DDFrom     string        `json:"dd_from"`
		TemplateID string        `json:"template_id"`
		UserID     string        `json:"userid"`
		ToChat     bool          `json:"to_chat"`
		ToCIDs     []string      `json:"to_cids"`
		ToUserIDs  []string      `json:"to_userids"`
	} `json:"create_report_param"`
}

type ContentItem struct {
	ContentType string `json:"content_type"`
	Sort        int    `json:"sort"`
	Type        int    `json:"type"`
	Content     string `json:"content"`
	Key         string `json:"key"`
}

type CreateReportResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	Result  struct {
		ReportID string `json:"report_id"`
		Status   string `json:"status"`
	} `json:"result"`
	RequestID string `json:"request_id"`
}

// 保存草稿
func (s *ReportService) SaveDraft(userId string, draftRequest *CreateReportRequest) (*CreateReportResponse, error) {
	accessToken, err := s.client.GetAccessToken()
	if err != nil {
		return nil, err
	}

	url := "https://oapi.dingtalk.com/topapi/report/create?access_token=" + accessToken.AccessToken

	// 设置草稿状态
	if draftRequest.CreateReportParam.ToChat {
		// 如果是草稿，不发送到聊天
		// 钉钉API中，to_chat=true表示发送，false表示草稿
		// 为了保存草稿，我们设置to_chat=false
		// 注意：这里可能需要根据实际需求调整
		// 保持原始设置，由调用者控制
	}

	jsonData, err := json.Marshal(draftRequest)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.httpClient.R().SetBody(jsonData).Post(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response CreateReportResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

type SaveReportParam struct {
	Contents   []ContentItem `json:"contents"`
	DDFrom     string        `json:"dd_from"`
	TemplateID string        `json:"template_id"`
	UserID     string        `json:"userid"`
}

type SaveReportResponse struct {
	ErrCode   int    `json:"errcode"`
	Result    string `json:"result"`
	RequestID string `json:"request_id"`
}

func (s *ReportService) SaveContent(userId string, param SaveReportParam) (*SaveReportResponse, error) {
	accessToken, err := s.client.GetAccessToken()
	if err != nil {
		return nil, err
	}

	url := "https://oapi.dingtalk.com/topapi/report/savecontent?access_token=" + accessToken.AccessToken

	jsonData, err := json.Marshal(param)

	resp, err := s.client.httpClient.R().SetBody(jsonData).Post(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response SaveReportResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
