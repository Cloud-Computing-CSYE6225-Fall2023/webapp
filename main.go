package main

import (
	"database/sql"
	"fmt"
	"github.com/shivasaicharanruthala/webapp/publish"
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
	"github.com/shivasaicharanruthala/webapp/log"
	"github.com/shivasaicharanruthala/webapp/mailer"
	"github.com/shivasaicharanruthala/webapp/middleware"
	accountSvc "github.com/shivasaicharanruthala/webapp/service/account"
	assignmentSvc "github.com/shivasaicharanruthala/webapp/service/assignment"
	"github.com/shivasaicharanruthala/webapp/store/account"
	"github.com/shivasaicharanruthala/webapp/store/assignment"
	"github.com/shivasaicharanruthala/webapp/types"
	"gopkg.in/alexcesaro/statsd.v2"
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
	//metricServerPort := os.Getenv("METRIC_SERVER_PORT")
	logFilePath := os.Getenv("LOG_FILE_PATH")
	migrationFilePath := os.Getenv("MIGRATION_FILE_PATH")

	connectionStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)

	// Initialize Logger
	logger, err := log.NewCustomLogger(logFilePath)
	if err != nil {
		lm := log.Message{Level: "ERROR", ErrorMessage: fmt.Sprintf("Initiating logger with error %v", err.Error())}
		logger.Log(&lm)
	}

	lm := log.Message{Level: "INFO", Msg: "Logger initialized successfully"}
	logger.Log(&lm)

	// Initialize Metric Server
	metricClient, err := statsd.New()
	if err != nil {
		lm = log.Message{Level: "ERROR", ErrorMessage: fmt.Sprintf("Initiating stasd metricClient with error %v", err.Error())}
		logger.Log(&lm)
	}

	lm = log.Message{Level: "INFO", Msg: "StatsD metric client initialized successfully"}
	logger.Log(&lm)

	defer metricClient.Close()

	// Initialize MailerSend Client
	mailerClient := mailer.New()

	// Initialize SNS Client
	snsClient, err := publish.New(logger)
	if err != nil {
		lm = log.Message{Level: "ERROR", ErrorMessage: fmt.Sprintf("Error Initilizing SNS client config with error %v", err.Error())}
		logger.Log(&lm)
	}

	// Initialize Context
	ctx := types.NewContext(logger, metricClient, mailerClient, snsClient)

	// Initialize DB connection
	db, err := sql.Open(dbDriver, connectionStr)
	defer func(db *sql.DB) {
		er := db.Close()
		if er != nil {
			lm = log.Message{Level: "ERROR", ErrorMessage: fmt.Sprintf("Closing Connection to database failed with error %v", er.Error())}
			logger.Log(&lm)
		}
	}(db)
	if err != nil {
		lm = log.Message{Level: "ERROR", ErrorMessage: fmt.Sprintf("Connection opening to Database failed with error %v", err.Error())}
		logger.Log(&lm)
	}

	err = db.Ping()
	if err != nil {
		lm = log.Message{Level: "ERROR", ErrorMessage: fmt.Sprintf("Ping to Database failed with error %v", err.Error())}
		logger.Log(&lm)
	}

	// Run Migrations
	if err == nil {
		driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			lm = log.Message{Level: "ERROR", ErrorMessage: fmt.Sprintf("Connecting to database for migration failed with error %v", err.Error())}
			logger.Log(&lm)
		}

		m, err := migrate.NewWithDatabaseInstance(migrationFilePath, dbName, driver)
		if err != nil {
			lm = log.Message{Level: "ERROR", ErrorMessage: fmt.Sprintf("Connecting to database for migration failed with error %v", err.Error())}
			logger.Log(&lm)
		}

		if m != nil && err == nil {
			err = m.Up()
			if err != nil {
				if err.Error() != "no change" {
					lm = log.Message{Level: "ERROR", ErrorMessage: fmt.Sprintf("Migration to schemas failed with error %v", err.Error())}
					logger.Log(&lm)
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
	assignmentHandler := aHandler.New(ctx, assgnmntSvc)
	//accountHandler := accntHandler.New(ctx, accntSvc)

	// Load test user accounts from a given file
	fileName := os.Getenv("USER_DATA_FILE_PATH")
	if err == nil {
		if err = accntSvc.BulkInsert(ctx, fileName); err != nil {
			lm = log.Message{Level: "ERROR", ErrorMessage: fmt.Sprintf("Loading users data to database failed with error %v", err.Error())}
			logger.Log(&lm)
		}
	}

	// Setup router using mux
	router := mux.NewRouter().StrictSlash(true)
	router.MethodNotAllowedHandler = http.HandlerFunc(MethodNotImplementedHandler)
	router.Use(middleware.Logging(logger))
	//router.Use(middleware.BasicAuths(ctx, accntSvc))
	router.Use(middleware.APICountMetrics(ctx))

	// Health Check Route
	router.HandleFunc("/healthz", HealthCheckHandler).Methods("GET")

	// Accounts Route
	//router.HandleFunc("/v1/account", accountHandler.Insert).Methods("POST")

	// Assignments Routes
	router.Handle("/v1/assignments", middleware.NewBasicAuth(ctx, assignmentHandler.Get, accntSvc)).Methods("GET")
	router.Handle("/v1/assignments", middleware.NewBasicAuth(ctx, assignmentHandler.Insert, accntSvc)).Methods("POST")
	router.Handle("/v1/assignments/{id}", middleware.NewBasicAuth(ctx, assignmentHandler.GetById, accntSvc)).Methods("GET")
	router.Handle("/v1/assignments/{id}", middleware.NewBasicAuth(ctx, assignmentHandler.Modify, accntSvc)).Methods("PUT")
	router.Handle("/v1/assignments/{id}", middleware.NewBasicAuth(ctx, assignmentHandler.Delete, accntSvc)).Methods("DELETE")
	router.Handle("/v1/assignments/{id}/submission", middleware.NewBasicAuth(ctx, assignmentHandler.PostAssignmentSubmission, accntSvc)).Methods("POST")

	// Start the server
	port := os.Getenv("PORT")
	server := fmt.Sprintf(":%s", port)

	lm = log.Message{Level: "INFO", Msg: fmt.Sprintf("Webapp Server starting to listen on port %v", port)}
	logger.Log(&lm)

	err = http.ListenAndServe(server, router)
	if err != nil {
		lm = log.Message{Level: "ERROR", ErrorMessage: fmt.Sprintf("Initializing weapp server to listen on port %v with error %v", port, err.Error())}
		logger.Log(&lm)
	}
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
			return
		}
	}(db)

	if err != nil {
		SetNoCacheResponseHeaders(w)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	err = db.Ping()
	if err != nil {
		SetNoCacheResponseHeaders(w)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	SetNoCacheResponseHeaders(w)
	w.WriteHeader(http.StatusOK)
	return
}

func MethodNotImplementedHandler(w http.ResponseWriter, r *http.Request) {
	SetNoCacheResponseHeaders(w)
	w.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func SetNoCacheResponseHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Pragma", "no-cache")
}
