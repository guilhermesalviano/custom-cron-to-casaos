package notifier

import (
	"fmt"
	lib "google-flights-crawler/lib"
	"log"
	"os"
)

func Notify(message string) {
	host, _ := os.Hostname()
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")

	log.Printf("🔔 [%s] %s\n", host, message)

	if webhookURL == "" {
		log.Println("⚠️ DISCORD_WEBHOOK_URL não configurada")
		return
	}

	lib.SendWebhook(webhookURL, map[string]interface{}{
		"content": fmt.Sprintf("[%s] %s", host, message),
	})
}