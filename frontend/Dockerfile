# 使用官方 Node.js 運行時作為父鏡像
FROM node:16

# 設置工作目錄
WORKDIR /app

# 複製 package.json 和 package-lock.json（如果可用）
COPY package*.json ./

# 安裝項目依賴
RUN npm install

# 將源代碼複製到鏡像中
COPY . .

RUN npm install react-table@7.8.0 react-toastify@8.2.0
RUN npm run build

# 開放端口 3000 供 React 開發服務器使用
EXPOSE 3000

# 運行應用
CMD ["npm", "start"]