# Wishlist service

This is the independently deployable Amazon wishlist crawler. It runs one scheduled crawl, saves results to MySQL, and optionally sends Discord or NTFY notifications.

Required environment variables:

```env
WISHLIST_ID=your_amazon_wishlist_id
WISHLIST_DAY=Saturday
WISHLIST_TIME=17:00
DB_HOST=your_mysql_host
DB_NAME=your_database_name
DB_USER=your_database_user
DB_PASSWORD=your_database_password
```

`DB_PORT` defaults to `3306`. `WISHLIST_TIMEZONE` defaults to `America/Sao_Paulo`. `DISCORD_WEBHOOK_URL` is optional.
Set `NTFY_TOPIC` to publish the same notifications to `https://ntfy.sh`.

Build and run locally:

```bash
go run .
docker build -t wishlist .
```
