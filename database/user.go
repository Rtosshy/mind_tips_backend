package database

import (
	"mind_tips_backend/models"
	"time"
)

func CreateOrUpdateUserByGoogle(googleUser *models.GoogleUserInfo) (*models.User, error) {
	user, err := GetUserByGoogleID(googleUser.ID)
	if err == nil {
		if needsUpdate(user, googleUser) {
			return updateUser(user, googleUser)
		}
		return user, nil
	}
	return createNewUser(googleUser)
}

// 更新が必要かどうかを判定する関数
func needsUpdate(user *models.User, googleUser *models.GoogleUserInfo) bool {
	return user.Email != googleUser.Email ||
		user.Name != googleUser.Name ||
		user.PictureURL != googleUser.Picture
}

func GetUserByGoogleID(googleID string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, name, google_id, picture_url, created_at, updated_at
		FROM users WHERE google_id = $1
	`

	err := DB.QueryRow(query, googleID).Scan(
		&user.ID, &user.Email, &user.Name, &user.GoogleID,
		&user.PictureURL, &user.CreatedAt, &user.UpdatedAt,
	)

	return user, err
}

// IDでユーザーを取得（自分の情報用）
func GetUserByID(userID int) (*models.User, error) {
	user := &models.User{}
	query := `
        SELECT id, email, name, google_id, picture_url, created_at, updated_at
        FROM users WHERE id = $1
    `

	err := DB.QueryRow(query, userID).Scan(
		&user.ID, &user.Email, &user.Name, &user.GoogleID,
		&user.PictureURL, &user.CreatedAt, &user.UpdatedAt,
	)

	return user, err
}

// 公開情報のみ取得（他のユーザー用）
func GetPublicUserByID(userID int) (*models.PublicUser, error) {
	user := &models.PublicUser{}
	query := `
        SELECT id, name, picture_url
        FROM users WHERE id = $1
    `

	err := DB.QueryRow(query, userID).Scan(
		&user.ID, &user.Name, &user.PictureURL,
	)

	return user, err
}

// ユーザー情報更新
func UpdateUserProfile(userID int, name string) (*models.User, error) {
	query := `
        UPDATE users
        SET name = $1, updated_at = $2
        WHERE id = $3
    `

	now := time.Now()
	_, err := DB.Exec(query, name, now, userID)
	if err != nil {
		return nil, err
	}

	return GetUserByID(userID)
}

func createNewUser(googleUser *models.GoogleUserInfo) (*models.User, error) {
	user := &models.User{
		Email:      googleUser.Email,
		Name:       googleUser.Name,
		GoogleID:   googleUser.ID,
		PictureURL: googleUser.Picture,
	}

	query := `
        INSERT INTO users (email, name, google_id, picture_url, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `

	now := time.Now()
	err := DB.QueryRow(query,
		user.Email, user.Name, user.GoogleID, user.PictureURL, now, now,
	).Scan(&user.ID)

	if err != nil {
		return nil, err
	}

	user.CreatedAt = now
	user.UpdatedAt = now
	return user, nil
}

func updateUser(user *models.User, googleUser *models.GoogleUserInfo) (*models.User, error) {
	query := `
        UPDATE users
        SET email = $1, name = $2, picture_url = $3, updated_at = $4
        WHERE google_id = $5
    `

	now := time.Now()
	_, err := DB.Exec(query,
		googleUser.Email, googleUser.Name, googleUser.Picture, now, user.GoogleID,
	)

	if err != nil {
		return nil, err
	}

	user.Email = googleUser.Email
	user.Name = googleUser.Name
	user.PictureURL = googleUser.Picture
	user.UpdatedAt = now

	return user, nil

}
