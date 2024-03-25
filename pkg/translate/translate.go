package translate

import (
	"fmt"

	"github.com/cmmarslender/edgefig/pkg/config"
	"github.com/cmmarslender/edgefig/pkg/edgeconfig"
)

// ConfigToEdgeConfig translates the friendly config to edgerouter config
// @TODO this should return a whole set of configs, not just router configs
func ConfigToEdgeConfig(cfg *config.Config) (*edgeconfig.Router, error) {
	if len(cfg.Routers) == 0 {
		return nil, fmt.Errorf("no routers configured")
	}

	// @TODO Deal with more than the 0th router
	router := cfg.Routers[0]

	defaultRouter := getDefaultRouterConfig()
	defaultRouter.Firewall.AllPing = edgeconfig.Enable
	defaultRouter.Firewall.SendRedirects = edgeconfig.Enable
	defaultRouter.Firewall.SynCookies = edgeconfig.Enable

	for intf, intCfg := range router.Interfaces {
		_iface := edgeconfig.Interface{
			Type:        edgeconfig.InterfaceTypeEthernet,
			Name:        intf,
			State:       edgeconfig.Enabled,
			Description: intCfg.Name,
			Address:     intCfg.Addresses,
			//Duplex:      "",
			//Speed:       "",
		}
		if intCfg.MTU != 0 {
			_iface.MTU = intCfg.MTU
		}

		// @TODO make some methods to keep references by key vs this hunting/replacing
		for replI, replInt := range defaultRouter.Interfaces.Interfaces {
			if replInt.Name == _iface.Name {
				defaultRouter.Interfaces.Interfaces[replI] = _iface
			}
		}
	}

	_dhcpServer := edgeconfig.DHCPServer{
		Disabled:       len(router.DHCP) == 0,
		HostfileUpdate: false,
		StaticARP:      false,
		UseDNSMASQ:     false,
	}
	for _, dhcpCfg := range router.DHCP {
		_dhcpNetwork := edgeconfig.DHCPNetwork{
			Name:          dhcpCfg.Name,
			Authoritative: true,
			Subnets: []edgeconfig.DHCPSubnet{
				{
					Subnet: dhcpCfg.Subnet,
					Router: dhcpCfg.Router,
					Lease:  dhcpCfg.Lease,
					DNS:    dhcpCfg.DNS,
					StartStop: edgeconfig.DHCPStartStop{
						Start: dhcpCfg.Start,
						Stop:  dhcpCfg.Stop,
					},
				},
			},
		}

		_dhcpServer.Networks = append(_dhcpServer.Networks, _dhcpNetwork)
	}
	defaultRouter.Service.DHCPServer = _dhcpServer

	_natService := edgeconfig.NatService{Rules: []edgeconfig.NatRule{}}
	for _, natRule := range router.NAT {
		_natService.Rules = append(_natService.Rules, edgeconfig.NatRule{
			Name:              natRule.Name,
			Type:              natRule.Type,
			InboundInterface:  natRule.InboundInterface,
			OutboundInterface: natRule.OutboundInterface,
			Protocol:          natRule.Protocol,
			Log:               edgeconfig.EnableDisable(natRule.Log),
			OutsideAddress:    natRule.OutsideAddress,
			InsideAddress:     natRule.InsideAddress,
		})
	}
	defaultRouter.Service.NAT = _natService

	defaultRouter.System.HostName = router.Name
	for _, user := range router.Users {
		defaultRouter.System.Login.Users = append(defaultRouter.System.Login.Users, edgeconfig.User{
			Username: user.Username,
			Authentication: edgeconfig.Authentication{
				EncryptedPassword: user.Password,
			},
			Level: user.Role,
		})
	}

	return defaultRouter, nil
}