package models

type Recommendation struct {
	IdBuku  int    `json:"id_buku"`
	Judul   string `json:"judul"`
	Penulis string `json:"penulis"`
	Tahun   int    `json:"tahun"`
	Status  string `json:"status"`
}