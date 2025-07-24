package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mind_tips_backend/config"
	"mind_tips_backend/database"
	"mind_tips_backend/handlers"
	"mind_tips_backend/middleware"
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
	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDB()

	// OAuth設定初期化
	config.InitOAuth()

	config.InitJWT()

	// Gorilla Muxルーターを初期化
	router := mux.NewRouter()

	// パブリックルート（認証不要）
	router.HandleFunc("/", homeHandler).Methods("GET")
	router.HandleFunc("/health", healthHandler).Methods("GET")
	router.HandleFunc("/api/tips", tipsHandler).Methods("GET")
	router.HandleFunc("/auth/google/login", handlers.GoogleLogin).Methods("GET")
	router.HandleFunc("/auth/google/callback", handlers.GoogleCallback).Methods("GET")

	// ユーザー公開情報（認証不要）
	router.HandleFunc("/api/users/{id}", handlers.GetUserByID).Methods("GET")

	// 認証が必要なルート
	protected := router.PathPrefix("/api/user").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/me", handlers.GetMyProfile).Methods("GET")
	protected.HandleFunc("/me", handlers.UpdateMyProfile).Methods("PUT")

	// ミドルウェア（ログ出力）
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	})

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
