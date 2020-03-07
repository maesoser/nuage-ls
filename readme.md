# nuage-nsg

This tool let you query VSD about an [NSG](http://bootstrap.nuagenetworks.net/) or an specific group of NSGs and returns to you the configuration of the matched NSGs.

A Network Services Gateway (NSG) is a SD-WAN capable router made by [Nuage Networks](https://www.nuagenetworks.net/) that is able to interconect the local networks of multiple sites belonging to the same company without the need to use expensive MPLS networks. It also simplifies its management by following the SDN principles.

This tool gives you information about:

- [VSCs](https://nuagenetworks.github.io/vsd-api-documentation/v5_0/infrastructurevscprofile.html) configured on an specific NSG
- TPM Version and status.
- NSG IP addresses and operational status.
- NSG phisical [location](https://nuagenetworks.github.io/vsd-api-documentation/v5_0/location.html).
- NSG enterprise name and description.
- Network sie interface configuration.
- Network side BGP configuration.
- Access (LAN) side interface configuration.
- Active [alarms](https://nuagenetworks.github.io/vsd-api-documentation/v5_0/alarm.html) on that NSG.


```bash
./nuage-nsg -h

Usage of nsg-ls:

-vsd
  VSD Url.
-usr
  VSD User.
-pwd
  VSD Password.
-org
  VSD Enterprise.

-f
  NSG name/description/systemID filter.
-e
  Enterprise Name/Description filter.

-detail
  Show details about NSG device(s).
-vsc
  Show details about configured NSG(s) VSC profiles.
-uplink
  Show details about configured NSG(s) uplinks.
-port
  Show details about configured NSG(s) vports.
-alarm
  Show last alarms generated on NSG(s).

-save
  Save NSG data to json file.
-log
  Save log file.
-session
  Connection details stored on /etc/nsgconfig.json

```

```bash
./nuage-nsg -f "Lab1" -vsc -uplink -port -detail -alarm

NSG: NSG_Lab1_Nokia
        Description: NSG de prueba
        Location: RIBERA DEL DUERO, 2 - MADRID - MADRID
        Enterprise: Empresa de prueba
        Enterprise Descr: Empresa de prueba
        Status: CONNECTED               Uptime: 194 days, 15:27
        TPM: ENABLED_OPERATIONAL        TPM Ver: 1.2.4.43
        SystemID: 104.145.233.147        Address: 8.5.17.2
        Model: NSG_X                    Version: Nuage NSG 5.4.1_148
        
        VSC: labvsc01
                Status: UP                      Version: C-5.4.1-148
                IPs: [10.10.10.11 10.0.0.11 192.168.1.11]
        VSC: labvsc02
                Status: UP                      Version: C-5.4.1-148
                IPs: [10.10.10.12 10.0.0.12 192.168.1.12]
        VSC: labvsc03
                Status: UP                      Version: C-5.4.1-148
                IPs: [10.10.11.11 10.0.1.11 192.168.2.11]
        VSC: labvsc04
                Status: UP                      Version: C-5.4.1-148
                IPs: [10.10.11.11 10.0.1.12 192.168.2.12]
                
        Uplink: port1.0
                Underlay: MPLS
                Mode: Static                    Role: PRIMARY
                Address: 128.0.20.20           ScnAddress: 128.0.10.20
                Mask: 255.255.255.240           Gateway: 128.0.10.1
                BGP: BGPNeighbor-test-lab1
                   Peer AS: 61100          Peer Addr: 128.0.20.21
        Uplink: port2.0
                Underlay: Internet
                Mode: Dynamic                   Role: SECONDARY
                
        Port: port3.0   (RG)
                Name: Test_Port_RG
                L3: Test_Domain_1
                    Description: Dominio de pruebas ACLs
                    ServiceID  : 130091033
                Zone: Test_Zone_1
                L2: Subnet_Zone_1
                    Description: Subred de pruebas ACLs
                    ServiceID  : 4º2000472
                    Network: 10.5.5.0/24       Gateway: 10.5.5.1
        Alarms:
                [2019-08-09 13:26:22] Port [port2] of gateway [NSG_Lab1_Nokia] with system-id [104.145.233.147] is down.
                [2019-08-01 08:17:15] Gateway Instance is not using [4] ports
```
