package cache

import (
	"encoding/json"
	"log"
	"sort"
	"strings"
)

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

func (c *Cache) GetDevice(name string) *Device {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	//name = b.publisher.EncodeKey(name)
	if strings.HasPrefix(name, "zigbee2mqtt/") {
		name = name[12:]
	}
	log.Println(name)
	if d, exists := c.devices[name]; exists {
		return d
	}
	return nil
}

func (b *Cache) addDevice(d *Device) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if d.FriendlyName != "" {
		/*if _, exists := b.devices[d.FriendlyName]; !exists {
			log.Printf("New device %s %q %q", d.FriendlyName, d.Model, d.Description)
		}*/
		b.devices[d.FriendlyName] = d
	}
}

// GetDevices returns a list of devices
func (b *Cache) GetDevices() []*Device {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	var r []*Device
	for _, v := range b.devices {
		r = append(r, v)
	}

	sort.SliceStable(r, func(i, j int) bool {
		return r[i].FriendlyName < r[j].FriendlyName
	})

	return r
}
