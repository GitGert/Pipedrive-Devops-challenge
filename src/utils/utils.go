package utils

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
)

var RESET = "\033[0m"
var RED = "\033[31m"
var GREEN = "\033[32m"

func Log_event(event string, message string) {
	time := time.Now()
	fmt.Println(time.Format("2006-01-02 15:04:05") + "\t" + event + "\t" + message)
}

func Log_request(r *http.Request, message string) {
	fullURL := GetUrl(r)
	Log_event(MakeGreen(fullURL), MakeGreen(message))
}

func LoadEnvFile() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(MakeRed("Please make sure your .env is in the project root"))
		fmt.Println(MakeRed("and that API_TOKEN and COMPANY_DOMAIN are set"))
		log.Fatal("Error loading .env file")
	}
}

func GetUrl(r *http.Request) string {
	scheme := "http" // Default scheme, change based on your application logic
	if r.TLS != nil {
		scheme = "https"
	}
	host := r.Host
	path := r.URL.Path
	fullURL := scheme + "://" + host + path
	return fullURL
}

func MakeRed(text string) string {
	return RED + text + RESET
}

func MakeGreen(text string) string {
	return GREEN + text + RESET
}
