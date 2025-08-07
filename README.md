# 智能报告助手

本工具是一款旨在帮助用户自动生成周报和月报草稿的Web应用，它能从钉钉中拉取并汇总您的日报。

## 🚀 项目架构

项目主要分为三个部分：

-   `/backend`: Go语言开发的后端服务，负责处理核心业务逻辑、与钉钉开放平台通信，并提供API接口。
-   `/frontend`: 基于Vue.js的单页应用，为用户提供友好的操作界面。

## 🛠️ 技术栈

-   **后端**: Go, Gorilla/Mux
-   **前端**: Vue.js, Vite, Tailwind CSS, Tiptap Editor
-   **部署**: Docker, Docker Compose

## 快速开始

### 环境准备

-   Docker 和 Docker Compose
-   一个钉钉应用 (`App Key` 和 `App Secret`)
-   Node.js 18+ 和 Go 1.18+ (用于本地开发)

### 运行方式 (推荐)

1.  **创建应用:**
    
    -   **钉钉:** 在 [钉钉开放平台](https://open-dev.dingtalk.com) 创建应用，并将重定向URI设置为 `http://localhost:5173/dingtalk/callback`。

2.  **配置环境变量:**
    
    在项目根目录创建一个 `.env` 文件，并填入您的应用凭证：

    ```env
    # 钉钉配置
    DINGTALK_APP_KEY=your_dingtalk_app_key
    DINGTALK_APP_SECRET=your_dingtalk_app_secret
    DINGTALK_REDIRECT_URI=http://localhost:5173/auth/callback
    DINGTALK_BASE_URL=https://oapi.dingtalk.com
    
    # 通用配置
    JWT_SECRET=your-jwt-secret-key-change-in-production
    FRONTEND_URL=http://localhost:5173
    ```

3.  **构建并运行:**

    ```bash
    docker-compose up --build
    ```

4.  **访问应用:**
    -   前端: `http://localhost:5173`
    -   后端 API: `http://localhost:8080`



## API 集成说明

### 认证接口
- **登录**: `GET /api/auth/dingtalk/login` - 获取OAuth登录URL
- **交换Code**: `POST /api/auth/dingtalk/exchange` - 用授权码换取JWT
- **当前用户**: `GET /api/auth/user` - 获取当前用户信息 (需认证)
- **登出**: `POST /api/auth/logout` - 退出登录

### 模板接口 (需认证)
- **URL**: `GET /api/dingtalk/templates/detail`
- **查询参数**:
  - `name`: 模板名称 (可选)

### 报告接口 (需认证)
- **URL**: `GET /api/reports`
- **查询参数**:
  - `template_name`: 模板名称 (可选)
  - `start_time`: 开始时间戳 (可选)
  - `end_time`: 结束时间戳 (可选)

## 开发说明

### 后端
后端服务位于 `backend` 目录，使用标准的Go项目结构。

### 前端
前端应用位于 `frontend` 目录，是一个基于Vite的Vue.js项目。

- **API服务**: `frontend/src/utils/apiService.js` 封装了对后端API的调用。
- **认证服务**: `frontend/src/utils/authService.js` 处理所有与用户认证相关的逻辑。
- **AI集成**: 前端直接调用 [DeepSeek](https://platform.deepseek.com/) API实现内容的智能生成与优化。

## 许可证

[MIT License](LICENSE)