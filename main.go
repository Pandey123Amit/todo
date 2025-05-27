package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

type EmailData struct {
	Name     string
	Subject  string
	Body     string
	Prevtask string
	Newtask  string
	Id       int
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	fmt.Println("Inside Main in todo app")

	todos := Todos{}
	storage := NewStorage[Todos]("todos.json")
	storage.Load(&todos)
	cmdFlags := NewCmdFlags()
	cmdFlags.Execute(&todos)
	// todos.print()
	storage.Save(todos)
	// Cron Job
	c := cron.New()
	c.AddFunc("@every 1m", func() {
		fmt.Println("⏰ Cron job running at:", time.Now())

		data := EmailData{
			Name:    "Amit",
			Subject: "Welcome to Todo App",
			Body:    "Thank you for signing up for our todo application!",
		}

		tmpl, err := template.ParseFiles("email_template.html")
		if err != nil {
			log.Fatal("Error parsing template:", err)
		}

		var body bytes.Buffer
		if err := tmpl.Execute(&body, data); err != nil {
			log.Fatal("Error executing template:", err)
		}

		err = sendEmail("emailcheck@yopmail.com", data.Subject, body.String())
		if err != nil {
			log.Fatalf("Failed to send email: %v", err)
		} else {
			fmt.Println("Email sent successfully!")
		}
	})
	c.Start()

	select {}
	c.Stop()

}

func sendEmail(to, subject, body string) error {
	from := os.Getenv("EMAIL_USER")
	fmt.Println(from)
	password := os.Getenv("EMAIL_PASSWORD")

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	msg := "From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n" +
		body

	auth := smtp.PlainAuth("", from, password, smtpHost)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
}
