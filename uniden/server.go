package uniden

import (
	"fmt"
	"log"

	socketio "github.com/googollee/go-socket.io"
	"github.com/smoke7385/smk-uniden-bluetooth/utils"

	"net/http"
	"strconv"
)

// SERVER
type UnidenInterfaceServer struct {
	socket *socketio.Server
	uniden *Uniden
	port   int
}

// https://github.com/googollee/go-socket.io/tree/master/_examples
func NewServer(uniden *Uniden, port int) *UnidenInterfaceServer {
	server := socketio.NewServer(nil)

	uis := UnidenInterfaceServer{
		socket: server,
		uniden: uniden,
		port:   port,
	}

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})

	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		return "recv " + msg
	})

	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		// server.Remove(s.ID())
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		// Add the Remove session id. Fixed the connection & mem leak
		// server.Remove(s.ID())
		fmt.Println("closed", reason)
	})

	return &uis
}

func (s *UnidenInterfaceServer) handleSettingsUpdate(settings *Settings) {
	s.socket.BroadcastToNamespace("/", "settingsUpdate", settings)
}

func (s *UnidenInterfaceServer) start() {
	go s.socket.Serve()
	defer s.socket.Close()

	http.Handle("/", s.socket)

	portString := utils.ConcatenateStrings(":", strconv.Itoa(s.port))
	s.uniden.println("Server started on port:", portString)

	log.Fatal(http.ListenAndServe(portString, nil))
}
