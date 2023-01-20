package use_case

import (
	"context"
	"fmt"
	"rental-vehicle-system/src/entity/booking"
	"rental-vehicle-system/src/entity/price_model"
)

func (u UseCase) CancelationBooking(ctx context.Context, bId string) (string, error) {
	bDetail, err := u.bookingRepository.GetBookingById(ctx, bId)
	if err != nil {
		return "", err
	}

	bStatus := bDetail.Booking.GetStatus()
	if bStatus != booking.RentStatus {
		return "", ErrBookingStatusIsNotCancel
	}

	u.bookingRepository.UpdateBookingStatus(ctx, bId, booking.CancelStatus)
	u.vehicleRepository.UpdateVehicleStatus(ctx, bDetail.Booking.VehicleId, price_model.ActiveStatus)

	successMessage := fmt.Sprintf("Cancelation success with bookindId %s", bId)
	return successMessage, nil
}
