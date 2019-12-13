package show

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/maesoser/nuage-nsg/pkg/fetch"
	"github.com/maesoser/nuage-nsg/pkg/util"
)

func NSGData(data fetch.NSGData) {
	status := "+"
	if data.NSG.OperationStatus != "CONNECTED" {
		status = "-"
	} else if data.NSG.OperationStatus == "UNKNOWN" {
		status = "?"
	}
	address := "unknown"
	if data.VRS != nil {
		address = data.VRS.Address
	}
	fmt.Printf("%s %-20s\t%-10s\t%-38s\t%-30s\n",
		status, data.NSG.SystemID, address,
		data.NSG.Description, data.EnterpriseDescription)
}

func DetailedNSGData(data fetch.NSGData) {
	uptime := "Unknown"
	address := "Unknown"
	if data.VRS != nil {
		uptime = util.UptimeToStr(data.VRS.Uptime)
		address = data.VRS.Address
	}

	fmt.Printf("\nNSG: %s\n", data.NSG.Name)
	fmt.Printf("\tDescription: %s\n", data.NSG.Description)
	for _, loc := range data.Locations {
		if loc.Address != "" {
			fmt.Printf("\tLocation: %s\n", loc.Address)
		}
	}
	fmt.Printf("\t%-30s\n", "Enterprise: "+data.EnterpriseName)
	fmt.Printf("\t%-30s\n", "Enterprise Descr: "+data.EnterpriseDescription)
	fmt.Printf("\t%-30s\t%-30s\n", "Status: "+data.NSG.OperationStatus, "Uptime: "+uptime)
	fmt.Printf("\t%-30s\t%-30s\n", "TPM: "+data.NSG.TPMStatus, "TPM Ver: "+data.NSG.TPMVersion)
	fmt.Printf("\t%-30s\t%-30s\n", "SystemID: "+data.NSG.SystemID, "Address: "+address)
	fmt.Printf("\t%-30s\t%-30s\n", "Model: "+data.NSG.Family, "Version: "+data.NSG.NSGVersion)

	if data.VSCs != nil {
		for _, vsc := range data.VSCs {
			fmt.Printf("\tVSC: %s (%s)\n", vsc.Name, vsc.Description)
			fmt.Printf("\t\t%-30s\t%-30s\n", "Status: "+vsc.Status, "Version: "+vsc.ProductVersion)
			vsc.Addresses = append(vsc.Addresses, vsc.ManagementIP)
			fmt.Printf("\t\tIPs: %s\n", vsc.Addresses)
		}
	}
	if data.Uplinks != nil {
		for _, uplink := range data.Uplinks {

			fmt.Printf("\tUplink: %s.%s\n", uplink.Info.PortName, strconv.Itoa(uplink.Info.Vlan))
			fmt.Printf("\t\t%-30s\n", "Underlay: "+uplink.Info.AssociatedUnderlayName)

			if uplink.Info.AuxiliaryLink {
				uplink.Info.Role = uplink.Info.Role + " AUX"
			}
			fmt.Printf("\t\t%-30s\t%-30s\n", "Mode: "+uplink.Info.Mode, "Role: "+uplink.Info.Role)
			if uplink.Info.Mode == "Static" {
				fmt.Printf("\t\t%-30s\t%-30s\n", "Address: "+uplink.Info.Address, "ScnAddress: "+uplink.Info.SecondaryAddress)
				fmt.Printf("\t\t%-30s\t%-30s\n", "Mask: "+uplink.Info.Netmask, "Gateway: "+uplink.Info.Gateway)
			}
			if uplink.BGP.ID != "" {
				fmt.Printf("\t\t%-30s\n", "BGP: "+uplink.BGP.Name)
				fmt.Printf("\t\t%s%d\t%-30s\n", "   Peer AS: ", uplink.BGP.PeerAS, "   Peer Addr: "+uplink.BGP.PeerIP)
			}
		}
	}
	if data.Ports != nil {
		for _, port := range data.Ports {
			mask, _ := net.IPMask(net.ParseIP(port.Subnet.Netmask).To4()).Size()
			network := port.Subnet.Address + "/" + strconv.Itoa(mask)

			fmt.Printf("\tPort: %s.%d", port.Vport.GatewayPortName, port.Vport.VLAN)
			if data.IsRedundant {
				fmt.Printf("\t(RG)")
			}
			fmt.Printf("\n\t\t%-30s\n", "Name: "+port.Vport.Name)
			fmt.Printf("\t\t%-30s\n", "L3: "+port.Domain.Name)
			fmt.Printf("\t\t%s%s\n", "    Description: ", port.Domain.Description)
			fmt.Printf("\t\t%s%d\n", "    ServiceID  : ", port.Domain.ServiceID)

			fmt.Printf("\t\t%-30s\t%-30s\n", "Zone: "+port.Zone.Name, port.Zone.Description)
			fmt.Printf("\t\t%-30s\n", "L2: "+port.Subnet.Name)
			fmt.Printf("\t\t%s%s\n", "    Description: ", port.Subnet.Description)
			fmt.Printf("\t\t%s%d\n", "    ServiceID  : ", port.Subnet.ServiceID)
			fmt.Printf("\t\t%s\t%s\n", "    Network: "+network, "Gateway: "+port.Subnet.Gateway)
		}
	}
	if data.Alarms != nil {
		fmt.Printf("\tAlarms:\n")
		for _, alarm := range data.Alarms {
			tm := time.Unix(int64(alarm.Timestamp/1000), 0)
			fmt.Printf("\t\t[%s] %s\n", tm.String(), alarm.Description)
		}
	}
}
