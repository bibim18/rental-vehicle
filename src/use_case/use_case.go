package use_case

import (
	"context"
	"errors"
	"rental-vehicle-system/src/entity/booking"
	"rental-vehicle-system/src/entity/customer"
	"rental-vehicle-system/src/entity/price_model"
)

type UseCase struct {
	vehicleRepository VehicleRepository
	bookingRepository BookingRepository
}

type BookingFullDetail struct {
	Booking  booking.Booking
	Customer customer.Customer
}

type VehicleRepository interface {
	GetVehicleById(ctx context.Context, vehicleId string) (price_model.PriceModel, error)
	GetVehicleByStatus(ctx context.Context, status string, limit, offset int64) ([]price_model.PriceModel, error)

	CreateVehicle(ctx context.Context, vehicle price_model.PriceModel) error
	CreateVehicles(ctx context.Context, vehicles []price_model.PriceModel) error

	UpdateVehicleStatus(ctx context.Context, vehicleId string, status string) error
	UpdateVehicle(ctx context.Context, vehicleId string, vehicle price_model.PriceModel) error
}

type BookingRepository interface {
	UpdateBookingStatus(ctx context.Context, bookingId string, status booking.BookingStatus) error
	UpdateReturnedBooking(ctx context.Context, bookingId string, booking booking.Booking) error
	GetBookingById(ctx context.Context, bookingId string) (BookingFullDetail, error)
	CreateBooking(ctx context.Context, booking BookingFullDetail) error
}

var (
	ErrVehicleNotReadyForRent   = errors.New("price_model not ready for rent")
	ErrBookingStatusIsNotRent   = errors.New("Booking cannot return")
	ErrBookingStatusIsNotCancel = errors.New("Booking cannot cancel")
)

func New(vehicleRepo VehicleRepository, bookingRepo BookingRepository) *UseCase {
	return &UseCase{
		vehicleRepository: vehicleRepo,
		bookingRepository: bookingRepo,
	}
}
