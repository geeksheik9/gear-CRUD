package db

import (
	"context"

	model "github.com/geeksheik9/gear-CRUD/models"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
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

// InsertWeapon is the database implementation to insert a weapon object
func (g *GearDB) InsertWeapon(weapon *model.Weapon) error {
	logrus.Debug("BEGIN - InsertWeapon")

	collection := g.client.Database(g.databaseName).Collection(g.weaponCollection)

	_, err := collection.InsertOne(context.Background(), weapon)

	return err
}
