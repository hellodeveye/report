import { GraphQLClient, gql } from 'graphql-request';
import { authService } from './authService';
const endpoint = (import.meta.env.VITE_API_BASE_URL || '') + '/api/graphql';
const client = new GraphQLClient(endpoint);

const graphqlService = {
  async request(query, variables = {}) {
    const token = authService.getToken();
    const headers = {};
    if (token) {
      headers.Authorization = `Bearer ${token}`;
    }
    return client.request(query, variables, headers);
  },

  async getTemplates(provider, userId) {
    let query;
    const variables = {};

    if (provider === 'dingtalk') {
      query = gql`
        query GetDingTalkTemplates($userId: String!) {
          dingtalkTemplates(userId: $userId) {
            name
            reportCode
            detail(userId: $userId) {
              id
              name
              fields {
                fieldName
                type
              }
            }
          }
        }
      `;
      variables.userId = userId;
    } else { // default to feishu
      query = gql`
        query GetFeishuTemplates {
          feishuTemplates {
            id
            name
            # Assuming the schema supports getting fields directly.
            # If not, the backend schema will need to be adjusted.
            fields {
              key
              title
              type
            }
          }
        }
      `;
    }
    
    const data = await this.request(query, variables);
    const templates = provider === 'dingtalk' 
      ? data.dingtalkTemplates.map(t => ({ ...t, ...t.detail }))
      : data.feishuTemplates;

    // We no longer need the nested detail object for dingtalk
    return templates.map(({ detail, ...rest }) => rest);
  },



  async getReports(provider, params) {
    let query;
    let variables;

    if (provider === 'dingtalk') {
      query = gql`
        query GetDingTalkReports($template_name: String!, $start_time: Int!, $end_time: Int!, $cursor: Int, $size: Int) {
          dingtalkReports(template_name: $template_name, start_time: $start_time, end_time: $end_time, cursor: $cursor, size: $size) {
            data_list {
              report_id
              template_name
              creator_name
              create_time
              contents {
                key
                value
              }
            }
          }
        }
      `;
      variables = {
        template_name: params.template_name,
        start_time: params.start_time,
        end_time: params.end_time,
        cursor: params.cursor,
        size: params.size,
      };
    } else { // feishu
      query = gql`
        query GetFeishuReports($rule_id: String!, $start_time: String!, $end_time: String!) {
          feishuReports(rule_id: $rule_id, start_time: $start_time, end_time: $end_time) {
            items {
              report_id
              title
              submitter_name
              submit_time
            }
          }
        }
      `;
      variables = {
        rule_id: params.rule_id,
        start_time: params.start_time,
        end_time: params.end_time,
      };
    }

    const data = await this.request(query, variables);
    return provider === 'dingtalk' ? data.dingtalkReports : data.feishuReports;
  },

  async createDingTalkReport(reportData) {
    const mutation = gql`
      mutation CreateDingTalkReport($template_name: String!, $template_id: String!, $contents: [ReportContentInput!]!) {
        createDingtalkReport(template_name: $template_name, template_id: $template_id, contents: $contents) {
          report_id
        }
      }
    `;
    return this.request(mutation, reportData);
  }
};

export default graphqlService;
