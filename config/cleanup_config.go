package config

import (
	"cocopen-backend/services"
	"database/sql"
	"log"
	"time"
)

func StartCleanupJob(db *sql.DB) {
    log.Println("Background job: Membersihkan akun belum diverifikasi sekali setiap 24 jam")

    ticker := time.NewTicker(24 * time.Hour)

    go func() {
        runCleanup(db)

        for range ticker.C {
            runCleanup(db)
        }
    }()
}

func runCleanup(db *sql.DB) {
    cutoff := time.Now().Add(-30 * time.Minute) 
    err := services.DeleteUnverifiedUsersBefore(db, cutoff)
    if err != nil {
        log.Printf("Gagal membersihkan akun belum diverifikasi: %v", err)
    } else {
        log.Printf("Pembersihan selesai: akun dibuat sebelum %v telah dicek", cutoff.Format("2006-01-02 15:04:05"))
    }
}