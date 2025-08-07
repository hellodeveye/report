import { authService } from './authService.js';
import graphqlService from './graphqlService.js';

class ApiService {
    constructor() {
        this.baseURL = '/api';
    }

    // 获取当前提供商
    getProvider() {
        const user = authService.getUser();
        return user ? user.provider : null;
    }

    // 获取模板列表
    async getTemplates() {
        const provider = this.getProvider();
        const user = authService.getUser();
        
        try {
            const templates = await graphqlService.getTemplates(provider, user?.id);
            // 钉钉和飞书返回的格式需要统一
            if (provider === 'dingtalk') {
                return templates.map(t => ({ id: t.reportCode, name: t.name }));
            }
            return templates; // 飞书的格式已经是 { id, name }
        } catch (error) {
            console.error(`获取模板列表失败:`, error);
            throw new Error(`获取模板列表失败: ${error.message}`);
        }
    }

    // 获取模板详情
    async getTemplateDetail(templateName, templateId) {
        const provider = this.getProvider();
        if (!provider) throw new Error("用户未登录");
        const user = authService.getUser();

        try {
            const detail = await graphqlService.getTemplateDetail(provider, templateName, user?.id);
            
            // 统一数据格式
            if (provider === 'dingtalk') {
                if (!detail) {
                    throw new Error('获取钉钉模板详情失败: 响应为空');
                }
                return {
                    id: templateId, // 钉钉详情接口不返回id，我们从列表传入
                    name: detail.name,
                    fields: detail.fields.map((field, index) => {
                        const fieldType = this.mapDingTalkFieldType(field.type);
                        const baseField = {
                            id: `field_${templateId}_${index}`,
                            label: field.fieldName,
                            type: fieldType,
                            placeholder: `请输入${field.fieldName}...`,
                        };
                         if (fieldType === 'dropdown' || fieldType === 'multiSelect') {
                            baseField.options = field.options || [
                                { value: 'option1', text: '选项1' },
                                { value: 'option2', text: '选项2' },
                            ];
                        }
                        if (fieldType === 'image') {
                            baseField.maxCount = 99;
                            baseField.maxSize = 20 * 1024 * 1024;
                        }
                         if (fieldType === 'attachment') {
                            baseField.maxCount = 9;
                            baseField.maxSize = 50 * 1024 * 1024;
                        }
                        return baseField;
                    })
                };
            }
    
            // 飞书的详情返回的是一个对象
            return {
                id: detail.id,
                name: detail.name,
                fields: (detail.form_schema || []).map((field, index) => ({
                    id: `field_${detail.id}_${index}`,
                    label: field.name,
                    type: this.mapFeishuFieldType(field.type),
                    placeholder: `请输入${field.name}...`,
                }))
            };
        } catch (error) {
            console.error(`获取模板详情失败:`, error);
            throw new Error(`获取模板详情失败: ${error.message}`);
        }
    }

    // 获取报告列表
    async getReports(params) {
        const provider = this.getProvider();
        if (!provider) throw new Error("用户未登录");

        try {
            const reportsData = await graphqlService.getReports(provider, params);

            // 统一数据格式
            if (provider === 'feishu') {
                if (!reportsData || !reportsData.items || !Array.isArray(reportsData.items)) {
                   console.warn('API返回的数据格式不正确:', reportsData);
                   return [];
               }
               return reportsData.items.map(report => ({
                   id: report.report_id,
                   title: `${report.title} - ${report.submitter_name} (${new Date(parseInt(report.submit_time, 10) * 1000).toLocaleString('zh-CN')})`,
                   isCollapsed: true,
                   fields: (report.form_contents || []).map(content => ({
                       name: content.field_name,
                       value: content.field_value,
                       type: 'tiptap' // 默认为富文本
                   }))
               }));
           } else if (provider === 'dingtalk') {
               if (!reportsData || !reportsData.data_list || !Array.isArray(reportsData.data_list)) {
                   console.warn('钉钉API返回的数据格式不正确:', reportsData);
                   return [];
               }
               return reportsData.data_list.map(report => ({
                   id: report.report_id,
                   title: `${report.template_name} - ${report.creator_name} (${new Date(report.create_time).toLocaleString('zh-CN')})`,
                   isCollapsed: true,
                   fields: (report.contents || []).map(content => ({
                       name: content.key,
                       value: content.value,
                       type: 'tiptap'
                   }))
               }));
           }
   
           return [];
        } catch (error) {
            console.error(`获取报告列表失败:`, error);
            throw new Error(`获取报告列表失败: ${error.message}`);
        }
    }

    async sendDingTalkReport(reportData) {
        const provider = this.getProvider();
        if (provider !== 'dingtalk') {
            throw new Error("This function is only available for DingTalk.");
        }

        try {
            return await graphqlService.createDingTalkReport(reportData);
        } catch (error) {
            console.error(`发送报告失败:`, error);
            throw new Error(`发送报告失败: ${error.message}`);
        }
    }

    // 字段类型映射
    mapFeishuFieldType(apiType) {
        const typeMap = {
            'text': 'tiptap',
            'number': 'number',
            'dropdown': 'dropdown',
            'image': 'image',
            'attachment': 'attachment',
            'multiSelect': 'multiSelect',
            'address': 'address',
            'datetime': 'datetime'
        };
        return typeMap[apiType] || 'tiptap';
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
