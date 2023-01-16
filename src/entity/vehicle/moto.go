package vehicle

type Moto struct {
	*Vehicle
}

func (m Moto) GetDeposit() Deposit {
	return m.RatePrice.Deposit / 100
}

func (m Moto) GetType() VType {
	return OrdinaryMotorcycleType
}

func NewOrdinaryMotorcycle(v Vehicle) *Moto {
	return &Moto{
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
