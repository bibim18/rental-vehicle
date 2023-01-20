package price_model

import (
	"testing"
	"time"
)

func TestCar_Validate(t *testing.T) {
	tests := []struct {
		name    string
		car     Basic
		partial bool
		wantErr bool
	}{
		{
			name: "valid car",
			car: Basic{
				priceModel: priceModel{
					Id: "1",
					//LicensePlate: "B 1234 ABC",
					//Brand:        "Toyota",
					//Model:        "Yaris",
					//Color:        "White",
					Status: ActiveStatus,
				},
				Upfront:       10000,
				PricePerDay:   10000,
				PricePerMonth: 100000,
				PricePerKm:    1000,
			},
			partial: false,
			wantErr: false,
		},
		{
			name: "valid car with empty upfront",
			car: Basic{
				priceModel: priceModel{
					Id: "1",
				},
			},
			partial: true,
			wantErr: false,
		},
		{
			name: "invalid car (no brand)",
			car: Basic{
				priceModel: priceModel{
					Id: "1",
					//LicensePlate: "B 1234 ABC",
					//Model:        "Yaris",
					//Color:        "White",
					Status: ActiveStatus,
				},
				Upfront:       10000,
				PricePerDay:   10000,
				PricePerMonth: 100000,
				PricePerKm:    1000,
			},
			partial: false,
			wantErr: true,
		},
		{
			name: "valid car one of date type",
			car: Basic{
				priceModel: priceModel{
					Id: "1",
					//LicensePlate: "B 1234 ABC",
					//Brand:        "Toyota",
					//Model:        "Yaris",
					//Color:        "White",
					Status: ActiveStatus,
				},
				Upfront:       10000,
				PricePerMonth: 100000,
				PricePerKm:    1000,
			},
			partial: false,
			wantErr: false,
		},
		{
			name: "valid car one of date type",
			car: Basic{
				priceModel: priceModel{
					Id: "1",
					//LicensePlate: "B 1234 ABC",
					//Brand:        "Toyota",
					//Model:        "Yaris",
					//Color:        "White",
					Status: ActiveStatus,
				},
				Upfront:     10000,
				PricePerDay: 100000,
				PricePerKm:  1000,
			},
			partial: false,
			wantErr: false,
		},
		{
			name: "invalid car no date type",
			car: Basic{
				priceModel: priceModel{
					Id: "1",
					//LicensePlate: "B 1234 ABC",
					//Brand:        "Toyota",
					//Model:        "Yaris",
					//Color:        "White",
					Status: ActiveStatus,
				},
				Upfront:    10000,
				PricePerKm: 1000,
			},
			partial: false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.car.Validate(tt.partial); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCar_GetUpfront(t *testing.T) {
	type args struct {
		qty      uint
		unit     DateUnit
		distance uint
	}
	tests := []struct {
		name    string
		car     Basic
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "simple",
			car: Basic{
				priceModel:    priceModel{},
				Upfront:       10000,
				PricePerDay:   10000,
				PricePerMonth: 100000,
				PricePerKm:    1000,
			},
			args: args{
				qty:      3,
				unit:     dailyUnit,
				distance: 100,
			},
			want:    40000,
			wantErr: false,
		},
		{
			name: "cannot rent daily",
			car: Basic{
				priceModel:    priceModel{},
				Upfront:       10000,
				PricePerMonth: 100000,
				PricePerKm:    1000,
			},
			args: args{
				qty:      3,
				unit:     dailyUnit,
				distance: 100,
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.car.GetUpfront(tt.args.qty, tt.args.unit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetExceedPrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetExceedPrice() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCar_GetExceedPrice(t *testing.T) {
	type args struct {
		unit       DateUnit
		exceedTime time.Duration
		distance   uint
	}
	tests := []struct {
		name    string
		car     Basic
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "simple",
			car: Basic{
				priceModel:    priceModel{},
				Upfront:       10000,
				PricePerDay:   10000,
				PricePerMonth: 100000,
				PricePerKm:    1000,
			},
			args: args{
				unit:       dailyUnit,
				exceedTime: 24 * time.Hour,
				distance:   100,
			},
			want: 110000,
		},
		{
			name: "simple",
			car: Basic{
				priceModel:    priceModel{},
				Upfront:       10000,
				PricePerDay:   10000,
				PricePerMonth: 100000,
				PricePerKm:    1000,
			},
			args: args{
				unit:       dailyUnit,
				exceedTime: 0,
				distance:   10,
			},
			want: 10000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.car.GetExceedPrice(tt.args.unit, tt.args.exceedTime, tt.args.distance)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetExceedPrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetExceedPrice() got = %v, want %v", got, tt.want)
			}
		})
	}
}
