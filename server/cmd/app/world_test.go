package app

import (
	"fmt"
	"log"
	"main/cmd/types"
	"main/cmd/utils"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"

	"testing"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func testHandler(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot upgrade: %v", err), http.StatusInternalServerError)
	}
	mt, p, err := conn.ReadMessage()
	if err != nil {
		log.Printf("cannot read message: %v", err)
		return
	}
	conn.WriteMessage(mt, []byte("hello "+string(p)))
}

func TestWorld(t *testing.T) {
	root, _ := utils.FindProjectRoot()
	root = filepath.Join(root, "data", "level.json")
	t.Setenv("LEVEL_FILE_PATH", root)

	server := httptest.NewServer(http.HandlerFunc(testHandler))
	u, _ := url.Parse(server.URL)
	u.Scheme = "ws"
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Error("Websocket connection failed")
	}

	world := NewWorld()
	newClient := types.NewClient(conn)
	newClient.Input.InputX = 1
	world.Clients["test_client"] = newClient
	world.SimulateOnce()

	if world.Clients["test_client"].Object.Position != world.level.Objects[0].Position {
		t.Error("Player should have collided with object")
	}
}
