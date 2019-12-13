package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type Session struct {
	Name         string `json:"name"`
	Address      string `json:"address"`
	User         string `json:"user"`
	Organization string `json:"org"`
	Password     string `json:"pass"`
}

type Sessions struct {
	Sessions []Session `json:"sessions"`
}

func ProfileFromFile(name string) (Session, error) {
	var s Session
	jsonFile, err := os.Open("/etc/vsdconfig.json")
	if err != nil {
		return s, err
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var sessions Sessions
	json.Unmarshal(byteValue, &sessions)
	for _, session := range sessions.Sessions {
		if session.Name == name {
			return session, nil
		}
	}
	return s, errors.New("Profile not found")
}
