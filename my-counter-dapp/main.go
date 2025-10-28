package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/gaojianshuai/mydapp_backend1/tree/main/my-counter-dapp/counter"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// 配置参数
	infuraURL := "https://sepolia.infura.io/v3/ad8d28671f5a47af902f71b5cf53d405"        // 替换为你的 Infura 项目 ID
	privateKeyHex := "9d21bb53117a73a76009ebe640e46fc023336de073e0085b8d4e33803c1172e0" // 替换为你的私钥（不带 0x 前缀）

	// 1. 连接到 Sepolia 测试网
	fmt.Println("连接到 Sepolia 测试网...")
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatal("连接失败: ", err)
	}
	defer client.Close()

	// 获取网络信息
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal("获取链ID失败: ", err)
	}
	fmt.Printf("连接成功! 链ID: %s\n", chainID.String())

	// 2. 配置账户和交易参数
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal("私钥解析失败: ", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("公钥类型断言失败")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Printf("账户地址: %s\n", fromAddress.Hex())

	// 获取账户余额
	balance, err := client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		log.Fatal("获取余额失败: ", err)
	}
	fmt.Printf("账户余额: %s ETH\n", weiToEther(balance))

	// 3. 部署合约
	fmt.Println("\n开始部署合约...")
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal("创建交易授权失败: ", err)
	}

	// 设置交易参数
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal("获取 nonce 失败: ", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasLimit = uint64(300000)

	// 获取推荐的 Gas 价格
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal("获取 Gas 价格失败: ", err)
	}
	auth.GasPrice = gasPrice
	fmt.Printf("Gas 价格: %s wei\n", gasPrice.String())

	// 设置初始计数器值为 10
	initialCount := big.NewInt(10)

	// 部署合约
	contractAddress, tx, instance, err := counter.DeployCounter(auth, client, initialCount)
	if err != nil {
		log.Fatal("合约部署失败: ", err)
	}

	fmt.Printf("部署交易已发送: %s\n", tx.Hash().Hex())
	fmt.Printf("合约地址: %s\n", contractAddress.Hex())

	// 等待交易确认
	fmt.Println("等待交易确认...")
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatal("等待交易确认失败: ", err)
	}

	if receipt.Status == 1 {
		fmt.Println("✅ 合约部署成功!")
	} else {
		log.Fatal("❌ 合约部署失败")
	}

	// 4. 与合约交互
	fmt.Println("\n开始与合约交互...")

	// 读取当前计数器值
	currentValue, err := instance.GetCount(&bind.CallOpts{})
	if err != nil {
		log.Fatal("读取计数器值失败: ", err)
	}
	fmt.Printf("当前计数器值: %d\n", currentValue)

	// 增加计数器
	fmt.Println("调用 increment 方法...")

	// 更新交易参数
	auth.Nonce = big.NewInt(int64(nonce + 1))
	auth.GasLimit = uint64(100000)

	incrementTx, err := instance.Increment(auth)
	if err != nil {
		log.Fatal("调用 increment 失败: ", err)
	}
	fmt.Printf("Increment 交易: %s\n", incrementTx.Hash().Hex())

	// 等待 increment 交易确认
	fmt.Println("等待 increment 交易确认...")
	_, err = bind.WaitMined(context.Background(), client, incrementTx)
	if err != nil {
		log.Fatal("等待 increment 交易确认失败: ", err)
	}

	// 再次读取计数器值
	time.Sleep(5 * time.Second) // 给节点一些时间更新状态
	newValue, err := instance.GetCount(&bind.CallOpts{})
	if err != nil {
		log.Fatal("读取新计数器值失败: ", err)
	}
	fmt.Printf("增加后的计数器值: %d\n", newValue)

	// 5. 验证操作成功
	if newValue.Cmp(big.NewInt(0).Add(currentValue, big.NewInt(1))) == 0 {
		fmt.Println("✅ 计数器增加操作成功!")
	} else {
		fmt.Printf("❌ 计数器增加操作失败: 期望 %d, 得到 %d\n",
			big.NewInt(0).Add(currentValue, big.NewInt(1)), newValue)
	}

	fmt.Println("\n🎉 所有操作完成!")
}

// 辅助函数：将 wei 转换为 ETH
func weiToEther(wei *big.Int) string {
	ether := new(big.Float).SetInt(wei)
	ether = ether.Quo(ether, big.NewFloat(1e18))
	return ether.Text('f', 6)
}
