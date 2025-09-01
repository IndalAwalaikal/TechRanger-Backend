package controllers

import (
	"cocopen-backend/dto"
	"cocopen-backend/models"
	"cocopen-backend/services"
	"cocopen-backend/utils"
	"cocopen-backend/middleware"
	"database/sql"
	"time"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

func CreatePendaftar(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.Error(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
			return
		}

		claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
		if !ok {
			utils.Error(w, http.StatusUnauthorized, "Akses ditolak")
			return
		}

		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "Gagal memproses form data")
			return
		}

		var req dto.CreatePendaftarRequest
		req.NamaLengkap = r.FormValue("nama_lengkap")
		req.AsalKampus = r.FormValue("asal_kampus")
		req.Prodi = r.FormValue("prodi")
		req.Semester = r.FormValue("semester")
		req.NoWA = r.FormValue("no_wa")
		req.Domisili = r.FormValue("domisili")
		req.AlamatSekarang = r.FormValue("alamat_sekarang")
		req.TinggalDengan = r.FormValue("tinggal_dengan")
		req.AlasanMasuk = r.FormValue("alasan_masuk")
		req.PengetahuanCoconut = r.FormValue("pengetahuan_coconut")

		if req.NamaLengkap == "" || req.AsalKampus == "" || req.Prodi == "" || req.NoWA == "" {
			utils.Error(w, http.StatusBadRequest, "Field wajib: nama_lengkap, asal_kampus, prodi, no_wa")
			return
		}

		file, header, err := r.FormFile("foto")
		if err != nil {
			if err == http.ErrMissingFile {
				utils.Error(w, http.StatusBadRequest, "Foto wajib diunggah")
				return
			}
			utils.Error(w, http.StatusBadRequest, "Gagal membaca file foto")
			return
		}
		defer file.Close()

		ext := strings.ToLower(filepath.Ext(header.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
			utils.Error(w, http.StatusBadRequest, "Format file tidak didukung (hanya .jpg, .jpeg, .png, .gif)")
			return
		}

		if header.Size > 2<<20 {
			utils.Error(w, http.StatusBadRequest, "Ukuran file maksimal 2 MB")
			return
		}

		latest, err := services.GetLatestPendaftarByUserID(db, claims.IDUser)
        if err != nil {
            utils.Error(w, http.StatusInternalServerError, "Gagal memeriksa riwayat pendaftaran")
            return
        }

        const duaBulan = 90 * 24 * time.Hour

        if latest != nil {
            now := time.Now()
            selisih := now.Sub(latest.CreatedAt)

            if selisih < duaBulan {
                hariSisa := int((duaBulan - selisih).Hours()/24) + 1
                if hariSisa < 1 {
                    hariSisa = 1
                }
                utils.Error(w, http.StatusTooEarly, "Anda baru bisa mendaftar kembali dalam " + strconv.Itoa(hariSisa) + " hari (batas 3 bulan antar pendaftaran)")
                return
            }
        }

		fotoName, err := utils.UploadFoto(file, header, utils.FotoPendaftarPath)
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, "Gagal mengunggah foto")
			return
		}

		pendaftar := models.Pendaftar{
			NamaLengkap:        req.NamaLengkap,
			AsalKampus:         req.AsalKampus,
			Prodi:              req.Prodi,
			Semester:           req.Semester,
			NoWA:               req.NoWA,
			Domisili:           req.Domisili,
			AlamatSekarang:     req.AlamatSekarang,
			TinggalDengan:      req.TinggalDengan,
			AlasanMasuk:        req.AlasanMasuk,
			PengetahuanCoconut: req.PengetahuanCoconut,
			FotoPath:           fotoName,
			Status:             "pending",
			UserID:             &claims.IDUser,
		}

		if err := services.CreatePendaftar(db, pendaftar); err != nil {
			utils.Error(w, http.StatusInternalServerError, "Gagal menambahkan pendaftar")
			return
		}

		utils.JSONResponse(w, http.StatusCreated, map[string]interface{}{
			"success": true,
			"message": "Pendaftar berhasil dibuat",
			"data": map[string]string{
				"foto_url": "/uploads/foto_pendaftar/" + fotoName,
			},
		})
	}
}

