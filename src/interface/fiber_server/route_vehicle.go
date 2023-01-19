package fiber_server

import (
	"rental-vehicle-system/src/entity/vehicle"
	"rental-vehicle-system/src/use_case"

	"github.com/gofiber/fiber/v2"
)

type createVehicleRequest struct {
	Vehicles  []vehicleDataRequest `json:"vehicles"`
	RatePrice ratePriceDataRequest `json:"rate_price"`
}

type updateVehicleRequest struct {
	Vehicle   vehicleDataRequest   `json:"vehicle"`
	RatePrice ratePriceDataRequest `json:"rate_price"`
}

type vehicleDataRequest struct {
	Brand        string        `json:"brand"`
	Model        string        `json:"model"`
	LicensePlate string        `json:"license_plate"`
	VehicleType  vehicle.VType `json:"vehicle_type"`
	Color        string        `json:"color"`
}

type ratePriceDataRequest struct {
	Daily   int             `json:"daily"`
	Monthly int             `json:"monthly"`
	Yearly  int             `json:"yearly"`
	Deposit vehicle.Deposit `json:"deposit"`
}

type ResponseData struct {
	Total int
	Data  []use_case.VehicleFullDetail
}

func (f FiberServer) addRouteVehicle(base fiber.Router) {
	r := base.Group("/vehicles")

	r.Post("/", f.createVehicle)

	r.Patch("/:vehicleId", f.updateVehicle)
	r.Patch("/active/:vehicleId", f.updateVehicleStatus)

	r.Get("/:status?", f.listVehicle)
	r.Delete("/:vehicleId", f.archiveVehicle)
}

func (cv createVehicleRequest) transferRequestToUsecase() (use_case.VehicleWithPrice, error) {
	vehicles := make([]vehicle.Vehicle, len(cv.Vehicles))
	for i, v := range cv.Vehicles {
		vehicles[i] = vehicle.Vehicle{
			Brand:        v.Brand,
			Model:        v.Model,
			Color:        v.Color,
			LicensePlate: v.LicensePlate,
			VehicleType:  v.VehicleType,
		}
	}

	return use_case.VehicleWithPrice{
		Vehicles: vehicles,
		RatePrice: vehicle.Price{
			Daily:   cv.RatePrice.Daily,
			Monthly: cv.RatePrice.Monthly,
			Yearly:  cv.RatePrice.Yearly,
			Deposit: cv.RatePrice.Deposit,
		},
	}, nil
}

func (uv updateVehicleRequest) transferRequestToUsecase() (use_case.VehicleFullDetail, error) {
	return use_case.VehicleFullDetail{
		Vehicle: vehicle.Vehicle{
			Brand:        uv.Vehicle.Brand,
			Model:        uv.Vehicle.Model,
			LicensePlate: uv.Vehicle.LicensePlate,
			VehicleType:  uv.Vehicle.VehicleType,
			Color:        uv.Vehicle.Color,
			RatePrice:    vehicle.Price(uv.RatePrice),
		},
	}, nil
}

func (f FiberServer) createVehicle(c *fiber.Ctx) error {
	var v createVehicleRequest
	if err := c.BodyParser(&v); err != nil {
		return f.errorHandler(c, err)
	}

	usecaseVehicle, err := v.transferRequestToUsecase()
	if err != nil {
		return f.errorHandler(c, err)
	}

	msg, err := f.useCase.CreateVehicle(getSpanContext(c), usecaseVehicle)

	if err != nil {
		return f.errorHandler(c, err)
	}

	return c.SendString(msg)
}

func (f FiberServer) listVehicle(c *fiber.Ctx) error {
	status := c.Params("status")

	list, err := f.useCase.GetVehicle(getSpanContext(c), status)
	if err != nil {
		return f.errorHandler(c, err)
	}

	response := ResponseData{Data: list, Total: len(list)}
	return c.JSON(response)
}

func (f FiberServer) updateVehicle(c *fiber.Ctx) error {
	vId := c.Params("vehicleId")
	var v updateVehicleRequest
	if err := c.BodyParser(&v); err != nil {
		return f.errorHandler(c, err)
	}

	usecaseVehicle, err := v.transferRequestToUsecase()
	msg, err := f.useCase.UpdateVehicleDetail(getSpanContext(c), vId, usecaseVehicle)
	if err != nil {
		return f.errorHandler(c, err)
	}

	return c.SendString(msg)
}

func (f FiberServer) updateVehicleStatus(c *fiber.Ctx) error {
	id := c.Params("vehicleId")
	msg, err := f.useCase.UpdateVehicleStatus(getSpanContext(c), id)

	if err != nil {
		return f.errorHandler(c, err)
	}

	return c.SendString(msg)
}

func (f FiberServer) archiveVehicle(c *fiber.Ctx) error {
	id := c.Params("vehicleId")
	msg, err := f.useCase.ArchiveVehicle(getSpanContext(c), id)

	if err != nil {
		return f.errorHandler(c, err)
	}

	return c.SendString(msg)
}
