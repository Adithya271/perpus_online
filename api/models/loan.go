
package models

type Loan struct {
	IdPeminjaman  int    `json:"id_peminjaman" db:"id_peminjaman"`
	IdUser        int    `json:"id_user" db:"id_user"`
	IdBuku        int    `json:"id_buku" db:"id_buku"`
	TanggalPinjam string `json:"tanggal_pinjam" db:"tanggal_pinjam"`
	TanggalKembali *string `json:"tanggal_kembali" db:"tanggal_kembali"`
}