package wishlist

import (
	"database/sql"
	"log"
	"time"

	"github.com/guilhermesalviano/casaos-cron/services/wishlist/internal/domain"
	"github.com/guilhermesalviano/casaos-cron/services/wishlist/internal/notify"
	"github.com/guilhermesalviano/casaos-cron/services/wishlist/internal/store"
)

func CrawlAndStore() {
	notify.Notify("Starting crawling Amazon wishlist...")

	results, err := ScrapeAmazonWishlist()
	if err != nil {
		log.Printf("Failed to scrape Amazon wishlist: %v", err)
		return
	}
	log.Printf("Amazon wishlist scraped successfully: %d items found", len(results))

	database, err := store.OpenFromEnvironment()
	if err != nil {
		log.Printf("Could not connect to database: %v", err)
		return
	}
	defer database.Close()

	for index := range results {
		if err := saveAmazonWishlistPrice(database, &results[index]); err != nil {
			log.Printf("Could not save wishlist item %d: %v", index, err)
			continue
		}
		log.Printf("Wishlist item %d saved to database", index)
	}
}

func saveAmazonWishlistPrice(database *sql.DB, item *domain.WishlistItem) error {
	_, err := database.Exec(
		"INSERT INTO wishlist_amazon (title, price, link, search_date) VALUES (?, ?, ?, ?)",
		item.Title,
		item.Price,
		item.Link,
		time.Now().Format("2006-01-02 15:04:05"),
	)
	return err
}
