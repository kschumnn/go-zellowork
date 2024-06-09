package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
	zwt "github.com/kschumnn/go-zellowork/types"
)

type (
	APIClient struct {
		apiClient *resty.Client
		host      string
		apiKey    string
		sessionID string
	}
)
type justStatusAndCode struct {
	Status string `json:"status"`
	Code   string `json:"code"`
}
type user_gettoken_resp struct {
	Status string `json:"status"`
	Code   string `json:"code"`
	SID    string `json:"sid"`
	Token  string `json:"token"`
}
type user_list_resp struct {
	justStatusAndCode
	Users []zwt.ZelloUser `json:"users"`
}
type channel_list_resp struct {
	justStatusAndCode
	Channels []zwt.ZelloChannel `json:"channels"`
}
type channel_roleslist_resp struct {
	justStatusAndCode
	Roles []zwt.ZelloChannelRole `json:"roles"`
}

func NewAPIClient(host string, apiKey string) (apiClient *APIClient) {
	return &APIClient{
		host:      host,
		apiKey:    apiKey,
		apiClient: resty.New(), //.EnableTrace().SetDebug(true),
	}
}
func NewAPIClientWithSessionID(host string, apiKey string, sessionID string) (apiClient *APIClient) {
	return &APIClient{
		host:      host,
		apiKey:    apiKey,
		sessionID: sessionID,
		apiClient: resty.New().
			SetJSONMarshaler(json.Marshal).
			SetJSONUnmarshaler(json.Unmarshal),
	}
}

func (ac *APIClient) getURL(command string) string {
	return fmt.Sprintf("%s/%s", ac.host, command)
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
func (ac *APIClient) Authenticate(username, password string) (res bool, err error) {
	res1 := user_gettoken_resp{}
	res2 := justStatusAndCode{}

	resp, err := ac.apiClient.R().
		SetHeader("Accept", "application/json").
		Get(ac.getURL("user/gettoken"))
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(resp.Body(), &res1)
	if err != nil {
		return res, err
	}
	if res1.Code != "200" {
		return res, fmt.Errorf(res1.Status)
	}
	if resp.IsSuccess() {
		ac.sessionID = res1.SID
		resp, err = ac.apiClient.R().
			SetQueryParams(map[string]string{
				"sid": ac.sessionID,
			}).
			SetFormData(map[string]string{
				"username": username,
				"password": GetMD5Hash(fmt.Sprintf("%s%s%s", GetMD5Hash(password), res1.Token, ac.apiKey)),
			}).
			SetHeader("Accept", "application/json").
			Post(ac.getURL("user/login"))
		if err != nil {
			return res, err
		}
		err = json.Unmarshal(resp.Body(), &res2)
		if err != nil {
			return res, err
		}
		if res2.Code != "200" {
			return res, fmt.Errorf(res2.Status)
		}
	}
	return res, nil
}
func (ac *APIClient) GetUsers(OnlyGateways bool) (users []zwt.ZelloUser, err error) {
	res1 := user_list_resp{}
	url := ac.getURL("user/get")
	if OnlyGateways {
		url = ac.getURL("user/get/gateway/true")
	}
	resp, err := ac.apiClient.R().
		SetQueryParams(map[string]string{
			"sid": ac.sessionID,
		}).
		SetHeader("Accept", "application/json").
		Get(url)
	if err != nil {
		return users, err
	}
	err = json.Unmarshal(resp.Body(), &res1)
	if err != nil {
		return users, err
	}
	return res1.Users, nil
}
func (ac *APIClient) AddUserToChannel(users []string, channel string) (err error) {
	res1 := justStatusAndCode{}

	resp, err := ac.apiClient.R().
		SetQueryParams(map[string]string{
			"sid": ac.sessionID,
		}).
		SetFormDataFromValues(url.Values{
			"login[]": users,
		}).
		//SetHeader("Accept", "application/json").
		Post(fmt.Sprintf("%s%s", ac.getURL("/user/addto/"), url.PathEscape(channel)))
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp.Body(), &res1)
	if err != nil {
		return err
	}

	if res1.Code != "200" {
		return fmt.Errorf(res1.Status)
	}
	return nil
}
func (ac *APIClient) RemoveUserToChannel(users []string, channel string) (err error) {
	res1 := justStatusAndCode{}

	resp, err := ac.apiClient.R().
		SetQueryParams(map[string]string{
			"sid": ac.sessionID,
		}).
		SetFormDataFromValues(url.Values{
			"login[]": users,
		}).
		//SetHeader("Accept", "application/json").
		Post(fmt.Sprintf("%s%s", ac.getURL("/user/removefrom/"), url.PathEscape(channel)))
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp.Body(), &res1)
	if err != nil {
		return err
	}

	if res1.Code != "200" {
		return fmt.Errorf(res1.Status)
	}
	return nil
}

