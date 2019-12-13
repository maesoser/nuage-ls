package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/nuagenetworks/go-bambou/bambou"
)

func Save(name string, b interface{}) error {
	data, err := json.MarshalIndent(b, "", "\t")
	if err != nil {
		return err
	}
	tstr := time.Now().Format("20060102_")
	name = strings.Trim(name, " ")
	err = ioutil.WriteFile(tstr+name+".json", []byte(data), 0644)
	if err != nil {
		return err
	}
	return nil
}

func UptimeToStr(uptime int) string {
	seconds := uptime / 1000.0
	minutes := seconds / 60.0
	hours := minutes / 60.0
	days := hours / 24.0
	return fmt.Sprintf("%d days, %d:%02d", days, hours%24, minutes%60)
}

func Filter(text string) *bambou.FetchingInfo {
	if text == "" {
		return nil
	}
	f := bambou.NewFetchingInfo()
	f.FilterType = "predicate"
	f.Filter = text
	return f
}

func GetEnvStr(name, value string) string {
	if os.Getenv(name) != "" {
		return os.Getenv(name)
	}
	return value
}
