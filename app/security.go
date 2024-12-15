package app

import (
	"bufio"
	"database/sql"
	"os"
	"strings"
)

const (
	blacklistFile        = "app/blacklisted_ips.txt"
	blacklistedSitesFile = "app/blacklisted_sites.txt"
)

// isIPBlacklisted checks if an IP is in the blacklist
func isIPBlacklisted(ip string) (bool, error) {
	var blacklisted bool
	err := db.QueryRow("SELECT blacklisted FROM user_info WHERE ip = ?", ip).Scan(&blacklisted)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // IP is not found in the database, so it's not blacklisted
		}
		return false, err // An actual error occurred
	}
	return blacklisted, nil
}

// addIPToBlacklist adds an IP to the blacklist
func addIPToBlacklist(ip string) error {
	file, err := os.OpenFile(blacklistFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(ip + "\n")
	return err
}

// loadBlacklistedSites loads blacklisted sites from a file into a map
func loadBlacklistedSites() (map[string]bool, error) {
	file, err := os.Open(blacklistedSitesFile)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]bool), nil
		}
		return nil, err
	}
	defer file.Close()

	sites := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		site := strings.TrimSpace(scanner.Text())
		site = strings.Trim(site, "\"") // Remove any surrounding quotes
		sites[site] = true
	}

	return sites, scanner.Err()
}
