package app

import (
    "log"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/oschwald/geoip2-golang"
)

// SetupGinRouter configures the Gin router and middleware
func SetupGinRouter(logFile *os.File, cityDB *geoip2.Reader, allowedCountries map[string]bool) *gin.Engine {
    r := gin.New()
    r.Use(gin.LoggerWithWriter(logFile), gin.Recovery())

    // Apply the blacklistSensitive middleware before verifyGeoIP
    r.Use(checkBlacklist)
    r.Use(blacklistSensitive) // This should come before verifyGeoIP
    r.Use(verifyGeoIP(cityDB, allowedCountries))

    r.GET("/robots.txt", func(c *gin.Context) {
        log.Println("robots.txt requested")
        c.File("./app/public/robots.txt")
    })

    // Serve static files under /public
    r.Static("/public", "./app/public")

    // Serve the index.html file for the root
    r.StaticFile("/", "./app/public/index.html")

    return r
}
