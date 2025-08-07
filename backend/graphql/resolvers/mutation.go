package resolvers

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/hellodeveye/report/pkg/auth"
	"github.com/hellodeveye/report/pkg/dingtalk"
)

func CreateDingTalkReportResolver(p graphql.ResolveParams) (interface{}, error) {
	userID, ok := auth.GetUserOpenID(p.Context)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	templateName, _ := p.Args["template_name"].(string)
	templateID, _ := p.Args["template_id"].(string)
	contents, _ := p.Args["contents"].([]interface{})

	templateDetail, err := dingtalkReportService.GetTemplateDetail(userID, templateName)
	if err != nil {
		return nil, fmt.Errorf("failed to get template details: %v", err)
	}
	fieldMap := make(map[string]dingtalk.Field)
	for _, field := range templateDetail.Result.Fields {
		fieldMap[field.FieldName] = field
	}
	var reportContents []dingtalk.ContentItem
	for _, c := range contents {
		contentMap := c.(map[string]interface{})
		key := contentMap["key"].(string)
		value := contentMap["value"].(string)
		if field, exists := fieldMap[key]; exists {
			reportContents = append(reportContents, dingtalk.ContentItem{
				Key: field.FieldName, Sort: field.Sort, Type: field.Type, Content: value, ContentType: "markdown",
			})
		}
	}
	createReq := dingtalk.CreateReportRequest{
		CreateReportParam: struct {
			Contents   []dingtalk.ContentItem `json:"contents"`
			DDFrom     string                 `json:"dd_from"`
			TemplateID string                 `json:"template_id"`
			UserID     string                 `json:"userid"`
			ToChat     bool                   `json:"to_chat"`
			ToCIDs     []string               `json:"to_cids"`
			ToUserIDs  []string               `json:"to_userids"`
		}{
			TemplateID: templateID, UserID: userID, Contents: reportContents,
		},
	}
	createResp, err := dingtalkReportService.Create(userID, &createReq)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"report_id": createResp.Result,
	}, nil
}
