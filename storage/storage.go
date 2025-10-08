package storage

import (
	"encoding/json"
	"github.com/TeseySTD/GoHospitalApi/models"
	"os"
	"sync"
)

type Storage struct {
	Patients     []models.Patient     `json:"patients"`
	Doctors      []models.Doctor      `json:"doctors"`
	Appointments []models.Appointment `json:"appointments"`
	mutex        sync.RWMutex
}

var Store = &Storage{
	Patients:     []models.Patient{},
	Doctors:      []models.Doctor{},
	Appointments: []models.Appointment{},
}

const storageFile = "hospital_data.json"

func (s *Storage) LoadFromFile() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	data, err := os.ReadFile(storageFile)
	if err != nil {
		return s.SaveToFile()
	}

	return json.Unmarshal(data, s)
}

func (s *Storage) SaveToFile() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(storageFile, data, 0644)
}

func (s *Storage) Lock() {
	s.mutex.Lock()
}

func (s *Storage) Unlock() {
	s.mutex.Unlock()
}

func (s *Storage) RLock() {
	s.mutex.RLock()
}

func (s *Storage) RUnlock() {
	s.mutex.RUnlock()
}
