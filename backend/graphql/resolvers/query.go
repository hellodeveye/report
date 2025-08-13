package resolvers

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/hellodeveye/report/graphql/types"
	"github.com/hellodeveye/report/pkg/auth"
	"github.com/hellodeveye/report/pkg/dingtalk"
)

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
	reports, err := dingtalkReportService.GetReports(userID, templateName, int64(startTime), int64(endTime), cursor, size)
	if err != nil {
		return nil, err
	}
	return reports.Result, nil
}

var dingtalkReportService *dingtalk.ReportService

func InitDingTalkResolvers(service *dingtalk.ReportService) {
	dingtalkReportService = service
	types.TemplateType.AddFieldConfig("detail", &graphql.Field{
		Type: types.TemplateDetailType,
		Args: graphql.FieldConfigArgument{
			"userId": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: GetTemplateDetailResolver,
	})
}

func GetTemplateDetailResolver(p graphql.ResolveParams) (interface{}, error) {
	template, _ := p.Source.(dingtalk.TemplateItem)
	userId, _ := p.Args["userId"].(string)
	templateDetail, err := dingtalkReportService.GetTemplateDetail(userId, template.Name)
	if err != nil {
		return nil, err
	}
	return templateDetail.Result, nil
}
