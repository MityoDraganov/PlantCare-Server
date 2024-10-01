package wsutils

import (
	"PlantCare/websocket/wsTypes"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/gorilla/websocket"
)



func sendResponse(connection *wsTypes.Connection, response wsTypes.WsResponse) {
	responseBytes, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error marshalling response:", err)
		return
	}

	select {
	case connection.Send <- responseBytes:
		// Message successfully queued for sending
	default:
		// Handle the case where sending is blocked or the channel is full
		fmt.Println("Error: Unable to send response, channel may be blocked.")
	}
}

// SendOkResponse creates and sends a response with Ok = true.
func SendValidResponse(connection *wsTypes.Connection, data interface{}) {
	response := wsTypes.WsResponse{
		Ok:   true,
		Status: 200,
		Data: toJSON(data),
	}

	sendResponse(connection, response)
}

// SendErrorResponse creates and sends a response with Ok = false and an optional status message.
func SendErrorResponse(connection *wsTypes.Connection, status int) {
	response := wsTypes.WsResponse{
		Ok:     false,
		Status: status,
		Data:   nil, // No data in error response
	}

	sendResponse(connection, response)
}

func SendMessages(connection *wsTypes.Connection) {
	for msg := range connection.Send {
		err := connection.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			fmt.Println("Error while sending message:", err)
			break
		}
	}
}

// func SendMessage(connection *wsTypes.Connection, statusResponse wsTypes.StatusResponse, event wsTypes.Event, data interface{}) error {
//     // Check if at least one of statusResponse or event is provided
//     if statusResponse == "" && event == "" {
//         return fmt.Errorf("at least one of statusResponse or event must be provided")
//     }

//     // Initialize IsRead to false by default
//     isRead := false
// 	timestamp := time.Now()

//     // If data is a map and contains "IsRead", extract and remove it
// 	fmt.Println("here")
//     if mapData, ok := data.(map[any]interface{}); ok {
// 		fmt.Printf("%+v\n", mapData)
//         if val, found := mapData["IsRead"]; found {
//             if boolVal, ok := val.(bool); ok {
//                 isRead = boolVal
//             } else {
//                 return fmt.Errorf("IsRead field is not a boolean")
//             }
//             // Remove "IsRead" from the map
//             delete(mapData, "IsRead")
//         }

// 		if val, found := mapData["CreatedAt"]; found {
//             if timeVal, ok := val.(time.Time); ok {
//                 timestamp = timeVal
//             } else {
//                 return fmt.Errorf("CreatedAt field is not a time.Time")
//             }
//             // Remove "IsRead" from the map
//             delete(mapData, "IsRead")
//         }
//     }

//     // Convert the data to JSON
//     jsonData, err := json.Marshal(data)
//     if err != nil {
//         fmt.Println("Error marshaling data:", err)
//         return err
//     }

//     // Create the message
//     message := wsTypes.Message{
//         StatusResponse: &statusResponse,
//         Event:          &event,
//         Data:           json.RawMessage(jsonData),
//         Timestamp:      timestamp,
//         IsRead:         isRead,  // Use the value extracted from data or default to false
//     }

//     // Marshal the message to JSON
//     messageBytes, err := json.Marshal(message)
//     if err != nil {
//         fmt.Println("Error marshaling message:", err)
//         return err
//     }

//     // Send the message through the WebSocket connection
//     err = connection.Conn.WriteMessage(1, messageBytes)
//     if err != nil {
//         fmt.Println("Error while sending message:", err)
//         return err
//     }

//     return nil
// }


func SendMessage(connection *wsTypes.Connection, statusResponse wsTypes.StatusResponse, event wsTypes.Event, data interface{}) error {
    // Check if at least one of statusResponse or event is provided
    if statusResponse == "" && event == "" {
        return fmt.Errorf("at least one of statusResponse or event must be provided")
    }

    // Initialize default values for IsRead and Timestamp
    isRead := false
    timestamp := time.Now()

    // Use reflection to check for fields "IsRead" and "Timestamp"
    val := reflect.ValueOf(data)
    typ := reflect.TypeOf(data)

    // If data is a pointer, get the value it points to
    if val.Kind() == reflect.Ptr {
        val = val.Elem()
        typ = typ.Elem()
    }

    // Ensure the underlying value is a struct before proceeding
    if val.Kind() != reflect.Struct {
        return fmt.Errorf("data must be a struct or a pointer to a struct")
    }

    // Create a new map to hold the fields, excluding "IsRead" and "Timestamp"
    resultMap := make(map[string]interface{})

    // Iterate over the struct fields
    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)
        fieldType := typ.Field(i)
        fieldName := fieldType.Name

        // Extract the JSON tag if it exists
        jsonTag := fieldType.Tag.Get("json")
        if jsonTag == "" {
            jsonTag = fieldName // fallback to field name if no JSON tag
        }

        // Handle IsRead and Timestamp separately, exclude them from the map
        if fieldName == "IsRead" && field.Kind() == reflect.Bool {
            isRead = field.Bool()
            continue
        }

        if fieldName == "Timestamp" && field.Type() == reflect.TypeOf(time.Time{}) {
            timestamp = field.Interface().(time.Time)
            continue
        }

        // Add all other fields to the result map using the JSON tag
        if fieldType.IsExported() {
            resultMap[jsonTag] = field.Interface()
        }
    }

    // Convert the result map (with IsRead and Timestamp removed) to JSON
    jsonData, err := json.Marshal(resultMap)
    if err != nil {
        fmt.Println("Error marshaling data:", err)
        return err
    }

    // Create the message
    message := wsTypes.Message{
        StatusResponse: &statusResponse,
        Event:          &event,
        Data:           json.RawMessage(jsonData),
        Timestamp:      timestamp,
        IsRead:         isRead,
    }

    // Marshal the message to JSON
    messageBytes, err := json.Marshal(message)
    if err != nil {
        fmt.Println("Error marshaling message:", err)
        return err
    }

    // Send the message through the WebSocket connection
    err = connection.Conn.WriteMessage(websocket.TextMessage, messageBytes)
    if err != nil {
        fmt.Println("Error while sending message:", err)
        return err
    }

    return nil
}



func SendValidRequest(connection *wsTypes.Connection, data interface{}){
	response := wsTypes.WsResponse{
		Ok:     true,
		Status: 200,
		Data:   toJSON(data),
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error marshalling response:", err)
		return
	}

	// Send the response to the WebSocket client
	err = connection.Conn.WriteMessage(websocket.TextMessage, responseBytes)
	if err != nil {
		fmt.Println("Error while sending message:", err)
		return
	}
}

func toJSON(data interface{}) json.RawMessage {
	if data == nil {
		return nil
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling data:", err)
		return nil
	}
	return dataBytes
}