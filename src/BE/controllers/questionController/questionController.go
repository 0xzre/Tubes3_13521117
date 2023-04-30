package questionController

import (
	"net/http"

	"BE/models"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func IndexQuestion(c *gin.Context) {

	var question []models.Question

	models.DB.Find(&question)
	c.JSON(http.StatusOK, gin.H{"question": question})

}

func ShowQuestion(c *gin.Context) {
	var question models.Question
	pertanyaan := c.Param("pertanyaan")

	if err := models.DB.First(&question, pertanyaan).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Pertanyaan tidak ditemukan"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"question": question})
}

func CreateQuestion(c *gin.Context) {

	var question models.Question

	if err := c.ShouldBindJSON(&question); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	models.DB.Create(&question)
	c.JSON(http.StatusOK, gin.H{"question": question})
}

func UpdateQuestion(c *gin.Context) {
	var question models.Question
	pertanyaan := c.Param("pertanyaan")

	if err := c.ShouldBindJSON(&question); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if models.DB.Model(&question).Where("pertanyaan = ?", pertanyaan).Updates(&question).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Tidak dapat mengupdate pertanyaan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pertanyaan berhasil diperbarui"})

}

func DeleteQuestion(c *gin.Context) {

	var question models.Question

	var input struct {
		Pertanyaan string `json:"pertanyaan"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	pertanyaan := input.Pertanyaan
	if models.DB.Delete(&question, pertanyaan).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Tidak dapat menghapus pertanyaan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pertanyaan berhasil dihapus"})
}
