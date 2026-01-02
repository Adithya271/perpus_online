package models

type Book struct {
	IdBuku  int    `json:"id_buku" db:"id_buku"`
	Judul   string `json:"judul" db:"judul"`
	Penulis string `json:"penulis" db:"penulis"`
	Tahun   int    `json:"tahun" db:"tahun"`
	Status  string `json:"status" db:"status"`
	Stok    int    `json:"stok"` 
}

type Loan struct {
	IdPeminjaman  int    `json:"id_peminjaman" db:"id_peminjaman"`
	IdUser        int    `json:"id_user" db:"id_user"`
	IdBuku        int    `json:"id_buku" db:"id_buku"`
	TanggalPinjam string `json:"tanggal_pinjam" db:"tanggal_pinjam"`
	TanggalKembali *string `json:"tanggal_kembali" db:"tanggal_kembali"`
}

type User struct {
	IdUser int    `json:"id_user" db:"id_user"`
	Nama   string `json:"nama" db:"nama"`
	Role   string `json:"role" db:"role"`
}

type Recommendation struct {
	IdBuku  int    `json:"id_buku"`
	Judul   string `json:"judul"`
	Penulis string `json:"penulis"`
	Tahun   int    `json:"tahun"`
	Status  string `json:"status"`
	Score   int    `json:"score"`
	Reason  string `json:"reason"`
}
