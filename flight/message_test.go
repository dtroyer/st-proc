// st-proc: Message Processor CLI
// SPDX-License-Identifier: MIT

// Test the flight data decoding

package flight

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"gotest.tools/v3/assert"
)

// Test case 1: original example

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

var testPacket1 = []byte{
	0x41, 0x49, 0x52, 0x00, 0x00, 0x00, 0x06, 0x4E, 0x32, 0x30, 0x39, 0x30, 0x34, 0x00, 0x00, 0x00,
	0x02, 0x00, 0x00, 0x00, 0x07, 0x47, 0x45, 0x6E, 0x78, 0x2D, 0x31, 0x42, 0x40, 0x43, 0x8E, 0xD6,
	0xEB, 0xFF, 0x60, 0x1D, 0xC0, 0x50, 0xD4, 0xC0, 0x91, 0x63, 0x01, 0x65, 0x40, 0xE2, 0x03, 0xF0,
	0x00, 0x00, 0x00, 0x00, 0xC0, 0x4A, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9A,
}

// Test case 2: extended fields and UTF-8

const testJson2 = `{
  "header": "QUlS",
  "tail_number": "N⏚04£",
  "engine_count": 3,
  "engine_name": "CF6-80C2D1F",
  "latitude": 39.11593389482025,
  "longitude": -67.32425341289998,
  "altitude": 36895.5,
  "temperature": -53.2
}`

var testPacket2 = []byte{
	0x41, 0x49, 0x52, 0x00, 0x00, 0x00, 0x08, 0x4E, 0xE2, 0x8F, 0x9A, 0x30, 0x34, 0xC2, 0xA3, 0x00,
	0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x0b, 0x43, 0x46, 0x36, 0x2D, 0x38, 0x30, 0x43, 0x32, 0x44,
	0x31, 0x46, 0x40, 0x43, 0x8E, 0xD6, 0xEB, 0xFF, 0x60, 0x1D, 0xC0, 0x50, 0xD4, 0xC0, 0x91, 0x63,
	0x01, 0x65, 0x40, 0xE2, 0x03, 0xF0, 0x00, 0x00, 0x00, 0x00, 0xC0, 0x4A, 0x99, 0x99, 0x99, 0x99,
	0x99, 0x9A,
}

// Test case 3: null strings and Lakeside

const testJson3 = `{
  "header": "QUlS",
  "tail_number": "",
  "engine_count": 0,
  "engine_name": "",
  "latitude": 39.198287657959135,
  "longitude": -94.80333019810332,
  "altitude": 36895.5,
  "temperature": -270
}`

var testPacket3 = []byte{
	0x41, 0x49, 0x52, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40,
	0x43, 0x99, 0x61, 0x7D, 0x6F, 0x11, 0x45, 0xC0, 0x57, 0xB3, 0x69, 0xC3, 0x10, 0x2F, 0x8D, 0x40,
	0xE2, 0x03, 0xF0, 0x00, 0x00, 0x00, 0x00, 0xC0, 0x70, 0xE0, 0x00, 0x00, 0x00, 0x00, 0x00,
}

type CaseTable struct {
	expected string
	bindata  []byte
}

// This test decodes the above JSON and binary data testing the JSON tags in
// the FlightMessage struct and DecodePacketBuffer()
func TestDecodePacketBuffer(t *testing.T) {
	data := []CaseTable{
		{testJson1, testPacket1},
		{testJson2, testPacket2},
		{testJson3, testPacket3},
	}

	for _, c := range data {
		// Decode the 'expected' values from the JSON strings
		var expect FlightMessage
		err := json.Unmarshal([]byte(c.expected), &expect)
		if err != nil {
			fmt.Println("json error: ", err)
		}

		// Decode the test packets
		buf := FlightMessage{}
		err = DecodePacketBuffer(bytes.NewBuffer(c.bindata), &buf)
		if err != nil {
			fmt.Println("error: ", err)
		}

		// See how it went
		assert.Equal(t, FlightMessageHeader, buf.Header)
		assert.Equal(t, expect.TailNumber, buf.TailNumber)
		assert.Equal(t, expect.EngineCount, buf.EngineCount)
		assert.Equal(t, expect.EngineName, buf.EngineName)
		assert.Equal(t, expect.Latitude, buf.Latitude)
		assert.Equal(t, expect.Longitude, buf.Longitude)
		assert.Equal(t, expect.Altitude, buf.Altitude)
		assert.Equal(t, expect.Temperature, buf.Temperature)
	}
}
