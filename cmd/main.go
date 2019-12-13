package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/maesoser/nuage-nsg/pkg/config"
	"github.com/maesoser/nuage-nsg/pkg/fetch"
	"github.com/maesoser/nuage-nsg/pkg/show"
	"github.com/maesoser/nuage-nsg/pkg/util"
	"github.com/nuagenetworks/vspk-go/vspk"
)

func usage() {
	flagSet := flag.CommandLine
	fmt.Printf("Usage of %s:\n", "nsg-ls")
	order := []string{"vsd", "usr", "pwd", "org", "f", "e", "detail", "vsc", "uplink", "port", "alarm", "save", "log", "session"}
	for _, name := range order {
		flag := flagSet.Lookup(name)
		fmt.Printf("-%s\n", flag.Name)
		fmt.Printf("  %s\n", flag.Usage)
	}
}

func main() {

	flag.Usage = usage
	// Login variables
	VSDUrl := flag.String("vsd", util.GetEnvStr("VSD_ADDRESS", "127.0.0.1:8443"), "VSD Url.")
	VSDUser := flag.String("usr", util.GetEnvStr("VSD_USER", ""), "VSD User.")
	VSDOrg := flag.String("org", util.GetEnvStr("VSD_ENTERPRISE", ""), "VSD Enterprise.")
	VSDPasswd := flag.String("pwd", util.GetEnvStr("VSD_PASSWD", ""), "VSD Password.")

	// Filtering variables
	nsgFilter := flag.String("f", "", "NSG name/description/systemID filter.")
	enterpriseFilter := flag.String("e", "", "Enterprise Name/Description filter.")

	// Show options
	showDetails := flag.Bool("detail", false, "Show details about NSG device(s).")
	showPorts := flag.Bool("port", false, "Show details about configured NSG(s) vports.")
	showUplinks := flag.Bool("uplink", false, "Show details about configured NSG(s) uplinks.")
	showAlarms := flag.Bool("alarm", false, "Show last alarms generated on NSG(s).")
	showVSCs := flag.Bool("vsc", false, "Show details about configured NSG(s) VSC profiles.")

	// Misc
	savePtr := flag.Bool("save", false, "Save NSG data to json file.")
	saveLog := flag.String("log", "", "Save log file.")
	profilePtr := flag.String("session", "", "Connection details stored on /etc/nsgconfig.json")
	flag.Parse()

	if *saveLog != "" {
		f, err := os.OpenFile(*saveLog, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Error opening file: %v", err)
		}
		defer f.Close()
		log.SetOutput(f)
		log.Println("[INFO] Executing nuage-nsg-ls")
	} else {
		log.SetOutput(ioutil.Discard)
	}

	if *showPorts || *showUplinks || *showAlarms || *showVSCs {
		*showDetails = true
	}
	*nsgFilter = strings.ToLower(*nsgFilter)

	if *profilePtr != "" {
		s, e := config.ProfileFromFile(*profilePtr)
		if e != nil {
			fmt.Println(e)
		} else {
			*VSDUrl = s.Address
			*VSDUser = s.User
			*VSDOrg = s.Organization
		}
	}

	s, root := vspk.NewSession(*VSDUser, *VSDPasswd, *VSDOrg, "https://"+*VSDUrl)

	if err := s.Start(); err != nil {
		fmt.Println("Unable to connect to Nuage VSD: " + err.Description)
		os.Exit(1)
	}

	if *showDetails == false {
		fmt.Printf("\n   %-20s\t%-10s\t%-38s\t%-30s\n", " SystemID", " Address", " Name", " Organization")
	}
	filter := util.Filter("name CONTAINS \"" + *enterpriseFilter + "\" OR description CONTAINS \"" + *enterpriseFilter + "\"")
	enterprises, err := root.Enterprises(filter)
	if err != nil {
		fmt.Println("Unable to fetch enterprises", err.Description)
	}
	for _, enterprise := range enterprises {
		data := make(chan fetch.NSGData)
		filter := util.Filter("systemID CONTAINS \"" + *nsgFilter + "\" OR name CONTAINS \"" + *nsgFilter +
			"\" OR description CONTAINS \"" + *nsgFilter + "\"")
		go fetch.GetNSGs(
			root,
			enterprise,
			filter,
			*showUplinks,
			*showPorts,
			*showAlarms,
			*showVSCs,
			data,
		)
		for nsg := range data {
			if *showDetails == true {
				show.DetailedNSGData(nsg)
			} else {
				show.NSGData(nsg)
			}
			if *savePtr == true {
				fileName := strings.Trim(enterprise.Name, " ") + "_" + strings.Trim(nsg.NSG.Name, " ")
				util.Save(fileName, nsg)
			}
		}
	}
}
