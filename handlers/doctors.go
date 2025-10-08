package handlers

import (
	"encoding/json"
	"github.com/TeseySTD/GoHospitalApi/models"
	"github.com/TeseySTD/GoHospitalApi/storage"
	"github.com/TeseySTD/GoHospitalApi/utils"
	"net/http"
	"strings"
)

func GetDoctorsHandler(w http.ResponseWriter, r *http.Request) {
	storage.Store.RLock()
	defer storage.Store.RUnlock()
	utils.RespondJSON(w, http.StatusOK, storage.Store.Doctors)
}

func GetDoctorHandler(w http.ResponseWriter, r *http.Request) {
	id := utils.ExtractID(r.URL.Path, "/doctors/")

	storage.Store.RLock()
	defer storage.Store.RUnlock()

	for _, doctor := range storage.Store.Doctors {
		if doctor.ID == id {
			utils.RespondJSON(w, http.StatusOK, doctor)
			return
		}
	}
	utils.RespondError(w, http.StatusNotFound, "Doctor not found")
}

func CreateDoctorHandler(w http.ResponseWriter, r *http.Request) {
	var doctor models.Doctor
	if err := json.NewDecoder(r.Body).Decode(&doctor); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	storage.Store.Lock()
	defer storage.Store.Unlock()

	doctor.ID = len(storage.Store.Doctors) + 1
	storage.Store.Doctors = append(storage.Store.Doctors, doctor)
	storage.Store.SaveToFile()

	utils.RespondJSON(w, http.StatusCreated, doctor)
}

func UpdateDoctorHandler(w http.ResponseWriter, r *http.Request) {
	id := utils.ExtractID(r.URL.Path, "/doctors/")

	var updatedDoctor models.Doctor
	if err := json.NewDecoder(r.Body).Decode(&updatedDoctor); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	storage.Store.Lock()
	defer storage.Store.Unlock()

	for i, doctor := range storage.Store.Doctors {
		if doctor.ID == id {
			updatedDoctor.ID = id
			storage.Store.Doctors[i] = updatedDoctor
			storage.Store.SaveToFile()
			utils.RespondJSON(w, http.StatusOK, updatedDoctor)
			return
		}
	}
	utils.RespondError(w, http.StatusNotFound, "Doctor not found")
}

func DeleteDoctorHandler(w http.ResponseWriter, r *http.Request) {
	id := utils.ExtractID(r.URL.Path, "/doctors/")

	storage.Store.Lock()
	defer storage.Store.Unlock()

	for i, doctor := range storage.Store.Doctors {
		if doctor.ID == id {
			storage.Store.Doctors = append(storage.Store.Doctors[:i], storage.Store.Doctors[i+1:]...)
			storage.Store.SaveToFile()
			utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Doctor deleted"})
			return
		}
	}
	utils.RespondError(w, http.StatusNotFound, "Doctor not found")
}

func DoctorsRouter(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if strings.HasPrefix(r.URL.Path, "/doctors/") && len(r.URL.Path) > len("/doctors/") {
			GetDoctorHandler(w, r)
		} else {
			GetDoctorsHandler(w, r)
		}
	case http.MethodPost:
		CreateDoctorHandler(w, r)
	case http.MethodPut:
		UpdateDoctorHandler(w, r)
	case http.MethodDelete:
		DeleteDoctorHandler(w, r)
	default:
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}
