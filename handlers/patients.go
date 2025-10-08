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

func GetPatientsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	patients, err := storage.Store.GetAllPatients(ctx)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "failed to fetch patients: "+err.Error())
		return
	}

	query := r.URL.Query()
	firstName := strings.ToLower(query.Get("first_name"))
	lastName := strings.ToLower(query.Get("last_name"))
	ageStr := query.Get("age")
	diagnosis := strings.ToLower(query.Get("diagnosis"))

	if firstName == "" && lastName == "" && ageStr == "" && diagnosis == "" {
		if patients == nil {
			patients = []models.Patient{}
		}
		utils.RespondJSON(w, http.StatusOK, patients)
		return
	}

	var filtered []models.Patient
	for _, p := range patients {
		match := true

		if firstName != "" && !strings.Contains(strings.ToLower(p.FirstName), firstName) {
			match = false
		}
		if lastName != "" && !strings.Contains(strings.ToLower(p.LastName), lastName) {
			match = false
		}
		if ageStr != "" {
			age, err := strconv.Atoi(ageStr)
			if err != nil || p.Age != age {
				match = false
			}
		}
		if diagnosis != "" && !strings.Contains(strings.ToLower(p.Diagnosis), diagnosis) {
			match = false
		}

		if match {
			filtered = append(filtered, p)
		}
	}

	if filtered == nil {
		filtered = []models.Patient{}
	}
	utils.RespondJSON(w, http.StatusOK, filtered)
}

func GetPatientHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := utils.ExtractID(r.URL.Path, "/patients/")

	patient, err := storage.Store.GetPatientByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") || strings.Contains(err.Error(), "not found") {
			utils.RespondError(w, http.StatusNotFound, "Patient not found")
		} else {
			utils.RespondError(w, http.StatusInternalServerError, "failed to fetch patient: "+err.Error())
		}
		return
	}
	utils.RespondJSON(w, http.StatusOK, patient)
}

func CreatePatientHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var patient models.Patient
	if err := json.NewDecoder(r.Body).Decode(&patient); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	created, err := storage.Store.CreatePatient(ctx, &patient)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "failed to create patient: "+err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusCreated, created)
}

func UpdatePatientHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := utils.ExtractID(r.URL.Path, "/patients/")

	var updated models.Patient
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	updated.ID = id

	if err := storage.Store.UpdatePatient(ctx, &updated); err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.RespondError(w, http.StatusNotFound, "Patient not found")
		} else {
			utils.RespondError(w, http.StatusInternalServerError, "update failed: "+err.Error())
		}
		return
	}
	utils.RespondJSON(w, http.StatusOK, updated)
}

func DeletePatientHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := utils.ExtractID(r.URL.Path, "/patients/")

	if err := storage.Store.DeletePatient(ctx, id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.RespondError(w, http.StatusNotFound, "Patient not found")
		} else {
			utils.RespondError(w, http.StatusInternalServerError, "delete failed: "+err.Error())
		}
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Patient deleted"})
}

func PatientsRouter(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if strings.HasPrefix(r.URL.Path, "/patients/") && len(r.URL.Path) > len("/patients/") {
			GetPatientHandler(w, r)
		} else {
			GetPatientsHandler(w, r)
		}
	case http.MethodPost:
		CreatePatientHandler(w, r)
	case http.MethodPut:
		UpdatePatientHandler(w, r)
	case http.MethodDelete:
		DeletePatientHandler(w, r)
	default:
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}
