package handlers

import (
	"encoding/json"
	"github.com/TeseySTD/GoHospitalApi/models"
	"github.com/TeseySTD/GoHospitalApi/storage"
	"github.com/TeseySTD/GoHospitalApi/utils"
	"net/http"
	"strings"
)

func GetAppointmentsHandler(w http.ResponseWriter, r *http.Request) {
	storage.Store.RLock()
	defer storage.Store.RUnlock()
	utils.RespondJSON(w, http.StatusOK, storage.Store.Appointments)
}

func GetAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	id := utils.ExtractID(r.URL.Path, "/appointments/")

	storage.Store.RLock()
	defer storage.Store.RUnlock()

	for _, appointment := range storage.Store.Appointments {
		if appointment.ID == id {
			utils.RespondJSON(w, http.StatusOK, appointment)
			return
		}
	}
	utils.RespondError(w, http.StatusNotFound, "Appointment not found")
}

func CreateAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	var appointment models.Appointment
	if err := json.NewDecoder(r.Body).Decode(&appointment); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	storage.Store.Lock()
	defer storage.Store.Unlock()

	appointment.ID = len(storage.Store.Appointments) + 1
	storage.Store.Appointments = append(storage.Store.Appointments, appointment)
	storage.Store.SaveToFile()

	utils.RespondJSON(w, http.StatusCreated, appointment)
}

func UpdateAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	id := utils.ExtractID(r.URL.Path, "/appointments/")

	var updatedAppointment models.Appointment
	if err := json.NewDecoder(r.Body).Decode(&updatedAppointment); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	storage.Store.Lock()
	defer storage.Store.Unlock()

	for i, appointment := range storage.Store.Appointments {
		if appointment.ID == id {
			updatedAppointment.ID = id
			storage.Store.Appointments[i] = updatedAppointment
			storage.Store.SaveToFile()
			utils.RespondJSON(w, http.StatusOK, updatedAppointment)
			return
		}
	}
	utils.RespondError(w, http.StatusNotFound, "Appointment not found")
}

func DeleteAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	id := utils.ExtractID(r.URL.Path, "/appointments/")

	storage.Store.Lock()
	defer storage.Store.Unlock()

	for i, appointment := range storage.Store.Appointments {
		if appointment.ID == id {
			storage.Store.Appointments = append(storage.Store.Appointments[:i], storage.Store.Appointments[i+1:]...)
			storage.Store.SaveToFile()
			utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Appointment deleted"})
			return
		}
	}
	utils.RespondError(w, http.StatusNotFound, "Appointment not found")
}

func AppointmentsRouter(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if strings.HasPrefix(r.URL.Path, "/appointments/") && len(r.URL.Path) > len("/appointments/") {
			GetAppointmentHandler(w, r)
		} else {
			GetAppointmentsHandler(w, r)
		}
	case http.MethodPost:
		CreateAppointmentHandler(w, r)
	case http.MethodPut:
		UpdateAppointmentHandler(w, r)
	case http.MethodDelete:
		DeleteAppointmentHandler(w, r)
	default:
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}
