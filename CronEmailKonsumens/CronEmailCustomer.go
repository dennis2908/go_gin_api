package main

import (
	"go-restapi-gin/models"

	"fmt"
	"log"
	"net/smtp"
	"runtime"
	"sync"
	"time"

	cron "github.com/robfig/cron/v3"
)

func main() {
	numberOfCores := runtime.NumCPU()
	fmt.Println(numberOfCores)
	runtime.GOMAXPROCS(numberOfCores)
	var wg sync.WaitGroup
	for i := 0; i < numberOfCores; i++ {
		wg.Add(1)

		config, err := models.LoadConfig("../")
		if err != nil {
			log.Fatal("? Could not load environment variables", err)
		}
		models.ConnectDB(&config)

		GetData()

	}
	wg.Wait()

}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func GetData() {
	fmt.Println("Email notif")
	jakartaTime, _ := time.LoadLocation("Asia/Jakarta")
	scheduler := cron.New(cron.WithLocation(jakartaTime))

	// stop scheduler tepat sebelum fungsi berakhir
	scheduler.Start()

	// set task yang akan dijalankan scheduler
	// gunakan crontab string untuk notifikasi harian tiap jam 9 malam ke email masing masing konsumen
	scheduler.AddFunc("*/1 * * * *", NotifyDailyNightNotif)
	time.Sleep(1 * time.Minute)
}

func Email(Name string, Email string) {

	if Name == "" {
		Name = "Konsumen"
	}
	username := "217b4a405ee6bf"
	password := "06a38497e6b5a3"
	host := "sandbox.smtp.mailtrap.io"
	port := "465"

	// Subject and body
	subject := "Good Night, " + Name
	body := "Hi, Good Night " + Name

	// Sender and receiver
	from := "michaeldenniseldima@gmail.com"
	to := []string{
		Email,
	}

	// Build the message
	message := fmt.Sprintf("From: %s\r\n", from)
	message += fmt.Sprintf("To: %s\r\n", to)
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += fmt.Sprintf("\r\n%s\r\n", body)

	// Authentication.
	auth := smtp.PlainAuth("", username, password, host)

	// Send email
	err := smtp.SendMail(host+":"+port, auth, from, to, []byte(message))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Email sent successfully.")
}

func NotifyDailyNightNotif() {

	fmt.Println("Email notif")
	var customers []models.Customer

	models.DB.Find(&customers)

	for _, data := range customers {
		fmt.Println(data.Email)
		if data.Email != "" {
			Email(data.Name, data.Email)
		}
	}
}
