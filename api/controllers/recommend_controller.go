

package controllers

import (
	"net/http"
	"strings"
	"sistem_perpus/config"
	"sistem_perpus/models"

	"github.com/gin-gonic/gin"
)

// GetRecommendationsBySearchHistory
func GetRecommendationsBySearchHistory(c *gin.Context) {
	db := config.ConnectDB()
	defer db.Close()

	userId := c.Param("id")

	keywordQuery := `
		SELECT keyword, MAX(search_date) as latest_date
		FROM search_history 
		WHERE id_user = ? 
		GROUP BY keyword
		ORDER BY latest_date DESC 
		LIMIT 20
	`
	
	keywordRows, err := db.Query(keywordQuery, userId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"recommendations": []models.Recommendation{},
			"total": 0,
			"message": "Belum ada riwayat pencarian",
		})
		return
	}
	defer keywordRows.Close()

	var keywords []string
	for keywordRows.Next() {
		var keyword string
		var latestDate string
		if err := keywordRows.Scan(&keyword, &latestDate); err == nil {
			keywords = append(keywords, keyword)
		}
	}

	if len(keywords) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"recommendations": []models.Recommendation{},
			"total": 0,
			"message": "Belum ada riwayat pencarian",
		})
		return
	}

	
	var whereClauses []string
	var args []interface{}
	
	for _, keyword := range keywords {
		whereClauses = append(whereClauses, "(b.judul LIKE ? OR b.penulis LIKE ?)")
		searchTerm := "%" + keyword + "%"
		args = append(args, searchTerm, searchTerm)
	}
	
	// userId di akhir
	args = append(args, userId)

	// Hanya exclude buku yang SEDANG dipinjam (tanggal_kembali IS NULL)
	query := `
		SELECT DISTINCT
			b.id_buku,
			b.judul,
			b.penulis,
			b.tahun,
			b.status
		FROM books b
		WHERE (` + strings.Join(whereClauses, " OR ") + `)
		AND b.status = 'tersedia'
		AND b.id_buku NOT IN (
			SELECT id_buku FROM loans 
			WHERE id_user = ? AND tanggal_kembali IS NULL
		)
		ORDER BY b.tahun DESC
		LIMIT 10
	`

	rows, err := db.Query(query, args...)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengambil rekomendasi",
			"detail": err.Error(),
		})
		return
	}
	defer rows.Close()

	var recommendations []models.Recommendation
	
	for rows.Next() {
		var rec models.Recommendation
		
		err := rows.Scan(
			&rec.IdBuku,
			&rec.Judul,
			&rec.Penulis,
			&rec.Tahun,
			&rec.Status,
		)
		
		if err != nil {
			continue
		}
		
		recommendations = append(recommendations, rec)
	}

	if recommendations == nil {
		recommendations = []models.Recommendation{}
	}

	c.JSON(http.StatusOK, gin.H{
		"recommendations": recommendations,
		"total": len(recommendations),
	})
}

// ========================================
// SAVE SEARCH HISTORY
// ========================================
func SaveSearchHistory(c *gin.Context) {
	db := config.ConnectDB()
	defer db.Close()

	var req struct {
		IdUser  int    `json:"id_user" binding:"required"`
		Keyword string `json:"keyword" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload tidak valid"})
		return
	}

	keyword := strings.ToLower(strings.TrimSpace(req.Keyword))
	if len(keyword) < 2 {
		c.JSON(http.StatusOK, gin.H{"message": "Keyword terlalu pendek"})
		return
	}

	// Insert langsung 
	_, err := db.Exec(`
		INSERT INTO search_history (id_user, keyword, search_date)
		VALUES (?, ?, NOW())
	`, req.IdUser, keyword)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menyimpan riwayat",
			"detail": err.Error(),
		})
		return
	}

	// Cleanup max 100 data
	db.Exec(`
		DELETE FROM search_history
		WHERE id_user = ?
		AND id_search NOT IN (
			SELECT id_search FROM (
				SELECT id_search
				FROM search_history
				WHERE id_user = ?
				ORDER BY search_date DESC
				LIMIT 100
			) t
		)
	`, req.IdUser, req.IdUser)

	c.JSON(http.StatusOK, gin.H{"message": "Riwayat pencarian disimpan"})
}