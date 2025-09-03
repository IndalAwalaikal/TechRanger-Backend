package services

import (
    "database/sql"
    "time"
    "log"
)

func Register(db *sql.DB, username, email, passwordHash, role string) (int, error) {
    log.Printf("Mencoba register: username=%s, email=%s, role=%s", username, email, role)

    res, err := db.Exec(`
        INSERT INTO users 
        (username, password, role, email, is_verified, full_name, profile_picture) 
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `, username, passwordHash, role, email, false, username, "default.jpg")

    if err != nil {
        log.Printf("ERROR DB: %v", err)
        return 0, err
    }

    lastID, err := res.LastInsertId()
    if err != nil {
        log.Printf("ERROR LastInsertId: %v", err)
        return 0, err
    }

    log.Printf("Register berhasil, ID: %d", lastID)
    return int(lastID), nil
}
func GetUserByUsername(db *sql.DB, username string) (int, string, string, string, string, error) {
    var id int
    var uname, email, hashed, role string

    err := db.QueryRow(`
        SELECT id_user, username, email, password, role
        FROM users
        WHERE username = ? AND is_verified = TRUE
    `, username).Scan(&id, &uname, &email, &hashed, &role)

    return id, uname, email, hashed, role, err
}

func GetUserByEmail(db *sql.DB, email string) (int, error) {
    var id int
    err := db.QueryRow("SELECT id_user FROM users WHERE email = ?", email).Scan(&id)
    return id, err
}

func DeleteUnverifiedUsersBefore(db *sql.DB, cutoff time.Time) error {
    tx, err := db.Begin()
    if err != nil {
        log.Printf("Gagal memulai transaksi: %v", err)
        return err
    }
    defer tx.Rollback()

    _, err = tx.Exec(`
        DELETE pr FROM password_resets pr
        INNER JOIN users u ON pr.user_id = u.id_user
        WHERE u.is_verified = 0 AND u.created_at < ?
    `, cutoff)
    if err != nil {
        log.Printf("Gagal hapus dari password_resets: %v", err)
        return err
    }

    _, err = tx.Exec(`
        DELETE evt FROM email_verification_tokens evt
        INNER JOIN users u ON evt.user_id = u.id_user
        WHERE u.is_verified = 0 AND u.created_at < ?
    `, cutoff)
    if err != nil {
        log.Printf("Gagal hapus dari email_verification_tokens: %v", err)
        return err
    }

    result, err := tx.Exec(`
        DELETE FROM users 
        WHERE is_verified = 0 AND created_at < ?
    `, cutoff)
    if err != nil {
        log.Printf("Gagal hapus dari users: %v", err)
        return err
    }

    rowsAffected, _ := result.RowsAffected()
    log.Printf("Berhasil hapus %d akun belum diverifikasi", rowsAffected)

    if err = tx.Commit(); err != nil {
        log.Printf("Gagal commit transaksi: %v", err)
        return err
    }

    return nil
}