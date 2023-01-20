package price_model

import "time"

type Exclusive struct {
	priceModel

	KmQuota    int
	PricePerKm int

	TimeQuota  time.Duration
	PricePerHr int
}

func (e Exclusive) GetUpfront(qty uint, unit DateUnit) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (e Exclusive) GetExceedPrice(unit DateUnit, exceedTime time.Duration, distance uint) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (e Exclusive) Validate(partial bool) error {
	//TODO implement me
	panic("implement me")
}

func (e Exclusive) Status() Status {
	//TODO implement me
	panic("implement me")
}

func NewExclusive(
	KmQuota int,
	PricePerKm int,
	TimeQuota time.Duration,
	PricePerHr int,
) PriceModel {
	return Exclusive{}
}
