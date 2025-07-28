import { authService } from './authService.js';

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
        if (!provider) {
            // 默认返回飞书的模板列表
            const url = `${this.baseURL}/feishu/templates`;
            const response = await fetch(url);
            if (!response.ok) {
                console.error("获取默认模板列表失败");
                return [];
            }
            return await response.json();
        }

        const url = `${this.baseURL}/${provider}/templates`;
        const response = await authService.authenticatedFetch(url);

        if (!response.ok) {
            throw new Error(`获取模板列表失败: ${response.status} ${response.statusText}`);
        }

        const data = await response.json();
        // 钉钉和飞书返回的格式需要统一
        if (provider === 'dingtalk') {
            // 钉钉返回的是 { result: { template_list: [...] } }
            return data.result.template_list.map(t => ({ id: t.report_code, name: t.name }));
        }
        // 飞书返回的是 [{ id: '...', name: '...' }]
        return data;
    }

    // 获取模板详情
    async getTemplateDetail(templateName, templateId) {
        const provider = this.getProvider();
        if (!provider) throw new Error("用户未登录");

        let url;
        if (provider === 'dingtalk') {
            url = `${this.baseURL}/${provider}/templates/detail?template_name=${encodeURIComponent(templateName)}`;
        } else {
            // 飞书使用name查询，但我们将其映射到templateId
            url = `${this.baseURL}/${provider}/templates/detail?name=${encodeURIComponent(templateName)}`;
        }

        const response = await authService.authenticatedFetch(url);
        if (!response.ok) {
            throw new Error(`获取模板详情失败: ${response.status} ${response.statusText}`);
        }
        
        const data = await response.json();

        // 统一数据格式
        if (provider === 'dingtalk') {
            if (!data.result) {
                throw new Error(`获取钉钉模板详情失败: ${data.errmsg || '响应格式不正确'}`);
            }
            return {
                id: templateId, // 钉钉详情接口不返回id，我们从列表传入
                name: data.result.name,
                fields: data.result.fields.map((field, index) => {
                    const fieldType = this.mapDingTalkFieldType(field.type);
                    const baseField = {
                        id: `field_${templateId}_${index}`,
                        label: field.field_name,
                        type: fieldType,
                        placeholder: `请输入${field.field_name}...`,
                    };
                    // 为钉钉的单选和多选添加占位选项
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

        // 飞书的详情返回的是一个数组，取第一个
        const rule = data[0];
        return {
            id: rule.rule_id,
            name: rule.name,
            fields: rule.form_schema.map((field, index) => ({
                id: `field_${rule.rule_id}_${index}`,
                label: field.name,
                type: this.mapFeishuFieldType(field.type),
                placeholder: `请输入${field.name}...`,
            }))
        };
    }

    // 获取报告列表
    async getReports(params) {
        const provider = this.getProvider();
        if (!provider) throw new Error("用户未登录");

        const queryParams = new URLSearchParams();
        if (provider === 'feishu') {
            if (params.rule_id) queryParams.append('rule_id', params.rule_id);
            if (params.start_time) queryParams.append('start_time', params.start_time);
            if (params.end_time) queryParams.append('end_time', params.end_time);
        } else if (provider === 'dingtalk') {
            if (params.template_name) queryParams.append('template_name', params.template_name);
            if (params.start_time) queryParams.append('start_time', params.start_time);
            if (params.end_time) queryParams.append('end_time', params.end_time);
        }

        const url = `${this.baseURL}/${provider}/reports${queryParams.toString() ? '?' + queryParams.toString() : ''}`;
        const response = await authService.authenticatedFetch(url);

        if (!response.ok) {
            throw new Error(`获取报告失败: ${response.status} ${response.statusText}`);
        }

        const data = await response.json();

        // 统一数据格式
        if (provider === 'feishu') {
             if (!data || !data.items || !Array.isArray(data.items)) {
                console.warn('API返回的数据格式不正确:', data);
                return [];
            }
            return data.items.map(report => ({
                id: report.task_id,
                title: `${report.rule_name} - ${report.from_user_name} (${new Date(report.commit_time * 1000).toLocaleString('zh-CN')})`,
                isCollapsed: true,
                fields: (report.form_contents || []).map(content => ({
                    name: content.field_name,
                    value: content.field_value,
                    type: 'tiptap' // 默认为富文本
                }))
            }));
        } else if (provider === 'dingtalk') {
            if (!data || !data.result || !Array.isArray(data.result.data_list)) {
                console.warn('钉钉API返回的数据格式不正确:', data);
                return [];
            }
            return data.result.data_list.map(report => ({
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

        // Fallback for any other case
        return [];
    }

    async sendDingTalkReport(reportData) {
        const provider = this.getProvider();
        if (provider !== 'dingtalk') {
            throw new Error("This function is only available for DingTalk.");
        }

        const url = `${this.baseURL}/dingtalk/reports`;
        const response = await authService.authenticatedFetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(reportData),
        });

        if (!response.ok) {
            const errorData = await response.json().catch(() => ({ message: `HTTP error! status: ${response.status}` }));
            throw new Error(errorData.message || '发送报告失败');
        }

        return await response.json();
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