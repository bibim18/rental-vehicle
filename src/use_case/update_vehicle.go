package use_case

import (
	"context"
	"fmt"
	"rental-vehicle-system/src/entity/vehicle"
)

func (u UseCase) UpdateVehicleDetail(ctx context.Context, vId string, v VehicleFullDetail) (string, error) {
	validateErr := v.Vehicle.Validate()
	if validateErr != nil {
		return "", validateErr
	}

	vehicleDetail, err := u.vehicleRepository.GetVehicleById(ctx, vId)
	if err != nil {
		return "", err
	}

	if vehicleDetail.Status != vehicle.ReadyStatus {
		return "", vehicle.ErrVehicleNotUpdate
	}

	u.vehicleRepository.UpdateVehicle(ctx, vId, v)
	successMessage := fmt.Sprintf("Update success with vehicleId %s", vId)
	return successMessage, nil
}
