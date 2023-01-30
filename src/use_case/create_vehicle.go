package use_case

import (
	"context"
	"fmt"
	"time"

	"rental-vehicle-system/src/entity/vehicle"
)

type VehicleWithPrice struct {
	Vehicles  []vehicle.Vehicle
	RatePrice vehicle.Price
}

func (u UseCase) CreateVehicle(ctx context.Context, v VehicleWithPrice) (string, error) {
	for _, vItem := range v.Vehicles {
		vItem.RatePrice = v.RatePrice
		err := vItem.Validate()

		if err != nil {
			return "", err
		}
		vFull := VehicleFullDetail{
			Vehicle: vehicle.Vehicle{
				Brand:        vItem.Brand,
				Model:        vItem.Model,
				LicensePlate: vItem.LicensePlate,
				VehicleType:  vItem.VehicleType,
				Color:        vItem.Color,
				RatePrice:    vItem.RatePrice,
			},
			Status:         vehicle.ActiveStatus,
			RegisteredDate: time.Now(),
		}
		u.vehicleRepository.CreateVehicle(ctx, vFull)
	}

	successMessage := fmt.Sprintf("Created vehicles %d items", len(v.Vehicles))
	return successMessage, nil
}
