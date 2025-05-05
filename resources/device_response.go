package resources

import (
	"ticatag_backend/models"
)

type DeviceResponse struct {
	//ID     primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Adress    string `json:"adress" bson:"adress"`
	Latitude  string `json:"latitude" bson:"latitude"`
	Longitude string `json:"longitude" bson:"longitude"`
}

func NewDeviceResponse(device models.Device) DeviceResponse {
	return DeviceResponse{
		// ID:    user.ID.Hex(),
		Adress:    device.Adress,
		Latitude:  device.Latitude,
		Longitude: device.Longitude,
	}
}

func NewDeviceListResponse(devices []models.Device) []DeviceResponse {
	responses := make([]DeviceResponse, 0, len(devices))
	for _, device := range devices {
		responses = append(responses, NewDeviceResponse(device))
	}
	return responses
}
