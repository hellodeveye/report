package resolvers

import (
	"github.com/graphql-go/graphql"
	"github.com/hellodeveye/report/graphql/types"
	"github.com/hellodeveye/report/pkg/dingtalk"
)

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
