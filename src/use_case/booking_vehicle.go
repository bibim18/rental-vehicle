package use_case

import (
	"context"

	"rental-vehicle-system/src/entity/booking"
	"rental-vehicle-system/src/entity/customer"
	"rental-vehicle-system/src/entity/vehicle"
)

type BookingDetail struct {
	Booking  booking.Booking
	Customer customer.Customer
}

func (u UseCase) BookingVehicle(ctx context.Context, b BookingDetail) (string, error) {
	validateCustomerErr := b.Customer.Validate()
	if validateCustomerErr != nil {
		return "", validateCustomerErr
	}

	validateBookingErr := b.Booking.Validate()
	if validateBookingErr != nil {
		return "", validateBookingErr
	}
	vehicleDetail, _ := u.vehicleRepository.GetVehicleById(ctx, b.Booking.VehicleId)
	if vehicleDetail.Status != vehicle.ReadyStatus {
		return "", ErrVehicleNotReadyForRent
	}

	vehicleMethod := vehicle.New(vehicleDetail.Vehicle.VehicleType, vehicleDetail.Vehicle)

	rentalCost := u.CalculateRentalPrice(vehicleDetail, b)
	deposit := vehicleMethod.GetDeposit()
	rentalDeposit, totalPrice := u.CalculatePrice(deposit, rentalCost)
	rentDate, dueDate := b.Booking.CalculateDate()

	bFull := booking.Booking{
		VehicleId:  b.Booking.VehicleId,
		Unit:       b.Booking.Unit,
		Qty:        int(b.Booking.Qty),
		Status:     "rent",
		RentPrice:  rentalCost,
		TotalPrice: booking.TotalPriceType(totalPrice),
		Deposit:    booking.DepositType(rentalDeposit),
		RentDate:   rentDate,
		DueDate:    dueDate,
	}
	u.bookingRepository.CreateBooking(ctx, BookingFullDetail{bFull, b.Customer})
	u.vehicleRepository.UpdateVehicleStatus(ctx, b.Booking.VehicleId, vehicle.InuseStatus)

	return "Create booking success", nil
}
