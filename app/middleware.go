package app

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/oschwald/geoip2-golang"
)

func checkBlacklist(c *gin.Context) {
	if !dbEnabled {
		log.Println("Database disabled. Skipping blacklist check.")
		c.Next()
		return
	}
	ip := c.ClientIP()

	// Check if the request is behind a proxy and use the forwarded IP
	if forwarded := c.Request.Header.Get("X-Forwarded-For"); forwarded != "" {
		ip = forwarded
	}

	log.Printf("Checking IP: %s for blacklist", ip)

	// Check if the IP is blacklisted
	if blacklisted, err := isIPBlacklisted(ip); err != nil {
		log.Printf("Blacklist error for IP %s: %v", ip, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	} else if blacklisted {
		log.Printf("Blocked blacklisted IP: %s", ip)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access restricted"})
		return
	}

	c.Next()
}

func blacklistSensitive(c *gin.Context) {
	sites, err := loadBlacklistedSites()
	if err != nil {
		log.Printf("Error loading blacklisted sites: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	requestPath := c.Request.URL.Path
	log.Printf("Requested path: %s", requestPath)

	for site := range sites {
		if strings.HasPrefix(requestPath, site) { // Check if the request path starts with a blacklisted site
			ip := c.ClientIP()
			log.Printf("Sensitive file access attempt %s by IP: %s", requestPath, ip)

			// Add IP to blacklist
			if err := addIPToBlacklist(ip); err != nil {
				log.Printf("Blacklist error: %v", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			} else {
				log.Printf("IP %s added to blacklist", ip)

				// Ensure the DB connection is active before updating the user_info table
				if err := ensureDatabaseAndTables(); err != nil {
					log.Printf("Database error while ensuring connection: %v", err)
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "database connection error"})
					return
				}

				_, err := db.Exec(`UPDATE user_info SET blacklisted = true WHERE ip = ?`, ip)
				if err != nil {
					log.Printf("Failed to update Blacklisted status in user_info: %v", err)
				}

				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access restricted"})
			}
			return
		}
	}

	log.Printf("Path %s not in blacklist", requestPath)
	c.Next()
}

func CreateMiddleware(r *gin.Engine, cityDB *geoip2.Reader, allowedCountries map[string]bool) {
	r.Use(checkBlacklist, verifyGeoIP(cityDB, allowedCountries), blacklistSensitive)
}
