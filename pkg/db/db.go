package db

import (
	"context"
	"errors"
	"net/url"
	"strconv"
	"time"

	model "github.com/geeksheik9/gear-CRUD/models"
	"github.com/geeksheik9/gear-CRUD/pkg/api"
	"github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// GearDB is the data access object for star wars FFG weapons and armor
type GearDB struct {
	client           *mongo.Client
	databaseName     string
	armorCollection  string
	weaponCollection string
}

//Ping checks that the database is running
func (g *GearDB) Ping() error {
	err := g.client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		logrus.Errorf("ERROR connectiong to database %v", err)
	}
	return err
}

// InsertArmor is the database implementation to insert an armor object
func (g *GearDB) InsertArmor(armor *model.Armor) error {
	logrus.Debug("BEGIN - InsertArmor")

	collection := g.client.Database(g.databaseName).Collection(g.armorCollection)

	_, err := collection.InsertOne(context.Background(), armor)

	return err
}

//GetArmor is the database implementation to get all armor objects
func (g *GearDB) GetArmor(queryParams url.Values) ([]model.Armor, error) {
	logrus.Debug("BEGIN - GetArmor")

	collection := g.client.Database(g.databaseName).Collection(g.armorCollection)

	pageNumber, pageCount, sort, filter := api.BuildFilter(queryParams)
	skip := 0
	if pageNumber > 0 {
		skip = (pageNumber - 1) * pageCount
	}

	opts := options.Find().
		SetMaxTime(30 * time.Second).
		SetSkip(int64(skip)).
		SetLimit(int64(pageCount)).
		SetSort(bson.D{{
			Key:   sort,
			Value: 1,
		}})

	cur, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}

	matches := []model.Armor{}

	for cur.Next(context.Background()) {
		elem := model.Armor{}
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		matches = append(matches, elem)
	}

	return matches, nil
}

//GetArmorByID is the database implementation to get a pspecific armor back from the database
func (g *GearDB) GetArmorByID(mongoID primitive.ObjectID) (*model.Armor, error) {
	logrus.Debugf("BEGIN - GetArmorByID: %v", mongoID)

	collection := g.client.Database(g.databaseName).Collection(g.armorCollection)
	query := api.BuildQuery(&mongoID, nil)

	armor := model.Armor{}

	err := collection.FindOne(context.Background(), query).Decode(&armor)
	if err != nil {
		return nil, err
	}

	return &armor, err
}

//UpdateArmorByID updates a specific armor in the armor database
func (g *GearDB) UpdateArmorByID(armor model.Armor, mongoID primitive.ObjectID) error {

	collection := g.client.Database(g.databaseName).Collection(g.armorCollection)

	result, err := collection.UpdateOne(context.Background(), bson.M{"_id": mongoID}, bson.D{{
		Key:   "$set",
		Value: armor,
	}})
	if err != nil {
		return err
	}

	matched := strconv.FormatInt(result.MatchedCount, 10)
	modified := strconv.FormatInt(result.ModifiedCount, 10)

	if result.MatchedCount != 1 {
		return errors.New("Could not update sheet. Tried to update " + mongoID.Hex() + " got " + matched + " matches instead of 1")
	}

	if result.ModifiedCount != 1 {
		return errors.New("Could not update sheet. Tried to updated " + mongoID.Hex() + " tried to update " + modified + " number of results instead of 1")
	}

	return nil
}

//DeleteArmorByID deletes a specific armor from the database
func (g *GearDB) DeleteArmorByID(mongoID primitive.ObjectID) error {
	logrus.Debugf("BEGIN - DeleteForceCharacterSheetByID: %v", mongoID)

	collection := g.client.Database(g.databaseName).Collection(g.armorCollection)

	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": mongoID})

	return err
}

// InsertWeapon is the database implementation to insert a weapon object
func (g *GearDB) InsertWeapon(weapon *model.Weapon) error {
	logrus.Debug("BEGIN - InsertWeapon")

	collection := g.client.Database(g.databaseName).Collection(g.weaponCollection)

	_, err := collection.InsertOne(context.Background(), weapon)

	return err
}

//GetWeapon is the database implementation to get all armor objects
func (g *GearDB) GetWeapon(queryParams url.Values) ([]model.Weapon, error) {
	logrus.Debug("BEGIN - GetWeapon")

	collection := g.client.Database(g.databaseName).Collection(g.weaponCollection)

	pageNumber, pageCount, sort, filter := api.BuildFilter(queryParams)

	skip := 0
	if pageNumber > 0 {
		skip = (pageNumber - 1) * pageCount
	}

	opts := options.Find().
		SetMaxTime(30 * time.Second).
		SetSkip(int64(skip)).
		SetLimit(int64(pageCount)).
		SetSort(bson.D{{
			Key:   sort,
			Value: 1,
		}})

	cur, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}

	matches := []model.Weapon{}

	for cur.Next(context.Background()) {
		elem := model.Weapon{}
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		matches = append(matches, elem)
	}

	return matches, nil
}

//GetWeaponByID is the database implementation to get a specific weapon back from the database
func (g *GearDB) GetWeaponByID(mongoID primitive.ObjectID) (*model.Weapon, error) {
	logrus.Debugf("BEGIN - GetWeaponByID: %v", mongoID)

	collection := g.client.Database(g.databaseName).Collection(g.weaponCollection)
	query := api.BuildQuery(&mongoID, nil)

	weapon := model.Weapon{}

	err := collection.FindOne(context.Background(), query).Decode(&weapon)
	if err != nil {
		return nil, err
	}

	return &weapon, err
}

//UpdateWeaponByID updates a specific weapon in the weapon database
func (g *GearDB) UpdateWeaponByID(weapon model.Weapon, mongoID primitive.ObjectID) error {

	collection := g.client.Database(g.databaseName).Collection(g.weaponCollection)

	result, err := collection.UpdateOne(context.Background(), bson.M{"_id": mongoID}, bson.D{{
		Key:   "$set",
		Value: weapon,
	}})
	if err != nil {
		return err
	}

	matched := strconv.FormatInt(result.MatchedCount, 10)
	modified := strconv.FormatInt(result.ModifiedCount, 10)

	if result.MatchedCount != 1 {
		return errors.New("Could not update sheet. Tried to update " + mongoID.Hex() + " got " + matched + " matches instead of 1")
	}

	if result.ModifiedCount != 1 {
		return errors.New("Could not update sheet. Tried to updated " + mongoID.Hex() + " tried to update " + modified + " number of results instead of 1")
	}

	return nil
}

//DeleteWeaponByID deletes a specific weapon from the database
func (g *GearDB) DeleteWeaponByID(mongoID primitive.ObjectID) error {
	logrus.Debugf("BEGIN - DeleteForceCharacterSheetByID: %v", mongoID)

	collection := g.client.Database(g.databaseName).Collection(g.weaponCollection)

	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": mongoID})

	return err
}
