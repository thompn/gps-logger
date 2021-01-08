package main

import (
    "fmt"
    "log"
    "net/http"
    "encoding/json"
    "os"

    "github.com/joho/godotenv"
    "github.com/gorilla/mux"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
)

type Track struct {
  Timestamp string `json:"timestamp"`
  Longitude string `json:"longitude"`
  Latitude string `json:"latitude"`
  Surface string `json:surface`
}

type gpsTracks []Track

var db *gorm.DB

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

// Function to handle the page requests and route them
func handleRequests() {
  myRouter := mux.NewRouter().StrictSlash(true)
  myRouter.HandleFunc("/api/v1/", homePage).Methods("GET")
  myRouter.HandleFunc("/api/v1/tracks", getTracks).Methods("GET")
  myRouter.HandleFunc("/api/v1/tracks", newTrack).Methods("POST")
  log.Fatal(http.ListenAndServe(":10000", myRouter))
}

// Function for main home page
func homePage(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Not Found!")
  fmt.Println("Endpoint Hit: homepage")
}

// Function to get all tracks from database
func getTracks(w http.ResponseWriter, r *http.Request) {
  var tracks []Track
  db.Preload("gpsTracks").Find(&tracks)
  json.NewEncoder(w).Encode(tracks)
}

// Function to post a track to the database
func newTrack(w http.ResponseWriter, r*http.Request) {
  var track Track
  json.NewDecoder(r.Body).Decode(&track)
  db.Create(&track)
  //w.Header().set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(track)
  fmt.Fprintf(w, "Added Track")
}

// Main Function
func main() {
  initDB()
  handleRequests()
}
