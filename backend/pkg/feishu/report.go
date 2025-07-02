package feishu

import (
	"context"
	"strconv"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkreport "github.com/larksuite/oapi-sdk-go/v3/service/report/v1"
)

// ReportService 飞书报告服务
type ReportService struct {
	client *lark.Client
}

// NewReportService 创建新的飞书报告服务
func NewReportService(client *lark.Client) *ReportService {
	return &ReportService{
		client: client,
	}
}

// QueryRules 查询报告规则
func (s *ReportService) QueryRules(name string) ([]*larkreport.Rule, error) {
	// 创建请求对象
	req := larkreport.NewQueryRuleReqBuilder().
		RuleName(name).
		IncludeDeleted(0).
		UserIdType(`open_id`).
		Build()

	// 发起请求
	resp, err := s.client.Report.V1.Rule.Query(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return resp.Data.Rules, nil
}

// QueryRuleDetail 查询报告规则详情
func (s *ReportService) QueryRuleDetail(ruleName string) ([]*larkreport.Rule, error) {
	// 创建请求对象
	req := larkreport.NewQueryRuleReqBuilder().
		RuleName(ruleName).
		IncludeDeleted(0).
		UserIdType(`open_id`).
		Build()

	// 发起请求
	resp, err := s.client.Report.V1.Rule.Query(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return resp.Data.Rules, nil
}

// QueryReports 查询报告
func (s *ReportService) QueryReports(ruleId, startTime, endTime string) (*larkreport.QueryTaskRespData, error) {
	// 解析时间戳
	var startTimeInt, endTimeInt int
	if startTime != "" {
		if parsed, err := strconv.ParseInt(startTime, 10, 64); err == nil {
			startTimeInt = int(parsed)
		}
	}
	if endTime != "" {
		if parsed, err := strconv.ParseInt(endTime, 10, 64); err == nil {
			endTimeInt = int(parsed)
		}
	}

	// 创建请求对象
	req := larkreport.NewQueryTaskReqBuilder().
		UserIdType(`open_id`).
		Body(larkreport.NewQueryTaskReqBodyBuilder().
			CommitStartTime(startTimeInt).
			CommitEndTime(endTimeInt).
			RuleId(ruleId).
			PageSize(10).
			PageToken("").
			Build()).
		Build()

	// 发起请求
	resp, err := s.client.Report.V1.Task.Query(context.Background(), req)
	if err != nil {
		return nil, err
	}

	if !resp.Success() {
		return nil, &ReportError{
			Code:    resp.Code,
			Message: resp.Msg,
		}
	}

	return resp.Data, nil
}

// ReportError 报告错误类型
type ReportError struct {
	Code    int
	Message string
}

func (e *ReportError) Error() string {
	return e.Message
}
