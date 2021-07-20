// +build eebus

package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/amp-x/eebus"
	"github.com/amp-x/eebus/cert"
	"github.com/amp-x/eebus/communication"
	"github.com/amp-x/eebus/mdns"
	"github.com/amp-x/eebus/server"
	"github.com/amp-x/eebus/ship"
	"github.com/amp-x/eebus/spine/model"
	"github.com/andig/evcc/util"
	"github.com/grandcat/zeroconf"
)

type EEBusClients struct {
	onConnect    func(string, ship.Conn) error
	onDisconnect func(string)
}

type EEBus struct {
	mux               sync.Mutex
	log               *util.Logger
	srv               *eebus.Server
	id                string
	clients           map[string]EEBusClients
	connectedClients  map[string]ship.Conn
	discoveredClients map[string]*zeroconf.ServiceEntry
}

var EEBusInstance *EEBus

func NewEEBus(other map[string]interface{}) (*EEBus, error) {
	cc := struct {
		Uri         string
		ShipID      string
		Certificate struct {
			Public, Private []byte
		}
	}{
		Uri: ":4712",
	}

	if err := util.DecodeOther(other, &cc); err != nil {
		return nil, err
	}

	// if !sponsor.IsAuthorized() {
	// 	return nil, errors.New("eebus requires evcc sponsorship, register at https://cloud.evcc.io")
	// }

	details := EEBusInstance.DeviceInfo()

	log := util.NewLogger("eebus")
	id := server.UniqueID{Prefix: details.BrandName}.String()
	if len(cc.ShipID) > 0 {
		id = cc.ShipID
	}

	cert, err := tls.X509KeyPair(cc.Certificate.Public, cc.Certificate.Private)
	if err != nil {
		return nil, err
	}

	srv := &eebus.Server{
		Log:         log.TRACE,
		Addr:        cc.Uri,
		Path:        "/ship/",
		Certificate: cert,
		ID:          id,
		Brand:       details.BrandName,
		Model:       details.DeviceCode,
		Type:        string(model.DeviceTypeEnumTypeEnergyManagementSystem),
		Register:    true,
	}

	if _, err = srv.Announce(); err != nil {
		return nil, err
	}

	c := &EEBus{
		log:               log,
		srv:               srv,
		id:                id,
		clients:           make(map[string]EEBusClients),
		connectedClients:  make(map[string]ship.Conn),
		discoveredClients: make(map[string]*zeroconf.ServiceEntry),
	}

	return c, nil
}

func (c *EEBus) DeviceInfo() communication.ManufacturerDetails {
	return communication.ManufacturerDetails{
		BrandName:     "EVCC",
		DeviceName:    "EVCC",
		DeviceCode:    "EVCC_HEMS_01",
		DeviceAddress: "EVCC_HEMS",
	}
}

func (c *EEBus) Register(ski string, shipConnectHandler func(string, ship.Conn) error, shipDisconnectHandler func(string)) {
	ski = strings.ReplaceAll(ski, "-", "")
	c.log.TRACE.Printf("registering ski: %s", ski)

	c.mux.Lock()
	c.clients[ski] = EEBusClients{onConnect: shipConnectHandler, onDisconnect: shipDisconnectHandler}
	c.mux.Unlock()

	// maybe the SKI is already discovered
	c.handleDiscoveredSKI(ski)
}

