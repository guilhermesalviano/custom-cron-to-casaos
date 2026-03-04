# Google Flights Crawler (Go)

A command-line tool to fetch and display the best flight prices from Google Flights via the **SerpApi** service.

> Google Flights is heavily JavaScript-rendered and actively blocks scrapers.
> SerpApi handles headless browsers, proxy rotation, and CAPTCHAs for you.
> A **free tier** gives you 100 searches/month — enough for personal use.

---

## Setup

### 1. Get a free SerpApi key
Sign up at https://serpapi.com/ — no credit card required for the free tier.

### 2. Build
```bash
go build -o flights ./main.go
```

### 3. Set your API key
```bash
export SERPAPI_KEY="your_key_here"
```
Or pass it directly with the `-key` flag.

---

## Usage

```bash
# Basic one-way search (GRU → JFK next month)
./flights -from GRU -to JFK -date 2025-04-10

# Round trip search
./flights -from GRU -to LIS -date 2025-04-10 -return 2025-04-20

# Business class, nonstop only, output to JSON
./flights -from GRU -to MIA -date 2025-04-15 -class 3 -stops 1 -output result.json

# All flags
./flights \
  -key  YOUR_API_KEY  \   # or use SERPAPI_KEY env var
  -from GRU           \   # departure IATA code
  -to   JFK           \   # arrival IATA code
  -date 2025-04-10    \   # outbound date
  -return 2025-04-20  \   # return date (omit for one-way)
  -adults 2           \   # number of adult passengers
  -class  1           \   # 1=Economy 2=Premium 3=Business 4=First
  -stops  0           \   # 0=Any 1=Nonstop 2=1stop 3=2stops
  -currency BRL       \   # currency code
  -lang     pt        \   # language (pt, en, es, fr...)
  -country  br        \   # country (br, us, fr...)
  -output result.json     # save JSON output (optional)
```

---

## Example Output

```
🔍 Searching flights GRU → JFK on 2025-04-10...

╔══════════════════════════════════════════════════════╗
║         Google Flights — Best Price Crawler          ║
╚══════════════════════════════════════════════════════╝
  Route     : GRU → JFK
  Outbound  : 2025-04-10
  Searched  : 2025-03-03 14:22:10 UTC
  ★ Best price: 2840 BRL

── Best Flights (3 results) ──────────────────────────────
  Airline                Dep      Arr      Duration   Stops  Price
  ─────────────────────────────────────────────────────────
  LATAM Airlines         23:10    07:35    10h25m     nonstop 2840 BRL
  American Airlines      18:00    06:20    13h20m     1 stop  3120 BRL
  United Airlines        19:45    09:15    14h30m     1 stop  3390 BRL
```

---

## JSON Output Structure

```json
{
  "searched_at": "2025-03-03T14:22:10Z",
  "origin": "GRU",
  "destination": "JFK",
  "outbound_date": "2025-04-10",
  "best_flights": [
    {
      "airline": "LATAM Airlines",
      "flight_number": "LA 8084",
      "departure_time": "2025-04-10 23:10",
      "arrival_time": "2025-04-11 07:35",
      "duration_minutes": 625,
      "stops": 0,
      "price": 2840,
      "currency": "BRL",
      "carbon_emissions_kg": 312
    }
  ],
  "other_flights": [...],
  "best_price": 2840,
  "currency": "BRL"
}
```

---

## Common IATA Codes (Brazil)

| Code | Airport |
|------|---------|
| GRU  | São Paulo (Guarulhos) |
| CGH  | São Paulo (Congonhas) |
| GIG  | Rio de Janeiro (Galeão) |
| BSB  | Brasília |
| CNF  | Belo Horizonte |
| SSA  | Salvador |
| REC  | Recife |
| FOR  | Fortaleza |
| CWB  | Curitiba |
| POA  | Porto Alegre |

---

## Notes

- **No third-party dependencies** — only the Go standard library is used.
- Results are sorted by price (cheapest first).
- The crawler respects Google's terms of service by routing through SerpApi.