import { authService } from './authService.js';
import graphqlService from './graphqlService.js';

class ApiService {
    constructor() {
        this.baseURL = '/api';
    }

    // 获取模板列表
    async getTemplates() {
        const user = authService.getUser();
        try {
            const templates = await graphqlService.getTemplates(user?.userid);
            
            return templates.map(t => ({
                id: t.reportCode,
                name: t.name,
                fields: t.fields.map((field, index) => {
                    const fieldType = this.mapDingTalkFieldType(field.type);
                    const baseField = {
                        id: `field_${t.reportCode}_${index}`,
                        label: field.fieldName,
                        type: fieldType,
                        placeholder: `请输入${field.fieldName}...`,
                    };
                    // Add other field properties as needed
                    return baseField;
                })
            }));
        } catch (error) {
            console.error(`获取模板列表失败:`, error);
            throw new Error(`获取模板列表失败: ${error.message}`);
        }
    }



    // 获取报告列表
    async getReports(params) {
        const user = authService.getUser();
        if (!user) throw new Error("用户未登录");

        try {
            const reportsData = await graphqlService.getReports(user.userid, params);

            if (!reportsData || !reportsData.data_list || !Array.isArray(reportsData.data_list)) {
                console.warn('钉钉API返回的数据格式不正确:', reportsData);
                return [];
            }
            return reportsData.data_list.map(report => ({
                id: report.report_id,
                title: `${report.template_name} - ${report.creator_name} (${new Date(parseInt(report.create_time)).toLocaleString('zh-CN')})`,
                isCollapsed: true,
                fields: (report.contents || []).map(content => ({
                    name: content.key,
                    value: content.value,
                    type: 'tiptap'
                }))
            }));
        } catch (error) {
            console.error(`获取报告列表失败:`, error);
            throw new Error(`获取报告列表失败: ${error.message}`);
        }
    }

    async sendDingTalkReport(reportData) {
        try {
            return await graphqlService.createDingTalkReport(reportData);
        } catch (error) {
            console.error(`发送报告失败:`, error);
            throw new Error(`发送报告失败: ${error.message}`);
        }
    }



    // 字段类型映射 for DingTalk
    mapDingTalkFieldType(apiType) {
        const typeMap = {
            1: 'tiptap',      // 文本
            2: 'number',      // 数字
            3: 'dropdown',    // 单选
            5: 'datetime',    // 日期
            7: 'multiSelect', // 多选
            8: 'image',       // 图片
            9: 'attachment',  // 附件
            12: 'user-picker', // 客户
            // 16 is table, ignored
        };
        return typeMap[apiType] || 'tiptap';
    }
}

export const apiService = new ApiService();
