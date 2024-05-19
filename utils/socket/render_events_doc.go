package socket

import (
	"fmt"
	"html/template"
	"net/http"

	errorhandling "github.com/chirag1807/task-management-system/error"
)

type SocketEventPayload interface {
	String() string // Method to convert payload to string (optional)
}

type StringPayload struct {
	Value string
}

func (p StringPayload) String() string {
	return p.Value
}

type JsonPayload struct {
	Data map[string]interface{}
}

func (p JsonPayload) String() string {
	// b, err := json.Marshal(p.Data)
	// if err != nil {
	// 	return ""
	// }
	// return string(b)

	var output string
	for key, value := range p.Data {
		output += fmt.Sprintf("%s: %v, ", key, value)
	}
	fmt.Println(output)
	return output
}

type SocketEvent struct {
	Name        string
	Description string
	// Payload      interface{}
	Payload      SocketEventPayload
	// Payload      map[string]interface{}
	ServerAction string
}

var events []SocketEvent = []SocketEvent{
	{
		Name:        "Connect",
		Description: "This Event is used to Connect Front end with Socket Server.",
		// Payload: request.User{FirstName: "Chirag", LastName: "Makwana"},
		Payload:      StringPayload{Value: "This is Payload"},
		ServerAction: "Connect Client to the Server",
	},
	{
		Name:         "message",
		Description:  "Client sends a message",
	 Payload:      JsonPayload{Data: map[string]interface{}{"text": "String content", "text1": "String Content 1"}},
		ServerAction: "Process message content and potentially broadcast to other clients",
	},
}

func RenderSocketEventsDoc(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("F:/GOLang Projects/Task Management System/utils/socket/events_doc.html")
	if err != nil {
		errorhandling.HandleJSONUnmarshlError(r, w, err)
		return
	}
	err = tmpl.Execute(w, struct{ SocketEvents []SocketEvent }{events})
	if err != nil {
		errorhandling.HandleJSONUnmarshlError(r, w, err)
		return
	}
}
