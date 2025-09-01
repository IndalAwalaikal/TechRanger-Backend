// dto/jadwal/jadwal_user_response.go
package dto

import "time"

type JadwalUserResponse struct {
    IDJadwal                 int        `json:"id_jadwal"`
    Tanggal                  time.Time  `json:"tanggal"`
    JamMulai                 string     `json:"jam_mulai"`
    JamSelesai               string     `json:"jam_selesai"`
    Tempat                   string     `json:"tempat"`
    KonfirmasiJadwal         string     `json:"konfirmasi_jadwal"`
    Catatan                  *string    `json:"catatan,omitempty"`
    PengajuanPerubahan       bool       `json:"pengajuan_perubahan"`
    AlasanPerubahan          *string    `json:"alasan_perubahan,omitempty"`
    TanggalDiajukan          *time.Time `json:"tanggal_diajukan,omitempty"`
    JamMulaiDiajukan         *string    `json:"jam_mulai_diajukan,omitempty"`
    JamSelesaiDiajukan       *string    `json:"jam_selesai_diajukan,omitempty"`
}

type JadwalAdminResponse struct {
    IDJadwal           int        `json:"id_jadwal"`
    UserID             int        `json:"user_id"`
    PendaftarID        *int       `json:"pendaftar_id,omitempty"`
    Tanggal            time.Time  `json:"tanggal"`
    JamMulai           string     `json:"jam_mulai"`
    JamSelesai         string     `json:"jam_selesai"`
    Tempat             string     `json:"tempat"`
    KonfirmasiJadwal   string     `json:"konfirmasi_jadwal"` 
    Catatan            *string    `json:"catatan,omitempty"`
    PengajuanPerubahan bool       `json:"pengajuan_perubahan"`
    AlasanPerubahan    *string    `json:"alasan_perubahan,omitempty"`
    TanggalDiajukan    *time.Time `json:"tanggal_diajukan,omitempty"`
    JamMulaiDiajukan   *string    `json:"jam_mulai_diajukan,omitempty"`
    JamSelesaiDiajukan *string    `json:"jam_selesai_diajukan,omitempty"`
    CreatedAt          time.Time  `json:"created_at"`
    UpdatedAt          time.Time  `json:"updated_at"`
    JenisJadwal        string     `json:"jenis_jadwal"`
    UserNama           string     `json:"user_nama,omitempty"`
    UserEmail          string     `json:"user_email,omitempty"`
    PendaftarNama      string     `json:"pendaftar_nama,omitempty"`
}

type JadwalCreateRequest struct {
    UserID      *int    `json:"user_id,omitempty"`
    PendaftarID *int    `json:"pendaftar_id,omitempty"`
    Tanggal     string  `json:"tanggal" validate:"required,datetime=2006-01-02"`
    JamMulai    string  `json:"jam_mulai" validate:"required,datetime=15:04:05"`
    JamSelesai  string  `json:"jam_selesai" validate:"required,datetime=15:04:05"`
    Tempat      string  `json:"tempat" validate:"required,min=3,max=255"`
    Catatan     *string `json:"catatan,omitempty"`
    JenisJadwal *string  `json:"jenis_jadwal,omitempty" validate:"omitempty,oneof=pribadi umum"`
}

type JadwalUpdateRequest struct {
    PendaftarID      *int    `json:"pendaftar_id,omitempty"`
    Tanggal          *string `json:"tanggal,omitempty" validate:"omitempty,datetime=2006-01-02"`
    JamMulai         *string `json:"jam_mulai,omitempty" validate:"omitempty,datetime=15:04:05"`
    JamSelesai       *string `json:"jam_selesai,omitempty" validate:"omitempty,datetime=15:04:05"`
    Tempat           *string `json:"tempat,omitempty" validate:"omitempty,min=3,max=255"`
    KonfirmasiJadwal *string `json:"konfirmasi_jadwal,omitempty" validate:"omitempty,oneof=belum dikonfirmasi ditolak"`
    Catatan          *string `json:"catatan,omitempty"`
    JenisJadwal      *string `json:"jenis_jadwal,omitempty" validate:"omitempty,oneof=pribadi umum"`
}

type JadwalAjukanPerubahanRequest struct {
    IDJadwal             int        `json:"id_jadwal" validate:"required"`
    TanggalDiajukan      *string    `json:"tanggal_diajukan,omitempty" validate:"omitempty,datetime=2006-01-02"`
    JamMulaiDiajukan     *string    `json:"jam_mulai_diajukan,omitempty" validate:"omitempty,datetime=15:04:05"`
    JamSelesaiDiajukan   *string    `json:"jam_selesai_diajukan,omitempty" validate:"omitempty,datetime=15:04:05"`
    AlasanPerubahan      string     `json:"alasan_perubahan" validate:"required,min=10,max=500"`
}