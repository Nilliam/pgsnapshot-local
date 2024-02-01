package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func SendWebhook(url string, message string) error {
	if url == "" {
		return nil
	}

	content := fmt.Sprintf(`{"content": "%s"}`, message)

	jsonStr := []byte(content)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	return nil
}
