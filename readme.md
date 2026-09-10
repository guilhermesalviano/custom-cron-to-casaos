# My Cron - selfhosted

A Go-based application that schedules automated web scraping tasks for Google Flights and Amazon Wishlists using cron-like scheduling. It integrates with AI analysis via Gemini, sends notifications via Discord or NTFY, and stores data in a MySQL database.

## Features

- **Scheduled Flights Crawling**: Automatically scrape Google Flights data based on predefined schedules from a CSV file.
- **Amazon Wishlist Monitoring**: Crawl Amazon wishlists for price tracking and availability.
- **AI-Powered Analysis**: Uses Google's Gemini AI to analyze scraped data.
- **Flexible Scheduling**: Uses gocron for precise scheduling (e.g., specific days and times).
- **Notification System**: Send alerts via Discord webhooks or NTFY.
- **Database Integration**: Stores results in MySQL database.
- **Docker Support**: Containerized deployment with Docker Compose.
- **Independent Wishlist Image**: Amazon wishlist can be deployed separately as `${DOCKER_USERNAME}/wishlist`.
- **Environment Configuration**: Uses .env files for sensitive configuration.

## Prerequisites

- Go 1.25.0 or later
- MySQL database
- API keys for:
  - SerpAPI (for Google Flights scraping)
  - Google Gemini AI
  - Discord webhook (optional)
- Docker and Docker Compose (for containerized deployment)

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/guilhermesalviano/google-flights-crawler.git
   cd google-flights-crawler
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up environment variables (see Configuration section).

## Configuration

Create a `.env` file in the root directory with the following variables:

```env
DB_HOST=your_mysql_host
DB_NAME=your_database_name
DB_PASSWORD=your_db_password
DB_PORT=your_db_port
DB_USER=your_db_username
DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/your_webhook_id/your_webhook_token
SERPAPI_KEY=your_serpapi_key
SCHEDULERS_FILE_PATH=path/to/schedulers.csv
```

### Schedulers CSV Format

The application uses a CSV file (`schedulers.csv`) to define scheduled tasks. The format includes:

- For flights: departure_id, arrival_id, outbound_date, return_date, adults, travel_class, stops, currency, language, country, day, time
- For wishlists: wishlist_url, day, time

Example `schedulers.csv`:
```
type,departure_id,arrival_id,outbound_date,return_date,adults,travel_class,stops,currency,language,country,day,time
flights,GRU,JFK,2024-12-01,2024-12-15,1,economy,0,USD,en,US,monday,09:00
wishlists,tuesday,10:00
```

## Usage

### Running Locally

1. Ensure your `.env` file is configured.
2. Run the application:
   ```bash
   go run main.go
   ```

The application will start the scheduler and begin executing tasks according to the CSV schedules.

### Docker Deployment

1. Update the `docker-compose.yml` with your environment variables.
2. Run with Docker Compose:
   ```bash
   docker-compose up -d
   ```

The independent wishlist service uses `WISHLIST_ID`, `WISHLIST_DAY`, and `WISHLIST_TIME` instead of the combined application's schedules CSV. Its timezone defaults to `America/Sao_Paulo`; set `WISHLIST_TIMEZONE` to override it. GitHub releases publish the image as `${DOCKER_USERNAME}/wishlist` with both the release version and `latest` tags.

See [services/wishlist/README.md](services/wishlist/README.md) for the standalone service configuration.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
