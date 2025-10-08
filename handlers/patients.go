package handlers

import (
	"encoding/json"
	"github.com/TeseySTD/GoHospitalApi/models"
	"github.com/TeseySTD/GoHospitalApi/storage"
	"github.com/TeseySTD/GoHospitalApi/utils"
	"net/http"
	"strconv"
	"strings"
)

func GetPatientsHandler(w http.ResponseWriter, r *http.Request) {
	storage.Store.RLock()
	defer storage.Store.RUnlock()
	
	query := r.URL.Query()
	firstName := query.Get("first_name")
	lastName := query.Get("last_name")
	ageStr := query.Get("age")
	diagnosis := query.Get("diagnosis")
	
	if firstName == "" && lastName == "" && ageStr == "" && diagnosis == "" {
		utils.RespondJSON(w, http.StatusOK, storage.Store.Patients)
		return
	}
	
	var filteredPatients []models.Patient
	
	for _, patient := range storage.Store.Patients {
		match := true
		
		if firstName != "" && !strings.Contains(strings.ToLower(patient.FirstName), strings.ToLower(firstName)) {
			match = false
		}
		if lastName != "" && !strings.Contains(strings.ToLower(patient.LastName), strings.ToLower(lastName)) {
			match = false
		}
		if ageStr != "" {
			age, err := strconv.Atoi(ageStr)
			if err != nil || patient.Age != age {
				match = false
			}
		}
		if diagnosis != "" && !strings.Contains(strings.ToLower(patient.Diagnosis), strings.ToLower(diagnosis)) {
			match = false
		}
		
		if match {
			filteredPatients = append(filteredPatients, patient)
		}
	}
	
	if filteredPatients == nil {
		filteredPatients = []models.Patient{}
	}
	
	utils.RespondJSON(w, http.StatusOK, filteredPatients)
}

func GetPatientHandler(w http.ResponseWriter, r *http.Request) {
	id := utils.ExtractID(r.URL.Path, "/patients/")

	storage.Store.RLock()
	defer storage.Store.RUnlock()

	for _, patient := range storage.Store.Patients {
		if patient.ID == id {
			utils.RespondJSON(w, http.StatusOK, patient)
			return
		}
	}
	utils.RespondError(w, http.StatusNotFound, "Patient not found")
}

func CreatePatientHandler(w http.ResponseWriter, r *http.Request) {
	var patient models.Patient
	if err := json.NewDecoder(r.Body).Decode(&patient); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	storage.Store.Lock()
	defer storage.Store.Unlock()

	patient.ID = len(storage.Store.Patients) + 1
	storage.Store.Patients = append(storage.Store.Patients, patient)
	storage.Store.SaveToFile()

	utils.RespondJSON(w, http.StatusCreated, patient)
}

func UpdatePatientHandler(w http.ResponseWriter, r *http.Request) {
	id := utils.ExtractID(r.URL.Path, "/patients/")

	var updatedPatient models.Patient
	if err := json.NewDecoder(r.Body).Decode(&updatedPatient); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	storage.Store.Lock()
	defer storage.Store.Unlock()

	for i, patient := range storage.Store.Patients {
		if patient.ID == id {
			updatedPatient.ID = id
			storage.Store.Patients[i] = updatedPatient
			storage.Store.SaveToFile()
			utils.RespondJSON(w, http.StatusOK, updatedPatient)
			return
		}
	}
	utils.RespondError(w, http.StatusNotFound, "Patient not found")
}

func DeletePatientHandler(w http.ResponseWriter, r *http.Request) {
	id := utils.ExtractID(r.URL.Path, "/patients/")

	storage.Store.Lock()
	defer storage.Store.Unlock()

	for i, patient := range storage.Store.Patients {
		if patient.ID == id {
			storage.Store.Patients = append(storage.Store.Patients[:i], storage.Store.Patients[i+1:]...)
			storage.Store.SaveToFile()
			utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Patient deleted"})
			return
		}
	}
	utils.RespondError(w, http.StatusNotFound, "Patient not found")
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