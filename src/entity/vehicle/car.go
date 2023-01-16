package vehicle

type Car struct {
	*Vehicle
}

func (c Car) GetDeposit() Deposit {
	return c.RatePrice.Deposit / 100
}

func (c Car) GetType() Type {
	return OrdinaryCarType
}

func NewOrdinaryCar(v Vehicle) *Car {
	return &Car{
		Vehicle: &Vehicle{
			Brand:        v.Brand,
			Model:        v.Model,
			LicensePlate: v.LicensePlate,
			VehicleType:  v.VehicleType,
			Color:        v.Color,
			RatePrice:    v.RatePrice,
		},
	}
}
