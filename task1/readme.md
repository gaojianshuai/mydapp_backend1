# 🚀 项目简介
这个项目演示了如何：

🔍 查询区块链信息 - 获取指定区块的详细信息

💸 发送交易 - 在 Sepolia 测试网络上进行 ETH 转账

🔐 安全交互 - 使用私钥签名交易并发送到网络

# 初始化 Go 模块
go mod init blockchain-interaction

# 安装必要的依赖
go get github.com/ethereum/go-ethereum
go get github.com/joho/godotenv

# 创建项目文件夹
mkdir blockchain-interaction
cd blockchain-interaction

# 初始化 Go 模块
go mod init blockchain-interaction



blockchain-interaction/
├── go.mod                 # Go 模块定义
├── go.sum                # 依赖校验和
├── .env                  # 环境变量（不要提交）
├── query_block.go        # 区块查询代码
├── send_transaction.go   # 交易发送代码
└── README.md            # 项目说明

# 获取 Infura API Key步骤：获取 Infura API Key
注册 Infura 账户

访问 Infura 官网

点击 "Get Started for Free" 注册账户

完成邮箱验证

创建新项目

登录后进入仪表板

点击 "Create New Project"

项目名称输入 "Sepolia-Test"

选择 "Ethereum" 作为产品

获取 API Key

在项目设置中，切换到 "Sepolia" 网络

复制 HTTPS 端点 URL

格式：https://sepolia.infura.io/v3/YOUR-PROJECT-ID

# 准备工作
注册 Infura 账户并创建项目，获取 Project ID

准备 Sepolia 测试网络的以太坊账户和私钥

从 Sepolia 水龙头获取测试 ETH

# 运行查询区块
go run query_block.go

# 运行发送交易
go run send_transaction.go


# 预期输出
=== 查询最新区块 ===
区块哈希: 0x4e3a3754410177e6937ef1f84bba68ea139e8d1a2258c5f85db9f1cd715a1bdd
区块号: 5687735
时间戳: 1698765432
交易数量: 124
难度: 12555445788935152
...

# 安全注意事项
私钥安全
✅ 正确做法：使用环境变量 .env 文件

✅ 正确做法：使用密钥管理服务

❌ 错误做法：硬编码在源代码中

❌ 错误做法：提交到版本控制系统

# 测试网络最佳实践
仅使用测试资金进行开发

在主网使用前充分测试

监控 Gas 价格避免过高费用

定期检查账户余额