# 贡献指南

感谢每位参与贡献的开发者！

## 快速开发

### 1. 环境准备

- [Go](https://go.dev/) (版本 >= 1.24)
- [mise](https://mise.jdx.dev/getting-started.html)
- [Node.js](https://nodejs.org/) (版本 >= 24)
- [Yarn](https://yarnpkg.com/)

### 2. 克隆项目

```bash
git clone https://github.com/shichen437/Gowlive.git
cd Gowlive
```

### 3. 后端启动

```bash
# 复制配置文件
cp manifest/config/config.yaml.example manifest/config/config.yaml

# 统一开发环境（需要安装 mise,可选）
mise trust
mise install

# 安装 Go 依赖
go mod tidy

# 启动后端服务 (使用 GoFrame CLI 工具可实现热更新)
# 首次运行需要安装 CLI
mise run cli.install

# 启动！
go run main.go # 热更新使用：gf run main.go
```

### 4. 前端启动

```bash
# 进入前端目录
cd web

# 统一开发环境（需要安装 mise,可选）
mise trust
mise install

# 安装 Node.js 依赖
npm install

# 启动前端开发服务
npm run dev # 或 yarn dev
```

## PR 规范

- 提交规范：[约定式提交](https://www.conventionalcommits.org/zh-hans/v1.0.0/)
- 提交数：同一 PR 建议不超过 2 个提交，超过请 squash
- 变更范围：尽量保持 PR 专注于单一主题；避免“大杂烩”式改动
- 必要说明：列出主要变更点、兼容性影响、迁移步骤
- 性能与安全：无明显性能回退与安全隐患
