package use_case

import (
	"context"

	"rental-vehicle-system/src/entity/booking"
	"rental-vehicle-system/src/entity/customer"
	"rental-vehicle-system/src/entity/price_model"
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

	v, err := u.vehicleRepository.GetVehicleById(ctx, b.Booking.VehicleId)
	if err != nil {
		return "", err
	}

	if v.Status() != price_model.ActiveStatus {
		return "", ErrVehicleNotReadyForRent
	}

	vehicleMethod := price_model.New(v.Vehicle.VehicleType, v.Vehicle)
	deposit := vehicleMethod.GetDeposit()

	rentalCost := u.calculateRentalPrice(v, b)
	rentalDeposit, totalPrice := u.CalculatePrice(deposit, rentalCost)
	rentDate, dueDate := b.Booking.CalculateDate()

	bFull := booking.Booking{
		VehicleId:  b.Booking.VehicleId,
		Unit:       b.Booking.Unit,
		Qty:        b.Booking.Qty,
		Status:     booking.RentStatus,
		RentPrice:  rentalCost,
		TotalPrice: booking.TotalPriceType(totalPrice),
		Deposit:    booking.DepositType(rentalDeposit),
		RentDate:   rentDate,
		DueDate:    dueDate,
	}
	u.bookingRepository.CreateBooking(ctx, BookingFullDetail{bFull, b.Customer})
	u.vehicleRepository.UpdateVehicleStatus(ctx, b.Booking.VehicleId, price_model.InuseStatus)

	return "Create booking success", nil
}
