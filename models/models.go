package models

type Patient struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
	Diagnosis string `json:"diagnosis"`
}

type Doctor struct {
	ID             int    `json:"id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Specialization string `json:"specialization"`
	Experience     int    `json:"experience"`
}

type Appointment struct {
	ID        int    `json:"id"`
	PatientID int    `json:"patient_id"`
	DoctorID  int    `json:"doctor_id"`
	Date      string `json:"date"`
	Time      string `json:"time"`
	Status    string `json:"status"`
}