func GetAllPendaftar(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.Error(w, http.StatusMethodNotAllowed, "Hanya metode GET yang diizinkan")
			return
		}

		rows, err := services.GetAllPendaftar(db)
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, "Gagal mengambil data pendaftar")
			return
		}
		defer rows.Close()

		var result []map[string]interface{}
		for rows.Next() {
			var p models.Pendaftar
			err := rows.Scan(
				&p.IDPendaftar, &p.NamaLengkap, &p.AsalKampus, &p.Prodi, &p.Semester,
				&p.NoWA, &p.Domisili, &p.AlamatSekarang, &p.TinggalDengan,
				&p.AlasanMasuk, &p.PengetahuanCoconut, &p.FotoPath,
				&p.CreatedAt, &p.UpdatedAt, &p.Status, &p.UserID,
			)
			if err != nil {
				utils.Error(w, http.StatusInternalServerError, "Gagal membaca data pendaftar")
				return
			}

			fotoURL := ""
			if p.FotoPath != "" {
				fotoURL = "/uploads/foto_pendaftar/" + p.FotoPath
			}

			result = append(result, map[string]interface{}{
				"id_pendaftar":         p.IDPendaftar,
				"nama_lengkap":         p.NamaLengkap,
				"asal_kampus":          p.AsalKampus,
				"prodi":                p.Prodi,
				"semester":             p.Semester,
				"no_wa":                p.NoWA,
				"domisili":             p.Domisili,
				"alamat_sekarang":      p.AlamatSekarang,
				"tinggal_dengan":       p.TinggalDengan,
				"alasan_masuk":         p.AlasanMasuk,
				"pengetahuan_coconut":  p.PengetahuanCoconut,
				"foto_path":            p.FotoPath,
				"foto_url":             fotoURL, // 👈 ini yang frontend butuhkan
				"created_at":           p.CreatedAt,
				"updated_at":           p.UpdatedAt,
				"status":               p.Status,
				"user_id":              p.UserID,
			})
		}

		utils.JSONResponse(w, http.StatusOK, result)
	}
}

func GetPendaftarByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.Error(w, http.StatusMethodNotAllowed, "Hanya metode GET yang diizinkan")
			return
		}

		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			utils.Error(w, http.StatusBadRequest, "Parameter id wajib diisi")
			return
		}

		idPendaftar, err := strconv.Atoi(idStr)
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "ID tidak valid")
			return
		}

		p, err := services.GetPendaftarByID(db, idPendaftar)
		if err != nil {
			if err == sql.ErrNoRows {
				utils.Error(w, http.StatusNotFound, "Pendaftar tidak ditemukan")
				return
			}
			utils.Error(w, http.StatusInternalServerError, "Gagal mengambil data pendaftar")
			return
		}

		utils.JSONResponse(w, http.StatusOK, p)
	}
}

func UpdatePendaftar(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			utils.Error(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
			return
		}

		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "Gagal memproses form")
			return
		}

		idStr := r.FormValue("id_pendaftar")
		id, err := strconv.Atoi(idStr)
		if err != nil || id == 0 {
			utils.Error(w, http.StatusBadRequest, "ID pendaftar tidak valid")
			return
		}

		status := r.FormValue("status")
		if status == "" {
			utils.Error(w, http.StatusBadRequest, "Status wajib diisi")
			return
		}

		if status != "pending" && status != "diterima" && status != "ditolak" {
			utils.Error(w, http.StatusBadRequest, "Status tidak valid")
			return
		}

		if err := services.UpdatePendaftar(db, id, status); err != nil {
			utils.Error(w, http.StatusInternalServerError, "Gagal memperbarui pendaftar")
			return
		}

		utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Pendaftar berhasil diperbarui",
		})
	}
}

func DeletePendaftar(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			utils.Error(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
			return
		}

		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			utils.Error(w, http.StatusBadRequest, "Parameter id wajib diisi")
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "ID tidak valid")
			return
		}

		oldData, err := services.GetPendaftarByID(db, id)
		if err == nil && oldData.FotoPath != "" {
			utils.HapusFoto(utils.FotoPendaftarPath, oldData.FotoPath)
		}

		if err := services.DeletePendaftar(db, id); err != nil {
			if err == sql.ErrNoRows {
				utils.Error(w, http.StatusNotFound, "Pendaftar tidak ditemukan")
				return
			}
			utils.Error(w, http.StatusInternalServerError, "Gagal menghapus pendaftar")
			return
		}

		utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Pendaftar berhasil dihapus",
		})
	}
}

func GetMyPendaftar(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            utils.Error(w, http.StatusMethodNotAllowed, "Hanya metode GET yang diizinkan")
            return
        }

        claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
        if !ok {
            utils.Error(w, http.StatusUnauthorized, "Akses ditolak")
            return
        }

        pendaftarList, err := services.GetPendaftarByUserID(db, claims.IDUser)
        if err != nil {
            utils.Error(w, http.StatusInternalServerError, "Gagal mengambil data pendaftaran")
            return
        }

        // Kirim response
        utils.JSONResponse(w, http.StatusOK, pendaftarList)
    }
}
