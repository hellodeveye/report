# 飞书报告助手 Backend

## 项目结构

```
backend/
├── main.go                     # 应用入口文件
├── go.mod                      # Go模块文件
├── go.sum                      # 依赖校验文件
├── Dockerfile                  # Docker构建文件
├── api/                        # API相关目录
│   ├── routes.go              # 路由配置
│   ├── handlers/              # HTTP处理器
│   │   └── feishu.go         # 飞书相关处理器
│   └── middleware/            # 中间件
│       └── cors.go           # CORS中间件
├── internal/                   # 内部包（不被其他项目引用）
│   ├── config/               # 配置管理
│   │   └── config.go
│   └── models/               # 数据模型
│       └── feishu.go
└── pkg/                       # 可被外部引用的包
    └── feishu/               # 飞书API客户端
        └── client.go
```

## 功能说明

### API端点

1. **健康检查**: `GET /`
   - 返回服务器运行状态

2. **飞书认证**: 
   - `GET /api/auth/feishu/login` - 飞书登录
   - `GET /api/auth/feishu/callback` - 飞书回调处理

3. **报告规则**: `GET /api/rules`
   - 查询飞书报告规则
   - 支持query参数: `rule_name`

4. **报告生成**: `POST /api/generate-draft`
   - 生成报告草稿

### 环境变量

```bash
# 飞书应用配置
FEISHU_APP_ID=your_app_id_here
FEISHU_APP_SECRET=your_app_secret_here
FEISHU_BASE_URL=https://open.feishu.cn

# 访问令牌（用于测试）
FEISHU_ACCESS_TOKEN=your_access_token_here

# 服务器端口
PORT=8080
```

## 运行方式

```bash
# 构建
go build

# 运行
go run main.go

# 或运行编译后的二进制文件
./github.com-hellodeveye-report
```

## API测试

```bash
# 健康检查
curl http://localhost:8080/

# 查询飞书规则（需要设置FEISHU_ACCESS_TOKEN环境变量）
curl "http://localhost:8080/api/rules?rule_name=工作月报"

# 生成报告草稿
curl -X POST http://localhost:8080/api/generate-draft
```

## 特性

- 🏗️ 清晰的项目结构分层
- 🔧 环境变量配置管理
- 🌐 CORS中间件支持
- 📝 飞书API客户端封装
- 🛣️ 集中化路由管理
- 📊 报告生成功能 