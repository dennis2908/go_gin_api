package customercontroller

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"go-restapi-gin/models"

	"go-restapi-gin/services"

	"go-restapi-gin/token"

	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/xuri/excelize/v2"

	redisconn "go-restapi-gin/redisconn"

	"github.com/joho/godotenv"

	"time"

	"crypto/aes"
	"crypto/cipher"
)

var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

type Customercontroller interface {
	ExcelCustomers()
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func init() {

	err := godotenv.Load("app.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func GenerateRefreshToken(c *gin.Context) {
	val, _ := c.Cookie("token")

	tokenRef, errx := Decrypt(val)

	if errx != nil {
		fmt.Println(errx)
	}

	fmt.Println(133, tokenRef)

	_, err := jwt.Parse(tokenRef, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_TOKEN_KEY")), nil
	})

	if err == nil {

		token, _ := token.GenerateToken()

		c.JSON(http.StatusOK, gin.H{"new token": token})

	} else {
		fmt.Println(err)
	}
}

func Login(c *gin.Context) {

	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.Customer{}

	u.UserName = input.Username
	u.Password = input.Password

	token, err := LoginCheck(u.UserName, u.Password, c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

}

func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Encrypt(keytext string) (string, error) {
	block, err := aes.NewCipher([]byte(os.Getenv(("SECRET_KEY"))))
	if err != nil {
		return "", err
	}
	plainText := []byte(keytext)
	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return Encode(cipherText), nil
}

func Decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

func Decrypt(text string) (string, error) {
	block, err := aes.NewCipher([]byte(os.Getenv(("SECRET_KEY"))))
	if err != nil {
		return "", err
	}
	cipherText := Decode(text)
	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}

func SetCookieToken(c *gin.Context) {
	refToken, err := token.GenerateRefreshToken()
	if err == nil {
		fmt.Println(121221, refToken)
		cookie, _ := Encrypt(refToken)
		fmt.Println(88, cookie)
		c.SetCookie("token", cookie, 3600000, "/", "localhost", false, true)
	}
}

func LoginCheck(username string, password string, c *gin.Context) (string, error) {

	var err error

	u := models.Customer{}

	err = models.DB.Where("user_name = ?", username).Take(&u).Error

	if err != nil {
		return "", err
	}

	err = VerifyPassword(password, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	token, err := token.GenerateToken()

	println(token)

	if err != nil {
		return "", err
	}

	SetCookieToken(c)

	return token, nil

}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func Index(c *gin.Context) {

	var customers []models.Customer

	models.DB.Find(&customers)

	rdb, _ := redisconn.Connect()
	urlsJson, _ := json.Marshal(customers)
	token, _ := GenerateRandomString(32)

	ttl := time.Duration(3) * time.Second

	op1 := rdb.Set(context.Background(), token, urlsJson, ttl)
	if err := op1.Err(); err != nil {
		fmt.Printf("unable to SET data. error: %v", err)
		return
	}
	op2 := rdb.Get(context.Background(), token)
	fmt.Printf("data", op2)
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

	rdb, _ := redisconn.Connect()
	urlsJson, _ := json.Marshal(customer)
	token, _ := GenerateRandomString(32)

	ttl := time.Duration(3) * time.Second

	op1 := rdb.Set(context.Background(), token, urlsJson, ttl)
	if err := op1.Err(); err != nil {
		fmt.Printf("unable to SET data. error: %v", err)
		return
	}
	op2 := rdb.Get(context.Background(), token)
	fmt.Printf("data", op2)

	c.JSON(http.StatusOK, gin.H{"customer": customer})
}

func ExcelCustomers(c *gin.Context) {
	services.ExcelCustomers(c)
}

func EmailCustomersMessage(customer models.Customer) string {
	conn, err := amqp.Dial(os.Getenv("rabbit_url"))
	msg := failOnError(err, "Failed to connect to RabbitMQ")
	if msg == "Error" {
		return "Error"
	}
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"sendEmailToCustomers", // name
		false,                  // durable
		false,                  // delete when unused
		false,                  // exclusive
		false,                  // no-wait
		nil,                    // arguments
	)
	msg = failOnError(err, "Failed to declare a queue")
	if msg == "Error" {
		return "Error"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "send data"

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

func ReadExcelKonsumens(c *gin.Context) {

	xlsx, err := excelize.OpenFile("static/excel/" + c.Param("filename") + ".xlsx")
	if err != nil {
		log.Fatal("ERROR", err.Error())
	}

	rowsExcel, _ := xlsx.GetRows("Sheet1")

	rows := make([]models.Customer, 0)
	for i, rowsExcel := range rowsExcel {
		if i == 0 {
			// Skip header row
			continue
		}
		rowEmail := rowsExcel[0]
		rowName := rowsExcel[1]
		rowPasword := rowsExcel[2]

		CreateCustomersMessage(models.Customer{Email: rowEmail, Name: rowName, Password: rowPasword})
		EmailCustomersMessage(models.Customer{Email: rowEmail, Name: rowName, Password: rowPasword})
	}

	fmt.Printf("%v \n", rows)

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
	retDataEmail := EmailCustomersMessage(customer)

	if retData == "success" && retDataEmail == "success" {
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

	fmt.Println(1212, os.Getenv("rabbit_url"))
	conn, err := amqp.Dial(os.Getenv("rabbit_url"))
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
