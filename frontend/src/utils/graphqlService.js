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

  async getTemplates(userId) {
    const variables = {};
    const query = gql`
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
    const data = await this.request(query, variables);
    const templates = data.dingtalkTemplates.map(t => ({ ...t, ...t.detail }));
    return templates.map(({ detail, ...rest }) => rest);
  },



  async getReports(params) {
    const query = gql`
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
    const variables = {
      template_name: params.template_name,
      start_time: params.start_time,
      end_time: params.end_time,
      cursor: params.cursor,
      size: params.size,
    };
    const data = await this.request(query, variables);
    return data.dingtalkReports;
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
