// models/jadwal.go
package models

import "time"

type Jadwal struct {
    IDJadwal           int        `json:"id_jadwal"`
    UserID             int        `json:"user_id"`
    PendaftarID        *int       `json:"pendaftar_id"`        
    Tanggal            time.Time  `json:"tanggal"`
    JamMulai           string     `json:"jam_mulai"`            
    JamSelesai         string     `json:"jam_selesai"`         
    Tempat             string     `json:"tempat"`
    KonfirmasiJadwal   string     `json:"konfirmasi_jadwal"`    
    Catatan            *string    `json:"catatan"`              
    PengajuanPerubahan bool       `json:"pengajuan_perubahan"`  
    AlasanPerubahan    *string    `json:"alasan_perubahan"`     

    TanggalDiajukan    *time.Time `json:"tanggal_diajukan,omitempty"`
    JamMulaiDiajukan   *string    `json:"jam_mulai_diajukan,omitempty"`
    JamSelesaiDiajukan *string    `json:"jam_selesai_diajukan,omitempty"`

    CreatedAt          time.Time  `json:"created_at"`
    UpdatedAt          time.Time  `json:"updated_at"`

    JenisJadwal        string     `json:"jenis_jadwal"`

    // Relasi (opsional)
    User      *User      `json:"user,omitempty"`
    Pendaftar *Pendaftar `json:"pendaftar,omitempty"`
    UserNama  string     `json:"user_nama,omitempty"`
}