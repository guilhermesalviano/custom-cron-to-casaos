package wishlist

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"

	"github.com/guilhermesalviano/casaos-cron/services/wishlist/internal/domain"
)

func ScrapeAmazonWishlist() ([]domain.WishlistItem, error) {
	var items []domain.WishlistItem

	collector := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"),
		colly.AllowedDomains("www.amazon.com.br", "amazon.com.br"),
	)

	collector.OnRequest(func(request *colly.Request) {
		request.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
		request.Headers.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
		request.Headers.Set("Sec-Fetch-Dest", "document")
		request.Headers.Set("Sec-Fetch-Mode", "navigate")
		request.Headers.Set("Sec-Fetch-Site", "none")
	})

	collector.OnHTML("li.g-item-sortable", func(element *colly.HTMLElement) {
		title := strings.TrimSpace(element.ChildText("a[id^='itemName_']"))
		link := element.ChildAttr("a[id^='itemName_']", "href")
		price := strings.TrimSpace(element.ChildText(".a-price .a-offscreen"))
		if price == "" {
			price = strings.TrimSpace(element.ChildText("span[id^='itemPrice_']"))
		}

		if title != "" {
			items = append(items, domain.WishlistItem{
				Title: title,
				Price: price,
				Link:  element.Request.AbsoluteURL(link),
			})
		}
	})

	collector.OnHTML("a.wl-see-more", func(element *colly.HTMLElement) {
		if nextPage := element.Attr("href"); nextPage != "" {
			if err := element.Request.Visit(element.Request.AbsoluteURL(nextPage)); err != nil {
				log.Printf("Could not visit wishlist next page: %v", err)
			}
		}
	})

	collector.OnError(func(response *colly.Response, err error) {
		log.Printf("Error on %s: %v (status: %d)", response.Request.URL, err, response.StatusCode)
	})

	wishlistID := strings.TrimSpace(os.Getenv("WISHLIST_ID"))
	if wishlistID == "" {
		return nil, fmt.Errorf("WISHLIST_ID is required")
	}

	url := fmt.Sprintf("https://www.amazon.com.br/hz/wishlist/ls/%s?_encoding=UTF8&sort=price-asc&filter=unpurchased", wishlistID)
	if err := collector.Visit(url); err != nil {
		return nil, fmt.Errorf("visit Amazon wishlist: %w", err)
	}
	collector.Wait()

	return items, nil
}
