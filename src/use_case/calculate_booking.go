package use_case

import (
	"rental-vehicle-system/src/entity/booking"
	"rental-vehicle-system/src/entity/price_model"
	"time"

	log "github.com/sirupsen/logrus"
)

type IsOverDue bool
type LateFines float32

func (u UseCase) calculateRentalPrice(v VehicleFullDetail, b BookingDetail) int {
	log.Info("calculate_booking.calculateRentalPrice")
	qty := b.Booking.Qty
	vehiclePrice := v.Vehicle.GetPrice(price_model.DateUnit(b.Booking.Unit)) * vehiclePrice
	return qty * vehiclePrice
}

func (u UseCase) CalculatePrice(deposit price_model.Deposit, price int) (price_model.Deposit, price_model.Deposit) {
	log.Info("calculate_booking.CalculatePrice")

	priceWithDepositType := price_model.Deposit(price)

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
		pricePerDate := vDetail.Vehicle.GetPrice(price_model.DateUnit(booking.DailyUnit))
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
