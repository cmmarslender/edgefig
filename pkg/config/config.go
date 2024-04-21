package config

import (
	"fmt"
	"net/netip"

	"github.com/cmmarslender/edgefig/pkg/types"
)

// Config is the top level config container
type Config struct {
	Routers []Router `yaml:"routers"`
	VLANs   []VLAN   `yaml:"vlans"`
}

// Connection common details for connecting to devices
type Connection struct {
	IP       netip.Addr `yaml:"ip"`
	Port     uint16     `yaml:"port"`
	Username string     `yaml:"username"`
	Password string     `yaml:"password"`
}

// User defines a common struct that represents a user across routers, switches, etc
type User struct {
	Username string          `yaml:"username"`
	Password string          `yaml:"password"`
	Role     types.UserLevel `yaml:"role"`
}

// Router is the top level config for a single router
type Router struct {
	Name string `yaml:"name"`
	Connection
	Interfaces map[string]RouterInterface `yaml:"interfaces"`
	Firewall   Firewall                   `yaml:"firewall"`
	BGP        []BGP                      `yaml:"bgp"`
	Routes     []StaticRoute              `yaml:"routes"`
	DHCP       []DHCP                     `yaml:"dhcp"`
	NAT        []NAT                      `yaml:"nat"`
	Users      []User                     `yaml:"users"`
}

// RouterInterface is a single physical interface on a router
type RouterInterface struct {
	Name      string         `yaml:"name"`
	Addresses []netip.Prefix `yaml:"addresses"`
	MTU       uint16         `yaml:"mtu"`
	Speed     uint32         `yaml:"speed"`
	Duplex    string         `yaml:"duplex"`
	VLANs     []string       `yaml:"vlans"`
}

// Firewall config for the router firewall
type Firewall struct {
	Groups interface{}    `yaml:"groups"`
	Zones  []FirewallZone `yaml:"zones"`
}

// FirewallZone a single firewall zone
type FirewallZone struct {
	Name          string              `yaml:"name"`
	IPType        types.IPAddressType `yaml:"ip-type"`
	DefaultAction string              `yaml:"default-action"`
	Description   string              `yaml:"description"`
	In            []string            `yaml:"in"`
	Out           []string            `yaml:"out"`
	Local         []string            `yaml:"local"`
	Rules         []FirewallRule      `yaml:"rules"`
}

// FirewallRule is a single rule within a firewall zone
type FirewallRule struct {
	Action      string              `yaml:"action"`
	Description string              `yaml:"description"`
	Destination types.NetworkPort   `yaml:"destination"`
	Log         types.EnableDisable `yaml:"log"`
	Protocol    types.Protocol      `yaml:"protocol"`
	Established types.EnableDisable `yaml:"established"`
	Invalid     types.EnableDisable `yaml:"invalid"`
	New         types.EnableDisable `yaml:"new"`
	Related     types.EnableDisable `yaml:"related"`
}

// BGP Defines a single BGP configuration for an AS
type BGP struct {
	ASN           uint32         `yaml:"asn"`
	IP            netip.Addr     `yaml:"ip"`
	Peers         []BGPPeer      `yaml:"peers"`
	Announcements []netip.Prefix `yaml:"announcements"`
}

// BGPPeer is a peer and its configuration for a given BGP session
type BGPPeer struct {
	IP              netip.Addr `yaml:"ip"`
	ASN             uint32     `yaml:"asn"`
	AnnounceDefault bool       `yaml:"announce-default"`
}

// StaticRoute is a statically configured route in the router
type StaticRoute struct {
	Description string       `yaml:"description"`
	Route       netip.Prefix `yaml:"route"`
	NextHop     netip.Addr   `yaml:"next-hop"`
	Distance    uint8        `yaml:"distance"`
}

// DHCP is a single DHCP config for a single subnet
type DHCP struct {
	Name          string            `yaml:"name"`
	Authoritative bool              `yaml:"authoritative"`
	Subnet        netip.Prefix      `yaml:"subnet"`
	Router        netip.Addr        `yaml:"router"`
	Start         netip.Addr        `yaml:"start"`
	Stop          netip.Addr        `yaml:"stop"`
	Lease         uint64            `yaml:"lease"`
	DNS           []netip.Addr      `yaml:"dns"`
	Reservations  []DHCPReservation `yaml:"reservations"`
}

// DHCPReservation is a reserved IP by MAC address for a DHCP server
type DHCPReservation struct {
	Name string     `yaml:"name"`
	MAC  string     `yaml:"mac"`
	IP   netip.Addr `yaml:"ip"`
}

// NAT configures NAT rules in a router
type NAT struct {
	Name              string            `yaml:"name"`
	Type              types.NATType     `yaml:"type"`
	InboundInterface  string            `yaml:"inbound_interface"`
	OutboundInterface string            `yaml:"outbound_interface"`
	Protocol          types.Protocol    `yaml:"protocol"`
	Log               bool              `yaml:"log"`
	InsideAddress     types.AddressPort `yaml:"inside_address"`
	OutsideAddress    types.AddressPort `yaml:"outside_address"`
}

// Switch is the top level config for a single switch
type Switch struct{}

// VLAN defines a single shared VLAN configuration
type VLAN struct {
	Name    string       `yaml:"name"`
	ID      uint16       `yaml:"id"`
	Address netip.Prefix `yaml:"address"`
	MTU     uint16       `yaml:"mtu"`
}

// GetVLANByName returns a VLAN by its name attribute
func (c *Config) GetVLANByName(name string) (VLAN, error) {
	for _, vlan := range c.VLANs {
		if vlan.Name == name {
			return vlan, nil
		}
	}

	return VLAN{}, fmt.Errorf("could not find requested VLAN %s in config", name)
}
