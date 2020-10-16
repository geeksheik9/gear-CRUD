package main

import (
	"fmt"
	"net/http"

	"github.com/geeksheik9/gear-CRUD/config"
	"github.com/geeksheik9/gear-CRUD/pkg/handler"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var version string

func main() {
	logrus.Info("INITIALIZING GEAR CRUD")

	accessor := viper.New()

	config, err := config.New(accessor)
	if err != nil {
		logrus.Fatalf("ERROR LOADING CONFIG: %v", err.Error())
	}

	gearService := handler.GearService{
		Version: version,
	}

	r := mux.NewRouter().StrictSlash(true)

	r = gearService.Routes(r)
	fmt.Printf("Server listen on port %v\n", config.Port)
	logrus.Info("END")
	logrus.Fatal(http.ListenAndServe(":"+config.Port, cors.Default().Handler(r)))
}
