package models

type Book struct {
	IdBuku  int    `json:"id_buku" db:"id_buku"`
	Judul   string `json:"judul" db:"judul"`
	Penulis string `json:"penulis" db:"penulis"`
	Tahun   int    `json:"tahun" db:"tahun"`
	Status  string `json:"status" db:"status"`
	Stok    int    `json:"stok"` 
}

