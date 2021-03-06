package main

import (
	"encoding/gob"
	"fmt"
	qrT "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/Rhymen/go-whatsapp"
	"os"
	"os/signal"
	"syscall"
	"time"
	"twentemilieu-whatsapp-bot/config"
	"twentemilieu-whatsapp-bot/handlers"
	"twentemilieu-whatsapp-bot/scheduler"
)

func main() {
	config := config.Config{}
	if err := config.SetFromBytes(); err != nil {
		fmt.Println(err)
		return
	}
	connect(&config)
}

func connect(conf *config.Config) {
	wac, err := whatsapp.NewConn(60 * time.Second)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create a connection to WhatsApp: %v\n", err)
		return
	}
	wac.SetClientVersion(2, 2021, 4)
	wac.AddHandler(&handlers.TwenteHandler{Conn: wac, StartTime: uint64(time.Now().Unix())})

	if err := login(wac); err != nil {
		fmt.Fprintf(os.Stderr, "Could not login to WhatsApp: %v\n", err)
		return
	}

	pong, err := wac.AdminTest()
	if !pong || err != nil {
		fmt.Fprintf(os.Stderr, "Cannot connect to phone: %v.\n", err)
		return
	}

	scheduler := scheduler.Scheduler{Wac: wac, Conf: conf}
	go scheduler.Watch(time.Hour * 24)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("Shutting down...")
	session, err := wac.Disconnect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to disconnect from WhatsApp: %v\n", err)
		return
	}

	if err := writeSession(session); err != nil {
		fmt.Fprintf(os.Stderr, "Could not store session: %v\n", err)
		return
	}
}

func login(conn *whatsapp.Conn) error {
	session, err := readSession()
	if err != nil {
		qr := make(chan string)
		go func() {
			term := qrT.New()
			term.Get(<-qr).Print()
		}()

		session, err = conn.Login(qr)
		if err != nil {
			return fmt.Errorf("cannot login: %v", err)
		}
	} else {
		session, err = conn.RestoreWithSession(session)
		if err != nil {
			return fmt.Errorf("cannot restore: %v", err)
		}
	}

	err = writeSession(session)
	return err
}

func readSession() (whatsapp.Session, error) {
	session := whatsapp.Session{}
	file, err := os.Open(".session.gob")
	defer file.Close()

	if err != nil {
		return session, err
	}

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&session)

	return session, err
}

func writeSession(session whatsapp.Session) error {
	file, err := os.OpenFile(".session.gob", os.O_CREATE|os.O_RDWR, 0600)
	defer file.Close()

	if err != nil {
		return err
	}

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(session)
	return err
}
