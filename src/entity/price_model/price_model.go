package price_model

import (
	"github.com/cockroachdb/errors"
	"time"
)

var (
	ErrInvalidVehicle    = errors.New("Invalid price_model")
	ErrVehicleNotFound   = errors.New("Cannot found price_model")
	ErrVehicleNotArchive = errors.New("price_model cannot disbled")
	ErrVehicleNotUpdate  = errors.New("price_model cannot update detail")
	ErrVehicleNotEnable  = errors.New("price_model already enabled")
)

type Status string
type DateUnit string

type priceModel struct {
	Id string
	//LicensePlate string
	//Brand        string `validate:"required,min=3"`
	//Model        string `validate:"required,min=3"`
	//Color        string
	Status Status
}

const (
	dailyUnit   DateUnit = "daily"
	weeklyUnit  DateUnit = "weekly"
	monthlyUnit DateUnit = "monthly"
	yearlyUnit  DateUnit = "yearly"
)

const (
	ActiveStatus   Status = "active"
	InuseStatus    Status = "inuse"
	InactiveStatus Status = "inactive"
)

type PriceModel interface {
	GetUpfront(qty uint, unit DateUnit) (int, error)
	GetExceedPrice(unit DateUnit, exceedTime time.Duration, distance uint) (int, error)
	Validate(partial bool) error
	Status() Status
}

func (u DateUnit) Duration() time.Duration {
	switch u {
	case dailyUnit:
		return time.Hour * 24
	case weeklyUnit:
		return time.Hour * 24 * 7
	case monthlyUnit:
		return time.Hour * 24 * 30
	case yearlyUnit:
		return time.Hour * 24 * 365
	default:
		return 0
	}
}

func NewStatus(status string) (Status, error) {
	switch status {
	case "active":
		return ActiveStatus, nil
	case "inuse":
		return InuseStatus, nil
	case "inactive":
		return InactiveStatus, nil
	default:
		return "", errors.Errorf("price_model status '%s' is invalid", status)
	}
}

func NewDateType(unit string) (DateUnit, error) {
	switch unit {
	case "daily":
		return dailyUnit, nil
	case "weekly":
		return weeklyUnit, nil
	case "monthly":
		return monthlyUnit, nil
	case "yearly":
		return yearlyUnit, nil
	default:
		return "", errors.Errorf("unit type '%s' is invalid", unit)
	}
}
