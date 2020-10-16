package handler

import (
	"encoding/json"
	"net/http"

	model "github.com/geeksheik9/gear-CRUD/models"
	"github.com/geeksheik9/gear-CRUD/pkg/api"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//GearService is the implementation of a service to access gear in a database
type GearService struct {
	Version string
}

//Routes sets up the routes for the RESTful interface
func (s *GearService) Routes(r *mux.Router) *mux.Router {
	r.HandleFunc("/armor", s.InsertArmor).Methods(http.MethodPost)

	return r
}

//InsertArmor is the handler function for inserting an armor object
func (s *GearService) InsertArmor(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("InsertArmor invoked with url: %v", r.URL)
	defer r.Body.Close()

	var armorModel model.Armor
	armorModel.ID = primitive.NewObjectID()

	err := json.NewDecoder(r.Body).Decode(&armorModel)
	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	api.RespondWithJSON(w, http.StatusOK, armorModel)

}
