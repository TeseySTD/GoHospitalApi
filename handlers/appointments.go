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

func GetAppointmentsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	appointments, err := storage.Store.GetAllAppointments(ctx)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "failed to fetch appointments: "+err.Error())
		return
	}

	query := r.URL.Query()
	patientIDStr := query.Get("patient_id")
	doctorIDStr := query.Get("doctor_id")
	date := query.Get("date")
	status := query.Get("status")

	if patientIDStr == "" && doctorIDStr == "" && date == "" && status == "" {
		if appointments == nil {
			appointments = []models.Appointment{}
		}
		utils.RespondJSON(w, http.StatusOK, appointments)
		return
	}

	var filtered []models.Appointment
	for _, ap := range appointments {
		match := true

		if patientIDStr != "" {
			patientID, err := strconv.Atoi(patientIDStr)
			if err != nil || ap.PatientID != patientID {
				match = false
			}
		}
		if doctorIDStr != "" {
			doctorID, err := strconv.Atoi(doctorIDStr)
			if err != nil || ap.DoctorID != doctorID {
				match = false
			}
		}
		if date != "" && ap.Date != date {
			match = false
		}
		if status != "" && !strings.EqualFold(ap.Status, status) {
			match = false
		}

		if match {
			filtered = append(filtered, ap)
		}
	}

	if filtered == nil {
		filtered = []models.Appointment{}
	}
	utils.RespondJSON(w, http.StatusOK, filtered)
}

func GetAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := utils.ExtractID(r.URL.Path, "/appointments/")

	appointment, err := storage.Store.GetAppointmentByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") || strings.Contains(err.Error(), "not found") {
			utils.RespondError(w, http.StatusNotFound, "Appointment not found")
		} else {
			utils.RespondError(w, http.StatusInternalServerError, "failed to fetch appointment: "+err.Error())
		}
		return
	}
	utils.RespondJSON(w, http.StatusOK, appointment)
}

func CreateAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var ap models.Appointment
	if err := json.NewDecoder(r.Body).Decode(&ap); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	created, err := storage.Store.CreateAppointment(ctx, &ap)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "failed to create appointment: "+err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusCreated, created)
}

func UpdateAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := utils.ExtractID(r.URL.Path, "/appointments/")

	var updated models.Appointment
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	updated.ID = id

	if err := storage.Store.UpdateAppointment(ctx, &updated); err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.RespondError(w, http.StatusNotFound, "Appointment not found")
		} else {
			utils.RespondError(w, http.StatusInternalServerError, "update failed: "+err.Error())
		}
		return
	}
	utils.RespondJSON(w, http.StatusOK, updated)
}

func DeleteAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := utils.ExtractID(r.URL.Path, "/appointments/")

	if err := storage.Store.DeleteAppointment(ctx, id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.RespondError(w, http.StatusNotFound, "Appointment not found")
		} else {
			utils.RespondError(w, http.StatusInternalServerError, "delete failed: "+err.Error())
		}
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Appointment deleted"})
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
