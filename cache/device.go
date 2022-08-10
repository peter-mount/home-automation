package cache

import (
	"encoding/json"
)

// Device represents a device reported by zigbee2mqtt
type Device struct {
	DateCode         string `json:"dateCode,omitempty"`
	FriendlyName     string `json:"friendly_name,omitempty"`
	IeeeAddr         string `json:"ieeeAddr,omitempty"`
	LastSeen         int64  `json:"lastSeen,omitempty"`
	NetworkAddress   int64  `json:"networkAddress,omitempty"`
	SoftwareBuildID  string `json:"softwareBuildID,omitempty"`
	Type             string `json:"type,omitempty"`
	Description      string `json:"description,omitempty"`
	HardwareVersion  int64  `json:"hardwareVersion,omitempty"`
	ManufacturerID   int64  `json:"manufacturerID,omitempty"`
	ManufacturerName string `json:"manufacturerName,omitempty"`
	Model            string `json:"model,omitempty"`
	ModelID          string `json:"modelID,omitempty"`
	PowerSource      string `json:"powerSource,omitempty"`
	Vendor           string `json:"vendor,omitempty"`
}

func (d *Device) ToJson() ([]byte, error) {
	return json.Marshal(d)
}
