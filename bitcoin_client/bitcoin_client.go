package bitcoinclient

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

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

func GetInstruments() []string {
	return []string{"TRANSACTIONS_COUNT"}
}

func LoadNext() {
	minDay := int64(19720)
	maxDay := time.Now().Unix()/86400 - 1
	tickers := GetInstruments()
	for _, t := range tickers {
		for day := minDay; day < maxDay; day++ {
			if !HasData(t, day) {
				tm := TimeByDayIndex(day)
				LoadData(t, tm, GetFileNameByDate(t, tm))
				return
			}
		}
	}
}

func ParseDate(value string) time.Time {
	t, _ := time.Parse("2006-01-02", value)
	return t
}

func TimeByDayIndex(dayIndex int64) time.Time {
	unixTime := dayIndex * 86400
	return time.Unix(int64(unixTime), 0)
}

func HasData(symbol string, dayIndex int64) bool {
	filePath := GetFileNameByDate(symbol, time.Unix(dayIndex*86400, 0))
	_, err := os.Lstat(filePath)
	return err == nil
}

func GetFileNameByDate(symbol string, t time.Time) string {
	p := logger.CurrentExePath() + "/data/btc/" + symbol + "/"
	filename := t.Format("2006-01-02") + ".txt"
	return p + filename
}

func LoadData(symbol string, date time.Time, fileName string) {
}

func CheckConnect() {
	// Настройка подключения к удаленному bitcoin-узлу через RPC
	connCfg := &rpcclient.ConnConfig{
		Host: "spb.u00.io:8332",
		//Host:         "localhost:8332",
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

	/*info, err := client.GetNetworkInfo()
	if err != nil {
		log.Fatalf("Error getting network info: %v", err)
	}*/

	//fmt.Println(info)

	for blockIndex := int64(400000); blockIndex < 834000; blockIndex++ {
		addrs := make(map[string]struct{})

		//blockNumber := int64(1234)

		//fmt.Println("\r\n-------------------")
		//fmt.Println("Block", blockIndex)

		blockHash, err := client.GetBlockHash(blockIndex)
		if err != nil {
			log.Fatalf("error getting block hash: %v", err)
		}

		//fmt.Println("gettings block")

		//header, err := client.GetBlockHeader(blockHash)
		//fmt.Println("BL:", blockIndex, "DT:", header.Timestamp.Format("2006-01-02 15-04-05"))

		block, err := client.GetBlock(blockHash)
		if err != nil {
			log.Fatalf("error getting block: %v", err)
		}

		//fmt.Println("BL:", blockIndex, "TRs:", len(block.Transactions), "DT:", block.Header.Timestamp.Format("2006-01-02 15-04-05"))
		//continue

		//fmt.Println("Block ", blockIndex, blockHash, block.Header.Timestamp)
		for _, tx := range block.Transactions {
			//fmt.Println("Transaction:", tx.TxHash())
			/*fmt.Println("input:")
			for _, inp := range tx.TxIn {
				fmt.Println("\t", inp.PreviousOutPoint)
			}*/
			//fmt.Println("output:")
			for _, out := range tx.TxOut {

				//receiverAddress, err := getAddressFromScriptPubKey(out.PkScript, &chaincfg.MainNetParams)

				_, addr, _, err := txscript.ExtractPkScriptAddrs(out.PkScript, &chaincfg.MainNetParams)
				for _, a := range addr {
					addrs[a.EncodeAddress()] = struct{}{}
					//fmt.Println("\t", a.EncodeAddress(), out.Value)
				}

				if err != nil {
					log.Println(err)
					continue
				}

				//fmt.Println(receiverAddress)
			}
		}

		fmt.Println("BL:", blockIndex, "Unique Addresses", len(addrs))

		addrsArray := make([]string, 0, len(addrs))
		for k := range addrs {
			addrsArray = append(addrsArray, k)
		}
		sort.Slice(addrsArray, func(i, j int) bool {
			return addrsArray[i] < addrsArray[j]
		})
		result := make([]byte, 0, 100000000)
		for i := 0; i < len(addrsArray); i++ {
			result = append(result, []byte(addrsArray[i])...)
			result = append(result, []byte("\r\n")...)
		}

		os.MkdirAll("data", 0777)

		//filePath := "data/" + block.Header.Timestamp.Format("2006.txt")
		filePath := "data/addrs.txt"

		f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			fmt.Println("File open error:", err)
		}
		f.Write(result)
		f.Close()
	}

}

