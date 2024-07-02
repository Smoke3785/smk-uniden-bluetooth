package uniden

import (
	"fmt"
	"log"

	"github.com/smoke7385/smk-uniden-bluetooth/utils"
	socket "github.com/zishang520/socket.io/v2/socket"

	"net/http"
	"strconv"
)

// SERVER
type UnidenInterfaceServer struct {
	clients []*socket.Socket
	socket  *socket.Server
	uniden  *Uniden
	port    int
}

// https://github.com/googollee/go-socket.io/tree/master/_examples
func NewServer(uniden *Uniden, port int) *UnidenInterfaceServer {
	server := socket.NewServer(nil, nil)

	uis := UnidenInterfaceServer{
		socket: server,
		uniden: uniden,
		port:   port,
	}

	return &uis
}

func (s *UnidenInterfaceServer) handleSettingsUpdate(settings *Settings) {
	fmt.Println("broadasting settings update")
	s.broadcast("settingsUpdate", s.uniden.Settings.Serialize())
}

func (s *UnidenInterfaceServer) broadcast(ev string, args ...any) {
	for _, client := range s.clients {
		client.Emit(ev, args...)
	}
}

func (s *UnidenInterfaceServer) listenForSocketEvents() {
	s.socket.On("connection", func(clients ...any) {
		client := clients[0].(*socket.Socket)
		s.clients = append(s.clients, client)

		client.Emit("settingsUpdate", s.uniden.Settings.Serialize())

		client.On("handshake", func(data ...any) {
			fmt.Println("handshake", data)
		})

		fmt.Println("connection", clients)
	})
}

func (s *UnidenInterfaceServer) start() {
	http.Handle("/", s.socket.ServeHandler(nil))

	// EVENT HANDLER
	s.listenForSocketEvents()

	portString := utils.ConcatenateStrings(":", strconv.Itoa(s.port))
	s.uniden.println("Server started on port:", portString)

	log.Fatal(http.ListenAndServe(portString, nil))
}
