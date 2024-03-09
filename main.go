package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/joho/godotenv"

	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"integrated-library-service/domain"
	"integrated-library-service/handlers"
	"integrated-library-service/middleware"
	"integrated-library-service/routes"

	_ "github.com/lib/pq" // pq driver.
)

var (
	// Derived from ldflags -X
	buildRevision string
	buildVersion  string
	buildTime     string

	// general options
	versionFlag bool
	helpFlag    bool

	// server port
	port string

	// other variables
	secretKey string

	// program controller
	done      = make(chan struct{})
	doneRetry = make(chan bool)
	errc      = make(chan error)
)

func init() {
	// Getting secret key
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	secretKey = os.Getenv("JWT_SECRET_KEY")
	flag.BoolVar(&versionFlag, "version", false, "show current version and exit")
	flag.BoolVar(&helpFlag, "help", false, "show usage and exit")
	flag.StringVar(&port, "port", ":8000", "server port")
}

func setBuildVariables() {
	if buildRevision == "" {
		buildRevision = "dev"
	}
	if buildVersion == "" {
		buildVersion = "dev"
	}
	if buildTime == "" {
		buildTime = time.Now().UTC().Format(time.RFC3339)
	}
}

func parseFlags() {
	flag.Parse()

	if helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	if versionFlag {
		fmt.Printf("%s %s %s\n", buildRevision, buildVersion, buildTime)
		os.Exit(0)
	}
}

func handleInterrupts() {
	log.Println("start handle interrupts")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	sig := <-interrupt
	log.Printf("caught sig: %v", sig)
	// close resource here
	done <- struct{}{}
}

func openDB() (*sql.DB, error) {
	var (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "postgres"
		dbname   = "integrated-library-management-service"
	)
	psqlInfo := os.Getenv("POSTGRESQL_CONN_STRING")
	if len(psqlInfo) == 0 {
		psqlInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	}
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	setBuildVariables()
	parseFlags()

	go handleInterrupts()

	server := gin.Default()
	ilmGroup := server.Group("ilm-service/v1")
	// initializing cors
	server.Use(middleware.CORS())
	authMiddleware := middleware.NewAuthMiddleware(secretKey)

	db, err := openDB()
	if err != nil {
		log.Printf("error connecting DB: %v", err)
		return
	}
	log.Println("DB connection is successful")
	defer db.Close()

	// create library service
	libraryService := domain.NewLibraryService(db)

	libraryHandler := handlers.NewLibraryHandler(libraryService, secretKey)
	apiRoutes := routes.NewRoutes(libraryHandler)
	routes.AttachRoutes(ilmGroup, apiRoutes, authMiddleware)

	go func() {
		errc <- server.Run(port)
	}()

	select {
	case err := <-errc:
		log.Printf("ListenAndServe error: %v", err)
	case <-done:
		log.Println("shutting down server ...")
	}

	time.AfterFunc(1*time.Second, func() {
		close(done)
		close(errc)
		close(doneRetry)
	})
}
