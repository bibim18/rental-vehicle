package price_model

import (
	"github.com/cockroachdb/errors"
	"time"
)

type Basic struct {
	priceModel

	Upfront       int
	PricePerDay   int
	PricePerMonth int
	PricePerKm    int

	// KmQuota int
	// PerDayStack int // 100 เช่า 1วันราคา 100, 2วันราคา 100+200, 3วันราคา 100+200+300...
	// KmStack int

}

func (c Basic) GetUpfront(qty uint, unit DateUnit) (int, error) {
	bp, err := c.getBasePrice(qty, unit)
	if err != nil {
		return 0, err
	}
	return c.Upfront + bp, nil
}

func (c Basic) GetExceedPrice(unit DateUnit, exceedTime time.Duration, distance uint) (int, error) {
	distancePrice := c.PricePerKm * int(distance)
	if exceedTime == 0 {
		return distancePrice, nil
	}

	bp, err := c.getBasePrice(1, unit)
	if err != nil {
		return 0, err
	}

	priceDuration := unit.Duration()

	return distancePrice + int(exceedTime/priceDuration)*bp, nil
}

func (c Basic) Status() Status {
	return c.priceModel.Status
}

func (c Basic) Validate(partial bool) error {
	if partial {
		if c.Id == "" {
			return errors.New("Id cannot be empty")
		}

		return nil
	}

	//if c.LicensePlate == "" {
	//	return errors.New("LicensePlate cannot be empty")
	//}
	//
	//if c.Brand == "" {
	//	return errors.New("Brand cannot be empty")
	//}
	//
	//if c.Model == "" {
	//	return errors.New("Model cannot be empty")
	//}
	//
	//if c.Color == "" {
	//	return errors.New("Color cannot be empty")
	//}

	if c.PricePerDay == 0 && c.PricePerMonth == 0 {
		return errors.New("PricePerDay and PricePerMonth cannot be empty")
	}

	if c.PricePerKm == 0 {
		return errors.New("PricePerKm cannot be empty")
	}

	return nil
}

func (c Basic) getBasePrice(qty uint, unit DateUnit) (int, error) {
	switch unit {
	case dailyUnit:
		if c.PricePerDay == 0 {
			return 0, errors.New("PricePerDay is empty")
		}

		return c.PricePerDay * int(qty), nil
	case monthlyUnit:
		if c.PricePerMonth == 0 {
			return 0, errors.New("PricePerMonth is empty")
		}

		return c.PricePerMonth * int(qty), nil
	case yearlyUnit:
		return 0, errors.New("yearly pricing not implemented")
	}

	return 0, errors.New("invalid unit")
}

func NewBasicPriceModel(
	LicensePlate string,
	Brand string,
	Model string,
	Color string,
	Deposit int,
	PricePerDay int,
	PricePerMonth int,
	PricePerKm int,
) PriceModel {
	return &Basic{
		priceModel: priceModel{
			//LicensePlate: LicensePlate,
			//Brand:        Brand,
			//Model:        Model,
			//Color:        Color,
			Status: ActiveStatus,
		},
		Upfront:       Deposit,
		PricePerDay:   PricePerDay,
		PricePerMonth: PricePerMonth,
		PricePerKm:    PricePerKm,
	}
}
