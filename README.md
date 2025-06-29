# Feishu Report Assistant

This project is a web application designed to help users automatically generate draft weekly and monthly reports by pulling and summarizing their daily reports from Feishu.

## Project Structure

The project is organized into three main directories:

-   `/backend`: The Go application that handles business logic, communicates with the Feishu Open Platform, and serves the API.
-   `/frontend`: The Vue.js single-page application that provides the user interface.
-   `/deployment`: Contains Docker, Docker Compose, and Kubernetes configuration files for deployment.

## Tech Stack

-   **Backend**: Go, Gorilla/Mux, Redis
-   **Frontend**: Vue.js, Vite, Tailwind CSS, Tiptap Editor
-   **Deployment**: Docker, Docker Compose, Kubernetes

## Getting Started

### Prerequisites

-   Docker and Docker Compose
-   `kubectl` for Kubernetes deployment (optional)
-   A Feishu App with `App ID` and `App Secret`
-   Node.js 18+ and Go 1.18+ (for local development)

### Running with Docker Compose (Recommended for Local Development)

1.  **Configure Apps:**
    
    **For Feishu:**
    Create a Feishu application in the [Feishu Developer Console](https://open.feishu.cn/app):
    - Set the redirect URI to: `http://localhost:5173/auth/callback`
    - Enable required permissions: user info access
    - Note down your App ID and App Secret
    
    **For DingTalk:**
    Create a DingTalk application in the [DingTalk Developer Console](https://open-dev.dingtalk.com):
    - Set the redirect URI to: `http://localhost:5173/auth/callback`
    - Enable required permissions: user info access
    - Note down your App Key and App Secret

2.  **Configure Environment Variables:**
    Create a `.env` file in the root directory and add your app details:

    ```env
    # 飞书配置
    FEISHU_APP_ID=your_app_id
    FEISHU_APP_SECRET=your_app_secret
    FEISHU_REDIRECT_URI=http://localhost:5173/auth/callback
    FEISHU_BASE_URL=https://open.feishu.cn
    
    # 钉钉配置
    DINGTALK_APP_KEY=your_dingtalk_app_key
    DINGTALK_APP_SECRET=your_dingtalk_app_secret
    DINGTALK_REDIRECT_URI=http://localhost:5173/auth/callback
    DINGTALK_BASE_URL=https://oapi.dingtalk.com
    
    # 通用配置
    JWT_SECRET=your-jwt-secret-key-change-in-production
    FRONTEND_URL=http://localhost:5173
    ```

3.  **Build and Run:**

    ```bash
    docker-compose up --build
    ```

4.  **Access the application:**
    -   Frontend: `http://localhost:5173`
    -   Backend API: `http://localhost:8080`

5.  **Login Process:**
    - Visit `http://localhost:5173` in your browser
    - Choose either "使用飞书登录" or "使用钉钉登录" button
    - You'll be redirected to the selected platform's authentication page
    - Scan the QR code with the corresponding mobile app or login with credentials
    - After successful authentication, you'll be redirected back to the main application

### Running with Kubernetes

Navigate to the `deployment` directory and apply the configurations:

```bash
# Deploy Redis
kubectl apply -f redis.yaml

# Deploy Backend
kubectl apply -f backend.yaml

# Deploy Frontend
kubectl apply -f frontend.yaml
```

Make sure to configure secrets for the Feishu App credentials within your Kubernetes cluster. 

## API集成说明

### 认证接口
- **Login**: `GET /api/auth/feishu/login` - 获取飞书OAuth登录URL
- **Exchange Code**: `POST /api/auth/feishu/exchange` - 用授权码换取JWT token
- **Current User**: `GET /api/auth/user` - 获取当前用户信息（需要认证）
- **Logout**: `POST /api/auth/logout` - 退出登录

### 模板接口（需要认证）
- **URL**: `GET /api/rules`
- **查询参数**:
  - `name`: 模板名称（可选）
- **响应格式**:
```json
[
  {
    "rule_id": "7519892154584432668",
    "name": "技术部-工作日报",
    "icon_name": "dailyReportIcon",
    "created_at": 1750861330,
    "creator_user_id": "ou_72c6d6ac0df0f229d7e19514388c7f83",
    "creator_user_name": "杨凯",
    "owner_user_id": "ou_72c6d6ac0df0f229d7e19514388c7f83",
    "owner_user_name": "杨凯",
    "form_schema": [
      {
        "name": "今日总结",
        "type": "text"
      },
      {
        "name": "明日计划",
        "type": "text"
      },
      {
        "name": "需要协调与帮助",
        "type": "text"
      }
    ],
    "need_report_user_ids": [
      "ou_72c6d6ac0df0f229d7e19514388c7f83"
    ],
    "manager_user_ids": [
      "ou_72c6d6ac0df0f229d7e19514388c7f83"
    ]
  }
]
```

### 报告接口
- **URL**: `GET /api/reports`
- **查询参数**:
  - `rule_id`: 模板ID（可选）
  - `start_time`: 开始时间戳（可选）
  - `end_time`: 结束时间戳（可选）
- **响应格式**:
```json
{
  "items": [
    {
      "task_id": "7519892643716136988",
      "rule_name": "技术部-工作日报",
      "from_user_id": "ou_72c6d6ac0df0f229d7e19514388c7f83",
      "from_user_name": "杨凯",
      "department_name": "",
      "commit_time": 1750861444,
      "form_contents": [
        {
          "field_id": "7519892155655176193",
          "field_name": "今日总结",
          "field_value": "test1"
        },
        {
          "field_id": "7519892155675131932",
          "field_name": "明日计划",
          "field_value": "test"
        },
        {
          "field_id": "7519892155700133916",
          "field_name": "需要协调与帮助",
          "field_value": ""
        }
      ],
      "rule_id": "7519892154584432668",
      "to_user_ids": [
        "ou_72c6d6ac0df0f229d7e19514388c7f83"
      ],
      "to_user_names": [
        "杨凯"
      ]
    }
  ],
  "has_more": false
}
```

### 支持的字段类型
- `text`: 文本（前端显示为富文本编辑器）
- `number`: 数字
- `dropdown`: 单选下拉框
- `multiSelect`: 多选框
- `image`: 图片上传
- `attachment`: 附件上传
- `address`: 地址输入
- `datetime`: 日期时间选择

### 前端集成特性
1. **自动加载模板**: 页面启动时自动从 `/api/rules` 加载模板列表
2. **动态报告获取**: 支持按模板ID和时间范围过滤报告
3. **字段类型映射**: 自动将API字段类型映射为前端组件类型
4. **数据格式转换**: 自动处理API响应数据格式转换
5. **错误处理**: 完整的错误提示和加载状态

### API服务类
创建了 `FeishuApiService` 类（位于 `frontend/src/utils/aiService.js`），提供：

#### 核心方法
- `getFixedTemplateList()`: 获取固定的模板名称列表
- `getAllTemplates()`: 获取所有模板（固定列表 + API内容）
- `getTemplateContent(templateName)`: 通过API获取指定模板的内容
- `getDefaultTemplate(templateId, templateName)`: 获取默认模板结构
- `getRuleByName(templateName)`: 获取指定名称的模板（带容错处理）

#### 报告相关
- `getReports(params, templateData)`: 获取报告列表
- `getRawRuleById(ruleId)`: 根据ID获取原始模板数据

#### 工具方法
- `mapFieldType(apiType)`: 字段类型映射
- `formatFieldValue(value, type)`: 字段值格式化
- `formatTime(timestamp)`: 时间格式化

### 模板架构设计
系统采用**固定模板列表 + 动态内容获取**的架构：

#### 固定模板列表
前端固定定义以下模板名称：
1. **工作月报**
2. **技术部-工作日报**

#### 动态内容获取
- 模板的具体字段内容通过 `/api/rules?name=模板名称` 接口获取
- 如果API获取失败，系统会使用预设的默认字段结构
- 支持完全离线使用（使用默认模板结构）

#### 默认模板结构

**工作月报默认字段**：
- 本月工作总结 (富文本)
- 主要成就 (富文本)
- 遇到的挑战 (富文本)
- 下月工作计划 (富文本)
- 关键指标数据 (文本)
- 团队反馈 (富文本)

**技术部-工作日报默认字段**：
- 今日总结 (富文本)
- 明日计划 (富文本)
- 需要协调与帮助 (富文本)

## 开发说明

### 修改内容概要
1. **新增API服务**: 创建 `FeishuApiService` 类处理飞书API调用
2. **数据格式转换**: 将API响应格式转换为前端组件所需格式
3. **动态模板加载**: 页面启动时从API加载模板，替换静态数据
4. **动态报告获取**: 支持按条件过滤获取报告数据
5. **字段类型映射**: 正确映射API字段类型到前端组件类型
6. **错误处理**: 添加完整的错误提示和加载状态管理

### 开发环境配置
确保后端API服务运行在 `http://localhost:8080`，前端会自动连接此地址。

## 许可证

MIT License 