package resources

import (
	"ticatag_backend/models"
)

type DeviceResponse struct {
	ID            string `bson:"_id,omitempty" json:"id"`
	Adress        string `json:"adress" bson:"adress"`
	Latitude      string `json:"latitude" bson:"latitude"`
	Longitude     string `json:"longitude" bson:"longitude"`
	Adresspostale string `json:"addresspostale" bson:"adresspostale"`
	CreatedAt     int64  `bson:"created_at" json:"created_at"`
}

func NewDeviceResponse(device models.Device) DeviceResponse {
	return DeviceResponse{
		ID:            device.ID.Hex(),
		Adress:        device.Adress,
		Latitude:      device.Latitude,
		Longitude:     device.Longitude,
		Adresspostale: device.Adresspostale,
		CreatedAt:     device.CreatedAt,
	}
}

func NewDeviceListResponse(devices []models.Device) []DeviceResponse {
	responses := make([]DeviceResponse, 0, len(devices))
	for _, device := range devices {
		responses = append(responses, NewDeviceResponse(device))
	}
	return responses
}
