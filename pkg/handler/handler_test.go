package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	model "github.com/geeksheik9/gear-CRUD/models"
	"github.com/geeksheik9/gear-CRUD/pkg/db/mocks"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func mockWeapon(id primitive.ObjectID, name string, price int64) model.Weapon {
	weapon := model.Weapon{
		ID:    id,
		Name:  name,
		Price: price,
	}
	return weapon
}

func mockWeapons(weapon model.Weapon) []model.Weapon {
	weapons := []model.Weapon{
		weapon,
	}
	return weapons
}

func mockSingleArmor(id primitive.ObjectID, armorType string, price int64) model.Armor {
	armor := model.Armor{
		ID:        id,
		ArmorType: armorType,
		Price:     price,
	}
	return armor
}

func mockArmor(armor model.Armor) []model.Armor {
	armors := []model.Armor{
		armor,
	}
	return armors
}

func TestGearService_PingCheck(t *testing.T) {
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

func TestGearService_HealthCheck(t *testing.T) {
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

func TestGearService_HealthCheck_error(t *testing.T) {
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

func TestGearService_InsertArmor_Success(t *testing.T) {
	id := primitive.NewObjectID()
	armor := mockSingleArmor(id, "test", 5)
	service := InitMockGearService(&armor, nil, nil, nil, nil)

	request, _ := json.Marshal(armor)

	r, err := http.NewRequest("POST", "/armor", bytes.NewBuffer(request))
	if err != nil {
		t.Errorf("InsertArmor() error creating request:\ngot: %v\nexpected: <no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("InsertArmor() error:\ngot: %v\nexpected: %v", w.Code, http.StatusOK)
	}
}

func TestGearService_InsertArmor_DBError(t *testing.T) {
	id := primitive.NewObjectID()
	armor := mockSingleArmor(id, "test", 5)
	service := InitMockGearService(nil, nil, nil, nil, errors.New("test error"))

	request, _ := json.Marshal(armor)

	r, err := http.NewRequest("POST", "/armor", bytes.NewBuffer(request))
	if err != nil {
		t.Errorf("InsertArmor() error creating request:\ngot: %v\nexpected: <no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("InsertArmor() error:\ngot: %v\nexpected:%v", w.Code, http.StatusInternalServerError)
	}
}

func TestGearService_InsertArmor_BadJSON(t *testing.T) {
	service := InitMockGearService(nil, nil, nil, nil, nil)

	request, _ := json.Marshal(`{bad json`)

	r, err := http.NewRequest("POST", "/armor", bytes.NewBuffer(request))
	if err != nil {
		t.Errorf("InsertArmor() error creating request:\ngot: %v\nexpected: <no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("InsertArmor() error:\ngot: %v\nexpected: %v", w.Code, http.StatusBadRequest)
	}
}

func TestGearService_InsertWeapon_Success(t *testing.T) {
	id := primitive.NewObjectID()
	weapon := mockWeapon(id, "test", 5)
	service := InitMockGearService(nil, nil, &weapon, nil, nil)

	request, _ := json.Marshal(weapon)

	r, err := http.NewRequest("POST", "/weapon", bytes.NewBuffer(request))
	if err != nil {
		t.Errorf("InsertWeapon() error creating request:\ngot: %v\nexpected: <no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("InsertWeapon() error:\ngot: %v\nexpected: %v", w.Code, http.StatusOK)
	}
}
func TestGearService_InsertWeapon_DBError(t *testing.T) {
	id := primitive.NewObjectID()
	weapon := mockWeapon(id, "test", 5)
	service := InitMockGearService(nil, nil, nil, nil, errors.New("test error"))

	request, _ := json.Marshal(weapon)

	r, err := http.NewRequest("POST", "/weapon", bytes.NewBuffer(request))
	if err != nil {
		t.Errorf("InsertWeapon() error creating request:\ngot: %v\nexpected: <no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("InsertWeapon() error:\ngot: %v\nexpected:%v", w.Code, http.StatusInternalServerError)
	}
}
func TestGearService_InsertWeapon_BadJSON(t *testing.T) {
	service := InitMockGearService(nil, nil, nil, nil, nil)

	request, _ := json.Marshal(`{bad json`)

	r, err := http.NewRequest("POST", "/weapon", bytes.NewBuffer(request))
	if err != nil {
		t.Errorf("InsertWeapon() error creating request:\ngot: %v\nexpected: <no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("InsertWeapon() error:\ngot: %v\nexpected: %v", w.Code, http.StatusBadRequest)
	}
}

func TestGearService_GetArmor_Success(t *testing.T) {
	id := primitive.NewObjectID()
	armor := mockSingleArmor(id, "test", 5)
	armors := mockArmor(armor)
	service := InitMockGearService(nil, armors, nil, nil, nil)

	r, err := http.NewRequest("GET", "/armor", nil)
	if err != nil {
		t.Errorf("GetArmor() error creating request:\ngot: %v\n expected: <no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("GetArmor() error:\ngot: %v\nexpected: %v", w.Code, http.StatusOK)
	}

	resp := []model.Armor{}
	err = json.NewDecoder(w.Body).Decode(&resp)
	if err != nil {
		t.Errorf("GetArmor() error:\n got: %v\n expected: <nil>", err)
	}
	if resp[0].ID != armor.ID || resp[0].ArmorType != armor.ArmorType || resp[0].Price != armor.Price {
		t.Errorf("GetArmor() error:\ngot: %v\nexpected: %v", resp[0], armor)
	}
}

func TestGearService_GetArmor_DBError(t *testing.T) {
	service := InitMockGearService(nil, nil, nil, nil, errors.New("test error"))

	r, err := http.NewRequest("GET", "/armor", nil)
	if err != nil {
		t.Errorf("GetArmor() error creating request:\ngot: %v\nexpected:<no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("GetForceCharacterSheets() error:\ngot: %v\nexpected: %v", w.Code, http.StatusInternalServerError)
	}
}

func TestGearService_GetWeapon_Success(t *testing.T) {
	id := primitive.NewObjectID()
	weapon := mockWeapon(id, "test", 5)
	weapons := mockWeapons(weapon)
	service := InitMockGearService(nil, nil, nil, weapons, nil)

	r, err := http.NewRequest("GET", "/weapon", nil)
	if err != nil {
		t.Errorf("GetWeapon() error creating request: \n got: %v\n expected: <no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("GetWeapon() error:\ngot: %v\nexpected: %v", w.Code, http.StatusOK)
	}

	resp := []model.Weapon{}
	err = json.NewDecoder(w.Body).Decode(&resp)
	if err != nil {
		t.Errorf("GetWeapon() error:\n got: %v\n expected: <nil>", err)
	}
	if resp[0].ID != weapon.ID || resp[0].WeaponType != weapon.WeaponType || resp[0].Price != weapon.Price {
		t.Errorf("GetWeapon() error:\ngot: %v\nexpected: %v", resp[0], weapon)
	}
}

func TestGearService_GetWeapon_DBError(t *testing.T) {
	service := InitMockGearService(nil, nil, nil, nil, errors.New("test error"))

	r, err := http.NewRequest("GET", "/weapon", nil)
	if err != nil {
		t.Errorf("GetWeapon() error creating request: \n got: %v \n expected: <no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("GetForceCharacterSheets() error: \n got: %v \n expected: %v", w.Code, http.StatusInternalServerError)
	}
}

func TestGearService_GetArmorByID_Success(t *testing.T) {
	id := primitive.NewObjectID()
	armor := mockSingleArmor(id, "test", 5)
	service := InitMockGearService(&armor, nil, nil, nil, nil)

	r, err := http.NewRequest("GET", "/armor/"+id.Hex(), nil)
	if err != nil {
		t.Errorf("GetArmorByID() error:\ngot: %v\nexpected: <no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("GetArmorByID() error:\ngot: %v\nexpected: %v", w.Code, http.StatusOK)
	}

	resp := model.Armor{}
	err = json.NewDecoder(w.Body).Decode(&resp)
	if err != nil {
		t.Errorf("GetArmorByID() error:\n got: %v\n expected: <nil>", err)
	}
	if resp.ID != armor.ID || resp.ArmorType != armor.ArmorType || resp.Price != armor.Price {
		t.Errorf("GetArmorByID() error:\ngot: %v\nexpected: %v", resp, armor)
	}
}

func TestGearService_GetArmorByID_DBError(t *testing.T) {
	id := primitive.NewObjectID()
	service := InitMockGearService(nil, nil, nil, nil, errors.New("test error"))

	r, err := http.NewRequest("GET", "/armor/"+id.Hex(), nil)
	if err != nil {
		t.Errorf("GetArmorByID() error:\ngot: %v\nexpected: <no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("GetArmorByID() error:\ngot: %v\nexpected: %v", w.Code, http.StatusInternalServerError)
	}
}

func TestGearService_GetArmorbyID_BadID(t *testing.T) {
	id := "this is a bad id"
	service := InitMockGearService(nil, nil, nil, nil, nil)

	r, err := http.NewRequest("GET", "/armor/"+id, nil)
	if err != nil {
		t.Errorf("GetArmorByID() error:\ngot: %v\nexpected: <no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("FindForceCharacterSheetByID() error:\ngot: %v\nexpected: %v", w.Code, http.StatusInternalServerError)
	}
}

func TestGearService_GetWeaponByID_Success(t *testing.T) {
	id := primitive.NewObjectID()
	weapon := mockWeapon(id, "test", 5)
	service := InitMockGearService(nil, nil, &weapon, nil, nil)

	r, err := http.NewRequest("GET", "/weapon/"+id.Hex(), nil)
	if err != nil {
		t.Errorf("GetWeaponByID() error:\n got: %v \n expected: <no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("GetArmorByID() error:\ngot: %v\nexpected: %v", w.Code, http.StatusOK)
	}
	resp := model.Weapon{}
	err = json.NewDecoder(w.Body).Decode(&resp)
	if err != nil {
		t.Errorf("GetWeaponByID() error: \n got: %v \n expected: <nil>", err)
	}
	if resp.ID != weapon.ID || resp.WeaponType != weapon.WeaponType || resp.Price != weapon.Price {
		t.Errorf("GetWeaponByID error: \n got: %v \n expected: %v", resp, weapon)
	}
}

func TestGearService_GetWeaponByID_DBError(t *testing.T) {
	id := primitive.NewObjectID()
	service := InitMockGearService(nil, nil, nil, nil, errors.New("test error"))

	r, err := http.NewRequest("GET", "/weapon/"+id.Hex(), nil)
	if err != nil {
		t.Errorf("GetWeaponByID error: \n got: %v \n expected: <no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("GetWeaponByID error: \n got: %v \n expected: %v", w.Code, http.StatusInternalServerError)
	}
}

func TestGearService_GetWeaponByID_BadID(t *testing.T) {
	id := "this is a bad id"
	service := InitMockGearService(nil, nil, nil, nil, nil)

	r, err := http.NewRequest("GET", "/weapon/"+id, nil)
	if err != nil {
		t.Errorf("GetWeaponByID error: \n got: %v \n expected: <no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("GetForceCharacterSheetByID error: \n got: %v \n expected: %v", w.Code, http.StatusInternalServerError)
	}
}

func TestGearService_UpdateArmorByID_Success(t *testing.T) {
	id := primitive.NewObjectID()
	armor := mockSingleArmor(id, "test", 5)
	service := InitMockGearService(&armor, nil, nil, nil, nil)

	request, _ := json.Marshal(armor)

	r, err := http.NewRequest("PUT", "/armor/"+id.Hex(), bytes.NewBuffer(request))
	if err != nil {
		t.Errorf("UpdateArmorByID() error creating request:\ngot: %v\nexpected: <no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("UpdatedArmorByID error:\ngot: %v\nexpected: %v", w.Code, http.StatusOK)
	}

	idString := "\"" + id.Hex() + "\""
	if w.Body.String() != idString {
		t.Errorf("UpdateArmorByID error:\ngot: %v\nexpected: %v", w.Body.String(), idString)
	}
}

func TestGearService_UpdateArmorByID_DBError(t *testing.T) {
	id := primitive.NewObjectID()
	armor := mockSingleArmor(id, "test", 5)
	service := InitMockGearService(nil, nil, nil, nil, errors.New("test error"))

	request, _ := json.Marshal(armor)

	r, err := http.NewRequest("PUT", "/armor/"+id.Hex(), bytes.NewBuffer(request))
	if err != nil {
		t.Errorf("UpdateArmorByID error: \n got: %v \n expected: <no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("UpdateArmorByID error: \n got: %v \n expected: %v", w.Code, http.StatusInternalServerError)
	}
}

func TestGearService_UpdateArmorByID_BadID(t *testing.T) {
	id := "this is a bad id"
	armor := mockSingleArmor(primitive.NewObjectID(), "test", 5)
	service := InitMockGearService(nil, nil, nil, nil, nil)

	request, _ := json.Marshal(armor)

	r, err := http.NewRequest("PUT", "/armor/"+id, bytes.NewBuffer(request))
	if err != nil {
		t.Errorf("UpdateArmorByID error: \n got: %v \n expected: <no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("GetForceCharacterSheetByID error: \n got: %v \n expected: %v", w.Code, http.StatusInternalServerError)
	}
}

func TestGearService_UpdateArmorByID_BadJSON(t *testing.T) {
	id := primitive.NewObjectID()
	service := InitMockGearService(nil, nil, nil, nil, nil)

	request, _ := json.Marshal(`{bad json`)

	r, err := http.NewRequest("PUT", "/armor/"+id.Hex(), bytes.NewBuffer(request))
	if err != nil {
		t.Errorf("UpdateArmorByID() error creating request:\ngot: %v\nexpected: <no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("UpdatedArmorByID() error:\ngot: %v\nexpected: %v", w.Code, http.StatusBadRequest)
	}
}

func TestGearService_UpdateWeaponByID_Success(t *testing.T) {
	id := primitive.NewObjectID()
	weapon := mockWeapon(id, "test", 5)
	service := InitMockGearService(nil, nil, &weapon, nil, nil)

	request, _ := json.Marshal(weapon)

	r, err := http.NewRequest("PUT", "/weapon/"+id.Hex(), bytes.NewBuffer(request))
	if err != nil {
		t.Errorf("UpdateWeaponByID() error creating request:\ngot: %v\nexpected: <no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("UpdatedWeaponByID error:\ngot: %v\nexpected: %v", w.Code, http.StatusOK)
	}

	idString := "\"" + id.Hex() + "\""
	if w.Body.String() != idString {
		t.Errorf("UpdateWeaponByID error:\ngot: %v\nexpected: %v", w.Body.String(), idString)
	}
}

func TestGearService_UpdateWeaponByID_DBError(t *testing.T) {
	id := primitive.NewObjectID()
	weapon := mockWeapon(id, "test", 5)
	service := InitMockGearService(nil, nil, nil, nil, errors.New("test error"))

	request, _ := json.Marshal(weapon)

	r, err := http.NewRequest("PUT", "/weapon/"+id.Hex(), bytes.NewBuffer(request))
	if err != nil {
		t.Errorf("UpdateWeaponByID error: \n got: %v \n expected: <no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("UpdateWeaponByID error: \n got: %v \n expected: %v", w.Code, http.StatusInternalServerError)
	}
}

func TestGearService_UpdateWeaponByID_BadID(t *testing.T) {
	id := "this is a bad id"
	weapon := mockWeapon(primitive.NewObjectID(), "test", 5)
	service := InitMockGearService(nil, nil, nil, nil, nil)

	request, _ := json.Marshal(weapon)

	r, err := http.NewRequest("PUT", "/weapon/"+id, bytes.NewBuffer(request))
	if err != nil {
		t.Errorf("UpdateWeaponByID error: \n got: %v \n expected: <no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("GetForceCharacterSheetByID error: \n got: %v \n expected: %v", w.Code, http.StatusInternalServerError)
	}
}

func TestGearService_UpdateWeaponByID_BadJSON(t *testing.T) {
	id := primitive.NewObjectID()
	service := InitMockGearService(nil, nil, nil, nil, nil)

	request, _ := json.Marshal(`{bad json`)

	r, err := http.NewRequest("PUT", "/weapon/"+id.Hex(), bytes.NewBuffer(request))
	if err != nil {
		t.Errorf("UpdateWeaponByID() error creating request:\ngot: %v\nexpected: <no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("UpdatedWeaponByID() error:\ngot: %v\nexpected: %v", w.Code, http.StatusBadRequest)
	}
}

func TestGearService_DeleteArmorByID_Success(t *testing.T) {
	id := primitive.NewObjectID()
	service := InitMockGearService(nil, nil, nil, nil, nil)

	r, err := http.NewRequest("DELETE", "/armor/"+id.Hex(), nil)
	if err != nil {
		t.Errorf("DeleteArmorByID() error creating request:\ngot: %v\nexpected: <no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusNoContent {
		t.Errorf("DeleteArmorByID() error:\n got:%v\nexpected: %v", w.Code, http.StatusNoContent)
	}
}

func TestGearService_DeleteArmorByID_DBError(t *testing.T) {
	id := primitive.NewObjectID()
	service := InitMockGearService(nil, nil, nil, nil, errors.New("test error"))

	r, err := http.NewRequest("DELETE", "/armor/"+id.Hex(), nil)
	if err != nil {
		t.Errorf("DeleteArmorByID error: \n got: %v \n expected: <no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("DeleteArmorByID error: \n got: %v \n expected: %v", w.Code, http.StatusInternalServerError)
	}
}

func TestGearService_DeleteArmorByID_BadID(t *testing.T) {
	id := "this is a bad id"
	service := InitMockGearService(nil, nil, nil, nil, nil)

	r, err := http.NewRequest("DELETE", "/armor/"+id, nil)
	if err != nil {
		t.Errorf("DeleteArmorByID error: \n got: %v \n expected: <no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("GetForceCharacterSheetByID error: \n got: %v \n expected: %v", w.Code, http.StatusInternalServerError)

	}
}

func TestGearService_DeleteWeaponByID_Success(t *testing.T) {
	id := primitive.NewObjectID()
	service := InitMockGearService(nil, nil, nil, nil, nil)

	r, err := http.NewRequest("DELETE", "/weapon/"+id.Hex(), nil)
	if err != nil {
		t.Errorf("DeleteWeaponByID() error creating request:\ngot: %v\nexpected: <no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusNoContent {
		t.Errorf("DeleteWeaponByID() error:\n got:%v\nexpected: %v", w.Code, http.StatusNoContent)
	}
}

func TestGearService_DeleteWeaponByID_DBError(t *testing.T) {
	id := primitive.NewObjectID()
	service := InitMockGearService(nil, nil, nil, nil, errors.New("test error"))

	r, err := http.NewRequest("DELETE", "/weapon/"+id.Hex(), nil)
	if err != nil {
		t.Errorf("DeleteWeaponByID error: \n got: %v \n expected: <no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("DeleteWeaponByID error: \n got: %v \n expected: %v", w.Code, http.StatusInternalServerError)
	}
}

func TestGearService_DeleteWeaponByID_BadID(t *testing.T) {
	id := "this is a bad id"
	service := InitMockGearService(nil, nil, nil, nil, nil)

	r, err := http.NewRequest("DELETE", "/weapon/"+id, nil)
	if err != nil {
		t.Errorf("DeleteWeaponByID error: \n got: %v \n expected: <no error>", err)
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	service.Routes(router).ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("GetForceCharacterSheetByID error: \n got: %v \n expected: %v", w.Code, http.StatusInternalServerError)

	}
}
