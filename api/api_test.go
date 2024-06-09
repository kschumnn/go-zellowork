package api

import (
	"fmt"
	"os"
	"testing"

	"github.com/kschumnn/go-zellowork/types"
)

func TestAuth(t *testing.T) {
	url, apikey, username, password := os.Getenv("HOST"), os.Getenv("APIKEY"), os.Getenv("USER"), os.Getenv("PASS")
	ac := NewAPIClient(url, apikey)
	_, err := ac.Authenticate(username, password)
	if err != nil {
		t.Error(err)
	}
}
func TestGetUsers(t *testing.T) {
	url, apikey, username, password := os.Getenv("HOST"), os.Getenv("APIKEY"), os.Getenv("USER"), os.Getenv("PASS")
	ac := NewAPIClient(url, apikey)
	_, err := ac.Authenticate(username, password)
	if err != nil {
		t.Error(err)
	}
	users, err := ac.GetUsers(false)
	if err != nil {
		t.Error(err)
	}
	for _, user := range users {
		fmt.Println(user.Name)
	}
}
func TestGetChannels(t *testing.T) {
	url, apikey, username, password := os.Getenv("HOST"), os.Getenv("APIKEY"), os.Getenv("USER"), os.Getenv("PASS")
	ac := NewAPIClient(url, apikey)
	_, err := ac.Authenticate(username, password)
	if err != nil {
		t.Error(err)
	}
	chans, err := ac.ChannelList()
	if err != nil {
		t.Error(err)
	}
	for _, ch := range chans {
		fmt.Println(ch.Name)
	}
}
func TestCreateChannel(t *testing.T) {
	url, apikey, username, password := os.Getenv("HOST"), os.Getenv("APIKEY"), os.Getenv("USER"), os.Getenv("PASS")
	ac := NewAPIClient(url, apikey)
	_, err := ac.Authenticate(username, password)
	if err != nil {
		t.Error(err)
	}
	err = ac.ChannelAdd("Test Channel", false, false)
	if err != nil {
		t.Error(err)
	}
}
func TestDeleteChannel(t *testing.T) {
	url, apikey, username, password := os.Getenv("HOST"), os.Getenv("APIKEY"), os.Getenv("USER"), os.Getenv("PASS")
	ac := NewAPIClient(url, apikey)
	_, err := ac.Authenticate(username, password)
	if err != nil {
		t.Error(err)
	}
	err = ac.ChannelDelete("Test Channel")
	if err != nil {
		t.Error(err)
	}
}
func TestAddUsersChannel(t *testing.T) {
	url, apikey, username, password := os.Getenv("HOST"), os.Getenv("APIKEY"), os.Getenv("USER"), os.Getenv("PASS")
	ac := NewAPIClient(url, apikey)
	_, err := ac.Authenticate(username, password)
	if err != nil {
		t.Error(err)
	}
	users := []string{}
	{
		usersList, err := ac.GetUsers(false)
		if err != nil {
			t.Error(err)
		}
		for _, user := range usersList {
			users = append(users, user.Name)
		}
	}

	err = ac.AddUserToChannel(users, "Test Channel")
	if err != nil {
		t.Error(err)
	}
}
func TestRoleListChannel(t *testing.T) {
	url, apikey, username, password := os.Getenv("HOST"), os.Getenv("APIKEY"), os.Getenv("USER"), os.Getenv("PASS")
	ac := NewAPIClient(url, apikey)
	_, err := ac.Authenticate(username, password)
	if err != nil {
		t.Error(err)
	}
	roles, err := ac.ChannelRolesList("Test Channel")
	if err != nil {
		t.Error(err)
	}
	for _, role := range roles {
		fmt.Println(role)
	}
}
func TestRoleSaveChannel(t *testing.T) {
	url, apikey, username, password := os.Getenv("HOST"), os.Getenv("APIKEY"), os.Getenv("USER"), os.Getenv("PASS")
	ac := NewAPIClient(url, apikey)
	_, err := ac.Authenticate(username, password)
	if err != nil {
		t.Error(err)
	}
	err = ac.ChannelRoleSave("Test Channel", types.ZelloChannelRole{
		Name: "LO",
		Settings: types.ZelloChannelRoleSettings{
			ListenOnly:   true,
			NoDisconnect: false,
		},
	})
	if err != nil {
		t.Error(err)
	}
	err = ac.ChannelRoleSave("Test Channel", types.ZelloChannelRole{
		Name: "GW",
		Settings: types.ZelloChannelRoleSettings{
			ListenOnly:   false,
			NoDisconnect: false,
		},
	})
	if err != nil {
		t.Error(err)
	}
}

func TestRoleAddUsersChannel(t *testing.T) {
	url, apikey, username, password := os.Getenv("HOST"), os.Getenv("APIKEY"), os.Getenv("USER"), os.Getenv("PASS")
	ac := NewAPIClient(url, apikey)
	_, err := ac.Authenticate(username, password)
	if err != nil {
		t.Error(err)
	}

	loUsers := []string{}
	{
		usersList, err := ac.GetUsers(false)
		if err != nil {
			t.Error(err)
		}
		for _, user := range usersList {
			loUsers = append(loUsers, user.Name)
		}
	}

	gwUsers := []string{}
	{
		usersList, err := ac.GetUsers(true)
		if err != nil {
			t.Error(err)
		}
		for _, user := range usersList {
			gwUsers = append(gwUsers, user.Name)
		}
	}

	err = ac.ChannelRoleAddUser("Test Channel", "LO", loUsers)
	if err != nil {
		t.Error(err)
	}
	err = ac.ChannelRoleAddUser("Test Channel", "GW", gwUsers)
	if err != nil {
		t.Error(err)
	}
}
