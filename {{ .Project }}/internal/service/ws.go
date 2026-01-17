package service

import (
	"net/http"

	"{{ .Computed.common_module_final }}/log"
	"github.com/gorilla/websocket"
)

const (
	wsReadBufferSize  = 1024
	wsWriteBufferSize = 1024
)

var upgrader = websocket.Upgrader{
	// Allow connections from any origin.
	CheckOrigin: func(*http.Request) bool { return true },
	ReadBufferSize:  wsReadBufferSize,
	WriteBufferSize: wsWriteBufferSize,
}

// Ws handles WebSocket connections at the /ws endpoint.
func (s *{{ .Computed.service_name_capitalized }}Service) Ws(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithError(err).Error("failed to upgrade connection to websocket")
		http.Error(w, "could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	log.Info("websocket connection established")

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.WithError(err).Warn("websocket unexpected close")
			}
			break
		}

		log.WithField("message", string(message)).Debug("received websocket message")

		if err := conn.WriteMessage(messageType, message); err != nil {
			log.WithError(err).Error("failed to write websocket message")
			break
		}
	}

	log.Info("websocket connection closed")
}

