package controllers

import (
	"cocopen-backend/dto"
	"cocopen-backend/models"
	"cocopen-backend/services"
	"cocopen-backend/utils"
	"cocopen-backend/middleware"
	"database/sql"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

func CreatePendaftar(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        // Hanya POST
        if r.Method != http.MethodPost {
            utils.Error(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
            return
        }

        // Tidak memakai JWT
        var userID *int = nil

        // Parse multipart (Wajib untuk file upload)
        if err := r.ParseMultipartForm(10 << 20); err != nil {
            utils.Error(w, http.StatusBadRequest, "Gagal memproses form-data: "+err.Error())
            return
        }

        // Ambil field form sesuai name di frontend
        req := dto.CreatePendaftarRequest{
            NamaLengkap:        strings.TrimSpace(r.FormValue("nama_lengkap")),
            AsalKampus:         strings.TrimSpace(r.FormValue("asal_kampus")),
            Prodi:              strings.TrimSpace(r.FormValue("prodi")),
            Semester:           strings.TrimSpace(r.FormValue("semester")),
            NoWA:               strings.TrimSpace(r.FormValue("no_wa")),
            Domisili:           r.FormValue("domisili"),
            AlamatSekarang:     r.FormValue("alamat_sekarang"),
            TinggalDengan:      r.FormValue("tinggal_dengan"),
            AlasanMasuk:        r.FormValue("alasan_masuk"),
            PengetahuanCoconut: r.FormValue("pengetahuan_coconut"),
        }

        // Validasi field wajib
        if req.NamaLengkap == "" || req.AsalKampus == "" || req.Prodi == "" || req.NoWA == "" {
            utils.Error(w, http.StatusBadRequest,
                "Field wajib belum lengkap: nama_lengkap, asal_kampus, prodi, no_wa")
            return
        }

        // --- Ambil file foto ---
        file, header, err := r.FormFile("foto")
        if err != nil {
            if err == http.ErrMissingFile {
                utils.Error(w, http.StatusBadRequest, "Foto wajib diunggah")
                return
            }
            utils.Error(w, http.StatusBadRequest, "Gagal membaca file foto: "+err.Error())
            return
        }
        defer file.Close()

        // Validasi ekstensi
        ext := strings.ToLower(filepath.Ext(header.Filename))
        allowed := map[string]bool{
            ".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
        }
        if !allowed[ext] {
            utils.Error(w, http.StatusBadRequest,
                "Format foto tidak didukung (hanya .jpg, .jpeg, .png, .gif)")
            return
        }

        // Validasi ukuran (maks 2MB)
        if header.Size > 2<<20 {
            utils.Error(w, http.StatusBadRequest, "Ukuran file maksimal 2 MB")
            return
        }

        // Upload foto memakai util
        fotoName, err := utils.UploadFoto(file, header, utils.FotoPendaftarPath)
        if err != nil {
            utils.Error(w, http.StatusInternalServerError,
                "Gagal mengunggah foto: "+err.Error())
            return
        }

        // Buat model untuk database
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
            UserID:             userID, // tidak perlu login
        }

        // Simpan ke database
        if err := services.CreatePendaftar(db, pendaftar); err != nil {
            utils.Error(w, http.StatusInternalServerError,
                "Gagal menambahkan pendaftar: "+err.Error())
            return
        }

        // Respon sukses
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
