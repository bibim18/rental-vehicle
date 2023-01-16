package booking_repository

import (
	"context"
	"rental-vehicle-system/src/entity/booking"
	"rental-vehicle-system/src/entity/customer"
	"rental-vehicle-system/src/use_case"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoDb struct {
	col *mongo.Collection
}

type mongoBooking struct {
	VehicleId  string                 `bson:"vehicleId"`
	Unit       booking.UnitType       `bson:"unit"`
	Qty        int                    `bson:"qty"`
	Status     booking.BookingStatus  `bson:"status"`
	TotalPrice booking.TotalPriceType `bson:"total_price"`
	Deposit    booking.DepositType    `bson:"deposit"`
	RentDate   time.Time              `bson:"rent_date"`
	RentPrice  int                    `bson:"rent_price"`
	DueDate    time.Time              `bson:"due_date"`
	Customer   mongoBookingCustomer   `bson:"customer"`
}

type mongoBookingCustomer struct {
	Name        string `bson:"name"`
	Lastname    string `bson:"last_name"`
	PhoneNumber string `bson:"phone_number"`
	Email       string `bson:"email"`
}

func transferBookingToMongoStruct(b use_case.BookingFullDetail) mongoBooking {
	return mongoBooking{
		VehicleId:  b.Booking.VehicleId,
		Unit:       b.Booking.Unit,
		Qty:        b.Booking.Qty,
		Status:     b.Booking.Status,
		RentPrice:  b.Booking.RentPrice,
		TotalPrice: b.Booking.TotalPrice,
		Deposit:    b.Booking.Deposit,
		RentDate:   b.Booking.RentDate,
		DueDate:    b.Booking.DueDate,
		Customer:   mongoBookingCustomer(b.Customer),
	}
}

func (m mongoDb) CreateBooking(ctx context.Context, b use_case.BookingFullDetail) error {
	log.Info("booking_repo.CreateBooking")
	prepareToMongo := transferBookingToMongoStruct(b)
	_, err := m.col.InsertOne(ctx, prepareToMongo)

	return err
}

func (m mongoDb) GetBookingById(ctx context.Context, bId string) (use_case.BookingFullDetail, error) {
	log.Info("booking_repo.GetBoookingById, bookingId is ", bId)
	objectId, err := primitive.ObjectIDFromHex(bId)
	if err != nil {
		return use_case.BookingFullDetail{}, err
	}

	cur := m.col.FindOne(ctx, bson.M{"_id": objectId})

	var b mongoBooking
	errDocode := cur.Decode(&b)
	if errDocode != nil {
		return use_case.BookingFullDetail{}, errDocode
	}

	return use_case.BookingFullDetail{
		Booking: booking.Booking{
			VehicleId: b.VehicleId,
			Unit:      b.Unit,
			Qty:       b.Qty,
			Status:    b.Status,
			DueDate:   b.DueDate,
			RentDate:  b.RentDate,
		},
		Customer: customer.Customer{
			Name:        b.Customer.Name,
			Lastname:    b.Customer.Lastname,
			PhoneNumber: b.Customer.PhoneNumber,
			Email:       b.Customer.Email,
		},
	}, nil
}

func (m mongoDb) UpdateReturnedBooking(ctx context.Context, bId string, b booking.Booking) error {
	log.Info("booking_repo.UpdateReturnedBooking")
	objectId, err := primitive.ObjectIDFromHex(bId)
	if err != nil {
		return err
	}
	res, err := m.col.UpdateOne(ctx, bson.M{"_id": objectId}, bson.M{
		"$set": bson.M{
			"status":      b.Status,
			"late_fines":  b.LateFines,
			"return_date": b.ReturnDate,
		},
	})
	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		return errors.Errorf("update booking at _id : %s not found: %w", bId, booking.ErrBookingNotFound)
	}

	return err
}

func (m mongoDb) UpdateBookingStatus(ctx context.Context, bId string, status booking.BookingStatus) error {
	log.Info("booking_repo.UpdateBookingStatus")
	objectId, err := primitive.ObjectIDFromHex(bId)
	if err != nil {
		return err
	}
	res, err := m.col.UpdateOne(ctx, bson.M{"_id": objectId}, bson.M{
		"$set": bson.M{
			"status": status,
		},
	})
	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		return errors.Errorf("update booking at _id : %s not found: %w", bId, booking.ErrBookingNotFound)
	}

	return err
}

func NewMongoDb(db *mongo.Database) use_case.BookingRepository {
	m := &mongoDb{col: db.Collection("bookings")}

	return m
}
