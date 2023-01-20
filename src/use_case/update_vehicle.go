package use_case

import (
	"context"
	"fmt"
	"rental-vehicle-system/src/entity/price_model"
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

	if vehicleDetail.Status != price_model.ActiveStatus {
		return "", price_model.ErrVehicleNotUpdate
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

	if vehicleDetail.Status != price_model.InactiveStatus {
		return "", price_model.ErrVehicleNotEnable
	}

	u.vehicleRepository.UpdateVehicleStatus(ctx, vehicleId, price_model.ActiveStatus)
	successMessage := fmt.Sprintf("Enabled price_model success with vehicleId %s", vehicleId)
	return successMessage, nil
}
