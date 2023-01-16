package vehicle

type CarLux struct {
	*Vehicle
}

func (c CarLux) GetDeposit() Deposit {
	return c.RatePrice.Deposit / 100
}

func (c CarLux) GetType() VType {
	return LuxuryCarType
}

func NewLuxuryCar(v Vehicle) *CarLux {
	return &CarLux{
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
