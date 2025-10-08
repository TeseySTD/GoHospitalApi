package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/TeseySTD/GoHospitalApi/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

var Store *Storage

// New creates storage wrapper
func New(pool *pgxpool.Pool) *Storage {
	return &Storage{pool: pool}
}

// ConnectPool - helper to create pool from connString
func ConnectPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}
	// optional: set some pool config
	config.MaxConns = 10
	config.MinConns = 1
	config.MaxConnIdleTime = time.Minute * 5

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

// Migrate creates tables if they do not exist
func (s *Storage) Migrate(ctx context.Context) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
CREATE TABLE IF NOT EXISTS doctors (
    id         integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    first_name text NOT NULL,
    last_name  text NOT NULL,
    specialization text NOT NULL,
    experience integer NOT NULL
);
`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
CREATE TABLE IF NOT EXISTS patients (
    id         integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    first_name text NOT NULL,
    last_name  text NOT NULL,
    age        integer NOT NULL,
    diagnosis  text
);
`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
CREATE TABLE IF NOT EXISTS appointments (
    id         integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    patient_id integer NOT NULL REFERENCES patients(id) ON DELETE CASCADE,
    doctor_id  integer NOT NULL REFERENCES doctors(id) ON DELETE SET NULL,
    date       date,
    time       time,
    status     text
);
`)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

//
// --- Patients CRUD ---
//

// CreatePatient inserts patient and returns created model (with ID)
func (s *Storage) CreatePatient(ctx context.Context, p *models.Patient) (*models.Patient, error) {
	row := s.pool.QueryRow(ctx, `
INSERT INTO patients (first_name, last_name, age, diagnosis)
VALUES ($1, $2, $3, $4)
RETURNING id
`, p.FirstName, p.LastName, p.Age, p.Diagnosis)

	var id int
	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	p.ID = id
	return p, nil
}

func (s *Storage) GetAllPatients(ctx context.Context) ([]models.Patient, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, first_name, last_name, age, diagnosis FROM patients ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Patient
	for rows.Next() {
		var p models.Patient
		err := rows.Scan(&p.ID, &p.FirstName, &p.LastName, &p.Age, &p.Diagnosis)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (s *Storage) GetPatientByID(ctx context.Context, id int) (*models.Patient, error) {
	row := s.pool.QueryRow(ctx, `SELECT id, first_name, last_name, age, diagnosis FROM patients WHERE id = $1`, id)
	var p models.Patient
	if err := row.Scan(&p.ID, &p.FirstName, &p.LastName, &p.Age, &p.Diagnosis); err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *Storage) UpdatePatient(ctx context.Context, p *models.Patient) error {
	ct, err := s.pool.Exec(ctx, `
UPDATE patients SET first_name=$1, last_name=$2, age=$3, diagnosis=$4 WHERE id=$5
`, p.FirstName, p.LastName, p.Age, p.Diagnosis, p.ID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("patient not found")
	}
	return nil
}

func (s *Storage) DeletePatient(ctx context.Context, id int) error {
	ct, err := s.pool.Exec(ctx, `DELETE FROM patients WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("patient not found")
	}
	return nil
}

//
// --- Doctors CRUD ---
//

func (s *Storage) CreateDoctor(ctx context.Context, d *models.Doctor) (*models.Doctor, error) {
	row := s.pool.QueryRow(ctx, `
INSERT INTO doctors (first_name, last_name, specialization, experience)
VALUES ($1, $2, $3, $4)
RETURNING id
`, d.FirstName, d.LastName, d.Specialization, d.Experience)
	var id int
	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	d.ID = id
	return d, nil
}

func (s *Storage) GetAllDoctors(ctx context.Context) ([]models.Doctor, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, first_name, last_name, specialization, experience FROM doctors ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Doctor
	for rows.Next() {
		var d models.Doctor
		if err := rows.Scan(&d.ID, &d.FirstName, &d.LastName, &d.Specialization, &d.Experience); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

func (s *Storage) GetDoctorByID(ctx context.Context, id int) (*models.Doctor, error) {
	row := s.pool.QueryRow(ctx, `SELECT id, first_name, last_name, specialization, experience FROM doctors WHERE id = $1`, id)
	var d models.Doctor
	if err := row.Scan(&d.ID, &d.FirstName, &d.LastName, &d.Specialization, &d.Experience); err != nil {
		return nil, err
	}
	return &d, nil
}

func (s *Storage) UpdateDoctor(ctx context.Context, d *models.Doctor) error {
	ct, err := s.pool.Exec(ctx, `
UPDATE doctors SET first_name=$1, last_name=$2, specialization=$3, experience=$4 WHERE id=$5
`, d.FirstName, d.LastName, d.Specialization, d.Experience, d.ID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("doctor not found")
	}
	return nil
}

func (s *Storage) DeleteDoctor(ctx context.Context, id int) error {
	ct, err := s.pool.Exec(ctx, `DELETE FROM doctors WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("doctor not found")
	}
	return nil
}

//
// --- Appointments CRUD ---
//

func (s *Storage) CreateAppointment(ctx context.Context, a *models.Appointment) (*models.Appointment, error) {
	row := s.pool.QueryRow(ctx, `
INSERT INTO appointments (patient_id, doctor_id, date, time, status)
VALUES ($1, $2, $3, $4, $5)
RETURNING id
`, a.PatientID, a.DoctorID, a.Date, a.Time, a.Status)
	var id int
	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	a.ID = id
	return a, nil
}

func (s *Storage) GetAllAppointments(ctx context.Context) ([]models.Appointment, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, patient_id, doctor_id, COALESCE(TO_CHAR(date,'YYYY-MM-DD'),'') , COALESCE(TO_CHAR(time,'HH24:MI:SS'),'') , status FROM appointments ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Appointment
	for rows.Next() {
		var a models.Appointment
		if err := rows.Scan(&a.ID, &a.PatientID, &a.DoctorID, &a.Date, &a.Time, &a.Status); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (s *Storage) GetAppointmentByID(ctx context.Context, id int) (*models.Appointment, error) {
	row := s.pool.QueryRow(ctx, `SELECT id, patient_id, doctor_id, COALESCE(TO_CHAR(date,'YYYY-MM-DD'),'') , COALESCE(TO_CHAR(time,'HH24:MI:SS'),'') , status FROM appointments WHERE id=$1`, id)
	var a models.Appointment
	if err := row.Scan(&a.ID, &a.PatientID, &a.DoctorID, &a.Date, &a.Time, &a.Status); err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *Storage) UpdateAppointment(ctx context.Context, a *models.Appointment) error {
	ct, err := s.pool.Exec(ctx, `
UPDATE appointments SET patient_id=$1, doctor_id=$2, date=$3, time=$4, status=$5 WHERE id=$6
`, a.PatientID, a.DoctorID, a.Date, a.Time, a.Status, a.ID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("appointment not found")
	}
	return nil
}

func (s *Storage) DeleteAppointment(ctx context.Context, id int) error {
	ct, err := s.pool.Exec(ctx, `DELETE FROM appointments WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("appointment not found")
	}
	return nil
}
