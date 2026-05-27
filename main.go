package main

import (
  "encoding/json"
  "fmt"
  "log"
  "math/rand"
  "net/http"
  "shortener/database"
  "shortener/users"
  "strconv"
  "strings"

  "github.com/gorilla/mux"
)

func main() {
  database.Connect()

  r := mux.NewRouter()

  r.HandleFunc("/shorturl", shorturlHandler).Methods("POST")
  r.HandleFunc("/{code}", redrectHandler).Methods("GET")
  r.HandleFunc("/api/stats/total", totalStatsHandler).Methods("GET")
  r.HandleFunc("/api/stats/user/{username}", userStatsHandler).Methods("GET")

  r.PathPrefix("/").Handler(http.FileServer(http.Dir(".")))

  fmt.Println("Server started on http://localhost:8080")
  http.ListenAndServe(":8080", r)
}

func totalStatsHandler(w http.ResponseWriter, r *http.Request) {
  var count int64
  if err := database.DB.Model(&users.Links{}).Count(&count).Error; err != nil {
    http.Error(w, "Database error", http.StatusInternalServerError)
    return
  }
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(map[string]int64{"total": count})
}

func userStatsHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  username := vars["username"]

  var count int64
  if err := database.DB.Model(&users.Links{}).Where("created_by = ?", username).Count(&count).Error; err != nil {
    http.Error(w, "Database error", http.StatusInternalServerError)
    return
  }
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(map[string]int64{"user_total": count})
}

func shorturlHandler(w http.ResponseWriter, r *http.Request) {
  var req users.Mix

  if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    http.Error(w, "Invalid JSON", http.StatusBadRequest)
    return
  }

  req.Old = strings.TrimSpace(req.Old)

  if !strings.HasPrefix(req.Old, "http://") && !strings.HasPrefix(req.Old, "https://") {
    req.Old = "http://" + req.Old
  }

  fmt.Println("Received URL:", req.Old)

  var a int = rand.Intn(900000) + 100000
  var b string = "http://localhost:8080/"

  str := strconv.Itoa(a)
  b = b + str

  l := users.Links{
    Old:       req.Old,
    Short:     b,
    CreatedBy: req.Username,
  }

  if err := database.DB.Create(&l).Error; err != nil {
    http.Error(w, "Database error", http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusCreated)

  response := map[string]string{
    "short_url": b,
  }
  json.NewEncoder(w).Encode(response)
}

func redrectHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  code := vars["code"]
  shortURL := "http://localhost:8080/" + code

  var link users.Links
  result := database.DB.Where("short = ?", shortURL).First(&link)

  if result.Error != nil {
    http.Error(w, "URL not found", http.StatusNotFound)
    return
  }

  log.Println(link.Old)

  http.Redirect(w, r, link.Old, http.StatusFound)
}
