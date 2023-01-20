package booking

import (
	"time"

	"github.com/cockroachdb/errors"
	"github.com/go-playground/validator/v10"
)

var (
	ErrInvalidBooking  = errors.New("invalid booking")
	ErrBookingNotFound = errors.New("Cannot found booking")
)

type Booking struct {
	VehicleId   string   `validate:"required"`
	Unit        UnitType `validate:"oneof='daily' 'monthly' 'yearly'"`
	Qty         int      `validate:"gt=0"`
	Status      BookingStatus
	DueDate     time.Time
	RentDate    time.Time
	ReturnDate  time.Time
	TotalPrice  TotalPriceType
	RentPrice   int
	LateFines   int
	Deposit     DepositType
	SummaryFine int
}
type UnitType string
type DepositType float32
type TotalPriceType float32
type BookingStatus string

const (
	RentStatus   BookingStatus = "rent"
	CancelStatus BookingStatus = "cancel"
	ReturnStatus BookingStatus = "return"
)

const (
	DailyUnit   UnitType = "daily"
	MonthlyUnit UnitType = "monthly"
	YearlyUnit  UnitType = "yearly"
)

func (b Booking) Validate() error {
	err := validator.New().Struct(b)
	if err != nil {
		return errors.Errorf("%s: %w", err, ErrInvalidBooking)
	}
	return nil
}

func (b Booking) GetStatus() BookingStatus {
	return b.Status
}

func (b Booking) CalculateDate() (time.Time, time.Time) {
	rentDate := time.Now()

	var addYear, addMonth, addDate int
	if b.Unit == DailyUnit {
		addDate = int(b.Qty)
	} else if b.Unit == MonthlyUnit {
		addMonth = int(b.Qty)
	} else if b.Unit == YearlyUnit {
		addYear = int(b.Qty)
	}
	dueDate := rentDate.AddDate(addYear, addMonth, addDate)

	return rentDate, dueDate
}

func (b Booking) GetTotalOverDueDays() int {
	return int(time.Now().Sub(b.DueDate).Hours() / 24)
}
