package uniden

import (
	"fmt"
	"log"

	"github.com/smoke7385/smk-uniden-bluetooth/utils"
	socket "github.com/zishang520/socket.io/v2/socket"

	"net/http"
	"strconv"
)

// Client wraps a socket.io socket with a name. Socket contains atomic state,
// so it must be referenced by pointer, never copied.
type Client struct {
	*socket.Socket
	Name string
}

// SERVER
type UnidenInterfaceServer struct {
	socket  *socket.Server
	clients []*Client
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

func (s *UnidenInterfaceServer) handleSettingsUpdate(settings *Settings, full bool) {
	fmt.Println("broadasting limited settings update")

	var bSettings *Settings

	if full {
		bSettings = &s.uniden.Settings
	} else {
		bSettings = settings
	}

	s.broadcast("settingsUpdate", bSettings.Serialize())

}

func (s *UnidenInterfaceServer) broadcastNewSettings() {
	fmt.Println("broadasting full settings update")
	s.broadcast("settingsUpdate", s.uniden.Settings.Serialize())
}

func (s *UnidenInterfaceServer) handleStatusUpdate(status *Status) {
	fmt.Println("broadasting status update")
	s.broadcast("statusUpdate", s.uniden.Status.Serialize())
}

func (s *UnidenInterfaceServer) handleRadarEvent(events *RadarEvents) {
	fmt.Println("broadasting radar events update")
	s.broadcast("radarEventUpdate", events.Serialize())
}

func (s *UnidenInterfaceServer) broadcast(ev string, args ...any) {
	for _, client := range s.clients {
		client.Emit(ev, args...)
	}
}

func (s *UnidenInterfaceServer) listenForSocketEvents() {
	s.socket.On("connection", func(clients ...any) {
		client := &Client{Socket: clients[0].(*socket.Socket)}
		s.clients = append(s.clients, client)

		// Send initial settings state
		client.Emit("settingsUpdate", s.uniden.Settings.Serialize())

		client.On("registerClient", func(data ...any) {
			client.Name = data[0].(string)
			fmt.Println("registerClient", client.Name)
		})

		client.On("updateSetting", func(data ...any) {
			requestUpdateSetting := NewRequestUpdateSetting(data...)
			setting, err := s.uniden.Settings.getByDeviceStorageIndex(requestUpdateSetting.DeviceStorageIndex)

			if err != nil {
				fmt.Println("Failed to find setting as requested by client.")
				return
			}

			setting.Update(requestUpdateSetting.ValueInt)
			fmt.Println("Updated setting upon client request")
		})

		client.On("mute", func(data ...any) {
			err := s.uniden.Mute()
			if err != nil {
				fmt.Println("Failed to mute device upon client request")
			} else {
				fmt.Println("Muted device upon client request")
			}

			s.broadcastNewSettings()
		})

		client.On("unmute", func(data ...any) {
			err := s.uniden.Unmute()
			if err != nil {
				fmt.Println("Failed to unmute device upon client request")
			} else {
				fmt.Println("Unmuted device upon client request")
			}

			s.broadcastNewSettings()
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
