package bot

import (
	"fmt"
	"github.com/Rhymen/go-whatsapp"
	"os"
)

type twenteHandler struct {
	Conn      *whatsapp.Conn
	StartTime uint64
}

func (twenteHandler) HandleError(err error) {
	fmt.Fprintf(os.Stderr, "%v", err)
}

func (twenteHandler) HandleTextMessage(message whatsapp.TextMessage) {
	fmt.Println(message)
}
