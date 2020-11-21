package handlers

import (
	"fmt"
	"github.com/Rhymen/go-whatsapp"
	"os"
)

type TwenteHandler struct {
	Conn      *whatsapp.Conn
	StartTime uint64
}

func (TwenteHandler) HandleError(err error) {
	fmt.Fprintf(os.Stderr, "%v", err)
}

func (TwenteHandler) HandleTextMessage(message whatsapp.TextMessage) {
	fmt.Println(message)
}
