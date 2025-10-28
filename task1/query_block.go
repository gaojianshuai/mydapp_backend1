package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// 连接到 Sepolia 测试网络
	// 替换为你的 Infura 项目 ID
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/ad8d28671f5a47af902f71b5cf53d405")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// 查询最新区块
	queryLatestBlock(client)

	// 查询指定区块（例如区块号 5000000）
	blockNumber := big.NewInt(5000000)
	queryBlockByNumber(client, blockNumber)
}

func queryLatestBlock(client *ethclient.Client) {
	fmt.Println("=== 查询最新区块 ===")

	// 获取最新区块号
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// 获取完整区块信息
	block, err := client.BlockByNumber(context.Background(), header.Number)
	if err != nil {
		log.Fatal(err)
	}

	printBlockInfo(block)
}

func queryBlockByNumber(client *ethclient.Client, blockNumber *big.Int) {
	fmt.Printf("\n=== 查询区块 #%d ===\n", blockNumber)

	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Printf("查询区块错误: %v", err)
		return
	}

	printBlockInfo(block)
}

func printBlockInfo(block *types.Block) {
	fmt.Printf("区块哈希: %s\n", block.Hash().Hex())
	fmt.Printf("区块号: %d\n", block.Number().Uint64())
	fmt.Printf("时间戳: %d\n", block.Time())
	fmt.Printf("交易数量: %d\n", len(block.Transactions()))
	fmt.Printf("难度: %d\n", block.Difficulty().Uint64())
	fmt.Printf("Nonce: %d\n", block.Nonce())
	fmt.Printf("Gas Limit: %d\n", block.GasLimit())
	fmt.Printf("Gas Used: %d\n", block.GasUsed())
	fmt.Printf("父区块哈希: %s\n", block.ParentHash().Hex())

	if len(block.Transactions()) > 0 {
		fmt.Println("\n前 3 笔交易哈希:")
		for i, tx := range block.Transactions() {
			if i >= 3 {
				break
			}
			fmt.Printf("  %d: %s\n", i+1, tx.Hash().Hex())
		}
	}
}

// === 查询最新区块 ===
// 区块哈希: 0x72069f64da06d0e29ec217f94d6fc5b0147c068cebdd2296edea2bf67d69d7e1
// 区块号: 9506872
// 时间戳: 1761633912
// 交易数量: 146
// 难度: 0
// Nonce: 0
// Gas Limit: 60000000
// Gas Used: 23661903
// 父区块哈希: 0x50d4c1394ee69e90179bebf4d37aa27a3c7e603b4a737cdb85be7c807f8d2314

// 前 3 笔交易哈希:
//   1: 0x4cf319a7610c70ae23690bda83d1a93e840018f238d1bf2ce0727d16bf56bb1a
//   2: 0xa7cc3573c7b52ecb77ba180cb98f016d172512a30048ece93ca945d00f61cca3
//   3: 0x6a8d50a7f99e6e0ae1fc76c7045b132b9fd7bcd25bb040041abb24f548072d81

// === 查询区块 #5000000 ===
// 区块哈希: 0x72247ea9191db039158939f7ba958e638a32df4f61a43edcae60cb7a686a2d55
// 区块号: 5000000
// 时间戳: 1704111828
// 交易数量: 81
// 难度: 0
// Nonce: 0
// Gas Limit: 30000000

// 前 3 笔交易哈希:
//   1: 0xe5a58e3d4460acf272945bc15bd91bad5b626b3f8ef33a36630cbce97d86b172
//   2: 0x2fd214d899e3b8e3ac957ab69927461f0c0a3afe9fc3b79a73a71a3f18c97e35
//   3: 0xab13feef2d739e8276da60663b9a58f3886eca8f7873093f9f31015b6c7a404f