func OptimizeFiles() {
	fmt.Println("Optimizing files")
	files, err := os.ReadDir("data/parts")
	if err != nil {
		fmt.Println("")
	}

	for _, fileInfo := range files {
		resultFile, err := os.OpenFile("data/result.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println("Open for write file error:", err)
		}
		addrsMap := make(map[string]struct{})
		fmt.Println(fileInfo.Name())
		content, err := os.ReadFile("data/parts/" + fileInfo.Name())
		if err != nil {
			fmt.Println("Read file error:", err)
		}
		fmt.Println("parsing ...")
		addrs := strings.FieldsFunc(string(content), func(r rune) bool {
			return r == '\r' || r == '\n'
		})
		fmt.Println("addresses:", len(addrs))
		for _, a := range addrs {
			addrsMap[a] = struct{}{}
		}
		addrList := make([]string, 0)
		for key := range addrsMap {
			addrList = append(addrList, key)
		}
		fmt.Println("map size:", len(addrList))
		fmt.Print("writing ... ")
		result := make([]byte, 0)
		for _, a := range addrList {
			result = append(result, []byte(a)...)
			result = append(result, []byte("\r\n")...)
		}
		resultFile.Write(result)
		fmt.Println("ok")
		resultFile.Close()
	}
}

func ParseRawFile() {
	filePath := "data/addrs.txt"
	os.MkdirAll("data/parts", 0777)
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0777)
	if err != nil {
		fmt.Println("File open error:", err)
	}

	buffer := make([]byte, 100*1024*1024)
	inputData := make([]byte, 0)

	//result := make(map[string]struct{})
	result := make([]string, 0)
	currentAddr := make([]byte, 0, 64)

	addItem := func() {
		if len(currentAddr) > 0 {
			//result[string(currentAddr)] = struct{}{}
			result = append(result, string(currentAddr))
			currentAddr = make([]byte, 0, 64)
		}
	}

	processedBytes := 0
	for {
		if len(inputData) < 200 {
			n, err := f.Read(buffer)
			if err != nil {
				fmt.Println("ERROR READING", err)
				break
			}
			processedBytes += n
			inputData = append(inputData, buffer[:n]...)
		}

		for i := 0; i < len(inputData)-100; i++ {
			b := inputData[0]
			inputData = inputData[1:]
			if b == 10 || b == 13 {
				addItem()
			}
			if b > 32 {
				currentAddr = append(currentAddr, b)
			}
		}

		fmt.Println("RESULT:", len(result), "processed MB:", processedBytes/1000000)

		sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
		currentFile := ""
		fmt.Println("writing count:", len(result))

		var partFile *os.File

		for _, a := range result {
			count := 2
			if strings.HasPrefix(a, "bc") {
				count = 5
			}

			prefix := a[:count]

			fileName := "data/parts/" + prefix + ".txt"
			if fileName != currentFile {
				currentFile = fileName
				if partFile != nil {
					partFile.Close()
				}
				//fmt.Println("Open file:", currentFile)
				partFile, err = os.OpenFile(currentFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
				if err != nil {
					fmt.Println("Open file error:", err)
					return
				}
			}

			partFile.Write([]byte(a + "\r\n"))
		}

		if partFile != nil {
			partFile.Close()
		}

		result = nil
	}

	f.Close()
}
