package bitcoinclient

import (
	"fmt"
	"log"
	"os"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/ipoluianov/bdata/logger"
	//"github.com/btcsuite/btcutil"
)

func password() string {
	filePath := logger.CurrentExePath() + "/password.txt"
	bs, err := os.ReadFile(filePath)
	if err != nil {
		return ""
	}
	return string(bs)
}

func CheckConnect() {
	// Настройка подключения к удаленному bitcoin-узлу через RPC
	connCfg := &rpcclient.ConnConfig{
		Host:         "spb.u00.io:8332",
		User:         "user",
		Pass:         password(),
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Если используется TLS, установите в false
	}

	// Создаем клиент
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		log.Fatalf("error creating new btc client: %v", err)
	}
	defer client.Shutdown()

	info, err := client.GetNetworkInfo()
	if err != nil {
		log.Fatalf("Error getting network info: %v", err)
	}

	fmt.Println(info)

	for blockIndex := int64(8066); blockIndex < 834066+1; blockIndex++ {
		//blockNumber := int64(1234)

		fmt.Println("gettings block hash", blockIndex)

		blockHash, err := client.GetBlockHash(blockIndex)
		if err != nil {
			log.Fatalf("error getting block hash: %v", err)
		}

		fmt.Println("gettings block")

		block, err := client.GetBlock(blockHash)
		if err != nil {
			log.Fatalf("error getting block: %v", err)
		}

		fmt.Printf("Block %d: %v\n", blockIndex, blockHash)
		for _, tx := range block.Transactions {
			for _, out := range tx.TxOut {

				//receiverAddress, err := getAddressFromScriptPubKey(out.PkScript, &chaincfg.MainNetParams)

				cl, addr, n, err := txscript.ExtractPkScriptAddrs(out.PkScript, &chaincfg.MainNetParams)
				fmt.Println(cl, addr, n)

				if err != nil {
					log.Println(err)
					continue
				}

				//fmt.Println(receiverAddress)
			}
		}
	}

}
