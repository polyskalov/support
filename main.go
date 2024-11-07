package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"support-server/internal"
)

// Send any text message to the bot after the bot has been started

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	dotenvError := godotenv.Load()
	if dotenvError != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("Token not set in .env file")
	}

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		panic(err)
	}

	dsn := "host=localhost user= password= dbname=support port=5432 sslmode=disable TimeZone=Asia/Yekaterinburg"
	db, dbErr := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if dbErr != nil {
		panic(err)
	}

	migrateError := db.AutoMigrate(&internal.MessageModel{})
	if migrateError != nil {
		panic(migrateError)
	}

	http.HandleFunc("/dialogs", getDialogs)
	httpErr := http.ListenAndServe(":80", nil)
	if httpErr != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	b.Start(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {

	message := internal.MessageModel{Text: update.Message.Text, ExternalId: update.Message.ID}

	db.Create(&message)
}

func getDialogs(w http.ResponseWriter, r *http.Request) {
	dialogs := &internal.Dialogs{}
	j, _ := json.Marshal(dialogs)
	_, _ = w.Write(j)
}
