# 客戶數據管理系統

這是一個基於 Go、React 和 MariaDB的應用，用於管理客戶資料和交易記錄。該項目包括自動化的 CI/CD 流程和 Kubernetes 部署配置。

## 功能特點

- 客戶數據的 CRUD 操作
- 交易記錄查詢和顯示
- 響應式前端設計
- Docker 容器化部署
- GitHub Actions 自動化 CI/CD
- Kubernetes 部署支持

## 技術棧

- 後端：Go (Gin 框架)
- 前端：React
- 資料庫：MariaDB
- 容器化：Docker & Docker Compose
- CI/CD：GitHub Actions
- 容器編排：Kubernetes

## 項目結構

```
.
├── backend/           # Go 後端代碼
├── frontend/          # React 前端代碼
├── k8s/               # Kubernetes 配置文件
├── .github/workflows/ # GitHub Actions 工作流配置
├── docker-compose.yml # Docker Compose 配置
└── README.md          # 項目說明文件
```

## 如何運行

1. clone倉庫：
   ```
   git clone https://github.com/Joyen09/devops_pre_test.git
   cd devops_pre_test
   ```

2. 使用 Docker Compose 啟動應用：
   ```
   docker-compose up -d
   ```

3. 訪問應用：
   - 前端：http://localhost:3000
   - 後端 API：http://localhost:8080
   - phpMyAdmin：http://localhost:9000

## CI/CD 流程

本項目使用 GitHub Actions 進行持續集成和部署。每次推送到 main 分支時，都會觸發以下流程：

1. 構建和測試後端和前端代碼
2. 構建 Docker 鏡像並推送到 Docker Hub
3. 部署到 Kubernetes 集群（使用 KinD 進行測試部署）

## Kubernetes 部署

Kubernetes 配置文件位於 `k8s/` 目錄中。包括：

- 後端部署
- 前端部署
- MariaDB 部署
- 服務配置
- 密鑰配置
