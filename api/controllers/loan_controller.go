package controllers

import (
	"database/sql"
	"net/http"
	"sistem_perpus/config"
	"time"

	"github.com/gin-gonic/gin"
)

func BorrowBook(c *gin.Context) {
	db := config.ConnectDB()
	defer db.Close()

	var req struct {
		IdUser int `json:"id_user" binding:"required"`
		IdBuku int `json:"id_buku" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	// Cek stok buku
	var stok int
	err := db.QueryRow("SELECT stok FROM books WHERE id_buku = ?", req.IdBuku).Scan(&stok)
	
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Buku tidak ditemukan"})
		return
	}

	// Validasi stok
	if stok <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stok buku habis"})
		return
	}

	// Cek apakah user sudah meminjam buku ini
	var count int
	db.QueryRow(`
		SELECT COUNT(*) FROM loans 
		WHERE id_user = ? AND id_buku = ? AND tanggal_kembali IS NULL
	`, req.IdUser, req.IdBuku).Scan(&count)

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Anda sudah meminjam buku ini"})
		return
	}

	// Mulai transaksi
	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memulai transaksi"})
		return
	}

	// Insert loan
	_, err = tx.Exec(
		"INSERT INTO loans (id_user, id_buku, tanggal_pinjam) VALUES (?, ?, ?)",
		req.IdUser, req.IdBuku, time.Now().Format("2006-01-02"),
	)

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal meminjam buku"})
		return
	}

// Kurangi stok
_, err = tx.Exec(`
    UPDATE books 
    SET stok = stok - 1
    WHERE id_buku = ?
`, req.IdBuku)
if err != nil {
    tx.Rollback()
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengurangi stok"})
    return
}

// Update status buku
_, err = tx.Exec(`
    UPDATE books
    SET status = CASE
        WHEN stok <= 0 THEN 'habis'
        ELSE 'dipinjam'
    END
    WHERE id_buku = ?
`, req.IdBuku)
if err != nil {
    tx.Rollback()
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update status buku"})
    return
}



	// Get data terbaru setelah update
	var stokBaru int
	var statusBaru string
	err = tx.QueryRow("SELECT stok, status FROM books WHERE id_buku = ?", req.IdBuku).Scan(&stokBaru, &statusBaru)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan data terbaru"})
		return
	}

	// Commit transaksi
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan transaksi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Peminjaman berhasil",
		"data": gin.H{
			"id_buku":        req.IdBuku,
			"tanggal_pinjam": time.Now().Format("2006-01-02"),
			"stok_tersisa":   stokBaru,
			"status":         statusBaru,
		},
	})
}

func ReturnBook(c *gin.Context) {
	db := config.ConnectDB()
	defer db.Close()

	var req struct {
		IdUser int `json:"id_user" binding:"required"`
		IdBuku int `json:"id_buku" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	// Cek apakah peminjaman ada
	var count int
	db.QueryRow(`
		SELECT COUNT(*) FROM loans 
		WHERE id_user = ? AND id_buku = ? AND tanggal_kembali IS NULL
	`, req.IdUser, req.IdBuku).Scan(&count)

	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Peminjaman tidak ditemukan"})
		return
	}

	// Mulai transaksi
	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memulai transaksi"})
		return
	}

	// Update loan
	_, err = tx.Exec(`
		UPDATE loans 
		SET tanggal_kembali = ? 
		WHERE id_user = ? AND id_buku = ? AND tanggal_kembali IS NULL
	`, time.Now().Format("2006-01-02"), req.IdUser, req.IdBuku)

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengembalikan buku"})
		return
	}

// Tambah stok
_, err = tx.Exec(`
    UPDATE books
    SET stok = stok + 1
    WHERE id_buku = ?
`, req.IdBuku)
if err != nil {
    tx.Rollback()
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambah stok"})
    return
}

// Update status jadi tersedia
_, err = tx.Exec(`
    UPDATE books
    SET status = 'tersedia'
    WHERE id_buku = ?
`, req.IdBuku)
if err != nil {
    tx.Rollback()
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update status"})
    return
}


	// Get stok terbaru setelah update
	var stokBaru int
	err = tx.QueryRow("SELECT stok FROM books WHERE id_buku = ?", req.IdBuku).Scan(&stokBaru)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan stok terbaru"})
		return
	}

	// Commit transaksi
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan transaksi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "Pengembalian berhasil",
		"tanggal_kembali": time.Now().Format("2006-01-02"),
		"stok_sekarang":   stokBaru,
	})
}

func GetUserLoans(c *gin.Context) {
	db := config.ConnectDB()
	defer db.Close()

	userId := c.Param("id")

	rows, err := db.Query(`
		SELECT l.id_peminjaman, l.id_buku, b.judul, b.penulis, 
		       l.tanggal_pinjam, l.tanggal_kembali
		FROM loans l
		JOIN books b ON l.id_buku = b.id_buku
		WHERE l.id_user = ?
		ORDER BY l.tanggal_pinjam DESC
	`, userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data peminjaman"})
		return
	}
	defer rows.Close()

	var loans []gin.H
	for rows.Next() {
		var id, idBuku int
		var judul, penulis, tanggalPinjam string
		var tanggalKembali sql.NullString

		rows.Scan(&id, &idBuku, &judul, &penulis, &tanggalPinjam, &tanggalKembali)

		loan := gin.H{
			"id_peminjaman":  id,
			"id_buku":        idBuku,
			"judul":          judul,
			"penulis":        penulis,
			"tanggal_pinjam": tanggalPinjam,
		}

		if tanggalKembali.Valid {
			loan["tanggal_kembali"] = tanggalKembali.String
			loan["status"] = "dikembalikan"
		} else {
			loan["tanggal_kembali"] = nil
			loan["status"] = "dipinjam"
		}

		loans = append(loans, loan)
	}

	if loans == nil {
		loans = []gin.H{}
	}

	c.JSON(http.StatusOK, loans)
}