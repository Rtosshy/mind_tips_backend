package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type HealthResponse struct {
	Status string `json:"status"`
}

type TipsResponse struct {
	Tips []string `json:"tips"`
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Mind Tips Backend!")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := HealthResponse{Status: "ok"}
	json.NewEncoder(w).Encode(response)
}

func tipsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tips := TipsResponse{
		Tips: []string{
			"Stay hydrated",
			"Take breaks",
			"Practice mindfulness",
		},
	}
	json.NewEncoder(w).Encode(tips)
}

func main() {
	// Gorilla Muxルーターを初期化
	r := mux.NewRouter()

	// ルートの定義
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/health", healthHandler).Methods("GET")
	r.HandleFunc("/api/tips", tipsHandler).Methods("GET")

	// ミドルウェア（ログ出力）
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	})

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
