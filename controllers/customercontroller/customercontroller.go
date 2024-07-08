package customercontroller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go-restapi-gin/models"

	"gorm.io/gorm"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {

	var customers []models.Customer

	models.DB.Find(&customers)
	c.JSON(http.StatusOK, gin.H{"customers": customers})

}

func Show(c *gin.Context) {
	var customer models.Customer
	id := c.Param("id")

	if err := models.DB.First(&customer, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Data tidak ditemukan"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"customer": customer})
}

func ExcelKonsumens(c *gin.Context) {

	var customers []models.Customer

	models.DB.Find(&customers)
	xlsx := excelize.NewFile()
	sheetName := "Sheet1"

	xlsx.SetSheetName("Sheet1", sheetName)

	// Add headers
	xlsx.SetCellValue(sheetName, "A1", "Email")
	xlsx.SetCellValue(sheetName, "B1", "Name")
	xlsx.SetCellValue(sheetName, "C1", "Password")
	// Create a new sheet.
	rowIndex := 2
	for _, data := range customers {
		xlsx.SetCellValue(sheetName, fmt.Sprintf("A%d", rowIndex), data.Email)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("B%d", rowIndex), data.Name)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("C%d", rowIndex), data.Password)

		rowIndex++
	}

	filename := "static/excel/" + c.Param("filename") + ".xlsx"

	if err := xlsx.SaveAs(filename); err != nil {
		log.Fatal(err)
	}
}

func Create(c *gin.Context) {

	var customer models.Customer

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	models.DB.Create(&customer)
	c.JSON(http.StatusOK, gin.H{"customer": customer})
}

func Update(c *gin.Context) {
	var customer models.Customer
	id := c.Param("id")

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if models.DB.Model(&customer).Where("id = ?", id).Updates(&customer).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "tidak dapat mengupdate customer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil diperbarui"})

}

func Delete(c *gin.Context) {

	var customer models.Customer

	var input struct {
		Id json.Number
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	id, _ := input.Id.Int64()
	if models.DB.Delete(&customer, id).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Tidak dapat menghapus customer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil dihapus"})
}
