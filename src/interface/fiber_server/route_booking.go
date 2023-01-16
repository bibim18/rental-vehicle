package fiber_server

import (
	"rental-vehicle-system/src/entity/booking"
	"rental-vehicle-system/src/entity/customer"
	"rental-vehicle-system/src/use_case"

	"github.com/gofiber/fiber/v2"
)

type bookingVehicleRequest struct {
	VehicleId string                     `json:"vehicle_id"`
	Unit      booking.UnitType           `json:"unit"`
	Qty       int                        `json:"qty"`
	Customer  bookingVehicleCustomerData `json:"customer"`
}

type bookingVehicleCustomerData struct {
	Name        string `json:"name"`
	Lastname    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}

func (f FiberServer) addRouteBooking(base fiber.Router) {
	r := base.Group("/bookings")

	r.Post("/", f.bookingVehicle)
	r.Patch("/return-vehicle/:bookingId", f.returnVehicle)
	r.Patch("/cancel/:bookingId", f.cancelBooking)

}

func (b bookingVehicleRequest) transferRequestToUsecase() use_case.BookingDetail {
	return use_case.BookingDetail{
		Booking: booking.Booking{VehicleId: b.VehicleId,
			Unit: b.Unit,
			Qty:  b.Qty},
		Customer: customer.Customer(b.Customer),
	}
}

func (f FiberServer) bookingVehicle(c *fiber.Ctx) error {
	var b bookingVehicleRequest
	if err := c.BodyParser(&b); err != nil {
		return f.errorHandler(c, err)
	}

	usecaseBooking := b.transferRequestToUsecase()
	msg, err := f.useCase.BookingVehicle(getSpanContext(c), usecaseBooking)
	if err != nil {
		return f.errorHandler(c, err)
	}

	return c.SendString(msg)
}

func (f FiberServer) returnVehicle(c *fiber.Ctx) error {
	bId := c.Params("bookingId")

	msg, err := f.useCase.ReturnVehicle(getSpanContext(c), bId)
	if err != nil {
		return f.errorHandler(c, err)
	}

	return c.SendString(msg)
}

func (f FiberServer) cancelBooking(c *fiber.Ctx) error {
	bId := c.Params("bookingId")

	msg, err := f.useCase.CancelationBooking(getSpanContext(c), bId)
	if err != nil {
		return f.errorHandler(c, err)
	}

	return c.SendString(msg)
}
