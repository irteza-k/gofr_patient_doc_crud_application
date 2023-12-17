package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"gofr.dev/examples/using-postgres/handler"
	"gofr.dev/examples/using-postgres/store"
	"gofr.dev/pkg/gofr"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Doctor struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Specialization string `json:"specialization"`
}

type Patient struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Age           int    `json:"age"`
	Gender        string `json:"gender"`
	Problem       string `json:"problem"`
	Contact       string `json:"contact"`
	AdmissionDate string `json:"admission_date"`
	Doctor        Doctor `json:"doctor"`
}
type Handler struct {
	DB *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{DB: db}
}

func main() {
	app := gofr.New()

	s := store.New()
	h := handler.New(s)
	//connect to database
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//create the table if it doesn't exist
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT, email TEXT)")

	if err != nil {
		log.Fatal(err)
	}

	// specifying the different routes supported by this service
	app.GET("/patient", h.Get)
	app.GET("/patient/{id}", h.GetByID)
	app.POST("/patient", h.Create)
	app.PUT("/patient/{id}", h.Update)
	app.DELETE("/patient/{id}", h.Delete)

	// specifying the different routes supported by this service
	app.GET("/Doctor", h.Get)
	app.GET("/Doctor/{id}", h.GetByID)
	app.POST("/Doctor", h.Create)
	app.PUT("/Doctor/{id}", h.Update)
	app.DELETE("/Doctor/{id}", h.Delete)

	//start server
	log.Fatal(http.ListenAndServe(":8000", jsonContentTypeMiddleware(app)))
}

func jsonContentTypeMiddleware(next *gofr.Gofr) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// fetch al patients

func (h *Handler) GetPatients(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM patients")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		patients := []Patient{}
		for rows.Next() {
			var p Patient
			if err := rows.Scan(&p.ID, &p.Name, &p.Age, &p.Gender, &p.Problem, &p.Contact, &p.AdmissionDate, &p.Doctor.ID, &p.Doctor.Name, &p.Doctor.Specialization); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			patients = append(patients, p)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Encode the list of patients as JSON and send it as the response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(patients)
	}
}

// patient by id
func (h *Handler) GetPatientByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var patient Patient
		err := db.QueryRow("SELECT * FROM patients WHERE id = $1", id).
			Scan(&patient.ID, &patient.Name, &patient.Age, &patient.Gender, &patient.Problem, &patient.Contact, &patient.AdmissionDate, &patient.Doctor.ID, &patient.Doctor.Name, &patient.Doctor.Specialization)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Patient not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Encode the patient information as JSON and send it as the response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(patient)
	}
}

// insert new patient
func (h *Handler) CreatePatient(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newPatient Patient
		if err := json.NewDecoder(r.Body).Decode(&newPatient); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Insert a new patient into the database
		_, err := db.Exec("INSERT INTO patients (name, age, gender, problem, contact, admission_date, doctor_id) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			newPatient.Name, newPatient.Age, newPatient.Gender, newPatient.Problem, newPatient.Contact, newPatient.AdmissionDate, newPatient.Doctor.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newPatient)
	}
}

//update patient

func (h *Handler) UpdatePatient(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the patient ID from the request URL
		vars := mux.Vars(r)
		id := vars["id"]

		var updatedPatient Patient
		if err := json.NewDecoder(r.Body).Decode(&updatedPatient); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Update the patient in the database
		_, err := db.Exec("UPDATE patients SET name=$1, age=$2, gender=$3, problem=$4, contact=$5, admission_date=$6, doctor_id=$7 WHERE id=$8",
			updatedPatient.Name, updatedPatient.Age, updatedPatient.Gender, updatedPatient.Problem, updatedPatient.Contact, updatedPatient.AdmissionDate, updatedPatient.Doctor.ID, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(updatedPatient)
	}
}

func (h *Handler) DeletePatient(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the patient ID from the request URL
		vars := mux.Vars(r)
		id := vars["id"]

		// Delete the patient from the database
		result, err := db.Exec("DELETE FROM patients WHERE id = $1", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check if the delete operation affected any rows
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			http.Error(w, "Patient not found", http.StatusNotFound)
			return
		}

		// Respond with a success message
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Patient deleted")
	}
}

// doctor crud operations
// fetch all doc
func (h *Handler) GetDoctors(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Fetch all doctors from the database
		rows, err := db.Query("SELECT * FROM doctors")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		doctors := []Doctor{}
		for rows.Next() {
			var d Doctor
			if err := rows.Scan(&d.ID, &d.Name, &d.Specialization); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			doctors = append(doctors, d)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Encode the list of doctors as JSON and send it as the response
		json.NewEncoder(w).Encode(doctors)
	}
}

// doc by id
func (h *Handler) GetDoctorByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Fetch a doctor by ID from the database
		vars := mux.Vars(r)
		id := vars["id"]

		var doctor Doctor
		err := db.QueryRow("SELECT * FROM doctors WHERE id = $1", id).Scan(&doctor.ID, &doctor.Name, &doctor.Specialization)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Doctor not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Encode the doctor information as JSON and send it as the response
		json.NewEncoder(w).Encode(doctor)
	}
}

//create doc

func (h *Handler) CreateDoctor(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var doctor Doctor
		err := json.NewDecoder(r.Body).Decode(&doctor)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := db.Exec("INSERT INTO doctors (name, specialization) VALUES ($1, $2)", doctor.Name, doctor.Specialization)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		newDoctorID, err := result.LastInsertId()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		doctor.ID = int(newDoctorID)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(doctor)
	}
}

// update doc
func (h *Handler) UpdateDoctor(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var doctor Doctor
		err := json.NewDecoder(r.Body).Decode(&doctor)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = db.Exec("UPDATE doctors SET name = $1, specialization = $2 WHERE id = $3", doctor.Name, doctor.Specialization, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		doctor.ID, _ = strconv.Atoi(id)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(doctor)
	}
}

// delete doc
func (h *Handler) DeleteDoctor(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		_, err := db.Exec("DELETE FROM doctors WHERE id = $1", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Doctor deleted")
	}
}
