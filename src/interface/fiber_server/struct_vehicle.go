package fiber_server

import (
	"github.com/cockroachdb/errors"
	"rental-vehicle-system/src/entity/price_model"
)

type vehicleDataRequest struct {
	VehicleType  string `json:"price_model"`
	Brand        string `json:"brand"`
	Model        string `json:"model"`
	LicensePlate string `json:"license_plate"`
	Color        string `json:"color"`

	// pricing
	Upfront       int `json:"upfront"`
	PricePerHour  int `json:"price_per_hour"`
	PricePerDay   int `json:"price_per_day"`
	PricePerWeek  int `json:"price_per_week"`
	PricePerMonth int `json:"price_per_month"`
	PricePerYear  int `json:"price_per_year"`
	PricePerKm    int `json:"price_per_km"`
}

func NewVehicle(v vehicleDataRequest) (price_model.PriceModel, error) {
	switch v.VehicleType {
	case "car":
		return price_model.NewBasicPriceModel(
			v.LicensePlate,
			v.Brand,
			v.Model,
			v.Color,
			v.Upfront,
			v.PricePerDay,
			v.PricePerMonth,
			v.PricePerKm,
		), nil
	}

	return nil, errors.New("invalid price_model type")

}
