package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: ./tokengen <user> <days>")
		os.Exit(1)
	}

	user := os.Args[1]
	days, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Ungültige Tage:", os.Args[2])
		os.Exit(1)
	}

	// Secret laden oder neu erstellen
	secret, err := os.ReadFile("secrets/" + user + ".key")
	if err != nil {
		secret = make([]byte, 32)
		if _, err := rand.Read(secret); err != nil {
			fmt.Println("Secret generieren Fehler:", err)
			os.Exit(1)
		}

		os.MkdirAll("secrets", 0700)
		if err := os.WriteFile("secrets/"+user+".key", secret, 0600); err != nil {
			fmt.Println("Secret speichern Fehler:", err)
			os.Exit(1)
		}

		fmt.Println("Neues Secret erstellt für:", user)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": user,
		"exp":  time.Now().Add(time.Duration(days) * 24 * time.Hour).Unix(),
	})

	signed, err := token.SignedString(secret)
	if err != nil {
		fmt.Println("Token Fehler:", err)
		os.Exit(1)
	}

	fmt.Printf("User:     %s\n", user)
	fmt.Printf("Läuft ab: %s\n", time.Now().Add(time.Duration(days)*24*time.Hour).Format("2006-01-02"))
	fmt.Printf("Token:    %s\n", signed)
}
