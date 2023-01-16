package use_case

import (
	"context"
	"errors"
	"time"

	"rental-vehicle-system/src/entity/booking"
	"rental-vehicle-system/src/entity/customer"
	"rental-vehicle-system/src/entity/vehicle"
)

type UseCase struct {
	vehicleRepository VehicleRepository
	bookingRepository BookingRepository
}

type VehicleFullDetail struct {
	Vehicle        vehicle.Vehicle
	Status         string
	RegisteredDate time.Time
}

type BookingFullDetail struct {
	Booking  booking.Booking
	Customer customer.Customer
}

type VehicleRepository interface {
	GetVehicleById(ctx context.Context, vehicleId string) (VehicleFullDetail, error)
	GetVehicleByStatus(ctx context.Context, status string, limit, offset int64) ([]VehicleFullDetail, error)
	CreateVehicle(ctx context.Context, vehicle VehicleFullDetail) error
	UpdateVehicleStatus(ctx context.Context, vehicleId string, status string) error
	UpdateVehicle(ctx context.Context, vehicleId string, vehicle VehicleFullDetail) error
}

type BookingRepository interface {
	UpdateBookingStatus(ctx context.Context, bookingId string, status booking.BookingStatus) error
	UpdateReturnedBooking(ctx context.Context, bookingId string, booking booking.Booking) error
	GetBookingById(ctx context.Context, bookingId string) (BookingFullDetail, error)
	CreateBooking(ctx context.Context, booking BookingFullDetail) error
}

var (
	ErrVehicleNotReadyForRent   = errors.New("vehicle not ready for rent")
	ErrBookingStatusIsNotRent   = errors.New("Booking cannot return")
	ErrBookingStatusIsNotCancel = errors.New("Booking cannot cancel")
)

func New(vehicleRepo VehicleRepository, bookingRepo BookingRepository) *UseCase {
	return &UseCase{
		vehicleRepository: vehicleRepo,
		bookingRepository: bookingRepo,
	}
}
