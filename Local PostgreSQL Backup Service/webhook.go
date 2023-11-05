package main

import (
	"bytes"
	"fmt"
	"html"
	"net/http"
)

func SendWebhook(url string, message string) error {
	if url == "" {
		return nil
	}

	sanitizedMessage := html.EscapeString(message)
	sanitizedMessage = fmt.Sprintf(`{"content": "%s"}`, sanitizedMessage)

	fmt.Printf(sanitizedMessage)

	jsonStr := []byte(sanitizedMessage)
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
