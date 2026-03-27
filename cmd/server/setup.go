package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/srmdn/foliocms/internal/config"
	"github.com/srmdn/foliocms/internal/db"
)

func runSetup(database *db.DB, cfg *config.Config) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== Folio First-Time Setup ===")
	fmt.Println()

	email := cfg.AdminEmail
	if email == "" {
		fmt.Print("Admin email: ")
		line, _ := reader.ReadString('\n')
		email = strings.TrimSpace(line)
	}
	if email == "" {
		return fmt.Errorf("admin email is required")
	}

	passwd := cfg.AdminPasswd
	if passwd == "" {
		fmt.Print("Admin password: ")
		line, _ := reader.ReadString('\n')
		passwd = strings.TrimSpace(line)
	}
	if len(passwd) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hashing password: %w", err)
	}

	_, err = database.Exec(
		`INSERT INTO users (email, passwd_hash) VALUES (?, ?)
		 ON CONFLICT(email) DO UPDATE SET passwd_hash = excluded.passwd_hash`,
		email, string(hash),
	)
	if err != nil {
		return fmt.Errorf("saving admin user: %w", err)
	}

	fmt.Printf("Admin user created: %s\n", email)
	return nil
}
