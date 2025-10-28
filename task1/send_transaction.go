package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 连接到 Sepolia 测试网络
	infuraURL := os.Getenv("INFURA_URL")
	if infuraURL == "" {
		log.Fatal("INFURA_URL not set in .env file")
	}

	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// 发送交易
	txHash, err := sendTransaction(client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("交易发送成功！交易哈希: %s\n", txHash)

	// // 等待交易确认
	// waitForTransaction(client, txHash)

	// 查询接收方的余额
	getBalance(client, common.HexToAddress(os.Getenv("RECIPIENT_ADDRESS")))

}

func sendTransaction(client *ethclient.Client) (string, error) {
	// 从环境变量获取私钥
	privateKeyStr := os.Getenv("PRIVATE_KEY")
	if privateKeyStr == "" {
		return "", fmt.Errorf("PRIVATE_KEY not set in .env file")
	}

	// 解析私钥
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return "", fmt.Errorf("解析私钥失败: %v", err)
	}

	// 从私钥获取公钥和地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("无法获取公钥")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 获取接收方地址
	recipientStr := os.Getenv("RECIPIENT_ADDRESS")
	if recipientStr == "" {
		return "", fmt.Errorf("RECIPIENT_ADDRESS not set in .env file")
	}
	toAddress := common.HexToAddress(recipientStr)

	// 获取 nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", fmt.Errorf("获取 nonce 失败: %v", err)
	}

	// 获取当前 gas 价格
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", fmt.Errorf("获取 gas 价格失败: %v", err)
	}

	// 获取链 ID
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return "", fmt.Errorf("获取链 ID 失败: %v", err)
	}

	// 设置转账金额 (0.001 ETH)
	value := big.NewInt(1000000000000000) // 0.001 ETH in wei

	// 设置 gas limit
	gasLimit := uint64(21000) // 标准转账的 gas limit

	// 创建交易
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	// 签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", fmt.Errorf("签名交易失败: %v", err)
	}

	// 发送交易
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", fmt.Errorf("发送交易失败: %v", err)
	}

	return signedTx.Hash().Hex(), nil
}

// 辅助函数：获取账户余额
func getBalance(client *ethclient.Client, address common.Address) {
	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		log.Printf("获取余额失败: %v", err)
		return
	}

	// 转换为 ETH
	ethValue := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18))
	fmt.Printf("地址 %s 的余额: %s ETH\n", address.Hex(), ethValue.String())
}
