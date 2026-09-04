# 食忆

私人菜品记忆 Web 应用：记录做过的菜、用照片和标签找回、从历史成功菜里获得复做推荐。

## 技术栈

- 后端：Go + Gin + GORM + PostgreSQL
- 前端：Vue 3 + Vite + TypeScript + Pinia（包管理使用 **pnpm**）
- 部署：Docker Compose

## 本地开发

需要：Go 1.23+、Node.js 20+、pnpm 10、Docker（仅用于开发库）。完整步骤见 [docs/07-本地开发.md](docs/07-本地开发.md)。

1. 启动开发数据库（映射宿主机 5432，与生产 Compose 的 `db` 不是同一份配置）：

```bash
docker compose -f deployments/docker-compose.dev.yml up -d
```

局域网 Linux 上跑库时，把 `backend/.env` 的 `DB_HOST` 设为该服务器内网 IP。若 5432 被拒绝，在服务器上对开发 compose 执行 `up -d --force-recreate`（需映射 `0.0.0.0:5432`）。

2. 配置并启动后端：

```bash
cd backend
cp .env.example .env   # PowerShell: Copy-Item .env.example .env
go run ./cmd/server
```

3. 启动前端：

```bash
cd frontend
pnpm install
pnpm dev
```

浏览器打开 http://localhost:5173 。Vite 会把 `/api` 和 `/uploads` 代理到 `http://127.0.0.1:8080`。

## Docker 部署

```bash
docker compose -f deployments/docker-compose.yml up --build -d
```

访问 http://localhost:8080 。

公网部署请设置访问口令：

```bash
APP_ACCESS_TOKEN=your-secret docker compose -f deployments/docker-compose.yml up -d
```

前端会跳转到 `/unlock` 输入口令。

## MVP 验收清单

- 可以新增、编辑、删除菜品
- 可以上传封面图和制作照片（列表使用缩略图）
- 可以记录多次制作，并自动更新制作次数/最近制作时间/平均评分
- 可按关键词、标签、状态、是否做过/成功筛选
- 首页展示统计、最近做过、很久没做但评分高、随机旧菜、想做清单
- PC 与手机端响应式布局
- `docker compose up` 可启动完整应用

## 目录

```
backend/     Go API
frontend/    Vue 3 前端（pnpm）
deployments/ Dockerfile、生产 compose、开发库 compose
uploads/     图片文件
docs/        产品与技术方案
```
