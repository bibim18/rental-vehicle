package vehicle_repository

import (
	"context"
	"time"

	"rental-vehicle-system/src/entity/price_model"
	"rental-vehicle-system/src/use_case"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDb struct {
	col *mongo.Collection
}

type mongoVehicleFullStruct struct {
	Brand          string                `bson:"brand"`
	Model          string                `bson:"model"`
	LicensePlate   string                `bson:"license_plate"`
	VehicleType    price_model.VType     `bson:"price_model"`
	Color          string                `bson:"color"`
	Status         string                `bson:"status"`
	RatePrice      mongoVehicleRatePrice `bson:"rate_price"`
	RegisteredDate time.Time             `bson:"registered_date"`
}

type mongoVehicle struct {
	Brand        string                `bson:"brand"`
	Model        string                `bson:"model"`
	LicensePlate string                `bson:"license_plate"`
	VehicleType  price_model.VType     `bson:"price_model"`
	Color        string                `bson:"color"`
	RatePrice    mongoVehicleRatePrice `bson:"rate_price"`
}

type mongoVehicleRatePrice struct {
	Daily   int                 `bson:"daily"`
	Monthly int                 `bson:"monthly"`
	Yearly  int                 `bson:"yearly"`
	Deposit price_model.Deposit `bson:"deposit"`
}

func transferVehicleToMongoFullStruct(v use_case.VehicleFullDetail) mongoVehicleFullStruct {
	return mongoVehicleFullStruct{
		Brand:          v.Vehicle.Brand,
		Model:          v.Vehicle.Model,
		LicensePlate:   v.Vehicle.LicensePlate,
		VehicleType:    v.Vehicle.VehicleType,
		Color:          v.Vehicle.Color,
		Status:         v.Status,
		RegisteredDate: v.RegisteredDate,
		RatePrice:      mongoVehicleRatePrice(v.Vehicle.RatePrice),
	}
}

func transferVehicleToMongoStruct(v use_case.VehicleFullDetail) mongoVehicle {
	return mongoVehicle{
		Brand:        v.Vehicle.Brand,
		Model:        v.Vehicle.Model,
		LicensePlate: v.Vehicle.LicensePlate,
		VehicleType:  v.Vehicle.VehicleType,
		Color:        v.Vehicle.Color,
		RatePrice:    mongoVehicleRatePrice(v.Vehicle.RatePrice),
	}
}

func transferMongoToVehicleStruct(mVehicle mongoVehicleFullStruct) use_case.VehicleFullDetail {
	return use_case.VehicleFullDetail{
		Vehicle: price_model.vehicle{
			Brand:        mVehicle.Brand,
			Model:        mVehicle.Model,
			LicensePlate: mVehicle.LicensePlate,
			VehicleType:  mVehicle.VehicleType,
			Color:        mVehicle.Color,
			RatePrice:    price_model.Pricing(mVehicle.RatePrice),
		},
		Status: mVehicle.Status,
	}
}

func (m mongoDb) CreateVehicle(ctx context.Context, vehicle use_case.VehicleFullDetail) error {
	log.Info("vehicle_repo.CreateVehicle")
	prepareToMongo := transferVehicleToMongoFullStruct(vehicle)
	_, err := m.col.InsertOne(ctx, prepareToMongo)

	return err
}

func (m mongoDb) UpdateVehicleStatus(ctx context.Context, vehicleId string, status string) error {
	log.Info("vehicle_repo.UpdateVehicleStatus")
	objectId, err := primitive.ObjectIDFromHex(vehicleId)
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
		return errors.Errorf("update price_model status at _id : %s not found: %w", vehicleId, price_model.ErrVehicleNotFound)
	}

	return err
}

func (m mongoDb) UpdateVehicle(ctx context.Context, vehicleId string, v use_case.VehicleFullDetail) error {
	log.Info("vehicle_repo.UpdateVehicle")
	objectId, err := primitive.ObjectIDFromHex(vehicleId)
	if err != nil {
		return err
	}

	newVehicleData := transferVehicleToMongoStruct(v)
	res, err := m.col.UpdateOne(ctx, bson.M{"_id": objectId}, bson.M{
		"$set": newVehicleData,
	})
	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		return errors.Errorf("update price_model at _id : %s not found: %w", vehicleId, price_model.ErrVehicleNotFound)
	}

	return err
}

func (m mongoDb) GetVehicleByStatus(ctx context.Context, status string, limit, offset int64) ([]use_case.VehicleFullDetail, error) {
	log.Info("vehicle_repo.GetVehicleByStatus, Query status is", status)
	opts := options.Find()
	opts.SetSkip(offset)
	opts.SetLimit(limit)
	var filter primitive.M

	if status != "" {
		filter = bson.M{"status": status}
	} else {
		filter = bson.M{"status": bson.M{"$ne": "unactive"}}
	}

	cur, err := m.col.Find(ctx, filter)

	var vehicles []mongoVehicleFullStruct
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var v mongoVehicleFullStruct
		err := cur.Decode(&v)
		if err != nil {
			return []use_case.VehicleFullDetail{}, err
		}
		vehicles = append(vehicles, v)
	}

	if err != nil {
		return []use_case.VehicleFullDetail{}, err
	}

	vehiclesWithFullDetail := make([]use_case.VehicleFullDetail, len(vehicles))
	for i, v := range vehicles {
		vehiclesWithFullDetail[i] = transferMongoToVehicleStruct(v)
	}
	return vehiclesWithFullDetail, nil
}

func (m mongoDb) GetVehicleById(ctx context.Context, vehicleId string) (use_case.VehicleFullDetail, error) {
	log.Info("vehicle_repo.GetVehicleById, vehicleId is ", vehicleId)
	objectId, err := primitive.ObjectIDFromHex(vehicleId)
	if err != nil {
		return use_case.VehicleFullDetail{}, err
	}

	cur := m.col.FindOne(ctx, bson.M{"_id": objectId})

	var v mongoVehicleFullStruct
	errDocode := cur.Decode(&v)
	if errDocode != nil {
		return use_case.VehicleFullDetail{}, errDocode
	}

	vFullDetail := transferMongoToVehicleStruct(v)
	return vFullDetail, nil
}

func NewMongoDb(db *mongo.Database) use_case.VehicleRepository {
	m := &mongoDb{col: db.Collection("vehicles")}

	return m
}
