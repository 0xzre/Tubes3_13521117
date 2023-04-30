package historyController

import (
	"net/http"

	"BE/models"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func IndexHistory(c *gin.Context) {

	var history []models.History

	models.DB.Find(&history)
	c.JSON(http.StatusOK, gin.H{"history": history})

}

func ShowHistory(c *gin.Context) {
	var history models.History
	prompt := c.Param("prompt")

	if err := models.DB.First(&history, prompt).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "History tidak ditemukan"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"history": history})
}

func CreateHistory(c *gin.Context) {

	var history models.History

	if err := c.ShouldBindJSON(&history); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	models.DB.Create(&history)
	c.JSON(http.StatusOK, gin.H{"history": history})
}

func UpdateHistory(c *gin.Context) {
	var history models.History
	prompt := c.Param("prompt")

	if err := c.ShouldBindJSON(&history); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if models.DB.Model(&history).Where("prompt = ?", prompt).Updates(&history).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Tidak dapat mengupdate history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "History berhasil diperbarui"})

}

func DeleteHistory(c *gin.Context) {

	var history models.History

	var input struct {
		Pertanyaan string `json:"prompt"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	prompt := input.Pertanyaan
	if models.DB.Delete(&history, prompt).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Tidak dapat menghapus history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "History berhasil dihapus"})
}
