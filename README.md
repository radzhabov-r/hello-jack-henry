# hello-jack-henry
Jack Henry intro assigment

# Weather Service

A simple web server that tells you the weather for any location using the National Weather Service.

## How it's organized

The code is split up into different folders to keep things tidy:

```
main.go             - Where the app starts
config.json         - Settings for the app (port, temperature ranges)
internal/config/    - Reads the config file
internal/router/    - Web page routes  
internal/handler/   - Handles web requests and does the main work
internal/models/    - Data structures
internal/service/   - Talks to the weather API
```

## Configuration

The app uses a `config.json` file to set things up. You can change these settings:

```json
{
  "server": {
    "port": 8080,
    "read_timeout_seconds": 15,
    "write_timeout_seconds": 15
  },
  "weather": {
    "temperature_ranges": {
      "hot_threshold": 80,
      "cold_threshold": 50
    }
  }
}
```

Want the app to run on a different port? Just change the port number in config.json. Want different temperature ranges? Update the hot and cold thresholds.

## Getting it running

### What you need
- Go 1.21 or newer

### Build it
```bash
go build -o weather-service main.go
```

### Start it up
```bash
./weather-service
```

Or just run it directly:
```bash
go run main.go
```

The server will start on whatever port is set in config.json (default is 8080). Want to change the port? Just edit the config.json file.

## How to use it

### Get the weather

**URL:** `GET /weather`

**What to send:**
- `lat` (required): Latitude (-90 to 90)
- `lon` (required): Longitude (-180 to 180)

**Try it:**
```bash
curl "http://localhost:8080/weather?lat=39.7456&lon=-97.0892"
```

**You'll get back:**
```json
{
  "forecast": "Partly Cloudy",
  "temperature": "moderate"
}
```

### Check if it's working

**URL:** `GET /health`

Just returns `OK` if everything's running fine.

## Try these places

```bash
# Denver, Colorado
curl "http://localhost:8080/weather?lat=39.7392&lon=-104.9903"

# New York City
curl "http://localhost:8080/weather?lat=40.7128&lon=-74.0060"

# Miami, Florida  
curl "http://localhost:8080/weather?lat=25.7617&lon=-80.1918"

# Barrow, AK (where it can snow in summer)    
curl "http://localhost:8080/weather?lat=71.2906&lon=-156.7887"
```

## Things I didn't add (to save time)

This is just a demo, so I kept it simple:
- Pretty basic error messages
- No tests yet

For a real app, you'd probably want:
- Caching so it doesn't hit the weather API every time
- Rate limiting so people can't spam it
- Health checks that actually test the weather API
- Monitoring and metrics if needed
