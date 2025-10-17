package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/TeseySTD/GoHospitalApi/handlers"
	"github.com/TeseySTD/GoHospitalApi/middleware"
	"github.com/TeseySTD/GoHospitalApi/storage"
)

func main() {
	ctx := context.Background()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5400/hospitaldb"
	}

	pool, err := storage.ConnectPool(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer pool.Close()

	st := storage.New(pool)
	storage.Store = st

	if err := st.Migrate(ctx); err != nil {
		log.Fatalf("migrate failed: %v", err)
	}
	log.Println("Database migrated/ready")

	//Public endpoints
	http.HandleFunc("/", middleware.LoggingMiddleware(handlers.RootHandler))
	http.HandleFunc("/login", middleware.LoggingMiddleware(handlers.LoginHandler))
	http.HandleFunc("/users", middleware.LoggingMiddleware(handlers.UsersListHandler))

	// Protected endpoints
	http.HandleFunc("/patients", middleware.Chain(
		handlers.PatientsRouter,
		middleware.LoggingMiddleware,
		middleware.JWTAuthMiddleware,
		middleware.RoleBasedAccess,
	))
	http.HandleFunc("/patients/", middleware.Chain(
		handlers.PatientsRouter,
		middleware.LoggingMiddleware,
		middleware.JWTAuthMiddleware,
		middleware.RoleBasedAccess,
	))

	http.HandleFunc("/doctors", middleware.Chain(
		handlers.DoctorsRouter,
		middleware.LoggingMiddleware,
		middleware.JWTAuthMiddleware,
		middleware.RoleBasedAccess,
	))
	http.HandleFunc("/doctors/", middleware.Chain(
		handlers.DoctorsRouter,
		middleware.LoggingMiddleware,
		middleware.JWTAuthMiddleware,
		middleware.RoleBasedAccess,
	))

	http.HandleFunc("/appointments", middleware.Chain(
		handlers.AppointmentsRouter,
		middleware.LoggingMiddleware,
		middleware.JWTAuthMiddleware,
		middleware.RoleBasedAccess,
	))
	http.HandleFunc("/appointments/", middleware.Chain(
		handlers.AppointmentsRouter,
		middleware.LoggingMiddleware,
		middleware.JWTAuthMiddleware,
		middleware.RoleBasedAccess,
	))

	port := ":8080"
	log.Printf("Server starting on port %s", port)

	srv := &http.Server{
		Addr:         port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
