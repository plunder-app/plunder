package apiserver

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

//FindFunctionEndpoint - will do a look up to find an exposed dynamic endpoint
func FindFunctionEndpoint(u *url.URL, c *http.Client, f, m string) (*EndPoint, *Response) {
	// Create local URL for the API call
	newURL := *u
	newURL.Path = fmt.Sprintf("%s/%s/%s", FunctionPath(), f, m)

	// Interact with the API server to find the endpoint
	response, err := ParsePlunderGet(&newURL, c)
	if err != nil {

		return nil, &Response{
			Warning: fmt.Sprintf("Unable to find method [%s] for function [%s]", m, f),
			Error:   err.Error(),
		}
	}
	var ep EndPoint
	err = json.Unmarshal(response.Payload, &ep)
	if err != nil {
		response.Error = err.Error()
		return nil, response
	}
	return &ep, response
}

//BuildEnvironmentFromConfig will use the apiserver pkg to parse a configuration file and create a http client with the correct authentication and URL
func BuildEnvironmentFromConfig(path, urlFlag string) (*url.URL, *http.Client, error) {
	log.Debugf("Parsing Configuration file [%s]", path)

	// Open the configuration
	c, err := openClientConfig(path)
	if err != nil {
		return nil, nil, err
	}
	// Retrieve the certificate
	cert, err := c.RetrieveClientCert()
	if err != nil {
		return nil, nil, err
	}

	// Build the certificate pool from the unencrypted cert
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cert)

	// Create a HTTPS client and supply the created CA pool
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}

	// Build the URL from the configuration
	serverURL := c.GetServerAddressURL()

	// Overwrite the configuration url if
	if urlFlag != "" {
		serverURL, err = url.Parse(urlFlag)
		if err != nil {
			return nil, nil, err
		}
	}

	return serverURL, client, nil
}

//ParsePlunderGet will attempt to retrieve data from the plunder API server
func ParsePlunderGet(u *url.URL, c *http.Client) (*Response, error) {
	var response Response

	log.Debugf("Querying the Plunder Server [%s]", u.String())

	resp, err := c.Get(u.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode > 200 {
		return nil, fmt.Errorf(resp.Status)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil

}

//ParsePlunderPost will attempt to retrieve data from the plunder API server
func ParsePlunderPost(u *url.URL, c *http.Client, data []byte) (*Response, error) {
	var response Response

	log.Debugf("Posting [%d] bytes to the Plunder Server [%s]", len(data), u.String())

	resp, err := c.Post(u.String(), "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode > 200 {
		return nil, fmt.Errorf(resp.Status)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil

}

//ParsePlunderDelete will attempt to retrieve data from the plunder API server
func ParsePlunderDelete(u *url.URL, c *http.Client) (*Response, error) {
	var response Response

	log.Debugf("Requesting DELETE method to [%s]", u.String())

	// Create request
	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode > 200 {
		return nil, fmt.Errorf(resp.Status)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil

}
