package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type UserInfo struct {
	ID                int            `json:"ID"`
	IPAddress         string         `json:"IP_Address"`
	AccessedParts     []AccessedPart `json:"Accessed_Parts"`
	TimeAccessed      string         `json:"Time_Accessed"`
	FirstTimeAccessed string         `json:"First_Time_Accessed"`
	LastTimeAccessed  string         `json:"Last_Time_Accessed"`
	Blacklisted       bool           `json:"Blacklisted"`
	// ClientData        string         `json:"Client_Data"`
	Country sql.NullString `json:"Country"`
	City    sql.NullString `json:"City"`
}

type AccessedPart struct {
	Part         string `json:"part"`
	TimeAccessed string `json:"time_accessed"`
}

func main() {
	// Database connection
	dsn := "users1234:User1234@tcp(127.0.0.1:3306)/user_logs"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Initialize Gin router
	router := gin.Default()

	// Serve static files
	router.Static("/public", "./public")

	// Define API route to fetch user data
	router.GET("/api/users", func(c *gin.Context) {
		users := fetchUserData(db)
		c.JSON(http.StatusOK, users)
	})

	// Define API route to fetch a single user's detailed information
	router.GET("/api/users/:id", func(c *gin.Context) {
		userID := c.Param("id")
		user, err := fetchSingleUser(db, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user details"})
			return
		}
		c.JSON(http.StatusOK, user)
	})

	// Define API route to delete a user
	router.DELETE("/api/users/:id", func(c *gin.Context) {
		userID := c.Param("id")

		// Delete related accessed_parts records
		_, err := db.Exec("DELETE FROM accessed_parts WHERE user_id = ?", userID)
		if err != nil {
			log.Printf("Error deleting accessed parts for user with ID %s: %v", userID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete related accessed parts"})
			return
		}

		// Delete the user
		result, err := db.Exec("DELETE FROM user_info WHERE id = ?", userID)
		if err != nil {
			log.Printf("Error deleting user with ID %s: %v", userID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Printf("Error fetching rows affected for user with ID %s: %v", userID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm user deletion"})
			return
		}

		if rowsAffected == 0 {
			log.Printf("No user found with ID %s", userID)
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		log.Printf("Successfully deleted user with ID %s", userID)
		c.Status(http.StatusOK)
	})

	// Define API route to toggle blacklist status
	router.PATCH("/api/users/:id/blacklist", func(c *gin.Context) {
		userID := c.Param("id")
		var payload struct {
			Blacklisted bool `json:"blacklisted"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		if _, err := db.Exec("UPDATE user_info SET blacklisted = ? WHERE id = ?", payload.Blacklisted, userID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update blacklist status"})
			return
		}
		c.Status(http.StatusOK)
	})

	// Serve the main HTML page
	router.GET("/", func(c *gin.Context) {
		c.File("./public/index.html")
	})

	// Serve the user details HTML page
	router.GET("/user/:id", func(c *gin.Context) {
		c.File("./public/user_details.html")
	})

	// Start the server on port 8081
	router.Run(":8080")
}

func fetchUserData(db *sql.DB) []UserInfo {
	var users []UserInfo

	// Updated query to fetch only the required fields
	rows, err := db.Query("SELECT id, ip, country, city, first_time_accessed, last_time_accessed, blacklisted FROM user_info")
	if err != nil {
		log.Fatalf("Error querying database: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var userInfo UserInfo
		var firstTimeAccessed, lastTimeAccessed []uint8

		err := rows.Scan(
			&userInfo.ID,
			&userInfo.IPAddress,
			&userInfo.Country,
			&userInfo.City,
			&firstTimeAccessed,
			&lastTimeAccessed,
			&userInfo.Blacklisted,
		)
		if err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}

		userInfo.FirstTimeAccessed = string(firstTimeAccessed)
		userInfo.LastTimeAccessed = string(lastTimeAccessed)

		// Append the userInfo to the users slice
		users = append(users, userInfo)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Error iterating over rows: %v", err)
	}

	return users
}

func fetchSingleUser(db *sql.DB, userID string) (UserInfo, error) {
	var userInfo UserInfo
	var timeAccessed, firstTimeAccessed, lastTimeAccessed []uint8

	userInfo.AccessedParts = []AccessedPart{}

	err := db.QueryRow("SELECT id, ip, time_accessed, first_time_accessed, last_time_accessed, blacklisted, country, city FROM user_info WHERE id = ?", userID).Scan(
		&userInfo.ID,
		&userInfo.IPAddress,
		&timeAccessed,
		&firstTimeAccessed,
		&lastTimeAccessed,
		&userInfo.Blacklisted,
		// &userInfo.ClientData,
		&userInfo.Country,
		&userInfo.City,
	)
	if err != nil {
		return userInfo, err
	}

	userInfo.TimeAccessed = string(timeAccessed)
	userInfo.FirstTimeAccessed = string(firstTimeAccessed)
	userInfo.LastTimeAccessed = string(lastTimeAccessed)

	partsRows, err := db.Query("SELECT part, time_accessed FROM accessed_parts WHERE user_id = ?", userInfo.ID)
	if err != nil {
		return userInfo, err
	}
	defer partsRows.Close()

	for partsRows.Next() {
		var part AccessedPart
		var partTimeAccessed []uint8
		err := partsRows.Scan(&part.Part, &partTimeAccessed)
		if err != nil {
			return userInfo, err
		}
		part.TimeAccessed = string(partTimeAccessed)
		userInfo.AccessedParts = append(userInfo.AccessedParts, part)
	}

	return userInfo, nil
}
