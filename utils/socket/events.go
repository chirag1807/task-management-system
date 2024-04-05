package socket

import (
	"log"

	"github.com/chirag1807/task-management-system/config"
	socketio "github.com/googollee/go-socket.io"
)

// SocketEvents defined various events like connect, disconnect, join and leave room that client will emit.
func SocketEvents(server *socketio.Server) {
	server.OnConnect("/", func(c socketio.Conn) error {
		config.LoggerInstance.Info("Connection Made Successfully." + c.ID())
		return nil
	})

	server.OnEvent("/", "join-room", func(s socketio.Conn, roomName string) {
		config.LoggerInstance.Info("User Joined the Room: " + roomName)
		server.JoinRoom("/", roomName, s)
	})

	server.OnEvent("/", "leave-room", func(s socketio.Conn, roomName string) {
		server.LeaveRoom("/", roomName, s)
	})

	server.OnError("/", func(c socketio.Conn, err error) {
		log.Fatal(err)
	})

	server.OnDisconnect("/", func(c socketio.Conn, s string) {
		config.LoggerInstance.Info("socket disconnected:" + s)
	})
}

// EmitCreateAndUpdateTaskEvents emits create-task and update-task event to either default namespace or specific room.
func EmitCreateAndUpdateTaskEvents(server *socketio.Server, event string, room string, msg interface{}, flag int) {
	//flag = 0 => broadcast to individual via id as event and flag = 1 => broadcast to room
	if flag == 0 {
		server.BroadcastToNamespace("/", event, msg)
	} else {
		server.BroadcastToRoom("/", room, event, msg)
	}
}
