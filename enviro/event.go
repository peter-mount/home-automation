package enviro

import "time"

// Event represents the fields the enviro sensors send to MQTT.
// Not all fields are populated by all devices
type Event struct {
	Device        string    `json:"device"`
	Timestamp     time.Time `json:"timestamp"`
	Pressure      float64   `json:"pressure"`
	Temperature   float64   `json:"temperature"`
	Humidity      float64   `json:"humidity"`
	WindDirection int       `json:"wind_direction"`
	Light         float64   `json:"light"`
	Rain          float64   `json:"rain"`
	WindSpeed     float64   `json:"wind_speed"`
	PM10          int       `json:"pm10"`
	PM25          int       `json:"pm2_5"`
	Noise         float64   `json:"noise"`
	PM1           int       `json:"pm1"`
}
