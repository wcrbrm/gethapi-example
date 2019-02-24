package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"strings"
)

func GetConnection() net.Conn {
	HOST, ok := os.LookupEnv("SERVER_HOST")
	if !ok {
		HOST = "127.0.0.1"
	}
	PORT, ok := os.LookupEnv("SERVER_PORT")
	if !ok {
		PORT = "9091"
	}
	address := HOST + ":" + PORT
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal("Dial Error:", err)
	}
	log.Println("Connected")
	return conn
}

func GetTestAccounts() map[string]string {
	FILEPATH, ok := os.LookupEnv("TEST_ACCOUNTS")
	if !ok {
		FILEPATH = "./testkeys/accounts10.txt"
	}
	file, err := os.Open(FILEPATH)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	accounts := make(map[string]string)

	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ":")
		accountName := strings.TrimSpace(parts[0])
		accountAddress := strings.TrimSpace(parts[1])
		accounts[accountName] = accountAddress
		log.Println("Loaded account alias", accountName)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return accounts
}

type SendEthRequest struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount string `json:"amount"`
	Key    string `json:"key"`
}

func GetShortcut(input string, accounts map[string]string) (*SendEthRequest, error) {
	words := strings.Fields(input)
	if len(words) != 3 {
		return nil, errors.New("Expected 3 words in the input")
	}
	pair1, exists1 := accounts[words[0]]
	if !exists1 {
		return nil, errors.New("Account 1 from should be in accounts list")
	}
	pair2, exists2 := accounts[words[1]]
	if !exists2 {
		return nil, errors.New("Account 2 to should be in accounts list")
	}
	pubpriv1 := strings.Split(pair1, ",")
	pubpriv2 := strings.Split(pair2, ",")

	from := pubpriv1[0]
	privateKey := pubpriv1[1]
	to := pubpriv2[0]

	n := new(big.Int)
	_, errVal := fmt.Sscan(words[2], n)
	if errVal != nil {
		return nil, errVal
	}
	value := n.Mul(n, big.NewInt(1e18))
	return &SendEthRequest{from, to, value.String(), privateKey}, nil
}

func main() {
	conn := GetConnection()
	accounts := GetTestAccounts()
	for {
		// read input from stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Message to send >> ")
		text, _ := reader.ReadString('\n')

		// if input string is in format of ACC#1 ACC#2 AMOUNT, then convert it
		sendEth, errShortcut := GetShortcut(text, accounts)
		if errShortcut != nil {
			// no shortcuts, send to tcp in raw format
			fmt.Fprintf(conn, text+"\n")
		} else {
			payloadSendEth, _ := json.Marshal(sendEth)
			log.Println("Sending )) ", string(payloadSendEth))
			fmt.Fprintf(conn, "SendEth"+string(payloadSendEth)+"\n")
		}

		// listen for reply
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Println("Message received: << " + message)
	}
}
