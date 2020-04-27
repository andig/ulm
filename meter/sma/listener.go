package sma

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync"

	"github.com/andig/evcc/api"
)

const (
	multicastAddr = "239.12.255.254:9522"
	udpBufferSize = 8192

	msgPreamble   = 28 // preamble size in bytes
	msgCodeLength = 4  // length in bytes

	ObisImportPower  = "1:1.4.0" // Wirkleistung (W)
	ObisImportEnergy = "1:1.8.0" // Wirkarbeit (Ws) +
	ObisExportPower  = "1:2.4.0" // Wirkleistung (W)
	ObisExportEnergy = "1:2.8.0" // Wirkarbeit (Ws) −
)

// obisDefinition defines the properties needed to parse the SMA multicast telegram values
type obisDefinition struct {
	length int     // data size in bytes of the return value
	factor float64 // the factor to multiply the value by to get the proper value in the given unit
}

// list of Obis codes and their properties as defined in the SMA EMETER-Protokoll-TI-de-10.pdf document
var knownObisCodes = map[string]obisDefinition{
	// Overall sums
	ObisImportPower: {4, 0.1}, ObisImportEnergy: {8, 1}, // Wirkleistung (W)/-arbeit (Ws) +
	ObisExportPower: {4, 0.1}, ObisExportEnergy: {8, 1}, // Wirkleistung (W)/-arbeit (Ws) −
	"1:3.4.0": {4, 0.1}, "1:3.8.0": {8, 1}, // Blindleistung (W)/-arbeit (Ws) +
	"1:4.4.0": {4, 0.1}, "1:4.8.0": {8, 1}, // Blindleistung (W)/-arbeit (Ws) −
	"1:9.4.0": {4, 0.1}, "1:9.8.0": {8, 1}, // Scheinleistung (W)/-arbeit (Ws) +
	"1:10.4.0": {4, 0.1}, "1:10.8.0": {8, 1}, // Scheinleistung (W)/-arbeit (Ws) −
	"1:13.4.0": {4, 0.001}, // Leistungsfaktor (Φ)
	// Phase 1: {
	"1:21.4.0": {4, 0.1}, "1:21.8.0": {8, 1}, // Wirkleistung (W)/-arbeit (Ws) +
	"1:22.4.0": {4, 0.1}, "1:22.8.0": {8, 1}, // Wirkleistung (W)/-arbeit (Ws) −
	"1:23.4.0": {4, 0.1}, "1:23.8.0": {8, 1}, // Blindleistung (W)/-arbeit (Ws) +
	"1:24.4.0": {4, 0.1}, "1:24.8.0": {8, 1}, // Blindleistung (W)/-arbeit (Ws) −
	"1:29.4.0": {4, 0.1}, "1:29.8.0": {8, 1}, // Scheinleistung (W)/-arbeit (Ws) +
	"1:30.4.0": {4, 0.1}, "1:30.8.0": {8, 1}, // Scheinleistung (W)/-arbeit (Ws) −
	"1:31.4.0": {4, 0.001}, // Strom (A)
	"1:32.4.0": {4, 0.001}, // Spannung (V
	// Phase 2: {
	"1:41.4.0": {4, 0.1}, "1:41.8.0": {8, 1}, // Wirkleistung (W)/-arbeit (Ws) +
	"1:42.4.0": {4, 0.1}, "1:42.8.0": {8, 1}, // Wirkleistung (W)/-arbeit (Ws) −
	"1:43.4.0": {4, 0.1}, "1:43.8.0": {8, 1}, // Blindleistung (W)/-arbeit (Ws) +
	"1:44.4.0": {4, 0.1}, "1:44.8.0": {8, 1}, // Blindleistung (W)/-arbeit (Ws) −
	"1:49.4.0": {4, 0.1}, "1:49.8.0": {8, 1}, // Scheinleistung (W)/-arbeit (Ws) +
	"1:50.4.0": {4, 0.1}, "1:50.8.0": {8, 1}, // Scheinleistung (W)/-arbeit (Ws) −
	"1:51.4.0": {4, 0.001}, // Strom (A)
	"1:52.4.0": {4, 0.001}, // Spannung (V)
	// Phase 3: {
	"1:61.4.0": {4, 0.1}, "1:61.8.0": {8, 1}, // Wirkleistung (W)/-arbeit (Ws) +
	"1:62.4.0": {4, 0.1}, "1:62.8.0": {8, 1}, // Wirkleistung (W)/-arbeit (Ws) −
	"1:63.4.0": {4, 0.1}, "1:63.8.0": {8, 1}, // Blindleistung (W)/-arbeit (Ws) +
	"1:64.4.0": {4, 0.1}, "1:64.8.0": {8, 1}, // Blindleistung (W)/-arbeit (Ws) −
	"1:69.4.0": {4, 0.1}, "1:69.8.0": {8, 1}, // Scheinleistung (W)/-arbeit (Ws) +
	"1:70.4.0": {4, 0.1}, "1:70.8.0": {8, 1}, // Scheinleistung (W)/-arbeit (Ws) −
	"1:71.4.0": {4, 0.001}, // Strom (A)
	"1:72.4.0": {4, 0.001}, // Spannung (V)
	// Others
	"144:0.0.0": {4, 1}, // SW Version
}

