package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/LampardNguyen234/go-incognito/incclient"
)

/*func (server *RPCServer) GetCurrentHeight() ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}

	method := getActiveShards

	params := make([]interface{}, 0)
	request := rpchandler.CreateJsonRequest("1.0", method, params, 1)

	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return server.SendPostRequestWithQuery(string(query))
}

func calculateIncreaseHeight(startHeight int64) int64 {

}*/

func fileToSlice(fileName string) []string {
	fileBytes, _ := ioutil.ReadFile(fileName)
	mySlice := strings.Split(string(fileBytes), "\n")
	return mySlice
}

func Init() ([]string, []string) {
	senderFile := "clientPrivateKey"
	receiverFile := "clientPaymentAddr"
	sendersPrivKey := fileToSlice(senderFile)
	receiversPaymentAddr := fileToSlice(receiverFile)
	return sendersPrivKey, receiversPaymentAddr
}

// dao gui tien cho 400 tk truoc khi bat dau giao dich
func initStresstest(receiversPaymentAddr []string) {
	DAOprivateKey := "112t8rnX8nLvrBmWg8hKQsYfRkf1oW7mKEjibM8PB8bBeKFnpPAbgRZNgwnAQBCFswGdZejnarPGKgPvHtjt1wfChEgAWojE2QeSK1hByQaJ"
	//total_receiver := len(receiversPaymentAddr)
	max_send_each := 200
	amount := 99000000000
	for i := 1; i < 400; i += max_send_each {
		receiverList := receiversPaymentAddr[i : i+max_send_each]
		txHash := sendTransaction(DAOprivateKey, receiverList, uint(max_send_each), uint64(amount))
		fmt.Println("Hash is: ", txHash)
		if i == 1 {
			time.Sleep(20 * time.Second)
		}
	}
}

func sendTransaction(sendersPrivKey string, receiversPaymentAddr []string, nRecv uint, amount uint64) string {

	if nRecv > 200 {
		return ""
	}

	client, err := incclient.NewTestNetClient()
	if err != nil {
		log.Fatal(err)
	}
	//version
	txVersion := int8(1)

	// chon danh sach nguoi nhan
	index := 0 //rand.Intn(800)
	receiverList := receiversPaymentAddr[index : uint(index)+nRecv]

	// tao slice amount
	var amountList []uint64
	for i := 0; uint(i) < nRecv; i++ {
		amountList = append(amountList, uint64(amount))
	}

	// gui txs
	txParam := incclient.NewTxParam(sendersPrivKey, receiverList, amountList, 0, nil, nil, nil)

	encodedTx, txHash, err := client.CreateRawTransaction(txParam, txVersion)
	if err != nil {
		log.Fatal(err)
	}
	err = client.SendRawTx(encodedTx)
	if err != nil {
		log.Fatal(err)

	}

	return txHash
}

func msendTransaction() string {
	client, err := incclient.NewTestNetClient()
	if err != nil {
		log.Fatal(err)
	}
	txVersion := int8(1)
	privateKey := "112t8rnX8nLvrBmWg8hKQsYfRkf1oW7mKEjibM8PB8bBeKFnpPAbgRZNgwnAQBCFswGdZejnarPGKgPvHtjt1wfChEgAWojE2QeSK1hByQaJ"
	fileName := "clientPaymentAddr"
	fileBytes, _ := ioutil.ReadFile(fileName)
	paymentAddrList := strings.Split(string(fileBytes), "\n")

	receiverList := paymentAddrList[0:]

	numRecv := len(receiverList)
	amount := 10000000
	var amountList []uint64
	for i := 0; i < numRecv; i++ {
		amountList = append(amountList, uint64(amount))
	}
	txParam := incclient.NewTxParam(privateKey, receiverList, amountList, 0, nil, nil, nil)

	encodedTx, txHash, err := client.CreateRawTransaction(txParam, txVersion)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendRawTx(encodedTx)
	if err != nil {
		log.Fatal(err)

	}

	//fmt.Printf("Create and send tx successfull, txhash: %v\n", txHash)
	return txHash
}

func main() {

	sendersPrivKey, receiversPaymentAddr := Init()

	// DAO gui tien cho 400 account khac
	initStresstest(receiversPaymentAddr)
	fmt.Println("DAO has send money to 400 account for initialization")

	amount := 1000
	nRecv := 100
	// gui trong sec_total giay, moi sec_interval giay gui 1 lan
	sec_total := 20
	sec_interval := 1
	ticker := time.Tick(time.Duration(sec_interval) * time.Second)
	start := time.Now()
	sender_index := 400
	for stay, timeout := true, time.After(time.Duration(sec_total)*time.Second); stay; {

		select {
		case <-timeout:
			elapsed := time.Since(start)
			fmt.Println("Elapse: ", elapsed)
			stay = false
			//calculate number of block here
		case <-ticker:

			txHash := sendTransaction(sendersPrivKey[sender_index], receiversPaymentAddr, uint(nRecv), uint64(amount))
			fmt.Printf("Create and send tx successfull, txhash: %v\n", txHash)
			if sender_index == 0 {
				sender_index = 400
			}
			sender_index--
		}
	}

}
