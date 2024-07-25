package services

import (
	"fmt"
	"go-restapi-gin/models"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/xuri/excelize/v2"
)

type customerService struct {
	customers []models.Customer
}

func ExcelCustomers(c *gin.Context) {
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