// Instance is the Listener singleton
var Instance *Listener

// Telegram defines the data structure of a SMA multicast data package
type Telegram struct {
	Addr   string
	Serial string
	Data   map[string]float64
}

// Listener for receiving SMA multicast data packages
type Listener struct {
	mux     sync.Mutex
	log     *api.Logger
	conn    *net.UDPConn
	clients map[string]chan<- Telegram
}

// New creates a Listener
func New(log *api.Logger, addr string) *Listener {
	// Parse the string address
	gaddr, err := net.ResolveUDPAddr("udp4", multicastAddr)
	if err != nil {
		log.FATAL.Fatalf("error resolving udp address: %s", err)
	}

	// Open up a connection
	conn, err := net.ListenMulticastUDP("udp4", nil, gaddr)
	if err != nil {
		log.FATAL.Fatalf("error opening connecting: %s", err)
	}

	if err := conn.SetReadBuffer(udpBufferSize); err != nil {
		log.FATAL.Fatalf("error setting read buffer: %s", err)
	}

	l := &Listener{
		log:  log,
		conn: conn,
	}

	go l.listen()

	return l
}

// processMessage converts a SMA multicast data package into Telegram
func (l *Listener) processMessage(src *net.UDPAddr, b []byte) (Telegram, error) {
	numBytes := len(b)

	if numBytes <= msgPreamble {
		return Telegram{}, errors.New("received data package is too small")
	}

	obisValues := make(map[string]float64)

	var obisDef obisDefinition
	for i := msgPreamble; i < numBytes-msgCodeLength; i = i + msgCodeLength + obisDef.length {
		// spec says value should be 1, but reading contains 0
		b0 := b[i+0]
		if b0 == 0 {
			b0 = 1
		}

		code := fmt.Sprintf("%d:%d.%d.%d", b0, b[i+1], b[i+2], b[i+3])
		if obisDef, ok := knownObisCodes[code]; ok {
			switch obisDef.length {
			case 4:
				obisValues[code] = obisDef.factor * float64(binary.BigEndian.Uint32(b[i+msgCodeLength:]))
			case 8:
				obisValues[code] = obisDef.factor * float64(binary.BigEndian.Uint64(b[i+msgCodeLength:]))
			}
		}
	}

	serial := strconv.FormatUint(uint64(binary.BigEndian.Uint32(b[20:24])), 10)

	msg := Telegram{
		Addr:   src.IP.String(),
		Serial: serial,
		Data:   obisValues,
	}

	return msg, nil
}

// listen for multicast data packages
func (l *Listener) listen() {
	buffer := make([]byte, udpBufferSize)

	for {
		read, src, err := l.conn.ReadFromUDP(buffer)
		if err != nil {
			l.log.WARN.Printf("udp read failed: %s", err)
			continue
		}

		if msg, err := l.processMessage(src, buffer[:read-1]); err == nil {
			l.send(msg)
		}
	}
}

// Subscribe adds a client address and message channel
func (l *Listener) Subscribe(addr string, c chan<- Telegram) {
	l.mux.Lock()
	defer l.mux.Unlock()

	if l.clients == nil {
		l.clients = make(map[string]chan<- Telegram)
	}

	l.clients[addr] = c
}

func (l *Listener) send(msg Telegram) {
	l.mux.Lock()
	defer l.mux.Unlock()

	for addr, client := range l.clients {
		if addr == msg.Addr {
			select {
			case client <- msg:
			default:
				l.log.TRACE.Println("listener: recv blocked")
			}
			break
		}
	}
}
