package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var config struct {
	Token     string `json:"token"`
	UID       string `json:"uid"`
	FirstTime int    `json:"first_time"`
}

var response struct {
	Status        int    `json:"status"`
	StatusMessage string `json:"status_message"`
	MessageToken  int64  `json:"message_token"`
	ChatHostname  string `json:"chat_hostname"`
}

func setWebhook() {
	payload := map[string]interface{}{
		"auth_token":  config.Token,
		"url":         "https://proton.me/",
		"event_types": []string{"failed"},
		"send_name":   false,
		"send_photo":  false,
	}
	payloadBytes, _ := json.Marshal(payload)
	resp, err := http.Post("https://chatapi.viber.com/pa/set_webhook", "application/json", bytes.NewReader(payloadBytes))
	if err != nil {
		fmt.Println("Error setting webhook:", err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func getAccountInfo() string {
	payload := map[string]interface{}{
		"auth_token": config.Token,
	}
	payloadBytes, _ := json.Marshal(payload)
	resp, err := http.Post("https://chatapi.viber.com/pa/get_account_info", "application/json", bytes.NewReader(payloadBytes))
	if err != nil {
		fmt.Println("Error getting account info:", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		fmt.Println(result)
		members := result["members"].([]interface{})
		if len(members) > 0 {
			return members[0].(map[string]interface{})["id"].(string)
		}
	}
	return ""
}

func sendContactMessage(phoneNumber string) {
	if config.UID == "" {
		fmt.Println("Error in uid: UID is empty")
		return
	}

	payload := map[string]interface{}{
		"auth_token": config.Token,
		"from":       config.UID,
		"type":       "contact",
		"contact": map[string]interface{}{
			"name":         "a",
			"phone_number": phoneNumber,
		},
	}
	payloadBytes, _ := json.Marshal(payload)
	resp, err := http.Post("https://chatapi.viber.com/pa/post", "application/json", bytes.NewReader(payloadBytes))
	if err != nil {
		fmt.Println("Error sending contact message:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))

	// Десериализуем данные JSON в структуру.
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		fmt.Println("JSON Unmarshal error:", err)
		return
	}

	if resp.StatusCode == 200 {
		fmt.Printf("Sent message for - %s response: %s\n", phoneNumber, response.StatusMessage)
	}
}

func sendBulkMessage(fname string) {
	fileContent, err := ioutil.ReadFile(fname)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	phones := bytes.Split(fileContent, []byte("\n"))
	for _, phone := range phones {
		sendContactMessage(string(phone))
	}
}

func sendCSVBulkMessage(fname string) {
	//fileContent, err := ioutil.ReadFile(fname)
	//if err != nil {
	//	fmt.Println("Error reading file:", err)
	//	return
	//}
	//phones := bytes.Split(fileContent, []byte("\n"))
	//for _, phone := range phones {
	//	sendContactMessage(string(phone))
	//}
	fmt.Printf("Reading file: %s\n", fname)
	file, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 5
	//	reader.Comment = '#'
	//	reader.Comma = '"'

	for {
		record, e := reader.Read()
		if e != nil {
			fmt.Println(e)
			break
		}
		phone := "+7" + record[0]
		sendContactMessage(phone)
	}
}

func loadConfig() {
	fileContent, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}
	json.Unmarshal(fileContent, &config)

	if config.Token == "" {
		fmt.Println("Please specify your TOKEN")
		return
	}

	if config.FirstTime == 1 {
		setWebhook()
		config.FirstTime = 0
		newConfigBytes, _ := json.Marshal(config)
		ioutil.WriteFile("config.json", newConfigBytes, 0644)
	}

	if config.UID == "" {
		uid := getAccountInfo()
		fmt.Println(uid)
		if uid == "" {
			fmt.Println("Something went wrong: UID is returned as None")
			return
		}
		config.UID = uid
		newConfigBytes, _ := json.Marshal(config)
		ioutil.WriteFile("config.json", newConfigBytes, 0644)
	}
}

func main() {
	phonePtr := flag.String("phone", "", "phone number +7...")
	listPtr := flag.String("list", "", "list with phone numbers")
	csvPtr := flag.String("csv", "", "csv list with phone numbers and other data")
	flag.Parse()

	loadConfig()

	if (*phonePtr == "" && *listPtr == "" && *csvPtr == "") || (*phonePtr != "" && *listPtr != "" && *csvPtr != "") {
		flag.PrintDefaults()
	} else {
		if *phonePtr != "" {
			sendContactMessage(*phonePtr)
			fmt.Println("########## All messages sent ###########")
		} else if *csvPtr != "" {
			sendCSVBulkMessage(*csvPtr)
			fmt.Println("########## All messages sent ###########")
		} else {
			sendBulkMessage(*listPtr)
			fmt.Println("########## All messages sent ###########")
		}
	}
}
