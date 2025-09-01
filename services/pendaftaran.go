// services/pendaftaran.go
package services

import (
	"cocopen-backend/models"
    "cocopen-backend/utils"
	"database/sql"
)

func CreatePendaftar(db *sql.DB, p models.Pendaftar) error {
    query := `
        INSERT INTO pendaftar (
            nama_lengkap, asal_kampus, prodi, semester, no_wa, domisili,
            alamat_sekarang, tinggal_dengan, alasan_masuk, pengetahuan_coconut,
            foto_path, status, user_id
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `

    _, err := db.Exec(query,
        p.NamaLengkap,
        p.AsalKampus,
        p.Prodi,
        p.Semester,
        p.NoWA,
        p.Domisili,
        p.AlamatSekarang,
        p.TinggalDengan,
        p.AlasanMasuk,
        p.PengetahuanCoconut,
        p.FotoPath,
        p.Status,
        p.UserID,
    )

    return err
}

func GetLatestPendaftarByUserID(db *sql.DB, userID int) (*models.Pendaftar, error) {
    var p models.Pendaftar
    query := `
        SELECT 
            id_pendaftar, nama_lengkap, asal_kampus, prodi, semester, no_wa, domisili,
            alamat_sekarang, tinggal_dengan, alasan_masuk, pengetahuan_coconut,
            foto_path, created_at, updated_at, status, user_id
        FROM pendaftar
        WHERE user_id = ?
        ORDER BY created_at DESC
        LIMIT 1
    `
    err := db.QueryRow(query, userID).Scan(
        &p.IDPendaftar, &p.NamaLengkap, &p.AsalKampus, &p.Prodi, &p.Semester,
        &p.NoWA, &p.Domisili, &p.AlamatSekarang, &p.TinggalDengan,
        &p.AlasanMasuk, &p.PengetahuanCoconut, &p.FotoPath,
        &p.CreatedAt, &p.UpdatedAt, &p.Status, &p.UserID,
    )

    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }

    return &p, nil
}

func GetAllPendaftar(db *sql.DB) (*sql.Rows, error) {
    query := `
        SELECT 
            id_pendaftar, nama_lengkap, asal_kampus, prodi, semester,
            no_wa, domisili, alamat_sekarang, tinggal_dengan,
            alasan_masuk, pengetahuan_coconut, foto_path,
            created_at, updated_at, status, user_id
        FROM pendaftar
        ORDER BY created_at DESC
    `
    return db.Query(query)
}

func GetPendaftarByID(db *sql.DB, idPendaftar int) (models.Pendaftar, error) {
    var p models.Pendaftar
    err := db.QueryRow(`
        SELECT id_pendaftar, nama_lengkap, asal_kampus, prodi, semester, no_wa, domisili,
               alamat_sekarang, tinggal_dengan, alasan_masuk, pengetahuan_coconut, foto_path,
               created_at, updated_at, status
        FROM pendaftar
        WHERE id_pendaftar = ?
    `).Scan(
        &p.IDPendaftar, &p.NamaLengkap, &p.AsalKampus, &p.Prodi, &p.Semester, &p.NoWA,
        &p.Domisili, &p.AlamatSekarang, &p.TinggalDengan, &p.AlasanMasuk,
        &p.PengetahuanCoconut, &p.FotoPath, &p.CreatedAt, &p.UpdatedAt, &p.Status,
    )
    return p, err
}


func UpdatePendaftar(db *sql.DB, idPendaftar int, status string) error {
    query := `
        UPDATE pendaftar
        SET status = ?, updated_at = NOW()
        WHERE id_pendaftar = ?
    `
    _, err := db.Exec(query, status, idPendaftar)
    return err
}



func DeletePendaftar(db *sql.DB, idPendaftar int) error {
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    _, err = tx.Exec("DELETE FROM hasil_test WHERE pendaftar_id = ?", idPendaftar)
    if err != nil {
        return err
    }

    var exists int
    err = tx.QueryRow("SELECT 1 FROM pendaftar WHERE id_pendaftar = ?", idPendaftar).Scan(&exists)
    if err != nil {
        if err == sql.ErrNoRows {
            return sql.ErrNoRows
        }
        return err
    }

    var fotoPath string
    err = tx.QueryRow("SELECT foto_path FROM pendaftar WHERE id_pendaftar = ?", idPendaftar).Scan(&fotoPath)
    if err == nil && fotoPath != "" {
        utils.HapusFoto(utils.FotoPendaftarPath, fotoPath)
    }

    result, err := tx.Exec("DELETE FROM pendaftar WHERE id_pendaftar = ?", idPendaftar)
    if err != nil {
        return err
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        return sql.ErrNoRows
    }

    return tx.Commit()
}

func GetPendaftarByUserID(db *sql.DB, userID int) ([]models.Pendaftar, error) {
    query := `
        SELECT 
            id_pendaftar, nama_lengkap, asal_kampus, prodi, semester,
            no_wa, domisili, alamat_sekarang, tinggal_dengan,
            alasan_masuk, pengetahuan_coconut, foto_path,
            created_at, updated_at, status
        FROM pendaftar
        WHERE user_id = ?
        ORDER BY created_at DESC
    `

    rows, err := db.Query(query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var pendaftarList []models.Pendaftar
    for rows.Next() {
        var p models.Pendaftar
        err := rows.Scan(
            &p.IDPendaftar,
            &p.NamaLengkap,
            &p.AsalKampus,
            &p.Prodi,
            &p.Semester,
            &p.NoWA,
            &p.Domisili,
            &p.AlamatSekarang,
            &p.TinggalDengan,
            &p.AlasanMasuk,
            &p.PengetahuanCoconut,
            &p.FotoPath,
            &p.CreatedAt,
            &p.UpdatedAt,
            &p.Status,
        )
        if err != nil {
            return nil, err
        }
        pendaftarList = append(pendaftarList, p)
    }

    return pendaftarList, nil
}