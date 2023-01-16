package use_case

import (
	"context"
	"fmt"
	"rental-vehicle-system/src/entity/vehicle"
)

func (u UseCase) ArchiveVehicle(ctx context.Context, vehicleId string) (string, error) {
	vehicleDetail, err := u.vehicleRepository.GetVehicleById(ctx, vehicleId)
	if err != nil {
		return "", err
	}

	if vehicleDetail.Status != vehicle.ReadyStatus {
		return "", vehicle.ErrVehicleNotArchive
	}

	u.vehicleRepository.UpdateVehicleStatus(ctx, vehicleId, vehicle.UnactiveStatus)
	successMessage := fmt.Sprintf("Archive success with vehicleId %s", vehicleId)
	return successMessage, nil
}
