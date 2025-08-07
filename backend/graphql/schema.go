package graphql

import (
	"log"

	"github.com/graphql-go/graphql"
	"github.com/hellodeveye/report/graphql/resolvers"
	"github.com/hellodeveye/report/graphql/types"
	"github.com/hellodeveye/report/internal/config"
	"github.com/hellodeveye/report/pkg/dingtalk"
	"github.com/hellodeveye/report/pkg/feishu"
)

func SetupGraphQLSchema() *graphql.Schema {
	// DingTalk Services
	dingTalkConfig := config.GetDingTalkConfig()
	dingtalkClient := dingtalk.NewClient(dingTalkConfig)
	dingtalkReportService := dingtalk.NewReportService(dingtalkClient)

	// Feishu Services
	feishuConfig := config.GetFeishuConfig()
	feishuClient := feishu.NewClient(feishuConfig)
	feishuReportService := feishu.NewReportService(feishuClient)

	// Initialize resolvers
	resolvers.InitDingTalkResolvers(dingtalkReportService)
	resolvers.InitFeishuResolvers(feishuReportService)

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: graphql.Fields{
		"dingtalkTemplates": &graphql.Field{
			Type: graphql.NewList(types.TemplateType),
			Args: graphql.FieldConfigArgument{
				"userId": &graphql.ArgumentConfig{Type: graphql.String},
				"name":   &graphql.ArgumentConfig{Type: graphql.String},
			},
			Resolve: resolvers.GetDingTalkTemplatesResolver,
		},
		"dingtalkReports": &graphql.Field{
			Type: types.ReportListType,
			Args: graphql.FieldConfigArgument{
				"template_name": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"start_time":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"end_time":      &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"cursor":        &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 0},
				"size":          &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 20},
			},
			Resolve: resolvers.GetDingTalkReportsResolver,
		},
		"feishuTemplates": &graphql.Field{
			Type:    graphql.NewList(types.FeishuTemplateType),
			Resolve: resolvers.GetFeishuTemplatesResolver,
		},
		"feishuTemplateDetail": &graphql.Field{
			Type: types.FeishuRuleType,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: resolvers.GetFeishuTemplateDetailResolver,
		},
		"feishuReports": &graphql.Field{
			Type: types.FeishuReportListType,
			Args: graphql.FieldConfigArgument{
				"rule_id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"start_time": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"end_time":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: resolvers.GetFeishuReportsResolver,
		},
	}}

	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"createDingtalkReport": &graphql.Field{
				Type: types.ReportType,
				Args: graphql.FieldConfigArgument{
					"template_name": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"template_id":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"contents":      &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.NewList(types.ReportContentInputType))},
				},
				Resolve: resolvers.CreateDingTalkReportResolver,
			},
		},
	})

	schemaConfig := graphql.SchemaConfig{
		Query:    graphql.NewObject(rootQuery),
		Mutation: rootMutation,
	}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}
	return &schema
}
