package scheduler

import (
	"fmt"
	"github.com/Rhymen/go-whatsapp"
	"time"
	"twentemilieu-whatsapp-bot/api"
	"twentemilieu-whatsapp-bot/config"
)

type Scheduler struct {
	Wac  *whatsapp.Conn
	Conf *config.Config
}

func (s *Scheduler) Watch(d time.Duration) {
	for {
		err := s.Fetch()
		if err != nil {
			fmt.Print(err)
		}
		time.Sleep(d)
	}
}

func (s *Scheduler) Fetch() error {
	message, ok := api.GetMessage(s.Conf)
	if ok != nil {
		return fmt.Errorf("There are no containers for today.")
	}
	SendTextMessage(s.Wac, *message, *s.Conf.Whatsapp)
	return nil
}

func SendTextMessage(wac *whatsapp.Conn, message string, jid string) {
	text := whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: jid,
		},
		Text: message,
	}
	err, _ := wac.Send(text)
	fmt.Println(err)
}
