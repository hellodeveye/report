package resolvers

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/hellodeveye/report/internal/models"
	"github.com/hellodeveye/report/pkg/auth"
	"github.com/hellodeveye/report/pkg/dingtalk"
	"github.com/hellodeveye/report/pkg/feishu"
)

var feishuReportService *feishu.ReportService

func InitFeishuResolvers(service *feishu.ReportService) {
	feishuReportService = service
}
func GetDingTalkTemplatesResolver(p graphql.ResolveParams) (interface{}, error) {
	userId, _ := p.Args["userId"].(string)
	templates, err := dingtalkReportService.GetTemplates(userId)
	if err != nil {
		return nil, err
	}
	if name, ok := p.Args["name"].(string); ok && name != "" {
		for _, template := range templates.Result.TemplateList {
			if template.Name == name {
				return []dingtalk.TemplateItem{template}, nil
			}
		}
		return nil, nil // Return empty if not found
	}
	return templates.Result.TemplateList, nil
}

func GetDingTalkReportsResolver(p graphql.ResolveParams) (interface{}, error) {
	userID, ok := auth.GetUserOpenID(p.Context)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	templateName, _ := p.Args["template_name"].(string)
	startTime, _ := p.Args["start_time"].(int)
	endTime, _ := p.Args["end_time"].(int)
	cursor, _ := p.Args["cursor"].(int)
	size, _ := p.Args["size"].(int)
	return dingtalkReportService.GetReports(userID, templateName, int64(startTime), int64(endTime), cursor, size)
}

func GetFeishuTemplatesResolver(p graphql.ResolveParams) (interface{}, error) {
	return []models.TemplateInfo{
		{ID: "monthly_report", Name: "工作月报"},
		{ID: "daily_report", Name: "技术部-工作日报"},
		{ID: "daily_scrum", Name: "每日站会"},
		{ID: "complex_template", Name: "复杂模板"},
	}, nil
}

func GetFeishuTemplateDetailResolver(p graphql.ResolveParams) (interface{}, error) {
	name, _ := p.Args["name"].(string)
	return feishuReportService.QueryRules(name)
}

func GetFeishuReportsResolver(p graphql.ResolveParams) (interface{}, error) {
	ruleID, _ := p.Args["rule_id"].(string)
	startTime, _ := p.Args["start_time"].(string)
	endTime, _ := p.Args["end_time"].(string)
	return feishuReportService.QueryReports(ruleID, startTime, endTime)
}
