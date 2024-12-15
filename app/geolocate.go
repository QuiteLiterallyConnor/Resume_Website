package app

import (
	"bufio"
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/oschwald/geoip2-golang"
)

// LoadGeoIPDatabases loads the GeoIP2 databases
func LoadGeoIPDatabases() (*geoip2.Reader, error) {
	cityDB, err := geoip2.Open("app/geolite/GeoLite2-City.mmdb")
	if err != nil {
		return nil, err
	}
	return cityDB, nil
}

// LoadWhitelistedCountries loads the list of allowed countries from a file
func LoadWhitelistedCountries(filePath string) (map[string]bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	countries := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			countries[line] = true
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return countries, nil
}

// verifyGeoIP verifies the country and logs the city for the client IP
func verifyGeoIP(cityDB *geoip2.Reader, allowedCountries map[string]bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := net.ParseIP(c.ClientIP())
		requestPath := c.Request.URL.Path
		clientData := c.Request.UserAgent()

		log.Printf("Client IP: %s, Request Path: %s", ip, requestPath)

		var country, city string
		var blacklisted bool

		if ip.IsLoopback() || ip.IsPrivate() {
			country, city = "US", "Local"
		} else {
			cityRecord, err := cityDB.City(ip)
			if err != nil {
				log.Printf("GeoIP city error for IP %s: %v", ip, err)
				city, country = "Unknown", "Unknown"
			} else {
				city = cityRecord.City.Names["en"]
				country = cityRecord.Country.IsoCode
			}
		}

		log.Printf("IP: %s, Country: %s, City: %s", ip, country, city)

		// Log the IP, country, city, and additional data
		logUserIP(ip.String(), country, city, requestPath, clientData, blacklisted)

		if !allowedCountries[country] {
			log.Printf("Rejected country: %s", country)
			blacklisted = true
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access restricted"})
		} else {
			c.Next()
		}
	}
}

// logUserIP logs the IP, country, city, accessed parts, and client data to the database
func logUserIP(ip, country, city, accessedParts, clientData string, blacklisted bool) {
	timeAccessed := time.Now()

	// Check if the IP already exists to update the last accessed time
	var userID int
	var firstTimeAccessedRaw []uint8

	err := db.QueryRow("SELECT id, first_time_accessed FROM user_info WHERE ip = ?", ip).Scan(&userID, &firstTimeAccessedRaw)
	if err == sql.ErrNoRows {
		// Insert new record with country and city
		result, err := db.Exec(`INSERT INTO user_info 
			(ip, country, city, time_accessed, first_time_accessed, last_time_accessed, blacklisted, client_data) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			ip, country, city, timeAccessed, timeAccessed, timeAccessed, blacklisted, clientData)
		if err != nil {
			log.Printf("Failed to insert user info into database: %v", err)
			return
		}
		userID64, err := result.LastInsertId()
		if err != nil {
			log.Printf("Failed to get last insert ID: %v", err)
			return
		}
		userID = int(userID64)
	} else if err == nil {
		// Update existing record with new last_time_accessed, country, and city
		_, err = db.Exec(`UPDATE user_info 
			SET last_time_accessed = ?, country = ?, city = ?, blacklisted = ?, client_data = ? 
			WHERE id = ?`,
			timeAccessed, country, city, blacklisted, clientData, userID)
		if err != nil {
			log.Printf("Failed to update user info in database: %v", err)
			return
		}
	} else {
		log.Printf("Database query error: %v", err)
		return
	}

	// Log the accessed parts in the accessed_parts table
	_, err = db.Exec(`INSERT INTO accessed_parts (user_id, part, time_accessed) 
		VALUES (?, ?, ?)`, userID, accessedParts, timeAccessed)
	if err != nil {
		log.Printf("Failed to log accessed part: %v", err)
	}
}
