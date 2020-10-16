package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Armor struct is used to create an Armor object
type Armor struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	ArmorType   string             `json:"type" bson:"type"`
	Defense     int64              `json:"defense" bson:"defense"`
	Soak        int64              `json:"soak" bson:"soak"`
	Price       int64              `json:"price" bson:"price"`
	Encumbrance int64              `json:"encumbrance" bson:"encumbrance"`
	HardPoitns  int64              `json:"hardPoints" bson:"hardPoints"`
	Rarity      int64              `json:"rarity" bson:"rarity"`
}
