package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	ngrok "golang.ngrok.com/ngrok"
	ngrokConfig "golang.ngrok.com/ngrok/config"
)

var (
	logFile     *os.File
	db          *sql.DB
	tunnel      ngrok.Tunnel
	stopChannel = make(chan bool)
	dbEnabled   = true // Flag to indicate whether DB functionality is active
)

func initLogger() (*os.File, error) {
	f, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	log.SetOutput(f)
	return f, nil
}

func initDBConnection() error {
	var err error
	serverDSN := "users1234:User1234@tcp(172.18.0.2:3306)/"
	db, err = sql.Open("mysql", serverDSN)
	if err != nil {
		dbEnabled = false
		log.Printf("Error opening database connection: %v. Database functionality disabled.", err)
		return nil // Return nil to avoid halting the server
	}

	if err = db.Ping(); err != nil {
		dbEnabled = false
		log.Printf("Error connecting to the database: %v. Database functionality disabled.", err)
		return nil
	}

	log.Println("Connected to MariaDB server successfully")
	return nil
}

func ensureDatabaseAndTables() error {
	// Ensure `user_logs` database exists
	_, err := db.Exec("CREATE DATABASE IF NOT EXISTS user_logs")
	if err != nil {
		return fmt.Errorf("error creating database: %v", err)
	}
	log.Println("Database `user_logs` ensured")

	// Reconnect to the `user_logs` database
	dsn := "users1234:User1234@tcp(172.18.0.2:3306)/user_logs"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("error connecting to `user_logs` database: %v", err)
	}

	// Set up connection pool settings
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection
	if err = db.Ping(); err != nil {
		return fmt.Errorf("error reconnecting to `user_logs` database: %v", err)
	}

	log.Println("Connected to `user_logs` database successfully")

	// Create the tables
	if err := createTables(); err != nil {
		return fmt.Errorf("error creating tables: %v", err)
	}

	return nil
}

func createTables() error {
	// `accessed_parts` table
	accessedPartsQuery := `
	CREATE TABLE IF NOT EXISTS accessed_parts (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_id INT NOT NULL,
		part VARCHAR(255) NOT NULL,
		time_accessed DATETIME NOT NULL
	)`
	if _, err := db.Exec(accessedPartsQuery); err != nil {
		return fmt.Errorf("error creating `accessed_parts` table: %v", err)
	}
	log.Println("Table `accessed_parts` ensured")

	// `user_info` table
	userInfoQuery := `
	CREATE TABLE IF NOT EXISTS user_info (
		id INT AUTO_INCREMENT PRIMARY KEY,
		ip VARCHAR(45) NOT NULL,
		accessed_parts TEXT,
		time_accessed DATETIME,
		first_time_accessed DATETIME,
		last_time_accessed DATETIME,
		blacklisted TINYINT(1) DEFAULT 0,
		client_data TEXT,
		country VARCHAR(100),
		city VARCHAR(100)
	)`
	if _, err := db.Exec(userInfoQuery); err != nil {
		return fmt.Errorf("error creating `user_info` table: %v", err)
	}
	log.Println("Table `user_info` ensured")

	return nil
}

func initDB() error {
	// Initialize database connection
	if err := initDBConnection(); err != nil {
		return err
	}

	// Ensure database and tables exist
	if err := ensureDatabaseAndTables(); err != nil {
		return err
	}

	return nil
}

func startNgrok(ctx context.Context) (ngrok.Tunnel, error) {
	authToken := os.Getenv("NGROK_AUTH_TOKEN")
	if authToken == "" {
		return nil, fmt.Errorf("NGROK_AUTH_TOKEN is not set in environment variables")
	}

	return ngrok.Listen(ctx,
		ngrokConfig.HTTPEndpoint(ngrokConfig.WithDomain("resume.connorisseur.com")),
		ngrok.WithAuthtoken(authToken),
	)
}

func StartServer() {
	var err error
	logFile, err = initLogger()
	if err != nil {
		log.Fatalf("Logger error: %v", err)
	}
	defer logFile.Close()

	if err = initDB(); err != nil {
		log.Fatalf("Database error: %v", err)
	}
	defer db.Close()

	cityDB, err := LoadGeoIPDatabases()
	if err != nil {
		log.Fatalf("GeoIP error: %v", err)
	}
	defer cityDB.Close()

	countries, err := LoadWhitelistedCountries("app/whitelisted_countries.txt")
	if err != nil {
		log.Fatalf("GeoIP error: %v", err)
	}

	r := SetupGinRouter(logFile, cityDB, countries)

	tunnel, err = startNgrok(context.Background())
	if err != nil {
		log.Fatalf("Ngrok error: %v", err)
	}

	log.Printf("Started Ngrok server at resume.connorisseur.com")
	go func() {
		if err := r.RunListener(tunnel); err != nil {
			log.Fatalf("Gin server error: %v", err)
		}
	}()

	<-stopChannel
}

func StopServer() {
	log.Println("Shutting down server...")

	if tunnel != nil {
		if err := tunnel.Close(); err != nil {
			log.Printf("Error closing Ngrok tunnel: %v", err)
		} else {
			log.Println("Ngrok tunnel closed")
		}
	}

	if db != nil {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		} else {
			log.Println("Database connection closed")
		}
	}

	// Closing the log file
	if logFile != nil {
		if err := logFile.Close(); err != nil {
			log.Printf("Error closing log file: %v", err)
		} else {
			log.Println("Log file closed")
		}
	}

	time.Sleep(2 * time.Second)

	log.Println("Server gracefully stopped")
	stopChannel <- true
}
