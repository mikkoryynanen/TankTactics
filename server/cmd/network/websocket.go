package network

import (
	"fmt"
	"log"
	"main/cmd/events"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Websocket struct {
	upgrader    websocket.Upgrader
	eventEmitter *events.EventEmitter
}

func NewWebsocket(eventEmitter events.EventEmitter) *Websocket {
	return &Websocket{
		eventEmitter: &eventEmitter,
	}
}

func (ws *Websocket) Connect() {
	r := mux.NewRouter()
	r.HandleFunc("/c", ws.handleConnection)
	r.HandleFunc("/c/room", ws.handleRoomConnection)

	// TODO users disabled for now
	// r.PathPrefix("/user").Handler(routes.UserRouter())

	log.Fatal(http.ListenAndServe(":8080", r))
}

func (ws *Websocket) handleConnection(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Connection received")

	c, err := ws.upgradeConnection(w, r)
	if err != nil {
		fmt.Println(err)
	}

	ws.eventEmitter.Emit(
		events.EventOnClientConnected, 
		events.ConnectionEvent{
			Conn: c,
		})
}

func (ws *Websocket) handleRoomConnection(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Connection received. Connecting to room...")

	queryParams := r.URL.Query()
	roomId, err := uuid.Parse(queryParams.Get("roomId"))
	if err != nil {
		fmt.Println("failed to parse roomId")
		return
	}
	c, err := ws.upgradeConnection(w, r)
	if err != nil {
		fmt.Println(err)
		return
	}

	ws.eventEmitter.Emit(
		events.EventOnClientConnected, 
		events.RoomConnectionEvent{
			RoomId: roomId,
			Conn: c,
		})
}

func (ws *Websocket) upgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	c, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Could not upgrade HTTP connection to Websocket")
		return nil, err
	}
	return c, nil
}