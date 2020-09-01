/*Copyright (c) 2019 Cisco and/or its affiliates.

This software is licensed to you under the terms of the Cisco Sample
Code License, Version 1.0 (the "License"). You may obtain a copy of the
License at

               https://developer.cisco.com/docs/licenses

All use of the material herein must be in accordance with the terms of
the License. All rights not expressly granted by the License are
reserved. Unless required by applicable law or agreed to separately in
writing, software distributed under the License is distributed on an "AS
IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
or implied.*/

package ccp

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"net/http"
)

// LoginCreds for provider
type LoginCreds struct {
	Username *string `json:"username" validate:"nonzero"`
	Password *string `json:"password" validate:"nonzero"`
}

// How the heck does this work?
// func (s *Client) Login(client *Client) error {
//       ^ attach method: can take on functions
//                        ^ input arg
//                                        ^ return

// Login updated for v3
func (s *Client) Login(client *Client) error {

	url := s.BaseURL + "/v3/system/login"

	loginCreds := LoginCreds{
		Username: String(client.Username),
		Password: String(client.Password),
	}

	// Set up custom transport with Proxy from CLI using HTTP_PROXY and TLS set not to verify certs
	var PTransport = &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	cl := http.Client{
		Transport: PTransport,
	}

	// Marshal the JSON payload to then send
	j, err := json.Marshal(loginCreds)
	if err != nil {
		return err
	}

	// Send the JSON payload
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))
	if err != nil {
		Debug(1, "Error logging in: "+err.Error())
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := cl.Do(req)
	if err == nil {
		Debug(1, "Logged in as user "+client.Username)
	} else {
		Debug(1, "Error logging in: "+err.Error())
		// Debug(1, "Response: "+ioutil.ReadAll(resp.Body))
		return err
	}
	// fmt.Println("Response: \n")
	// fmt.Println(resp)
	// fmt.Println("Geting X-Auth-Token")

	var xauthtoken = resp.Header.Get("X-Auth-Token")
	Debug(1, "X-Auth-Token = "+xauthtoken)
	// set xauth
	s.XAuthToken = xauthtoken

	defer resp.Body.Close()

	return nil
}
