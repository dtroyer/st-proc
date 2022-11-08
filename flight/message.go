// st-proc: Message Processor CLI
// SPDX-License-Identifier: MIT

// Message handling for flight data
// * FlightMessage Go struct representing the message packat, with JSON tags
// * Decode on-the-wire packet into FlightMessage
// * JSON Marshal and Unmarshal support for the 3 byte message header

package flight

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"reflect"
	"strings"
)

// Define a specific type for the header to simplify decoding
type MessageHeader [3]byte

// The MessageHeader is base64 encoded in JSON

// Encode a MessageHeader to base64
func (b *MessageHeader) MarshalJSON() ([]byte, error) {
	str := base64.StdEncoding.EncodeToString((*b)[0:3])
	return []byte(str), nil
}

// Decode a base64-encoded string into MessageHeader
func (b *MessageHeader) UnmarshalJSON(input []byte) error {
	dst, err := base64.StdEncoding.DecodeString(string(input[1 : len(input)-1]))
	if err != nil {
		return err
	}
	*b = MessageHeader{dst[0], dst[1], dst[2]}
	return nil
}

// The flight data struct
type FlightMessage struct {
	Header      MessageHeader `json:"header"`
	TailNumber  string        `json:"tail_number"`
	EngineCount uint32        `json:"engine_count"`
	EngineName  string        `json:"engine_name"`
	Latitude    float64       `json:"latitude"`
	Longitude   float64       `json:"longitude"`
	Altitude    float64       `json:"altitude"`
	Temperature float64       `json:"temperature"`
}

// The MessageHeader we are concerned with
var FlightMessageHeader = MessageHeader{'A', 'I', 'R'}

// Decode the binary message into a FlightMessage struct
func DecodePacketBuffer(buffer *bytes.Buffer, data interface{}) error {
	// Validate passed destination buffer is a pointer and FlightMessage struct
	dataType := reflect.TypeOf(data).Kind()
	if dataType != reflect.Ptr {
		return errors.New("data is not a pointer")
	}

	dataValue := reflect.ValueOf(data).Elem()
	if dataValue.Type().Name() != "FlightMessage" {
		return errors.New("data is not a FlightMessage struct")
	}

	// Loop through the FlightMessage struct fields and extract bytes from the buffer
	for i := 0; i < dataValue.NumField(); i++ {
		switch dataValue.Field(i).Type().Kind() {
		case reflect.Array:
			// Look for our MessageHeader type
			switch dataValue.Field(i).Type().Name() {
			case "MessageHeader":
				err := binary.Read(buffer, binary.BigEndian, dataValue.Field(i).Addr().Interface())
				if err != nil {
					return err
				}
				if dataValue.Field(i) == reflect.ValueOf(FlightMessageHeader) {
					return errors.New("Invalid message header")
				}
				break
			}
			break
		case reflect.Float64:
			err := binary.Read(buffer, binary.BigEndian, dataValue.Field(i).Addr().Interface())
			if err != nil {
				return err
			}
			break
		case reflect.String:
			// Get 4 byte string size
			var strsize uint32
			err := binary.Read(buffer, binary.BigEndian, &strsize)
			if err != nil {
				return err
			}
			// Get string value
			dataValue.Field(i).SetString(
				strings.ToValidUTF8(string(buffer.Next(int(strsize))), ""),
			)
			break
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			err := binary.Read(buffer, binary.BigEndian, dataValue.Field(i).Addr().Interface())
			if err != nil {
				return err
			}
			break
		default:
			return errors.New("unsupported type: " + dataValue.Field(i).Type().Name())
		}
	}
	return nil
}
