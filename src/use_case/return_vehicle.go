package use_case

import (
	"context"
	"fmt"
	"rental-vehicle-system/src/entity/booking"
	"rental-vehicle-system/src/entity/price_model"

	"github.com/cockroachdb/errors"
)

func (u UseCase) ReturnVehicle(ctx context.Context, bookingId string) (string, error) {
	bDetail, err := u.bookingRepository.GetBookingById(ctx, bookingId)
	if err != nil {
		return "", err
	}

	vehicleDetail, err := u.vehicleRepository.GetVehicleById(ctx, bDetail.Booking.VehicleId)
	if err != nil {
		return "", err
	}

	bStatus := bDetail.Booking.GetStatus()
	if bStatus != booking.RentStatus {
		return "", errors.Errorf("Booking status is '%s': %w", bStatus, ErrBookingStatusIsNotRent)
	}

	returnDate, lateFines := u.CalculateLateFines(bDetail, vehicleDetail)
	summaryFine := u.CalculateSummaryLateFines(bDetail.Booking.Deposit, lateFines)

	bookingReturnedDetail := booking.Booking{
		Status:      booking.ReturnStatus,
		ReturnDate:  returnDate,
		LateFines:   lateFines,
		Deposit:     bDetail.Booking.Deposit,
		SummaryFine: summaryFine,
	}

	u.bookingRepository.UpdateReturnedBooking(ctx, bookingId, bookingReturnedDetail)
	u.vehicleRepository.UpdateVehicleStatus(ctx, bDetail.Booking.VehicleId, price_model.ActiveStatus)

	successMessage := fmt.Sprintf("Returned success with bookindId %s", bookingId)
	return successMessage, nil
}
