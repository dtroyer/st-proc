// st-proc: Message Processor CLI
// SPDX-License-Identifier: MIT

package flight

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"gotest.tools/v3/assert"
)

const testJson1 = `{
  "header": "QUlS",
  "tail_number": "N20904",
  "engine_count": 2,
  "engine_name": "GEnx-1B",
  "latitude": 39.11593389482025,
  "longitude": -67.32425341289998,
  "altitude": 36895.5,
  "temperature": -53.2
}`

var testPkt1 = []byte{
	0x41, 0x49, 0x52, 0x00, 0x00, 0x00, 0x06, 0x4E, 0x32, 0x30, 0x39, 0x30, 0x34, 0x00, 0x00, 0x00,
	0x02, 0x00, 0x00, 0x00, 0x07, 0x47, 0x45, 0x6E, 0x78, 0x2D, 0x31, 0x42, 0x40, 0x43, 0x8E, 0xD6,
	0xEB, 0xFF, 0x60, 0x1D, 0xC0, 0x50, 0xD4, 0xC0, 0x91, 0x63, 0x01, 0x65, 0x40, 0xE2, 0x03, 0xF0,
	0x00, 0x00, 0x00, 0x00, 0xC0, 0x4A, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9A,
}

func TestReadPktBuf(t *testing.T) {
	var expect FlightMessage
	err := json.Unmarshal([]byte(testJson1), &expect)
	// fmt.Printf("tn: %s\n", expect.TailNumber)
	if err != nil {
		fmt.Println("json error: ", err)
	}
	buf := FlightMessage{}

	err = DecodePacketBuffer(bytes.NewBuffer(testPkt1), &buf)
	if err != nil {
		fmt.Println("error: ", err)
	}
	assert.Equal(t, FlightMessageHeader, buf.Header)
	assert.Equal(t, "N20904", buf.TailNumber)
	assert.Equal(t, uint32(2), buf.EngineCount)
	assert.Equal(t, "GEnx-1B", buf.EngineName)
	assert.Equal(t, 39.11593389482025, buf.Latitude)
	assert.Equal(t, -67.32425341289998, buf.Longitude)
	assert.Equal(t, 36895.5, buf.Altitude)
	assert.Equal(t, -53.2, buf.Temperature)
}