func (c *EEBus) Run() {
	entries := make(chan *zeroconf.ServiceEntry)
	go c.discoverDNS(entries, func(entry *zeroconf.ServiceEntry) {
		c.addDisoveredEntry(entry)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// discover all services on the network (e.g. _workstation._tcp)
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		panic(fmt.Errorf("mDNS: failed initializing resolver: %w", err))
	}

	if err = resolver.Browse(ctx, ship.ZeroconfType, ship.ZeroconfDomain, entries); err != nil {
		panic(fmt.Errorf("failed to browse: %w", err))
	}

	ln := &server.Listener{
		Log:          c.log.TRACE,
		AccessMethod: c.id,
		Handler:      c.shipHandler,
	}

	if err := c.srv.Listen(ln, c.certificateHandler); err != nil {
		c.log.ERROR.Println("eebus listen:", err)
	}
}

func (c *EEBus) addDisoveredEntry(entry *zeroconf.ServiceEntry) {
	// we need to get the SKI only
	svc, err := mdns.NewFromDNSEntry(entry)

	if err == nil {
		c.mux.Lock()
		c.discoveredClients[svc.SKI] = entry
		c.mux.Unlock()

		// maybe the SKI is already registered
		c.handleDiscoveredSKI(svc.SKI)
	} else {
		c.log.TRACE.Printf("%s: could not create ship service from DNS entry: %v", entry.HostName, err)
	}
}

func (c *EEBus) handleDiscoveredSKI(ski string) {
	c.mux.Lock()

	_, connected := c.connectedClients[ski]
	_, registered := c.clients[ski]
	entry, discovered := c.discoveredClients[ski]

	c.log.TRACE.Printf("client %s connected %t, registered %t, discovered %t ", ski, connected, registered, discovered)

	if !connected && discovered && registered {
		c.mux.Unlock()
		c.connectDiscoveredEntry(entry)
		return
	}

	c.mux.Unlock()
}

func (c *EEBus) connectDiscoveredEntry(entry *zeroconf.ServiceEntry) {
	svc, err := mdns.NewFromDNSEntry(entry)

	var conn ship.Conn
	if err == nil {
		c.log.TRACE.Printf("%s: client connect", entry.HostName)
		conn, err = svc.Connect(c.log.TRACE, c.id, c.srv.Certificate, c.shipCloseHandler)
	}

	if err != nil {
		c.log.TRACE.Printf("%s: client done: %v", entry.HostName, err)
		return
	}

	err = c.shipHandler(svc.SKI, conn)
	if err != nil {
		log.FATAL.Fatalf("%s: error calling shipHandler: %v", entry.HostName, err)
		return
	}
}

func (c *EEBus) discoverDNS(results <-chan *zeroconf.ServiceEntry, connector func(*zeroconf.ServiceEntry)) {
	for entry := range results {
		c.log.TRACE.Println("mDNS:", entry.HostName, entry.AddrIPv4, entry.Text)

		for _, typ := range entry.Text {
			if strings.HasPrefix(typ, "type=") && typ == "type=EVSE" {
				connector(entry)
			}
		}
	}
}

func (c *EEBus) certificateHandler(leaf *x509.Certificate) error {
	ski, err := cert.SkiFromX509(leaf)
	if err != nil {
		return err
	}

	c.log.TRACE.Printf("verifying client ski: %s", ski)

	c.mux.Lock()
	defer c.mux.Unlock()

	for client := range c.clients {
		if client == ski {
			return nil
		}
	}

	return fmt.Errorf("client ski not allowed: %s", ski)
}

func (c *EEBus) shipHandler(ski string, conn ship.Conn) error {
	c.mux.Lock()

	for client, cb := range c.clients {
		if client == ski {
			currentConnection, found := c.connectedClients[ski]
			connect := true
			c.log.TRACE.Printf("client %s found? %t", ski, found)
			if found {
				if currentConnection.IsConnectionClosed() {
					c.log.TRACE.Printf("client has closed connection")
					delete(c.connectedClients, ski)
				} else {
					c.log.TRACE.Printf("client has no closed connection")
					connect = false
				}
			}
			c.log.TRACE.Printf("client %s connect? %t", ski, connect)
			if connect {
				c.connectedClients[ski] = conn
				c.mux.Unlock()
				err := cb.onConnect(ski, conn)
				if err != nil {
					c.mux.Lock()
					delete(c.connectedClients, ski)
					c.mux.Unlock()
				}
				return err
			}
		}
	}

	c.mux.Unlock()

	return errors.New("client not registered")
}

// handles connection closed
func (c *EEBus) shipCloseHandler(ski string) {
	c.mux.Lock()

	_, found := c.connectedClients[ski]
	if found {
		for client, cb := range c.clients {
			if client == ski {
				cb.onDisconnect(ski)
				break
			}
		}
		delete(c.connectedClients, ski)
	}

	c.mux.Unlock()
}
