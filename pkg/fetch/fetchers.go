package fetch

import (
	"fmt"
	"log"

	"github.com/maesoser/nuage-nsg/pkg/util"
	"github.com/nuagenetworks/go-bambou/bambou"
	"github.com/nuagenetworks/vspk-go/vspk"
)

type NSGData struct {
	EnterpriseName        string
	EnterpriseDescription string
	IsRedundant           bool
	NSG                   *vspk.NSGateway
	Locations             vspk.LocationsList
	VRS                   *vspk.VRS
	VSCs                  []vspk.VSC
	Uplinks               []UplinkData
	Ports                 []PortData
	Alarms                vspk.AlarmsList
}

type NSGList []NSGData

type PortData struct {
	Domain vspk.Domain  `json:"domain"`
	Zone   vspk.Zone    `json:"zone"`
	Subnet *vspk.Subnet `json:"subnet"`
	Vport  *vspk.VPort  `json:"vport"`
}

type PortList []PortData

type UplinkData struct {
	Info *vspk.UplinkConnection
	BGP  vspk.BGPNeighbor
}

type UplinkList []UplinkData

func GetPorts(nsg *vspk.NSGateway, nsgid string) ([]PortData, error) {
	var ports []PortData
	subnets, err := nsg.Subnets(nil)
	if err != nil {
		return ports, err
	}
	for _, subnet := range subnets {
		vports, err := subnet.VPorts(nil)
		if err != nil {
			return ports, err
		}
		for _, vport := range vports {
			// Obtengo la VLAN por Vport
			var vlan vspk.VLAN
			vlan.SetIdentifier(vport.VLANID)
			vlan.Fetch()
			if vlan.GatewayID == nsgid {
				var zone vspk.Zone
				zone.SetIdentifier(vport.ZoneID)
				zone.Fetch()

				var domain vspk.Domain
				domain.SetIdentifier(vport.DomainID)
				domain.Fetch()

				var eVPort PortData
				eVPort.Vport = vport
				eVPort.Domain = domain
				eVPort.Subnet = subnet
				eVPort.Zone = zone

				ports = append(ports, eVPort)
			}
		}
	}
	return ports, nil
}

func GetVSCs(vrs *vspk.VRS) []vspk.VSC {
	var vscs []vspk.VSC
	for _, id := range vrs.ParentIDs {
		var vsc vspk.VSC
		vsc.SetIdentifier(id.(string))
		vsc.Fetch()
		vscs = append(vscs, vsc)
	}
	return vscs
}

func GetVRS(root *vspk.Me, systemID string) (*vspk.VRS, error) {
	filter := util.Filter("description CONTAINS \"" + systemID + "\"")
	vrss, err := root.VRSs(filter)
	if err != nil {
		return nil, err
	} else if len(vrss) == 0 {
		return nil, fmt.Errorf("No VSCs found")
	}
	return vrss[0], nil
}

func GetNSGs(root *vspk.Me,
	enterprise *vspk.Enterprise,
	filter *bambou.FetchingInfo,
	uplinkInfo bool,
	portInfo bool,
	alarmInfo bool,
	vscInfo bool,
	data chan NSGData,
) {
	var nsgsdata NSGList

	nsgateways, err := enterprise.NSGateways(filter)
	if err != nil {
		log.Printf("[ERROR] Unable to fetch NSGateways for enterprise %s: %s\n",
			enterprise.Name,
			err.Error())
	}
	rgroups, err := enterprise.NSRedundantGatewayGroups(nil)
	if err != nil {
		log.Printf("[ERROR] Unable to fetch NSG Redundant Groups for enterprise %s: %s\n",
			enterprise.Name,
			err.Error())
	}
	for _, nsg := range nsgateways {
		var nsgdata NSGData

		nsgdata.EnterpriseName = enterprise.Name
		nsgdata.EnterpriseDescription = enterprise.Description

		nsgdata.NSG = nsg
		vrs, err := GetVRS(root, nsg.SystemID)
		nsgdata.VRS = vrs
		if err != nil {
			log.Printf("[ERROR] Unable to fetch VRS Object for NSG %s: %s\n", nsg.SystemID, err.Error())
		}
		loc, err := nsg.Locations(nil)
		nsgdata.Locations = loc
		if err != nil && loc == nil {
			log.Printf("[ERROR] Unable to gather Location for NSG %s: %v\n", nsg.SystemID, err)
		}
		if vscInfo {
			nsgdata.VSCs = GetVSCs(vrs)
		}
		if uplinkInfo {
			uplinks, err := nsg.UplinkConnections(nil)
			if err != nil {
				log.Printf("[ERROR] Unable to fetch Uplinks for NSG %s: %s\n", nsg.SystemID, err.Error())
			}
			for _, uplink := range uplinks {
				var bgp vspk.BGPNeighbor
				if uplink.AssociatedBGPNeighborID != "" {
					bgp.SetIdentifier(uplink.AssociatedBGPNeighborID)
					bgp.Fetch()
				}
				udata := UplinkData{
					Info: uplink,
					BGP:  bgp,
				}
				nsgdata.Uplinks = append(nsgdata.Uplinks, udata)
			}
		}
		if portInfo {
			ID, redundant := GetPortID(rgroups, nsg.ID)
			nsgdata.IsRedundant = redundant
			ports, err := GetPorts(nsg, ID)
			nsgdata.Ports = ports
			if err != nil {
				log.Printf("[ERROR] Unable to fetch Ports for NSG %s: %s\n", nsg.SystemID, err.Error())
			}
		}
		if alarmInfo {
			alarms, err := nsg.Alarms(nil)
			nsgdata.Alarms = alarms
			if err != nil {
				log.Printf("[ERROR] Unable to fetch Alarms for NSG %s: %s\n", nsg.SystemID, err.Error())
			}
		}
		log.Println(nsgdata.NSG.Name)
		nsgsdata = append(nsgsdata, nsgdata)
		data <- nsgdata
	}
	close(data)
}

func GetPortID(groups vspk.NSRedundantGatewayGroupsList, nsgID string) (string, bool) {
	for _, rgroup := range groups {
		if rgroup.GatewayPeer1ID == nsgID || rgroup.GatewayPeer2ID == nsgID {
			return rgroup.ID, true
		}
	}
	return nsgID, false
}
