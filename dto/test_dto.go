package dto

import "time"

type SoalCreateRequest struct {
    Nomor        int    `json:"nomor" validate:"required,min=1"`
    Pertanyaan   string `json:"pertanyaan" validate:"required,min=10"`
    PilihanA     string `json:"pilihan_a" validate:"required,min=1"`
    PilihanB     string `json:"pilihan_b" validate:"required,min=1"`
    PilihanC     string `json:"pilihan_c" validate:"required,min=1"`
    PilihanD     string `json:"pilihan_d" validate:"required,min=1"`
    JawabanBenar string `json:"jawaban_benar" validate:"required,oneof=A B C D"`
}

type SoalUpdateRequest struct {
    Nomor        *int    `json:"nomor,omitempty" validate:"omitempty,min=1"`
    Pertanyaan   *string `json:"pertanyaan,omitempty" validate:"omitempty,min=10"`
    PilihanA     *string `json:"pilihan_a,omitempty" validate:"omitempty,min=1"`
    PilihanB     *string `json:"pilihan_b,omitempty" validate:"omitempty,min=1"`
    PilihanC     *string `json:"pilihan_c,omitempty" validate:"omitempty,min=1"`
    PilihanD     *string `json:"pilihan_d,omitempty" validate:"omitempty,min=1"`
    JawabanBenar *string `json:"jawaban_benar,omitempty" validate:"omitempty,oneof=A B C D"`
}

type SubmitJawabanRequest struct {
    Jawaban map[int]string `json:"jawaban" validate:"required,dive,keys,required,endkeys,oneof=A B C D"`
}

type HasilResponse struct {
    IDHasil      int     `json:"id_hasil"`
    UserID       int     `json:"user_id"`
    PendaftarID  int     `json:"pendaftar_id"`
    SkorBenar    int     `json:"skor_benar"`
    SkorSalah    int     `json:"skor_salah"`
    Nilai        float64 `json:"nilai"`
    WaktuMulai   string  `json:"waktu_mulai"`
    WaktuSelesai *string `json:"waktu_selesai,omitempty"`
    DurasiMenit  *int    `json:"durasi_menit,omitempty"`
}

type SoalResponse struct {
    IDSoal     int    `json:"id_soal"`
    Nomor      int    `json:"nomor"`
    Pertanyaan string `json:"pertanyaan"`
    PilihanA   string `json:"pilihan_a"`
    PilihanB   string `json:"pilihan_b"`
    PilihanC   string `json:"pilihan_c"`
    PilihanD   string `json:"pilihan_d"`
}

type TestConfigUpdateRequest struct {
    Judul           *string    `json:"judul,omitempty" validate:"omitempty,min=1"`
    Deskripsi       *string    `json:"deskripsi,omitempty"`
    DurasiMenit     *int       `json:"durasi_menit,omitempty" validate:"omitempty,min=1"`
    WaktuMulai      *time.Time `json:"waktu_mulai,omitempty"`
    WaktuSelesai    *time.Time `json:"waktu_selesai,omitempty"`
    Aktif           *bool      `json:"aktif,omitempty"`
}