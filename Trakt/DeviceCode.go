package Trakt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type DeviceCodes struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURL string `json:"verification_url"`
}

func DeviceCode() DeviceCodes {
	client := &http.Client{}
	clientID := os.Getenv("TRAKT_CLIENT_ID")

	bodyStruct := struct {
		ClientID string `json:"client_id"`
	}{
		ClientID: clientID,
	}

	bodyBytes, err := json.Marshal(bodyStruct)

	if err != nil {
		fmt.Println("Error marshalling JSON")
		fmt.Println(err)
	}

	body := []byte(bodyBytes)

	req, _ := http.NewRequest("POST", "https://api.trakt.tv/oauth/device/code", bytes.NewBuffer(body))

	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Errored when sending request to the server")
	}

	defer resp.Body.Close()
	resp_body, _ := io.ReadAll(resp.Body)

	var deviceCode DeviceCodes
	err = json.Unmarshal(resp_body, &deviceCode)
	if err != nil {
		fmt.Println("Error unmarshalling the JSON values")
		fmt.Println(err)
	}

	fmt.Printf("Opening the default web browser, please input the code: %s\n", deviceCode.UserCode)

	return deviceCode
}
