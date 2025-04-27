package Trakt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"runtime"
)

type DeviceCodes struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURL string `json:"verification_url"`
}

func DeviceCode() DeviceCodes {
	client := &http.Client{}

	body := []byte("{\n  \"client_id\": \"9ddb03b27e491a11ac4447bd5f37b1246e78f29e3df1e2ccd8ae1d802cbf466e\"\n}")

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
	URL := deviceCode.VerificationURL

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", URL).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", URL).Start()
	case "darwin":
		err = exec.Command("open", URL).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		fmt.Println("Error while opening the URL")
		fmt.Println(err)
	}

	return deviceCode
}
