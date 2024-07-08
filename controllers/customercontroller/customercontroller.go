package customercontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go-restapi-gin/models"

	"gorm.io/gorm"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
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

func Create(c *gin.Context) {

	var customer models.Customer

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	fmt.Println(2222, &customer)

	// models.DB.Create(&customer)

	retData := CreateCustomersMessage(customer)

	if retData == "success" {
		c.JSON(http.StatusOK, gin.H{"customer": customer})
	} else {
		c.JSON(http.StatusInternalServerError, "Error Save Data")
	}

}

func failOnError(err error, msg string) string {
	if err != nil {
		return "Error"
	} else {
		return ""
	}
}

func CreateCustomersMessage(customer models.Customer) string {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	msg := failOnError(err, "Failed to connect to RabbitMQ")
	if msg == "Error" {
		return "Error"
	}
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"createCustomers", // name
		false,             // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	msg = failOnError(err, "Failed to declare a queue")
	if msg == "Error" {
		return "Error"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "send data"

	fmt.Println(1111, &customer)

	dataSend, err := json.Marshal(&customer)

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(dataSend),
		})
	msg = failOnError(err, "Failed to publish a message")
	if msg == "Error" {
		return "Error"
	}
	log.Printf(" [x] Sent %s\n", body)

	return "success"

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
