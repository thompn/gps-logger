package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
  "os"
  "context"
  "time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
  "github.com/rs/cors"
  "github.com/google/uuid"
  "github.com/shaj13/go-guardian/auth"
	"github.com/shaj13/go-guardian/auth/strategies/basic"
	"github.com/shaj13/go-guardian/auth/strategies/bearer"
	"github.com/shaj13/go-guardian/store"
)

type Track struct {
	Timestamp string `json:"timestamp"`
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
	Surface   string `json:"surface"`
	Tracking  string `json:"tracking"`
}

type gpsTracks []Track

var db *gorm.DB

// variables for auth
var authenticator auth.Authenticator
var cache store.Cache

// Function to initialise the database connection
func initDB() {
	creds := godotenv.Load(".env")
	if creds != nil {
		log.Fatalf("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	var err error
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUsername, dbPassword, dbHost, dbPort, dbName)
	db, err = gorm.Open("mysql", dataSourceName)

	if err != nil {
		fmt.Println(err)
		panic("failed to connect to database")
	}

}

// Function to create auth token
func createToken(w http.ResponseWriter, r *http.Request) {
  creds := godotenv.Load(".env")
	if creds != nil {
		log.Fatalf("Error loading .env file")
	}
	uname := os.Getenv("AUTH_USER")
	token := uuid.New().String()
	user := auth.NewDefaultUser(uname, "1", nil, nil)
	tokenStrategy := authenticator.Strategy(bearer.CachedStrategyKey)
	auth.Append(tokenStrategy, token, user, r)
	body := fmt.Sprintf("token: %s \n", token)
	w.Write([]byte(body))
}

// Auth validation function
func validateUser(ctx context.Context, r *http.Request, userName, password string) (auth.Info, error) {
  creds := godotenv.Load(".env")
	if creds != nil {
		log.Fatalf("Error loading .env file")
	}
  uname := os.Getenv("AUTH_USER")
  pass := os.Getenv("AUTH_PASS")
  if userName == uname && password == pass {
      return auth.NewDefaultUser(uname, "1", nil, nil), nil
  }
  return nil, fmt.Errorf("Invalid credentials")
}

// Function for go-guardian
func setupGoGuardian() {
  authenticator = auth.New()
  cache = store.NewFIFO(context.Background(), time.Minute*10)
  basicStrategy := basic.New(validateUser, cache) 
  tokenStrategy := bearer.New(bearer.NoOpAuthenticate, cache)
  authenticator.EnableStrategy(basic.StrategyKey, basicStrategy)
  authenticator.EnableStrategy(bearer.CachedStrategyKey,    tokenStrategy)
}

// Middleware function to intercept requests to auth
func middleware(next http.Handler) http.HandlerFunc {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      log.Println("Executing Auth Middleware")
      user, err := authenticator.Authenticate(r)
      if err != nil {
          code := http.StatusUnauthorized
          http.Error(w, http.StatusText(code), code)
          log.Println(user)
          return
      }
      next.ServeHTTP(w, r)
  })
}

// Function to handle the page requests and route them
func handleRequests() {
  router := mux.NewRouter().StrictSlash(true)
  router.HandleFunc("/api/v1/auth/token", middleware(http.HandlerFunc(createToken))).Methods("GET")
	router.HandleFunc("/api/v1/", middleware(http.HandlerFunc(homePage))).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/tracks", middleware(http.HandlerFunc(getTracks))).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/tracks", middleware(http.HandlerFunc(newTrack))).Methods("POST", "OPTIONS")
	c := cors.New(cors.Options{
		AllowedMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowedOrigins:     []string{"*"},
		AllowCredentials:   true,
		AllowedHeaders:     []string{"Content-Type", "Bearer", "Bearer ", "content-type", "Origin", "Accept"},
		OptionsPassthrough: true,
	})
	handler := c.Handler(router)
	log.Fatal(http.ListenAndServe(":10000", handler))
}

// Function for main home page
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Not Found!")
}

// Function to get all tracks from database
func getTracks(w http.ResponseWriter, r *http.Request) {
	var tracks []Track
	result := db.Find(&tracks)
	json.NewEncoder(w).Encode(result)
}

// Function to post a track to the database
func newTrack(w http.ResponseWriter, r *http.Request) {
	var track Track
	json.NewDecoder(r.Body).Decode(&track)
	fmt.Println(track)
	db.Create(&track)
	json.NewEncoder(w).Encode(track)
	fmt.Fprintf(w, "Added Track")
}

// Main Function
func main() {
  initDB()
  setupGoGuardian()
  handleRequests()
}
