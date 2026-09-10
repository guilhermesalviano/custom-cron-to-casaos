package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func Notify(message string) {
	host, _ := os.Hostname()
	log.Printf("🔔 [%s] %s", host, message)

	webhookURL := strings.TrimSpace(os.Getenv("DISCORD_WEBHOOK_URL"))
	if webhookURL != "" {
		if err := SendWebhook(webhookURL, map[string]string{
			"content": fmt.Sprintf("[%s] %s", host, message),
		}); err != nil {
			log.Printf("Could not send Discord notification: %v", err)
		}
	}

	topic := strings.TrimSpace(os.Getenv("NTFY_TOPIC"))
	if topic != "" {
		if err := SendNtfy(topic, message); err != nil {
			log.Printf("Could not send NTFY notification: %v", err)
		}
	}
}

func SendWebhook(url string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal webhook payload: %w", err)
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("create webhook request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("send webhook: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("webhook returned status %d", response.StatusCode)
	}
	return nil
}

func SendNtfy(topic, message string) error {
	request, err := http.NewRequest(http.MethodPost, "https://ntfy.sh/"+topic, strings.NewReader(message))
	if err != nil {
		return fmt.Errorf("create NTFY request: %w", err)
	}
	request.Header.Set("Title", "Wishlist service")
	request.Header.Set("Priority", "default")

	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("send NTFY notification: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("NTFY returned status %d", response.StatusCode)
	}
	return nil
}
