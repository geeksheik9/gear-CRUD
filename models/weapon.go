package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Weapon struct is used to create a Weapon object
type weapon struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	WeaponType   string             `json:"type" bson:"type"`
	Name         string             `json:"name" bson:"name"`
	Skill        string             `json:"skill" bson:"skill"`
	Damage       string             `json:"damage" bson:"damage"`
	Critical     int64              `json:"critical" bson:"critical"`
	Range        string             `json:"range" bson:"range"`
	Encumberence int64              `json:"encumberence" bson:"encumberence"`
	HP           int64              `json:"hp" bson:"hp"`
	Price        int64              `json:"price" bson:"price"`
	Rarity       int64              `json:"rarity" bson:"rarity"`
	Special      string             `json:"special" bson:"special"`
}
