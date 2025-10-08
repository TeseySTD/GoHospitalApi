package handlers

import (
	"github.com/TeseySTD/GoHospitalApi/utils"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		utils.RespondError(w, http.StatusNotFound, "Not found")
		return
	}

	info := map[string]interface{}{
		"message": "Hospital REST API",
		"endpoints": map[string]string{
			"GET /patients":             "Get all patients",
			"GET /patients/{id}":        "Get patient by ID",
			"POST /patients":            "Create a new patient",
			"PUT /patients/{id}":        "Update patient",
			"DELETE /patients/{id}":     "Delete patient",
			"GET /doctors":              "Get all doctors",
			"GET /doctors/{id}":         "Get doctor by ID",
			"POST /doctors":             "Create a new doctor",
			"PUT /doctors/{id}":         "Update doctor",
			"DELETE /doctors/{id}":      "Delete doctor",
			"GET /appointments":         "Get all appointments",
			"GET /appointments/{id}":    "Get appointment by ID",
			"POST /appointments":        "Create a new appointment",
			"PUT /appointments/{id}":    "Update appointment",
			"DELETE /appointments/{id}": "Delete appointment",
		},
	}
	utils.RespondJSON(w, http.StatusOK, info)
}
