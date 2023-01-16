package vehicle

type MotoLux struct {
	*Vehicle
}

func (m MotoLux) GetDeposit() Deposit {
	return m.RatePrice.Deposit / 100
}

func (m MotoLux) GetType() VType {
	return LuxuryMotorcycleType
}

func NewLuxuryMotorcycle(v Vehicle) *MotoLux {
	return &MotoLux{
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
