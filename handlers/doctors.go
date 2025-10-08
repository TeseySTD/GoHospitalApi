package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/TeseySTD/GoHospitalApi/models"
	"github.com/TeseySTD/GoHospitalApi/storage"
	"github.com/TeseySTD/GoHospitalApi/utils"
)

func GetDoctorsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	doctors, err := storage.Store.GetAllDoctors(ctx)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "failed to fetch doctors: "+err.Error())
		return
	}

	query := r.URL.Query()
	firstName := strings.ToLower(query.Get("first_name"))
	lastName := strings.ToLower(query.Get("last_name"))
	specialization := strings.ToLower(query.Get("specialization"))
	experienceStr := query.Get("experience")
	minExpStr := query.Get("min_experience")

	if firstName == "" && lastName == "" && specialization == "" && experienceStr == "" && minExpStr == "" {
		if doctors == nil {
			doctors = []models.Doctor{}
		}
		utils.RespondJSON(w, http.StatusOK, doctors)
		return
	}

	var filtered []models.Doctor
	for _, doctor := range doctors {
		match := true

		if firstName != "" && !strings.Contains(strings.ToLower(doctor.FirstName), firstName) {
			match = false
		}
		if lastName != "" && !strings.Contains(strings.ToLower(doctor.LastName), lastName) {
			match = false
		}
		if specialization != "" && !strings.Contains(strings.ToLower(doctor.Specialization), specialization) {
			match = false
		}
		if experienceStr != "" {
			experience, err := strconv.Atoi(experienceStr)
			if err != nil || doctor.Experience != experience {
				match = false
			}
		}
		if minExpStr != "" {
			minExp, err := strconv.Atoi(minExpStr)
			if err != nil || doctor.Experience < minExp {
				match = false
			}
		}

		if match {
			filtered = append(filtered, doctor)
		}
	}

	if filtered == nil {
		filtered = []models.Doctor{}
	}
	utils.RespondJSON(w, http.StatusOK, filtered)
}

func GetDoctorHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := utils.ExtractID(r.URL.Path, "/doctors/")

	doctor, err := storage.Store.GetDoctorByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") || strings.Contains(err.Error(), "not found") {
			utils.RespondError(w, http.StatusNotFound, "Doctor not found")
		} else {
			utils.RespondError(w, http.StatusInternalServerError, "failed to fetch doctor: "+err.Error())
		}
		return
	}
	utils.RespondJSON(w, http.StatusOK, doctor)
}

func CreateDoctorHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var doctor models.Doctor
	if err := json.NewDecoder(r.Body).Decode(&doctor); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	created, err := storage.Store.CreateDoctor(ctx, &doctor)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "failed to create doctor: "+err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusCreated, created)
}

func UpdateDoctorHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := utils.ExtractID(r.URL.Path, "/doctors/")

	var updatedDoctor models.Doctor
	if err := json.NewDecoder(r.Body).Decode(&updatedDoctor); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	updatedDoctor.ID = id

	if err := storage.Store.UpdateDoctor(ctx, &updatedDoctor); err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.RespondError(w, http.StatusNotFound, "Doctor not found")
		} else {
			utils.RespondError(w, http.StatusInternalServerError, "update failed: "+err.Error())
		}
		return
	}
	utils.RespondJSON(w, http.StatusOK, updatedDoctor)
}

func DeleteDoctorHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := utils.ExtractID(r.URL.Path, "/doctors/")

	if err := storage.Store.DeleteDoctor(ctx, id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.RespondError(w, http.StatusNotFound, "Doctor not found")
		} else {
			utils.RespondError(w, http.StatusInternalServerError, "delete failed: "+err.Error())
		}
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Doctor deleted"})
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
