package use_case

import (
	"rental-vehicle-system/src/entity/booking"
	"rental-vehicle-system/src/entity/vehicle"
	"time"

	log "github.com/sirupsen/logrus"
)

type IsOverDue bool
type LateFines float32

func (u UseCase) CalculateRentalPrice(v VehicleFullDetail, b BookingDetail) int {
	log.Info("calculate_booking.CalculateRentalPrice")
	qty := b.Booking.Qty
	vehiclePrice := v.Vehicle.GetPrice(vehicle.UnitType(b.Booking.Unit))
	return qty * vehiclePrice
}

func (u UseCase) CalculatePrice(deposit vehicle.Deposit, price int) (vehicle.Deposit, vehicle.Deposit) {
	log.Info("calculate_booking.CalculatePrice")

	priceWithDepositType := vehicle.Deposit(price)

	rentalDeposit := deposit * priceWithDepositType
	totalPrice := rentalDeposit + priceWithDepositType

	return rentalDeposit, totalPrice
}

func (u UseCase) CalculateLateFines(bDetail BookingFullDetail, vDetail VehicleFullDetail) (time.Time, int) {
	// คำนวนค่าปรับถ้าวันคืนเกิดกำหนด
	log.Info("calculate_booking.CalculateLateFines")

	var lateFines int = 0
	returnDate := time.Now()
	isOverDue := returnDate.After(bDetail.Booking.DueDate)
	if isOverDue {
		totalOverDue := bDetail.Booking.GetTotalOverDueDays()
		pricePerDate := vDetail.Vehicle.GetPrice(vehicle.UnitType(booking.DailyUnit))
		lateFines = totalOverDue * pricePerDate
	}

	return returnDate, lateFines
}

func (u UseCase) CalculateSummaryLateFines(deposit booking.DepositType, lateFines int) int {
	// หักค่าปรับกับมัดจำ
	log.Info("calculate_booking.CalculateSummaryLateFines")
	var totalLateFines int = 0
	if lateFines > 0 {
		totalLateFines = int(deposit) - lateFines
	}

	return totalLateFines
}
