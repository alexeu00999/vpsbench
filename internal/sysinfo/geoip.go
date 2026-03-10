package sysinfo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

const geoIPURL = "http://ip-api.com/json/?fields=query,country,city,status"

type geoIPResponse struct {
	Status  string `json:"status"`
	Country string `json:"country"`
	City    string `json:"city"`
	Query   string `json:"query"`
}

// detectLocation определяет локацию сервера через GeoIP API.
func detectLocation(ctx context.Context) (ip, location string) {
	slog.Debug("[sysinfo] detecting location via GeoIP", "url", geoIPURL)

	// Таймаут 5 секунд для GeoIP запроса
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, geoIPURL, nil)
	if err != nil {
		slog.Warn("[sysinfo] failed to create GeoIP request", "error", err)
		return "", "Unknown"
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Warn("[sysinfo] GeoIP request failed", "error", err)
		return "", "Unknown"
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Warn("[sysinfo] GeoIP returned non-200 status", "status", resp.StatusCode)
		return "", "Unknown"
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Warn("[sysinfo] failed to read GeoIP response", "error", err)
		return "", "Unknown"
	}

	slog.Debug("[sysinfo] GeoIP raw response", "body", string(body))

	var geo geoIPResponse
	if err := json.Unmarshal(body, &geo); err != nil {
		slog.Warn("[sysinfo] failed to parse GeoIP response", "error", err)
		return "", "Unknown"
	}

	if geo.Status != "success" {
		slog.Warn("[sysinfo] GeoIP returned non-success status", "status", geo.Status)
		return "", "Unknown"
	}

	ip = geo.Query
	if geo.City != "" && geo.Country != "" {
		location = fmt.Sprintf("%s, %s", geo.City, geo.Country)
	} else if geo.Country != "" {
		location = geo.Country
	} else {
		location = "Unknown"
	}

	slog.Info("[sysinfo] location detected", "ip", ip, "location", location)
	return ip, location
}
