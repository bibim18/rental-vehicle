package use_case

import (
	"context"
	"rental-vehicle-system/src/entity/price_model"
)

func (u UseCase) GetVehicle(ctx context.Context, status string) ([]VehicleFullDetail, error) {
	var list []VehicleFullDetail
	var err error

	if status != "" {
		err = price_model.VerifyStatus(status)
		if err != nil {
			return []VehicleFullDetail{}, err
		}
	}

	list, err = u.vehicleRepository.GetVehicleByStatus(ctx, status, 5, 10)

	return list, nil
}
