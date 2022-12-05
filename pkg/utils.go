package pkg

import (
	"context"
	"hash/fnv"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/telegram"
)

func HashShortening(s []byte) uint32 {
	/*
		Simple hash func.
		!!! It is NOT a cryptographic hash-func !!!
		return: positive num
	*/
	hash := fnv.New32a()
	if _, err := hash.Write(s); err != nil {
		log.Fatalf("ERROR : %s", err)
	}
	return hash.Sum32()
}

func URLValidation(inpURL string) bool {
	/*
		URL validation.
	*/
	_, err := url.ParseRequestURI(inpURL)
	if err != nil {
		log.Println(err)
	}
	return nil == err
}

func SendMessage(msg string) {
	/*
		Send message via telegram bot. Need BOT_TOKEN and recipient id.
		param msg: message to be sent
	*/
	if err := godotenv.Load(".env"); err != nil {
		log.Println(err)
	}

	telegramService, _ := telegram.New(os.Getenv("BOT_TOKEN"))
	// Write correct telegram/chat id (var idToSend)
	var idToSend int64
	idToSend = 88005553535
	telegramService.AddReceivers(idToSend)
	notify.UseServices(telegramService)

	_ = notify.Send(
		context.Background(),
		"Shortener service",
		msg,
	)
}

func StopNotifyAdmin() {
	/*
		Notifies via telegram about the exit.
	*/
	signalCancel := make(chan os.Signal, 1)
	signal.Notify(signalCancel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for {
			s := <-signalCancel
			switch s {
			case os.Interrupt:
				fallthrough
			case syscall.SIGINT:
				fallthrough
			case syscall.SIGTERM:
				SendMessage("STOPPED")
				os.Exit(1)
			}
		}
	}()
}
