package main

import (
	"log"
	"net/http"

	"github.com/TeseySTD/GoHospitalApi/handlers"
	"github.com/TeseySTD/GoHospitalApi/middleware"
	"github.com/TeseySTD/GoHospitalApi/storage"
)

func main() {
	if err := storage.Store.LoadFromFile(); err != nil {
		log.Println("Creating new storage file...")
	}

	http.HandleFunc("/patients", middleware.Chain(
		handlers.PatientsRouter,
		middleware.LoggingMiddleware,
		middleware.AuthMiddleware,
	))
	http.HandleFunc("/patients/", middleware.Chain(
		handlers.PatientsRouter,
		middleware.LoggingMiddleware,
		middleware.AuthMiddleware,
	))

	http.HandleFunc("/doctors", middleware.Chain(
		handlers.DoctorsRouter,
		middleware.LoggingMiddleware,
		middleware.AuthMiddleware,
	))
	http.HandleFunc("/doctors/", middleware.Chain(
		handlers.DoctorsRouter,
		middleware.LoggingMiddleware,
		middleware.AuthMiddleware,
	))

	http.HandleFunc("/appointments", middleware.Chain(
		handlers.AppointmentsRouter,
		middleware.LoggingMiddleware,
		middleware.AuthMiddleware,
	))
	http.HandleFunc("/appointments/", middleware.Chain(
		handlers.AppointmentsRouter,
		middleware.LoggingMiddleware,
		middleware.AuthMiddleware,
	))

	http.HandleFunc("/", middleware.LoggingMiddleware(handlers.RootHandler))

	port := ":8080"
	log.Printf("Server starting on port %s", port)
	log.Printf("API Key for authorization: %s", middleware.ValidAPIKey)

	log.Fatal(http.ListenAndServe(port, nil))
}
