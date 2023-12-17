This project implements a REST API for managing patients and doctors using Golang, Docker, PostgreSQL, and Postman. It allows CRUD operations (Create, Read, Update, Delete) on both patients and doctors.

Technology Stack:

    Backend: Golang
    Database: PostgreSQL
    Containerization: Docker
    Testing: Postman

Features:

    Create, read, update, and delete patients
    Manage patient details including name, age, gender, problem, contact, admission date, and assigned doctor
    Create, read, update, and delete doctors
    Manage doctor details including name and specialization
    JSON-based API endpoints for easy integration

Getting Started:

    Clone the repository:

git clone https://github.com/irteza-k/gofr_patient_doc_crud_application.git

    Set up the environment:

    Docker: Install Docker Desktop or Engine if not already installed.
    PostgreSQL:
        Run docker-compose up -d to start the PostgreSQL container.
        Alternatively, follow the instructions in docker-compose.yml to manually configure and run the container.
    Postman: Download and install Postman.

    Build and run the API:

make build
make run

    Test with Postman:

    Import the provided patient-and-doctor-api.postman_collection.json file into Postman.
    This collection contains pre-defined requests for all API endpoints.
    Follow the instructions in each request to test the API functionality.

API Endpoints:
Method	Endpoint	Description
POST	/patients	Create a new patient
GET	/patients	Get all patients
GET	/patients/:id	Get a specific patient by ID
PUT	/patients/:id	Update a patient
DELETE	/patients/:id	Delete a patient
POST	/doctors	Create a new doctor
GET	/doctors	Get all doctors
GET	/doctors/:id	Get a specific doctor by ID
PUT	/doctors/:id	Update a doctor
DELETE	/doctors/:id	Delete a doctor

Docker and Database:

    This project utilizes Docker to simplify deployment and isolation.
    The docker-compose.yml file defines two services:
        go-app: Builds and runs the Golang app.
        go-db: Runs a PostgreSQL container with pre-configured settings.
    You can easily launch and manage both services using docker-compose.