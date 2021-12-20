package sockets

import (
	"fmt"

	socketio "github.com/googollee/go-socket.io"
)

var ServerVisita *socketio.Server

func gestionarServerVisitas() {
	ServerVisita = socketio.NewServer(nil)
	ServerVisita.OnConnect("/", func(s socketio.Conn) error {
		print("Conexion exitosa")
		s.Emit("conexion", "Conexion exitosa.")
		return nil
	})
	ServerVisita.OnEvent("/", "entrar", func(s socketio.Conn, msg string) error {
		s.Join(msg)
		return nil
	})
	ServerVisita.OnEvent("/", "contestar", func(s socketio.Conn, msg string) {
		roomId := s.Rooms()[0]
		if roomId != "" {
			ServerVisita.BroadcastToRoom("/", roomId, "respuesta", func() string {
				return msg
			})
		} else {
			s.Emit("error", "No esta en observando ninguna solicitud")
		}
	})
	ServerVisita.OnDisconnect("/", func(c socketio.Conn, s string) {
		fmt.Printf("Conexion perdido con %s", c.ID())
	})

}
