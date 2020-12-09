package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	model "github.com/geeksheik9/gear-CRUD/models"
	"github.com/geeksheik9/gear-CRUD/pkg/db/mocks"
	"github.com/gorilla/mux"
)

func InitMockGearService(armor *model.Armor, armors []model.Armor, weapon *model.Weapon, weapons []model.Weapon, errors error) (s GearService) {
	db := mocks.MockGearDatabase{
		ArmorToReturn:   armor,
		ArmorsToReturn:  armors,
		WeaponToReturn:  weapon,
		WeaponsToReturn: weapons,
		ErrorToReturn:   errors,
	}

	return GearService{
		Version:  "test",
		Database: &db,
	}
}

func mockWeapon() {

}

func mockArmor() {

}

func TestCharacterService_PingCheck(t *testing.T) {
	service := InitMockGearService(nil, nil, nil, nil, nil)

	r, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		t.Errorf("Ping() error creating request:\ngot: %v\nexpected:<no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("Ping() error:\ngot: %v\n expected: %v", w.Code, http.StatusOK)
	}
}

func TestCharacterService_HealthCheck(t *testing.T) {
	service := InitMockGearService(nil, nil, nil, nil, nil)
	r, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Errorf("Ping() error creating request:\ngot: %v\nexpected:<no error>", err)
	}
	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("Ping() error:\ngot: %v\n expected: %v", w.Code, http.StatusOK)
	}
}

func TestCharacterService_HealthCheck_error(t *testing.T) {
	service := InitMockGearService(nil, nil, nil, nil, errors.New("test error"))
	r, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Errorf("Ping() error creating request:\ngot: %v\nexpected:<no error>", err)
	}
	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)
	if w.Code != http.StatusFailedDependency {
		t.Errorf("Ping() error:\ngot: %v\n expected: %v", w.Code, http.StatusFailedDependency)
	}
}
