package handler

import (
	"encoding/json"
	"net/http"
	"net/url"

	model "github.com/geeksheik9/gear-CRUD/models"
	"github.com/geeksheik9/gear-CRUD/pkg/api"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GearDatabase is the interface for the actual database object
type GearDatabase interface {
	//Armor methods
	InsertArmor(armor *model.Armor) error
	GetArmor(query url.Values) ([]model.Armor, error)
	GetArmorByID(mongoID primitive.ObjectID) (*model.Armor, error)
	UpdateArmorByID(armor model.Armor, mongoID primitive.ObjectID) error
	DeleteArmorByID(mongoID primitive.ObjectID) error
	//Weapon methods
	InsertWeapon(weapon *model.Weapon) error
	GetWeapon(query url.Values) ([]model.Weapon, error)
	GetWeaponByID(mongoID primitive.ObjectID) (*model.Weapon, error)
	UpdateWeaponByID(weapon model.Weapon, mongoID primitive.ObjectID) error
	DeleteWeaponByID(mongoID primitive.ObjectID) error
	//Helper methods
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

	r.HandleFunc("/armor", s.GetArmor).Methods(http.MethodGet)
	r.HandleFunc("/weapon", s.GetWeapon).Methods(http.MethodGet)

	r.HandleFunc("/armor/{ID}", s.GetArmorByID).Methods(http.MethodGet)
	r.HandleFunc("/weapon/{ID}", s.GetWeaponByID).Methods(http.MethodGet)

	r.HandleFunc("/armor/{ID}", s.UpdateArmorByID).Methods(http.MethodPut)
	r.HandleFunc("/weapon/{ID}", s.UpdateWeaponByID).Methods(http.MethodPut)

	r.HandleFunc("/armor/{ID}", s.DeleteArmorByID).Methods(http.MethodDelete)
	r.HandleFunc("/weapon/{ID}", s.DeleteWeaponByID).Methods(http.MethodDelete)

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

//GetArmor is the handler function to return all armor in the database
func (s *GearService) GetArmor(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("GetArmor invoked with url: %v", r.URL)

	armor, err := s.Database.GetArmor(r.URL.Query())
	if err != nil {
		api.RespondWithError(w, api.CheckError(err), err.Error())
		return
	}

	api.RespondWithJSON(w, http.StatusOK, armor)
}

//GetWeapon is the hanblder function to return all weapons in the database
func (s *GearService) GetWeapon(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("GetWeapon invoked with url: %v", r.URL)

	weapon, err := s.Database.GetWeapon(r.URL.Query())
	if err != nil {
		api.RespondWithError(w, api.CheckError(err), err.Error())
		return
	}

	api.RespondWithJSON(w, http.StatusOK, weapon)
}

//GetArmorByID is the handler function to return a specific armor in the database
func (s *GearService) GetArmorByID(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("GetArmorByID invoked with url: %v", r.URL)

	vars := mux.Vars(r)
	ID := vars["ID"]

	objectID, err := api.StringToObjectID(ID)
	if err != nil {
		api.RespondWithError(w, api.CheckError(err), err.Error())
		return
	}

	armor, err := s.Database.GetArmorByID(objectID)
	if err != nil {
		api.RespondWithError(w, api.CheckError(err), err.Error())
		return
	}

	api.RespondWithJSON(w, http.StatusOK, armor)
}

//GetWeaponByID is the handler function to return a specific weapon in the database
func (s *GearService) GetWeaponByID(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("GetWeaponByID invoked with url: %v", r.URL)

	vars := mux.Vars(r)
	ID := vars["ID"]

	objectID, err := api.StringToObjectID(ID)
	if err != nil {
		api.RespondWithError(w, api.CheckError(err), err.Error())
		return
	}

	weapon, err := s.Database.GetWeaponByID(objectID)
	if err != nil {
		api.RespondWithError(w, api.CheckError(err), err.Error())
		return
	}

	api.RespondWithJSON(w, http.StatusOK, weapon)
}

//UpdateArmorByID is the handler function to update a specific armor in the database
func (s *GearService) UpdateArmorByID(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("UpdateArmorByID invoked with url: %v", r.URL)
	defer r.Body.Close()

	vars := mux.Vars(r)
	ID := vars["ID"]

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		api.RespondWithError(w, api.CheckError(err), err.Error())
		return
	}

	armor := model.Armor{}
	err = json.NewDecoder(r.Body).Decode(&armor)
	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = s.Database.UpdateArmorByID(armor, objectID)
	if err != nil {
		api.RespondWithError(w, api.CheckError(err), err.Error())
		return
	}

	api.RespondWithJSON(w, http.StatusOK, objectID)
}

//UpdateWeaponByID is the handler function to update a specific weapon in the database
func (s *GearService) UpdateWeaponByID(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("UpdateWeaponByID invoked with url: %v", r.URL)
	defer r.Body.Close()

	vars := mux.Vars(r)
	ID := vars["ID"]

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		api.RespondWithError(w, api.CheckError(err), err.Error())
		return
	}

	weapon := model.Weapon{}
	err = json.NewDecoder(r.Body).Decode(&weapon)
	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = s.Database.UpdateWeaponByID(weapon, objectID)
	if err != nil {
		api.RespondWithError(w, api.CheckError(err), err.Error())
		return
	}

	api.RespondWithJSON(w, http.StatusOK, objectID)
}

//DeleteArmorByID is the handler function to remove a specific armor in the database
func (s *GearService) DeleteArmorByID(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("DeleterArmorByID invoked with url: %v", r.URL)

	vars := mux.Vars(r)
	ID := vars["ID"]

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		api.RespondWithError(w, api.CheckError(err), err.Error())
		return
	}

	err = s.Database.DeleteArmorByID(objectID)
	if err != nil {
		api.RespondWithError(w, api.CheckError(err), err.Error())
	}

	api.RespondNoContent(w, http.StatusNoContent)
}

//DeleteWeaponByID is the handler function to remove a specific weapon in the database
func (s *GearService) DeleteWeaponByID(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("DeleterWeaponByID invoked with url: %v", r.URL)

	vars := mux.Vars(r)
	ID := vars["ID"]

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		api.RespondWithError(w, api.CheckError(err), err.Error())
		return
	}

	err = s.Database.DeleteWeaponByID(objectID)
	if err != nil {
		api.RespondWithError(w, api.CheckError(err), err.Error())
	}

	api.RespondNoContent(w, http.StatusNoContent)
}
