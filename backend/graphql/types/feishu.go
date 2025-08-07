package types

import "github.com/graphql-go/graphql"

// FeishuTemplateType defines the GraphQL type for a Feishu template
var FeishuTemplateType = graphql.NewObject(graphql.ObjectConfig{
	Name: "FeishuTemplate",
	Fields: graphql.Fields{
		"id":   &graphql.Field{Type: graphql.String},
		"name": &graphql.Field{Type: graphql.String},
	},
})

// FeishuRuleType defines the GraphQL type for a Feishu rule
var FeishuRuleType = graphql.NewObject(graphql.ObjectConfig{
	Name: "FeishuRule",
	Fields: graphql.Fields{
		"id":         &graphql.Field{Type: graphql.String},
		"name":       &graphql.Field{Type: graphql.String},
		"is_default": &graphql.Field{Type: graphql.Boolean},
	},
})

// FeishuReportItemType defines the GraphQL type for a Feishu report item
var FeishuReportItemType = graphql.NewObject(graphql.ObjectConfig{
	Name: "FeishuReportItem",
	Fields: graphql.Fields{
		"report_id":      &graphql.Field{Type: graphql.String},
		"title":          &graphql.Field{Type: graphql.String},
		"cycle":          &graphql.Field{Type: graphql.String},
		"create_time":    &graphql.Field{Type: graphql.String},
		"submitter_name": &graphql.Field{Type: graphql.String},
		"submit_time":    &graphql.Field{Type: graphql.String},
	},
})

// FeishuReportListType defines the GraphQL type for a list of Feishu reports
var FeishuReportListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "FeishuReportList",
	Fields: graphql.Fields{
		"items":       &graphql.Field{Type: graphql.NewList(FeishuReportItemType)},
		"page_token":  &graphql.Field{Type: graphql.String},
		"has_more":    &graphql.Field{Type: graphql.Boolean},
		"total_count": &graphql.Field{Type: graphql.Int},
	},
})

var ReportContentInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "ReportContentInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"key":   &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"value": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
	},
})
