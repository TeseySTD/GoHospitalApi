package main

import (
	"log"
	"net/http"

	"github.com/TeseySTD/GoHospitalApi/handlers"
	"github.com/TeseySTD/GoHospitalApi/storage"
)

func main() {
	if err := storage.Store.LoadFromFile(); err != nil {
		log.Println("Creating new storage file...")
	}

	// Routers
	http.HandleFunc("/patients", handlers.PatientsRouter)
	http.HandleFunc("/patients/", handlers.PatientsRouter)
	http.HandleFunc("/doctors", handlers.DoctorsRouter)
	http.HandleFunc("/doctors/", handlers.DoctorsRouter)
	http.HandleFunc("/appointments", handlers.AppointmentsRouter)
	http.HandleFunc("/appointments/", handlers.AppointmentsRouter)

	// Root page with api description
	http.HandleFunc("/", handlers.RootHandler)

	port := ":8080"

	log.Fatal(http.ListenAndServe(port, nil))
}