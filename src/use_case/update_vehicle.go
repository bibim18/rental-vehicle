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

	if vehicleDetail.Status != vehicle.ActiveStatus {
		return "", vehicle.ErrVehicleNotUpdate
	}

	u.vehicleRepository.UpdateVehicle(ctx, vId, v)
	successMessage := fmt.Sprintf("Update success with vehicleId %s", vId)
	return successMessage, nil
}

func (u UseCase) UpdateVehicleStatus(ctx context.Context, vehicleId string) (string, error) {
	vehicleDetail, err := u.vehicleRepository.GetVehicleById(ctx, vehicleId)
	if err != nil {
		return "", err
	}

	if vehicleDetail.Status != vehicle.UnactiveStatus {
		return "", vehicle.ErrVehicleNotEnable
	}

	u.vehicleRepository.UpdateVehicleStatus(ctx, vehicleId, vehicle.ActiveStatus)
	successMessage := fmt.Sprintf("Enabled vehicle success with vehicleId %s", vehicleId)
	return successMessage, nil
}
