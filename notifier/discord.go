package notifier

import (
	"fmt"
	"os"
	lib "google-flights-crawler/lib"
)

func Notify(message string) {
	host, _ := os.Hostname()
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")

	fmt.Printf("🔔 [%s] %s\n", host, message)

	if webhookURL == "" {
		fmt.Println("⚠️ DISCORD_WEBHOOK_URL não configurada")
		return
	}

	lib.SendWebhook(webhookURL, map[string]interface{}{
		"content": fmt.Sprintf("[%s] %s", host, message),
	})
}