// Base Channel Stuff
func (ac *APIClient) ChannelList() (chans []zwt.ZelloChannel, err error) {
	res1 := channel_list_resp{}
	resp, err := ac.apiClient.R().
		SetQueryParams(map[string]string{
			"sid": ac.sessionID,
		}).
		SetHeader("Accept", "application/json").
		Get(ac.getURL("channel/get"))
	if err != nil {
		return chans, err
	}
	err = json.Unmarshal(resp.Body(), &res1)
	if err != nil {
		return chans, err
	}
	return res1.Channels, nil
}

/*
shared	GET	(Optional) "true" or "false". Set to "true" to create group channel, set to "false" to create dynamic channel. Default is "false"
invisible	GET	(Optional) "true" or "false". Set to "true" in combination with shared=true to create a hidden group channel. When combined with shared=false the behavior is not defined. Default is "false"
*/
func (ac *APIClient) ChannelAdd(name string, shared bool, hidden bool) (err error) {
	res1 := justStatusAndCode{}
	url := fmt.Sprintf("%s/name/%s/shared/%s/invisible/%s",
		ac.getURL("channel/add"),
		url.PathEscape(name),
		url.PathEscape(fmt.Sprint(shared)),
		url.PathEscape(fmt.Sprint(hidden)),
	)
	resp, err := ac.apiClient.R().
		SetQueryParams(map[string]string{
			"sid": ac.sessionID,
		}).
		SetHeader("Accept", "application/json").
		Get(url)
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp.Body(), &res1)
	if err != nil {
		return err
	}
	if res1.Code != "200" {
		return fmt.Errorf(res1.Status)
	}
	return nil
}
func (ac *APIClient) ChannelDelete(name string) (err error) {
	res1 := justStatusAndCode{}
	resp, err := ac.apiClient.R().
		SetQueryParams(map[string]string{
			"sid": ac.sessionID,
		}).
		SetFormData(map[string]string{
			"name[]": name,
		}).
		SetHeader("Accept", "application/json").
		Post(ac.getURL("channel/delete"))
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp.Body(), &res1)
	if err != nil {
		return err
	}
	if res1.Code != "200" {
		return fmt.Errorf(res1.Status)
	}
	return nil
}
func (ac *APIClient) ChannelRolesList(name string) (chans []zwt.ZelloChannelRole, err error) {
	res1 := channel_roleslist_resp{}
	resp, err := ac.apiClient.R().
		SetQueryParams(map[string]string{
			"sid": ac.sessionID,
		}).
		SetHeader("Accept", "application/json").
		Get(fmt.Sprintf("%s%s", ac.getURL("channel/roleslist/name/"), url.PathEscape(name)))
	if err != nil {
		return chans, err
	}
	err = json.Unmarshal(resp.Body(), &res1)
	if err != nil {
		return chans, err
	}
	return res1.Roles, nil
}
func (ac *APIClient) ChannelRoleSave(name string, role zwt.ZelloChannelRole) (err error) {
	res1 := justStatusAndCode{}
	urls := fmt.Sprintf("%s/channel/%s/name/%s",
		ac.getURL("channel/saverole"),
		url.PathEscape(name),
		url.PathEscape(role.Name),
	)
	jsonSettings, err := json.Marshal(role.Settings)
	if err != nil {
		return err
	}
	resp, err := ac.apiClient.R().
		SetQueryParams(map[string]string{
			"sid": ac.sessionID,
		}).
		SetFormDataFromValues(url.Values{
			"settings": []string{string(jsonSettings)},
		}).
		SetHeader("Accept", "application/json").
		Post(urls)
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp.Body(), &res1)
	if err != nil {
		return err
	}
	if res1.Code != "200" {
		return fmt.Errorf(res1.Status)
	}
	return nil
}
func (ac *APIClient) ChannelRoleAddUser(name string, roleName string, users []string) (err error) {
	res1 := justStatusAndCode{}
	urls := fmt.Sprintf("%s/channel/%s/name/%s",
		ac.getURL("channel/addtorole"),
		url.PathEscape(name),
		url.PathEscape(roleName),
	)
	resp, err := ac.apiClient.R().
		SetQueryParams(map[string]string{
			"sid": ac.sessionID,
		}).
		SetFormDataFromValues(url.Values{
			"login[]": users,
		}).
		SetHeader("Accept", "application/json").
		Post(urls)
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp.Body(), &res1)
	if err != nil {
		return err
	}
	if res1.Code != "200" {
		return fmt.Errorf(res1.Status)
	}
	return nil
}

//
//
