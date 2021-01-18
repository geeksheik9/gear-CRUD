package mocks

import (
	"net/url"

	model "github.com/geeksheik9/gear-CRUD/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//MockGearDatabase is a mock struct for testing
type MockGearDatabase struct {
	ArmorToReturn   *model.Armor
	ArmorsToReturn  []model.Armor
	WeaponToReturn  *model.Weapon
	WeaponsToReturn []model.Weapon
	ErrorToReturn   error
}

//InsertArmor is the mock method for testing
func (db *MockGearDatabase) InsertArmor(armor *model.Armor) error {
	return db.ErrorToReturn
}

//GetArmor is the mock method for testing
func (db *MockGearDatabase) GetArmor(query url.Values) ([]model.Armor, error) {
	return db.ArmorsToReturn, db.ErrorToReturn
}

//GetArmorByID is the mock method for testing
func (db *MockGearDatabase) GetArmorByID(mongoID primitive.ObjectID) (*model.Armor, error) {
	return db.ArmorToReturn, db.ErrorToReturn
}

//UpdateArmorByID is the mock method for testing
func (db *MockGearDatabase) UpdateArmorByID(armor model.Armor, mongoID primitive.ObjectID) error {
	return db.ErrorToReturn
}

//DeleteArmorByID is the mock method for testing
func (db *MockGearDatabase) DeleteArmorByID(mongoID primitive.ObjectID) error {
	return db.ErrorToReturn
}

//InsertWeapon is the mock method for testing
func (db *MockGearDatabase) InsertWeapon(weapon *model.Weapon) error {
	return db.ErrorToReturn
}

//GetWeapon is the mock method for testing
func (db *MockGearDatabase) GetWeapon(query url.Values) ([]model.Weapon, error) {
	return db.WeaponsToReturn, db.ErrorToReturn
}

//GetWeaponByID is the mock method for testing
func (db *MockGearDatabase) GetWeaponByID(mongoID primitive.ObjectID) (*model.Weapon, error) {
	return db.WeaponToReturn, db.ErrorToReturn
}

//UpdateWeaponByID is the mock method for testing
func (db *MockGearDatabase) UpdateWeaponByID(weapon model.Weapon, mongoID primitive.ObjectID) error {
	return db.ErrorToReturn
}

//DeleteWeaponByID is the mock method for testing
func (db *MockGearDatabase) DeleteWeaponByID(mongoID primitive.ObjectID) error {
	return db.ErrorToReturn
}

//Ping is the mock implementation for testing
func (db *MockGearDatabase) Ping() error {
	return db.ErrorToReturn
}
