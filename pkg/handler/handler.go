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

// GearDatabase is the interface for the actual database object
type GearDatabase interface {
	InsertArmor(armor *model.Armor) error
	InsertWeapon(weapon *model.Weapon) error
	Ping() error
}

//GearService is the implementation of a service to access gear in a database
type GearService struct {
	Version  string
	Database GearDatabase
}

//Routes sets up the routes for the RESTful interface
func (s *GearService) Routes(r *mux.Router) *mux.Router {
	r.HandleFunc("/ping", s.PingCheck).Methods(http.MethodGet)
	r.Handle("/health", s.healthCheck(s.Database)).Methods(http.MethodGet)

	//Inserts
	r.HandleFunc("/armor", s.InsertArmor).Methods(http.MethodPost)
	r.HandleFunc("/weapon", s.InsertWeapon).Methods(http.MethodPost)

	//TODO: GETALL

	//TODO: GETBYID

	//TODO: UPDATEBYID

	//TODO: DELETEBYID
	return r
}

//PingCheck checks that the app is running and returns 200, OK, version
func (s *GearService) PingCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK, " + s.Version))
}

func (s *GearService) healthCheck(database GearDatabase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dbErr := database.Ping()
		var stringDBErr string

		if dbErr != nil {
			stringDBErr = dbErr.Error()
		}

		response := model.HealthCheckResponse{
			APIVersion: s.Version,
			DBError:    stringDBErr,
		}

		if dbErr != nil {
			api.RespondWithJSON(w, http.StatusFailedDependency, response)
			return
		}

		api.RespondWithJSON(w, http.StatusOK, "Object Created")
	})
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

	err = s.Database.InsertArmor(&armorModel)
	if err != nil {
		api.RespondWithError(w, api.CheckError(err), err.Error())
		return
	}

	api.RespondWithJSON(w, http.StatusOK, "Armor Object Created")
}

//InsertWeapon is the handler function for inserting a weapon object
func (s *GearService) InsertWeapon(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("InsertWeapon invoked with url: %v", r.URL)
	defer r.Body.Close()

	var weaponModel model.Weapon
	weaponModel.ID = primitive.NewObjectID()

	err := json.NewDecoder(r.Body).Decode(&weaponModel)
	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}
	err = s.Database.InsertWeapon(&weaponModel)
	if err != nil {
		api.RespondWithError(w, api.CheckError(err), err.Error())
		return
	}

	api.RespondWithJSON(w, http.StatusOK, "Weapon Object Created")
}
