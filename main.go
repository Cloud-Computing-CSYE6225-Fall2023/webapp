package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	aHandler "github.com/shivasaicharanruthala/webapp/handler/assignment"
	"github.com/shivasaicharanruthala/webapp/middleware"
	accountSvc "github.com/shivasaicharanruthala/webapp/service/account"
	assignmentSvc "github.com/shivasaicharanruthala/webapp/service/assignment"
	"github.com/shivasaicharanruthala/webapp/store/account"
	"github.com/shivasaicharanruthala/webapp/store/assignment"
)

func init() {
	for _, envFilePath := range EnvFilePaths {
		err := godotenv.Load(envFilePath)
		if err != nil {
			fmt.Printf("Error loading %v file\n", envFilePath)
		} else {
			break
		}
	}
}

func main() {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbDriver := os.Getenv("DRIVER_NAME")
	migrationFilePath := os.Getenv("MIGRATION_FILE_PATH")

	connectionStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := sql.Open(dbDriver, connectionStr)
	defer func(db *sql.DB) {
		er := db.Close()
		if er != nil {
			fmt.Printf("ERROR: Closing Connection to database failed with error %v", er.Error())
		}
	}(db)
	if err != nil {
		fmt.Printf("ERROR: Connection opening to Database failed with error %v", err.Error())
	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("ERROR: Ping to Database failed with error %v", err.Error())
	}

	// Run Migrations
	if err == nil {
		driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			fmt.Printf("ERROR: Connecting to database for migration failed with error %v", err.Error())
		}

		m, err := migrate.NewWithDatabaseInstance(migrationFilePath, dbName, driver)
		if err != nil {
			fmt.Printf("ERROR: Connecting to database for migration failed with error %v", err.Error())
		}

		if m != nil && err == nil {
			err = m.Up()
			if err != nil {
				if err.Error() != "no change" {
					fmt.Printf("ERROR: Migration to schemas failed with error %v", err.Error())
				}
			}
		}
	}

	// Store Layer
	assignmentStore := assignment.New(db)
	accountStore := account.New(db)

	// Service Layer
	assgnmntSvc := assignmentSvc.New(assignmentStore)
	accntSvc := accountSvc.New(accountStore)

	// Handler Layer
	assignmentHandler := aHandler.New(assgnmntSvc)
	//accountHandler := accntHandler.New(accntSvc)

	// Load test user accounts from a given file
	fileName := os.Getenv("USER_DATA_FILE_PATH")
	if err == nil {
		if err = accntSvc.BulkInsert(fileName); err != nil {
			fmt.Printf("ERROR: Loading users data to database failed with error %v", err.Error())
			//panic(errors.NewCustomError(err))
		}
	}

	// Setup router using mux
	router := mux.NewRouter().StrictSlash(true)
	router.MethodNotAllowedHandler = http.HandlerFunc(MethodNotImplementedHandler)

	// Health Check Route
	router.HandleFunc("/healthz", HealthCheckHandler).Methods("GET")

	// Accounts Route
	//router.HandleFunc("/v1/account", accountHandler.Insert).Methods("POST")

	// Assignments Routes
	router.Handle("/v1/assignments", middleware.NewBasicAuth(assignmentHandler.Get, accntSvc)).Methods("GET")
	router.Handle("/v1/assignments", middleware.NewBasicAuth(assignmentHandler.Insert, accntSvc)).Methods("POST")
	router.Handle("/v1/assignments/{id}", middleware.NewBasicAuth(assignmentHandler.GetById, accntSvc)).Methods("GET")
	router.Handle("/v1/assignments/{id}", middleware.NewBasicAuth(assignmentHandler.Modify, accntSvc)).Methods("PUT")
	router.Handle("/v1/assignments/{id}", middleware.NewBasicAuth(assignmentHandler.Delete, accntSvc)).Methods("DELETE")

	// Start the server
	port := os.Getenv("PORT")
	server := fmt.Sprintf(":%s", port)

	_ = http.ListenAndServe(server, router)
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbDriver := os.Getenv("DRIVER_NAME")

	connectionStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)

	if len(r.URL.Query()) > 0 {
		SetNoCacheResponseHeaders(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		SetNoCacheResponseHeaders(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the request body is not empty
	if len(body) > 0 {
		SetNoCacheResponseHeaders(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db, err := sql.Open(dbDriver, connectionStr)
	defer func(db *sql.DB) {
		er := db.Close()
		if er != nil {
			SetNoCacheResponseHeaders(w)
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	}(db)

	if err != nil {
		SetNoCacheResponseHeaders(w)
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	err = db.Ping()
	if err != nil {
		SetNoCacheResponseHeaders(w)
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	SetNoCacheResponseHeaders(w)
	w.WriteHeader(http.StatusOK)
	return
}

func MethodNotImplementedHandler(w http.ResponseWriter, r *http.Request) {
	SetNoCacheResponseHeaders(w)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func SetNoCacheResponseHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Pragma", "no-cache")
}
