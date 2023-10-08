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
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}
}

func main() {
	connectionStr := os.Getenv("POSTGRESQL_CONNECTION_STRING")
	driverName := os.Getenv("DRIVER_NAME")
	migrationFilePath := os.Getenv("MIGRATION_FILE_PATH")
	databaseName := os.Getenv("DB_NAME")

	db, err := sql.Open(driverName, connectionStr)
	defer func(db *sql.DB) {
		er := db.Close()
		if er != nil {
			fmt.Printf("ERROR: Closing Connection to database failed with error %v", er.Error())
			//panic(errors.NewCustomError(er))
		}
	}(db)
	if err != nil {
		fmt.Printf("ERROR: Connection opening to Database failed with error %v", err.Error())
		//panic(errors.NewCustomError(err))
	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("ERROR: Ping to Database failed with error %v", err.Error())
		//panic(errors.NewCustomError(err))
	}

	// Run Migrations
	if err == nil {
		driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			fmt.Printf("ERROR: Connecting to database for migration failed with error %v", err.Error())
			//panic(errors.NewCustomError(err))
		}

		m, err := migrate.NewWithDatabaseInstance(migrationFilePath, databaseName, driver)
		if err != nil {
			fmt.Printf("ERROR: Connecting to database for migration failed with error %v", err.Error())
			//panic(errors.NewCustomError(err))
		}

		if m != nil && err == nil {
			err = m.Up()
			if err != nil {
				if err.Error() != "no change" {
					fmt.Printf("ERROR: Migration to schemas failed with error %v", err.Error())
					//panic(errors.NewCustomError(err))
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
	serverAddr := os.Getenv("SERV_ADDR")
	port := os.Getenv("PORT")
	server := fmt.Sprintf("%s:%s", serverAddr, port)

	_ = http.ListenAndServe(server, router)
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	connectionStr := os.Getenv("POSTGRESQL_CONNECTION_STRING")
	driverName := os.Getenv("DRIVER_NAME")

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

	db, err := sql.Open(driverName, connectionStr)
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
