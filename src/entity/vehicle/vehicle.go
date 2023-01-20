package vehicle

import (
	"github.com/cockroachdb/errors"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

var (
	ErrInvalidVehicle    = errors.New("Invalid vehicle")
	ErrVehicleNotFound   = errors.New("Cannot found vehicle")
	ErrVehicleNotArchive = errors.New("Vehicle cannot disbled")
	ErrVehicleNotUpdate  = errors.New("Vehicle cannot update detail")
	ErrVehicleNotEnable  = errors.New("Vehicle already enabled")
)

type Deposit float32
type Status string

type Vehicle struct {
	Brand        string `validate:"required,min=3"`
	Model        string `validate:"required,min=3"`
	LicensePlate string `validate:"required,min=3"`
	VehicleType  VType  `validate:"oneof='car-lux' 'car' 'moto' 'moto-lux'"`
	Color        string
	RatePrice    Price
}

type Price struct {
	Daily   int     `validate:"required,gt=0"`
	Monthly int     `validate:"required,gt=0"`
	Yearly  int     `validate:"required,gt=0"`
	Deposit Deposit `validate:"required,min=1,max=100"`
}

type VType string
type UnitType string

const (
	dailyUnit   string = "daily"
	monthlyUnit string = "monthly"
	yearlyUnit  string = "yearly"
)

const (
	LuxuryCarType          VType = "car-lux"
	OrdinaryCarType        VType = "car"
	LuxuryMotorcycleType   VType = "moto-lux"
	OrdinaryMotorcycleType VType = "moto"
)

const (
	ActiveStatus   string = "active"
	InuseStatus    string = "inuse"
	UnactiveStatus string = "unactive"
)

type VehicleMethod interface {
	GetType() VType
	GetDeposit() Deposit
}

func (v Vehicle) Validate() error {
	err := validator.New().Struct(v)
	if err != nil {
		return errors.Errorf("%s: %w", err, ErrInvalidVehicle)
	}
	return nil
}

func (v Vehicle) GetPrice(unit UnitType) int {
	if string(unit) == dailyUnit {
		return v.RatePrice.Daily
	}
	if string(unit) == monthlyUnit {
		return v.RatePrice.Monthly
	}
	if string(unit) == yearlyUnit {
		return v.RatePrice.Yearly
	}
	return 0
}

func VerifyStatus(status string) error {
	log.Info("vehicle_entity.VerifyStatus", status)
	if status == ActiveStatus || status == InuseStatus || status == UnactiveStatus {
		return nil
	}
	return errors.Errorf("vehicle status '%s' is invalid", status)
}

func New(vType VType, v Vehicle) VehicleMethod {
	vehicleClass := map[VType]VehicleMethod{
		LuxuryCarType:          NewLuxuryCar(v),
		OrdinaryCarType:        NewOrdinaryCar(v),
		LuxuryMotorcycleType:   NewLuxuryMotorcycle(v),
		OrdinaryMotorcycleType: NewOrdinaryMotorcycle(v),
	}
	return vehicleClass[vType]
}
