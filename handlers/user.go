package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"mind_tips_backend/config"
	"mind_tips_backend/database"
	"mind_tips_backend/middleware"
	"mind_tips_backend/models"

	"github.com/gorilla/mux"
)

// Google OAuth ログインページへのリダイレクト
func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	// CSRF対策のためのstate parameter
	state := generateState()

	// セッションまたはクッキーにstateを保存（本番環境では要実装）
	// TODO: stateの検証機能を実装

	url := config.GoogleOauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// 自分の情報を取得
func GetMyProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	user, err := database.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// 自分の情報を更新
func UpdateMyProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	var updateReq models.UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// 簡単なバリデーション
	if updateReq.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	user, err := database.UpdateUserProfile(userID, updateReq.Name)
	if err != nil {
		log.Printf("Failed to update user: %v", err)
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// 特定ユーザーの公開情報取得
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["id"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := database.GetPublicUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Googleからのコールバック処理
func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	// 認可コードの取得
	code := r.FormValue("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	// アクセストークンの取得
	token, err := config.GoogleOauthConfig.Exchange(r.Context(), code)
	if err != nil {
		log.Printf("Failed to exchange token: %v", err)
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	// ユーザー情報の取得
	client := config.GoogleOauthConfig.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Printf("Failed to get user info: %v", err)
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var googleUser models.GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		log.Printf("Failed to decode user info: %v", err)
		http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
		return
	}

	// データベースにユーザーを作成/更新
	user, err := database.CreateOrUpdateUserByGoogle(&googleUser)
	if err != nil {
		log.Printf("Failed to create/update user: %v", err)
		http.Error(w, "Failed to process user", http.StatusInternalServerError)
		return
	}

	// JWTトークンを生成
	jwtToken, err := config.GenerateToken(user.ID)
	if err != nil {
		log.Printf("Failed to generate JWT token: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// レスポンス
	response := models.LoginResponse{
		User:  *user,
		Token: jwtToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func generateState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
