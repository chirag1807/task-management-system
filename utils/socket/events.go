package socket

import (
	"fmt"
	"log"

	socketio "github.com/googollee/go-socket.io"
)

func SocketEvents(server *socketio.Server) {
	server.OnConnect("/", func(c socketio.Conn) error {
		log.Println("Connection Made Successfully.", c.ID())
		return nil
	})

	server.OnEvent("/", "join-room", func(s socketio.Conn, roomName string) {
		fmt.Println(roomName)
		server.JoinRoom("/", roomName, s)
	})

	server.OnEvent("/", "leave-room", func(s socketio.Conn, roomName string) {
		server.LeaveRoom("/", roomName, s)
	})

	server.OnError("/", func(c socketio.Conn, err error) {
		log.Fatal(err)
	})

	server.OnDisconnect("/", func(c socketio.Conn, s string) {
		log.Println("socket disconnected:", s)
	})
}

func EmitCreateAndUpdateTaskEvents(server *socketio.Server, event string, room string, msg interface{}, flag int) {
	//flag = 0 => broadcast to individual via id as event and flag = 1 => broadcast to room
	if flag == 0 {
		server.BroadcastToNamespace("/", event, msg)
	} else {
		server.BroadcastToRoom("/", room, event, msg)
	}
}
