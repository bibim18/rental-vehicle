package use_case

import (
	"context"
	"rental-vehicle-system/src/entity/price_model"
)

func (u UseCase) CreateVehicle(ctx context.Context, v price_model.PriceModel) error {
	return u.vehicleRepository.CreateVehicle(ctx, v)
}
