package services

import (
	"cocopen-backend/models"
	"database/sql"
)

func CreateSoal(db *sql.DB, soal models.SoalTest) error {
	query := `
		INSERT INTO soal_test (nomor, pertanyaan, pilihan_a, pilihan_b, pilihan_c, pilihan_d, jawaban_benar)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.Exec(
		query,
		soal.Nomor,
		soal.Pertanyaan,
		soal.PilihanA,
		soal.PilihanB,
		soal.PilihanC,
		soal.PilihanD,
		soal.JawabanBenar,
	)
	return err
}

func GetAllSoal(db *sql.DB) ([]models.SoalTest, error) {
	query := `
		SELECT 
			id_soal, nomor, pertanyaan, pilihan_a, pilihan_b, pilihan_c, pilihan_d, 
			jawaban_benar, created_at, updated_at
		FROM soal_test
		ORDER BY nomor ASC
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var soals []models.SoalTest
	for rows.Next() {
		var s models.SoalTest
		err := rows.Scan(
			&s.IDSoal,
			&s.Nomor,
			&s.Pertanyaan,
			&s.PilihanA,
			&s.PilihanB,
			&s.PilihanC,
			&s.PilihanD,
			&s.JawabanBenar,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		soals = append(soals, s)
	}
	return soals, rows.Err()
}

func GetSoalByID(db *sql.DB, idSoal int) (*models.SoalTest, error) {
	query := `
		SELECT 
			id_soal, nomor, pertanyaan, pilihan_a, pilihan_b, pilihan_c, pilihan_d, 
			jawaban_benar, created_at, updated_at
		FROM soal_test
		WHERE id_soal = ?
	`
	row := db.QueryRow(query, idSoal)

	var s models.SoalTest
	err := row.Scan(
		&s.IDSoal,
		&s.Nomor,
		&s.Pertanyaan,
		&s.PilihanA,
		&s.PilihanB,
		&s.PilihanC,
		&s.PilihanD,
		&s.JawabanBenar,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func UpdateSoal(db *sql.DB, soal *models.SoalTest) error {
	query := `
		UPDATE soal_test SET
			nomor = ?,
			pertanyaan = ?,
			pilihan_a = ?,
			pilihan_b = ?,
			pilihan_c = ?,
			pilihan_d = ?,
			jawaban_benar = ?,
			updated_at = NOW()
		WHERE id_soal = ?
	`
	_, err := db.Exec(
		query,
		soal.Nomor,
		soal.Pertanyaan,
		soal.PilihanA,
		soal.PilihanB,
		soal.PilihanC,
		soal.PilihanD,
		soal.JawabanBenar,
		soal.IDSoal,
	)
	return err
}

func DeleteSoal(db *sql.DB, idSoal int) error {
	query := `DELETE FROM soal_test WHERE id_soal = ?`
	_, err := db.Exec(query, idSoal)
	return err
}

func GetTestConfig(db *sql.DB) (*models.Test, error) {
	query := `
		SELECT 
			id_test, judul, deskripsi, durasi_menit, 
			waktu_mulai, waktu_selesai, aktif, created_at, updated_at
		FROM test
		WHERE aktif = TRUE
		LIMIT 1
	`
	row := db.QueryRow(query)

	var t models.Test
	var deskripsi sql.NullString
	var waktuMulai, waktuSelesai sql.NullTime

	err := row.Scan(
		&t.IDTest,
		&t.Judul,
		&deskripsi,
		&t.DurasiMenit,
		&waktuMulai,
		&waktuSelesai,
		&t.Aktif,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if deskripsi.Valid {
		t.Deskripsi = &deskripsi.String
	}
	if waktuMulai.Valid {
		t.WaktuMulai = &waktuMulai.Time
	}
	if waktuSelesai.Valid {
		t.WaktuSelesai = &waktuSelesai.Time
	}

	return &t, nil
}

func GetUserPendaftarID(db *sql.DB, userID int) (int, error) {
	var pendaftarID int
	err := db.QueryRow("SELECT id_pendaftar FROM pendaftar WHERE user_id = ?", userID).Scan(&pendaftarID)
	return pendaftarID, err
}

func HasUserTakenTest(db *sql.DB, userID, idTest int) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1 FROM hasil_test 
			WHERE user_id = ? AND id_test = ?
		)
	`
	err := db.QueryRow(query, userID, idTest).Scan(&exists)
	return exists, err
}

func GetHasilIDByUserAndTest(db *sql.DB, userID, idTest int) (int, error) {
	var idHasil int
	query := `
		SELECT id_hasil 
		FROM hasil_test 
		WHERE user_id = ? AND id_test = ? 
		ORDER BY created_at DESC 
		LIMIT 1
	`
	err := db.QueryRow(query, userID, idTest).Scan(&idHasil)
	return idHasil, err
}

func CreateHasilTest(db *sql.DB, hasil models.HasilTest) error {
	query := `
		INSERT INTO hasil_test (user_id, pendaftar_id, id_test, waktu_mulai, waktu_selesai, skor_benar, skor_salah, nilai)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.Exec(
		query,
		hasil.UserID,
		hasil.PendaftarID,
		hasil.IDTest,
		hasil.WaktuMulai,
		hasil.WaktuSelesai,
		hasil.SkorBenar,
		hasil.SkorSalah,
		hasil.Nilai,
	)
	return err
}

func GetHasilByUserID(db *sql.DB, userID, idTest int) (*models.HasilTest, error) {
	query := `
		SELECT 
			id_hasil, user_id, pendaftar_id, id_test,
			skor_benar, skor_salah, nilai,
			waktu_mulai, waktu_selesai,
			durasi_menit,
			created_at, updated_at
		FROM hasil_test
		WHERE user_id = ? AND id_test = ?
	`
	row := db.QueryRow(query, userID, idTest)

	var h models.HasilTest
	var waktuSelesai sql.NullTime
	var durasiMenit sql.NullInt64

	err := row.Scan(
		&h.IDHasil,
		&h.UserID,
		&h.PendaftarID,
		&h.IDTest,
		&h.SkorBenar,
		&h.SkorSalah,
		&h.Nilai,
		&h.WaktuMulai,
		&waktuSelesai,
		&durasiMenit,
		&h.CreatedAt,
		&h.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if waktuSelesai.Valid {
		h.WaktuSelesai = &waktuSelesai.Time
	}
	if durasiMenit.Valid {
		min := int(durasiMenit.Int64)
		h.DurasiMenit = &min
	}

	return &h, nil
}

func CreateJawabanUser(db *sql.DB, jawaban models.JawabanUser) error {
	query := `
		INSERT INTO jawaban_user (id_hasil, id_soal, jawaban_user, is_benar)
		VALUES (?, ?, ?, ?)
	`
	_, err := db.Exec(
		query,
		jawaban.IDHasil,
		jawaban.IDSoal,
		jawaban.JawabanUser,
		jawaban.IsBenar,
	)
	return err
}

func GetJawabanByHasilID(db *sql.DB, idHasil int) ([]models.JawabanUser, error) {
	query := `
		SELECT id_jawaban, id_hasil, id_soal, jawaban_user, is_benar, created_at
		FROM jawaban_user
		WHERE id_hasil = ?
		ORDER BY id_soal
	`
	rows, err := db.Query(query, idHasil)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jawabans []models.JawabanUser
	for rows.Next() {
		var j models.JawabanUser
		err := rows.Scan(
			&j.IDJawaban,
			&j.IDHasil,
			&j.IDSoal,
			&j.JawabanUser,
			&j.IsBenar,
			&j.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		jawabans = append(jawabans, j)
	}
	return jawabans, rows.Err()
}

func UpdateHasilTest(db *sql.DB, hasil *models.HasilTest) error {
	query := `
		UPDATE hasil_test SET
			skor_benar = ?,
			skor_salah = ?,
			nilai = ?,
			waktu_selesai = ?,
			updated_at = NOW()
		WHERE id_hasil = ?
	`
	_, err := db.Exec(
		query,
		hasil.SkorBenar,
		hasil.SkorSalah,
		hasil.Nilai,
		hasil.WaktuSelesai,
		hasil.IDHasil,
	)
	return err
}

func GetAllHasilTestByTestID(db *sql.DB, idTest int) (*sql.Rows, error) {
    query := `
    SELECT 
        id_hasil, user_id, pendaftar_id,
        skor_benar, skor_salah, nilai,
        waktu_mulai, waktu_selesai,
        durasi_menit
    FROM hasil_test
    WHERE id_test = ?
    ORDER BY nilai DESC, waktu_mulai ASC
`
    return db.Query(query, idTest)
}


func UpsertTestConfig(db *sql.DB, config models.Test) error {
    existing, err := GetTestConfig(db)
    if err != nil {
        if err == sql.ErrNoRows {
            return CreateTestConfig(db, config)
        }
        return err
    }

    query := `
        UPDATE test 
        SET 
            judul = ?, 
            deskripsi = ?, 
            durasi_menit = ?, 
            waktu_mulai = ?, 
            waktu_selesai = ?, 
            aktif = ?, 
            updated_at = CURRENT_TIMESTAMP 
        WHERE id_test = ?
    `

    _, err = db.Exec(
        query,
        config.Judul,
        config.Deskripsi,
        config.DurasiMenit,
        config.WaktuMulai,
        config.WaktuSelesai,
        config.Aktif,
        existing.IDTest,
    )
    return err
}


func CreateTestConfig(db *sql.DB, config models.Test) error {
    query := `
        INSERT INTO test 
        (judul, deskripsi, durasi_menit, waktu_mulai, waktu_selesai, aktif) 
        VALUES (?, ?, ?, ?, ?, ?)
    `
    _, err := db.Exec(
        query,
        config.Judul,
        config.Deskripsi,
        config.DurasiMenit,
        config.WaktuMulai,
        config.WaktuSelesai,
        config.Aktif,
    )
    return err
}