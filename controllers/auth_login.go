// controllers/login.go
package controllers

import (
	"cocopen-backend/dto"
	"cocopen-backend/utils"
	"encoding/json"
	"net/http"
	"os"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.Error(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
		return
	}

	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	staticUsername := os.Getenv("STATIC_USERNAME")
	staticPassword := os.Getenv("STATIC_PASSWORD")

	// Validasi env tersedia
	if staticUsername == "" || staticPassword == "" {
		utils.Error(w, http.StatusInternalServerError, "Kredensial admin belum dikonfigurasi")
		return
	}

	// Cek kredensial
	if req.Username != staticUsername || req.Password != staticPassword {
		utils.Error(w, http.StatusUnauthorized, "Username atau password salah")
		return
	}

	// 🔑 Eksplisit: hanya role "admin" yang diizinkan
	role := "admin"
	userID := 1
	fullName := "Administrator"
	username := staticUsername
	profilePicture := ""

	// (Opsional) Tambahkan validasi role jika kamu punya logika dinamis nanti
	if role != "admin" {
		utils.Error(w, http.StatusForbidden, "Akses ditolak: hanya admin yang diizinkan")
		return
	}

	token := utils.GenerateToken(userID, username, fullName, role, profilePicture)

	utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"status":  http.StatusOK,
		"message": "Login berhasil",
		"token":   token,
		"user": map[string]interface{}{
			"id":              userID,
			"username":        username,
			"full_name":       fullName,
			"role":            role,
			"profile_picture": profilePicture,
		},
	})
}