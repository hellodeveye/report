package types

import "github.com/graphql-go/graphql"

// ConversationType 定义了钉钉会话的GraphQL类型
var ConversationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Conversation",
	Fields: graphql.Fields{
		"conversationId": &graphql.Field{Type: graphql.String},
		"title":          &graphql.Field{Type: graphql.String},
	},
})

// ReceiverType 定义了钉钉接收者的GraphQL类型
var ReceiverType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Receiver",
	Fields: graphql.Fields{
		"userName": &graphql.Field{Type: graphql.String},
		"userId":   &graphql.Field{Type: graphql.String},
	},
})

// FieldType 定义了钉钉模板字段的GraphQL类型
var FieldType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Field",
	Fields: graphql.Fields{
		"fieldName": &graphql.Field{Type: graphql.String},
		"sort":      &graphql.Field{Type: graphql.Int},
		"type":      &graphql.Field{Type: graphql.Int},
	},
})

// TemplateDetailType 定义了钉钉模板详情的GraphQL类型
var TemplateDetailType = graphql.NewObject(graphql.ObjectConfig{
	Name: "TemplateDetail",
	Fields: graphql.Fields{
		"defaultReceivedConvs": &graphql.Field{Type: graphql.NewList(ConversationType)},
		"defaultReceivers":     &graphql.Field{Type: graphql.NewList(ReceiverType)},
		"fields":               &graphql.Field{Type: graphql.NewList(FieldType)},
		"id":                   &graphql.Field{Type: graphql.String},
		"name":                 &graphql.Field{Type: graphql.String},
		"userName":             &graphql.Field{Type: graphql.String},
		"userId":               &graphql.Field{Type: graphql.String},
	},
})

// ReportContentType 定义了钉钉报告内容的GraphQL类型
var ReportContentType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ReportContent",
	Fields: graphql.Fields{
		"key":   &graphql.Field{Type: graphql.String},
		"value": &graphql.Field{Type: graphql.String},
		"sort":  &graphql.Field{Type: graphql.Int},
		"type":  &graphql.Field{Type: graphql.Int},
	},
})

// ReportType 定义了钉钉报告的GraphQL类型
var ReportType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Report",
	Fields: graphql.Fields{
		"report_id":      &graphql.Field{Type: graphql.String},
		"template_name":  &graphql.Field{Type: graphql.String},
		"creator_name":   &graphql.Field{Type: graphql.String},
		"creator_id":     &graphql.Field{Type: graphql.String},
		"dept_name":      &graphql.Field{Type: graphql.String},
		"remark":         &graphql.Field{Type: graphql.String},
		"create_time":    &graphql.Field{Type: graphql.Int},
		"contents":       &graphql.Field{Type: graphql.NewList(ReportContentType)},
		"read_user_list": &graphql.Field{Type: graphql.NewList(graphql.String)},
	},
})

// ReportListType 定义了钉钉报告列表的GraphQL类型
var ReportListType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ReportList",
	Fields: graphql.Fields{
		"data_list":   &graphql.Field{Type: graphql.NewList(ReportType)},
		"next_cursor": &graphql.Field{Type: graphql.Int},
		"has_more":    &graphql.Field{Type: graphql.Boolean},
		"size":        &graphql.Field{Type: graphql.Int},
	},
})

// TemplateType a a
var TemplateType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Template",
	Fields: graphql.Fields{
		"iconUrl":    &graphql.Field{Type: graphql.String},
		"name":       &graphql.Field{Type: graphql.String},
		"reportCode": &graphql.Field{Type: graphql.String},
		"url":        &graphql.Field{Type: graphql.String},
	},
})
