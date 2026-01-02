package controllers

import (
	"database/sql"
	"net/http"
	"strconv"
	"sistem_perpus/config"
	"sistem_perpus/models"

	"github.com/gin-gonic/gin"
)

func GetBooks(c *gin.Context) {
	db := config.ConnectDB()
	defer db.Close()

	keyword := c.Query("keyword")
	status := c.Query("status")
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "20")

	
	query := "SELECT id_buku, judul, penulis, tahun, status, stok FROM books WHERE 1=1"
	args := []interface{}{}

	if keyword != "" {
		query += " AND (judul LIKE ? OR penulis LIKE ?)"
		searchTerm := "%" + keyword + "%"
		args = append(args, searchTerm, searchTerm)
	}

	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}

	query += " ORDER BY judul ASC LIMIT ? OFFSET ?"
	
	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)
	
	var offset int
	if pageInt == 1 {
		offset = 0
	} else {
		offset = (pageInt - 1) * limitInt
	}
	
	args = append(args, limitInt, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengambil data buku",
		})
		return
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var b models.Book
		
		if err := rows.Scan(&b.IdBuku, &b.Judul, &b.Penulis, &b.Tahun, &b.Status, &b.Stok); err != nil {
			continue
		}
		books = append(books, b)
	}

	if books == nil {
		books = []models.Book{}
	}

	c.JSON(http.StatusOK, books)
}

func GetBookDetail(c *gin.Context) {
	db := config.ConnectDB()
	defer db.Close()

	id := c.Param("id")
	
	var book models.Book
	
	err := db.QueryRow(
		"SELECT id_buku, judul, penulis, tahun, status, stok FROM books WHERE id_buku = ?",
		id,
	).Scan(&book.IdBuku, &book.Judul, &book.Penulis, &book.Tahun, &book.Status, &book.Stok)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Buku tidak ditemukan"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil detail buku"})
		return
	}

	c.JSON(http.StatusOK, book)
}