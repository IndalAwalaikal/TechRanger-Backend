package controllers

import (
	"cocopen-backend/dto"
	"cocopen-backend/middleware"
	"cocopen-backend/models"
	"cocopen-backend/services"
	"cocopen-backend/utils"
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetUserSoalHandler(db *sql.DB) http.HandlerFunc {
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

        pendaftar, err := services.GetLatestPendaftarByUserID(db, claims.IDUser)
        if err != nil {
            utils.Error(w, http.StatusInternalServerError, "Gagal memeriksa pendaftaran: "+err.Error())
            return
        }
        if pendaftar == nil {
            utils.Error(w, http.StatusForbidden, "Anda belum melakukan pendaftaran")
            return
        }

        testConfig, err := services.GetTestConfig(db)
        if err != nil {
            if err == sql.ErrNoRows {
                utils.Error(w, http.StatusNotFound, "Tes belum tersedia")
                return
            }
            utils.Error(w, http.StatusInternalServerError, "Gagal memuat konfigurasi tes")
            return
        }

        now := time.Now()
        if testConfig.WaktuMulai != nil && now.Before(*testConfig.WaktuMulai) {
            utils.Error(w, http.StatusForbidden, "Tes belum dimulai")
            return
        }
        if testConfig.WaktuSelesai != nil && now.After(*testConfig.WaktuSelesai) {
            utils.Error(w, http.StatusForbidden, "Tes telah ditutup")
            return
        }

        alreadyTaken, err := services.HasUserTakenTest(db, claims.IDUser, testConfig.IDTest)
        if err != nil {
            utils.Error(w, http.StatusInternalServerError, "Gagal memeriksa riwayat tes")
            return
        }
        if alreadyTaken {
            utils.Error(w, http.StatusForbidden, "Anda sudah pernah mengikuti tes")
            return
        }

        soals, err := services.GetAllSoal(db)
        if err != nil {
            utils.Error(w, http.StatusInternalServerError, "Gagal mengambil soal: "+err.Error())
            return
        }

        var response []dto.SoalResponse
        for _, s := range soals {
            response = append(response, dto.SoalResponse{
                IDSoal:     s.IDSoal,
                Nomor:      s.Nomor,
                Pertanyaan: s.Pertanyaan,
                PilihanA:   s.PilihanA,
                PilihanB:   s.PilihanB,
                PilihanC:   s.PilihanC,
                PilihanD:   s.PilihanD,
            })
        }

        utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
            "soal":           response,
            "durasi_menit":   testConfig.DurasiMenit,
            "judul":          testConfig.Judul,
            "deskripsi":      testConfig.Deskripsi,
            "waktu_mulai":    testConfig.WaktuMulai,
            "waktu_selesai":  testConfig.WaktuSelesai,
        })
    }
}

func SubmitJawabanHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            utils.Error(w, http.StatusMethodNotAllowed, "Hanya metode POST yang diizinkan")
            return
        }

        claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
        if !ok {
            utils.Error(w, http.StatusUnauthorized, "Akses ditolak")
            return
        }

        var req dto.SubmitJawabanRequest
        if err := utils.ParseAndValidate(r, &req); err != nil {
            utils.Error(w, http.StatusBadRequest, err.Error())
            return
        }

        pendaftar, err := services.GetLatestPendaftarByUserID(db, claims.IDUser)
        if err != nil {
            utils.Error(w, http.StatusInternalServerError, "Gagal memeriksa pendaftaran")
            return
        }
        if pendaftar == nil {
            utils.Error(w, http.StatusForbidden, "Anda belum mendaftar")
            return
        }
        pendaftarID := pendaftar.IDPendaftar

        testConfig, err := services.GetTestConfig(db)
        if err != nil {
            utils.Error(w, http.StatusNotFound, "Tes tidak tersedia")
            return
        }

        alreadyTaken, err := services.HasUserTakenTest(db, claims.IDUser, testConfig.IDTest)
        if err != nil {
            utils.Error(w, http.StatusInternalServerError, "Gagal memeriksa riwayat tes")
            return
        }
        if alreadyTaken {
            utils.Error(w, http.StatusForbidden, "Anda sudah pernah mengikuti tes")
            return
        }

        soals, err := services.GetAllSoal(db)
        if err != nil {
            utils.Error(w, http.StatusInternalServerError, "Gagal memuat soal untuk penilaian")
            return
        }

        jawabanBenarMap := make(map[int]string)
        for _, s := range soals {
            jawabanBenarMap[s.IDSoal] = s.JawabanBenar
        }

        for idSoal := range req.Jawaban {
            if _, exists := jawabanBenarMap[idSoal]; !exists {
                utils.Error(w, http.StatusBadRequest, "Soal dengan ID "+strconv.Itoa(idSoal)+" tidak ditemukan")
                return
            }
        }

        skorBenar := 0
        var jawabans []models.JawabanUser

        for idSoal, jawabUser := range req.Jawaban {
            isBenar := strings.ToUpper(jawabUser) == jawabanBenarMap[idSoal]
            if isBenar {
                skorBenar++
            }
            jawabans = append(jawabans, models.JawabanUser{
                IDSoal:      idSoal,
                JawabanUser: strings.ToUpper(jawabUser),
                IsBenar:     isBenar,
            })
        }

        totalSoal := len(soals)
        skorSalah := totalSoal - skorBenar
        nilai := float64(skorBenar) / float64(totalSoal) * 100

        waktuSekarang := time.Now()

        hasilBaru := models.HasilTest{
            UserID:         claims.IDUser,
            PendaftarID:    pendaftarID,
            IDTest:         testConfig.IDTest,
            WaktuMulai:     waktuSekarang,
            WaktuSelesai:   &waktuSekarang,
            SkorBenar:      skorBenar,
            SkorSalah:      skorSalah,
            Nilai:          nilai,
        }

        err = services.CreateHasilTest(db, hasilBaru)
        if err != nil {
            utils.Error(w, http.StatusInternalServerError, "Gagal membuat hasil tes: "+err.Error())
            return
        }

        idHasil, err := services.GetHasilIDByUserAndTest(db, claims.IDUser, testConfig.IDTest)
        if err != nil {
            utils.Error(w, http.StatusInternalServerError, "Gagal ambil ID hasil tes")
            return
        }

        for _, j := range jawabans {
            j.IDHasil = idHasil
            err = services.CreateJawabanUser(db, j)
            if err != nil {
                utils.Error(w, http.StatusInternalServerError, "Gagal simpan jawaban: "+err.Error())
                return
            }
        }

        utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
            "success":    true,
            "message":    "Jawaban berhasil dikirim dan dinilai",
            "skor_benar": skorBenar,
            "skor_salah": skorSalah,
            "nilai":      nilai,
        })
    }
}

func CreateSoalHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.Error(w, http.StatusMethodNotAllowed, "Hanya POST yang diizinkan")
			return
		}

		claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
		if !ok || claims.Role != "admin" {
			utils.Error(w, http.StatusForbidden, "Akses ditolak")
			return
		}

		var req dto.SoalCreateRequest
		if err := utils.ParseAndValidate(r, &req); err != nil {
			utils.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		soal := models.SoalTest{
			Nomor:        req.Nomor,
			Pertanyaan:   req.Pertanyaan,
			PilihanA:     req.PilihanA,
			PilihanB:     req.PilihanB,
			PilihanC:     req.PilihanC,
			PilihanD:     req.PilihanD,
			JawabanBenar: req.JawabanBenar,
		}

		err := services.CreateSoal(db, soal)
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, "Gagal membuat soal: "+err.Error())
			return
		}

		utils.JSONResponse(w, http.StatusCreated, map[string]interface{}{
			"success": true,
			"message": "Soal berhasil dibuat",
		})
	}
}

func UpdateSoalHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			utils.Error(w, http.StatusMethodNotAllowed, "Hanya PUT yang diizinkan")
			return
		}

		claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
		if !ok || claims.Role != "admin" {
			utils.Error(w, http.StatusForbidden, "Akses ditolak")
			return
		}

		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			utils.Error(w, http.StatusBadRequest, "Parameter id wajib")
			return
		}

		idSoal, err := strconv.Atoi(idStr)
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "ID tidak valid")
			return
		}

		var req dto.SoalUpdateRequest
		if err := utils.ParseAndValidate(r, &req); err != nil {
			utils.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		soal, err := services.GetSoalByID(db, idSoal)
		if err != nil {
			if err == sql.ErrNoRows {
				utils.Error(w, http.StatusNotFound, "Soal tidak ditemukan")
				return
			}
			utils.Error(w, http.StatusInternalServerError, "Gagal ambil soal")
			return
		}

		if req.Nomor != nil {
			soal.Nomor = *req.Nomor
		}
		if req.Pertanyaan != nil {
			soal.Pertanyaan = *req.Pertanyaan
		}
		if req.PilihanA != nil {
			soal.PilihanA = *req.PilihanA
		}
		if req.PilihanB != nil {
			soal.PilihanB = *req.PilihanB
		}
		if req.PilihanC != nil {
			soal.PilihanC = *req.PilihanC
		}
		if req.PilihanD != nil {
			soal.PilihanD = *req.PilihanD
		}
		if req.JawabanBenar != nil {
			soal.JawabanBenar = *req.JawabanBenar
		}

		err = services.UpdateSoal(db, soal)
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, "Gagal update soal: "+err.Error())
			return
		}

		utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Soal berhasil diperbarui",
		})
	}
}

func DeleteSoalHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			utils.Error(w, http.StatusMethodNotAllowed, "Hanya DELETE yang diizinkan")
			return
		}

		claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
		if !ok || claims.Role != "admin" {
			utils.Error(w, http.StatusForbidden, "Akses ditolak")
			return
		}

		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			utils.Error(w, http.StatusBadRequest, "Parameter id wajib")
			return
		}

		idSoal, err := strconv.Atoi(idStr)
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "ID tidak valid")
			return
		}

		_, err = services.GetSoalByID(db, idSoal)
		if err != nil {
			if err == sql.ErrNoRows {
				utils.Error(w, http.StatusNotFound, "Soal tidak ditemukan")
				return
			}
			utils.Error(w, http.StatusInternalServerError, "Gagal cek soal")
			return
		}

		err = services.DeleteSoal(db, idSoal)
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, "Gagal hapus soal: "+err.Error())
			return
		}

		utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Soal berhasil dihapus",
		})
	}
}

func GetHasilTesUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.Error(w, http.StatusMethodNotAllowed, "Hanya GET yang diizinkan")
			return
		}

		claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
		if !ok {
			utils.Error(w, http.StatusUnauthorized, "Akses ditolak")
			return
		}

		testConfig, err := services.GetTestConfig(db)
		if err != nil {
			utils.Error(w, http.StatusNotFound, "Konfigurasi tes tidak ditemukan")
			return
		}

		hasil, err := services.GetHasilByUserID(db, claims.IDUser, testConfig.IDTest)
		if err != nil {
			if err == sql.ErrNoRows {
				utils.Error(w, http.StatusNotFound, "Anda belum mengikuti tes")
				return
			}
			utils.Error(w, http.StatusInternalServerError, "Gagal ambil hasil")
			return
		}

		var waktuSelesai *string
		if hasil.WaktuSelesai != nil {
			t := hasil.WaktuSelesai.Format("2006-01-02 15:04:05")
			waktuSelesai = &t
		}

		response := dto.HasilResponse{
			IDHasil:      hasil.IDHasil,
			UserID:       hasil.UserID,
			PendaftarID:  hasil.PendaftarID,
			SkorBenar:    hasil.SkorBenar,
			SkorSalah:    hasil.SkorSalah,
			Nilai:        hasil.Nilai,
			WaktuMulai:   hasil.WaktuMulai.Format("2006-01-02 15:04:05"),
			WaktuSelesai: waktuSelesai,
			DurasiMenit:  hasil.DurasiMenit,
		}

		utils.JSONResponse(w, http.StatusOK, response)
	}
}

func GetAllHasilTesHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            utils.Error(w, http.StatusMethodNotAllowed, "Hanya GET yang diizinkan")
            return
        }

        claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
        if !ok || claims.Role != "admin" {
            utils.Error(w, http.StatusForbidden, "Akses ditolak")
            return
        }

        testConfig, err := services.GetTestConfig(db)
        if err != nil {
            utils.Error(w, http.StatusNotFound, "Konfigurasi tes tidak ditemukan")
            return
        }

        // 🔁 Query: JOIN dengan pendaftar saja
        query := `
            SELECT 
                ht.id_hasil,
                ht.user_id,
                ht.pendaftar_id,
                p.nama_lengkap AS pendaftar_name,
                ht.skor_benar,
                ht.skor_salah,
                ht.nilai,
                ht.waktu_mulai,
                ht.waktu_selesai,
                ht.durasi_menit
            FROM hasil_test ht
            LEFT JOIN pendaftar p ON ht.pendaftar_id = p.id_pendaftar
            WHERE ht.id_test = ?
            ORDER BY ht.nilai DESC, ht.waktu_mulai ASC
        `

        rows, err := db.Query(query, testConfig.IDTest)
        if err != nil {
            utils.Error(w, http.StatusInternalServerError, "Gagal ambil hasil tes: "+err.Error())
            return
        }
        defer rows.Close()

        var hasilList []map[string]interface{}
        for rows.Next() {
            var (
                idHasil       int
                userID        int
                pendaftarID   sql.NullInt64
                pendaftarName sql.NullString
                skorBenar     int
                skorSalah     int
                nilai         float64
                waktuMulai    time.Time
                waktuSelesai  sql.NullTime
                durasiMenit   sql.NullInt64
            )

            err := rows.Scan(
                &idHasil,
                &userID,
                &pendaftarID,
                &pendaftarName,
                &skorBenar,
                &skorSalah,
                &nilai,
                &waktuMulai,
                &waktuSelesai,
                &durasiMenit,
            )
            if err != nil {
                utils.Error(w, http.StatusInternalServerError, "Gagal scan hasil: "+err.Error())
                return
            }

            item := map[string]interface{}{
                "id_hasil":       idHasil,
                "user_id":        userID,
                "pendaftar_id":   nil,
                "pendaftar_name": "Belum daftar",
                "skor_benar":     skorBenar,
                "skor_salah":     skorSalah,
                "nilai":          nilai,
                "waktu_mulai":    waktuMulai.Format("2006-01-02 15:04:05"),
                "waktu_selesai":  nil,
                "durasi_menit":   nil,
            }

            if pendaftarID.Valid {
                item["pendaftar_id"] = int(pendaftarID.Int64)
            }
            if pendaftarName.Valid && pendaftarName.String != "" {
                item["pendaftar_name"] = pendaftarName.String
            }
            if waktuSelesai.Valid {
                item["waktu_selesai"] = waktuSelesai.Time.Format("2006-01-02 15:04:05")
            }
            if durasiMenit.Valid {
                item["durasi_menit"] = int(durasiMenit.Int64)
            }

            hasilList = append(hasilList, item)
        }

        utils.JSONResponse(w, http.StatusOK, hasilList)
    }
}

func GetAllSoalAdminHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.Error(w, http.StatusMethodNotAllowed, "Hanya GET yang diizinkan")
			return
		}

		claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
		if !ok || claims.Role != "admin" {
			utils.Error(w, http.StatusForbidden, "Akses ditolak")
			return
		}

		soals, err := services.GetAllSoal(db)
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, "Gagal ambil soal: "+err.Error())
			return
		}

		var response []dto.SoalResponse
		for _, s := range soals {
			response = append(response, dto.SoalResponse{
				IDSoal:     s.IDSoal,
				Nomor:      s.Nomor,
				Pertanyaan: s.Pertanyaan,
				PilihanA:   s.PilihanA,
				PilihanB:   s.PilihanB,
				PilihanC:   s.PilihanC,
				PilihanD:   s.PilihanD,
			})
		}

		utils.JSONResponse(w, http.StatusOK, response)
	}
}

func UpdateTestConfigHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
        if !ok || claims.Role != "admin" {
            utils.Error(w, http.StatusForbidden, "Akses ditolak")
            return
        }

        switch r.Method {
        case http.MethodGet:
            // 🔹 Ambil konfigurasi saat ini
            testConfig, err := services.GetTestConfig(db)
            if err != nil {
                if err == sql.ErrNoRows {
                    utils.Error(w, http.StatusNotFound, "Konfigurasi tes belum dibuat")
                    return
                }
                utils.Error(w, http.StatusInternalServerError, "Gagal ambil konfigurasi: "+err.Error())
                return
            }

            utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
                "id_test":       testConfig.IDTest,
                "judul":         testConfig.Judul,
                "deskripsi":     testConfig.Deskripsi,
                "durasi_menit":  testConfig.DurasiMenit,
                "waktu_mulai":   testConfig.WaktuMulai,
                "waktu_selesai": testConfig.WaktuSelesai,
                "aktif":         testConfig.Aktif,
            })

        case http.MethodPut:
            // 🔹 Update konfigurasi
            var req dto.TestConfigUpdateRequest
            if err := utils.ParseAndValidate(r, &req); err != nil {
                utils.Error(w, http.StatusBadRequest, err.Error())
                return
            }

            existing, _ := services.GetTestConfig(db)

            config := models.Test{
                Judul:       "Tes Seleksi",
                DurasiMenit: 30,
                Aktif:       true,
            }

            if existing != nil {
                config.IDTest = existing.IDTest
            }

            if req.Judul != nil {
                config.Judul = *req.Judul
            }
            if req.Deskripsi != nil {
                config.Deskripsi = req.Deskripsi
            }
            if req.DurasiMenit != nil {
                config.DurasiMenit = *req.DurasiMenit
            }
            if req.WaktuMulai != nil {
                config.WaktuMulai = req.WaktuMulai
            }
            if req.WaktuSelesai != nil {
                config.WaktuSelesai = req.WaktuSelesai
            }
            if req.Aktif != nil {
                config.Aktif = *req.Aktif
            }

            err := services.UpsertTestConfig(db, config)
            if err != nil {
                utils.Error(w, http.StatusInternalServerError, "Gagal simpan konfigurasi: "+err.Error())
                return
            }

            utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
                "success": true,
                "message": "Konfigurasi tes berhasil disimpan",
            })

        default:
            utils.Error(w, http.StatusMethodNotAllowed, "Hanya GET dan PUT yang diizinkan")
        }
    }
}

func GetTestConfigHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            utils.Error(w, http.StatusMethodNotAllowed, "Hanya GET yang diizinkan")
            return
        }

        claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
        if !ok || claims.Role != "admin" {
            utils.Error(w, http.StatusForbidden, "Akses ditolak")
            return
        }

        testConfig, err := services.GetTestConfig(db)
        if err != nil {
            if err == sql.ErrNoRows {
                utils.Error(w, http.StatusNotFound, "Konfigurasi tes belum dibuat")
                return
            }
            utils.Error(w, http.StatusInternalServerError, "Gagal ambil konfigurasi: "+err.Error())
            return
        }

        utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
            "id_test":       testConfig.IDTest,
            "judul":         testConfig.Judul,
            "deskripsi":     testConfig.Deskripsi,
            "durasi_menit":  testConfig.DurasiMenit,
            "waktu_mulai":   testConfig.WaktuMulai,
            "waktu_selesai": testConfig.WaktuSelesai,
            "aktif":         testConfig.Aktif,
        })
    }
}