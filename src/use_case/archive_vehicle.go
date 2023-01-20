package use_case

import (
	"context"
	"fmt"
	"rental-vehicle-system/src/entity/price_model"
)

func (u UseCase) ArchiveVehicle(ctx context.Context, vehicleId string) (string, error) {
	vehicleDetail, err := u.vehicleRepository.GetVehicleById(ctx, vehicleId)
	if err != nil {
		return "", err
	}

	if vehicleDetail.Status != price_model.ActiveStatus {
		return "", price_model.ErrVehicleNotArchive
	}

	u.vehicleRepository.UpdateVehicleStatus(ctx, vehicleId, price_model.InactiveStatus)
	successMessage := fmt.Sprintf("Disabled price_model success with vehicleId %s", vehicleId)
	return successMessage, nil
}
