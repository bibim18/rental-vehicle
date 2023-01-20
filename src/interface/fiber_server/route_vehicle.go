package fiber_server

import (
	"github.com/gofiber/fiber/v2"
	"rental-vehicle-system/src/entity/price_model"
	"rental-vehicle-system/src/use_case"
)

func (f FiberServer) addRouteVehicle(base fiber.Router) {
	r := base.Group("/vehicles")

	r.Post("/", f.createVehicle)

	r.Patch("/:vehicleId", f.updateVehicle)
	r.Patch("/active/:vehicleId", f.updateVehicleStatus)

	r.Get("/:status?", f.listVehicle)
	r.Delete("/:vehicleId", f.archiveVehicle)
}

func (cv createVehicleRequest) transferRequestToUsecase() (use_case.VehicleWithPrice, error) {
	vehicles := make([]price_model.vehicle, len(cv.Vehicles))
	for i, v := range cv.Vehicles {
		vehicles[i] = price_model.vehicle{
			Brand:        v.Brand,
			Model:        v.Model,
			Color:        v.Color,
			LicensePlate: v.LicensePlate,
			VehicleType:  v.VehicleType,
		}
	}

	return use_case.VehicleWithPrice{
		Vehicles: vehicles,
		RatePrice: price_model.Pricing{
			Daily:   cv.RatePrice.Daily,
			Monthly: cv.RatePrice.Monthly,
			Yearly:  cv.RatePrice.Yearly,
			Deposit: cv.RatePrice.Deposit,
		},
	}, nil
}

func (uv updateVehicleRequest) transferRequestToUsecase() (use_case.VehicleFullDetail, error) {
	return use_case.VehicleFullDetail{
		Vehicle: price_model.vehicle{
			Brand:        uv.Vehicle.Brand,
			Model:        uv.Vehicle.Model,
			LicensePlate: uv.Vehicle.LicensePlate,
			VehicleType:  uv.Vehicle.VehicleType,
			Color:        uv.Vehicle.Color,
			RatePrice:    price_model.Pricing(uv.RatePrice),
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
