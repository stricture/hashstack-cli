package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func getTotal(path string) (int, error) {
	var total int
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", flServerURL, path), nil)
	if err != nil {
		return total, errors.New("there was an error creating the request")
	}
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", flToken))
	req.Header.Set("Range", "1-1")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return total, errors.New("there was an error completing the request")
	}
	switch resp.StatusCode {
	case 401:
		return total, errors.New("authentication failed: you may need to run auth again")
	case 400:
		return total, errors.New("there were some validation errors for your request")
	case 404:
		return total, errors.New("item not found on server")
	case 500:
		return total, errors.New("there was an internal server error")
	}
	contentRange := resp.Header.Get("Content-Range")
	if contentRange == "" {
		return total, errors.New("invalid response from server")
	}
	parts := strings.Split(contentRange, "/")
	if len(parts) != 2 {
		return total, errors.New("invalid response from server")
	}
	total, err = strconv.Atoi(parts[1])
	if err != nil {
		return total, errors.New("invalid response from server")
	}
	return total, nil
}

func getRangeJSON(path string, data interface{}) error {
	total, err := getTotal(path)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", flServerURL, path), nil)
	if err != nil {
		return errors.New("there was an error creating the request")
	}
	req.Header.Set("Range", fmt.Sprintf("1-%d", total))
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", flToken))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New("there was an error completing the request")
	}
	switch resp.StatusCode {
	case 401:
		return errors.New("authentication failed: you may need to run auth again")
	case 400:
		return errors.New("there were some validation errors for your request")
	case 404:
		return errors.New("item not found on server")
	case 500:
		return errors.New("there was an internal server error")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("there was an error reading the response body")
	}
	if err := json.Unmarshal(body, data); err != nil {
		return errors.New("there was an error decoding JSON returned from the server")
	}
	return nil
}

func getJSON(path string, data interface{}) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", flServerURL, path), nil)
	if err != nil {
		return errors.New("there was an error creating the request")
	}
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", flToken))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New("there was an error completing the request")
	}
	switch resp.StatusCode {
	case 401:
		return errors.New("authentication failed: you may need to run auth again")
	case 400:
		return errors.New("there were some validation errors for your request")
	case 404:
		return errors.New("item not found on server")
	case 500:
		return errors.New("there was an internal server error")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("there was an error reading the response body")
	}
	if err := json.Unmarshal(body, data); err != nil {
		return errors.New("there was an error decoding JSON returned from the server")
	}
	return nil
}

func postJSON(path string, data interface{}) ([]byte, error) {
	var body []byte
	buff, err := json.Marshal(data)
	if err != nil {
		return body, errors.New("there was an error encoding JSON")
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", flServerURL, path), bytes.NewBuffer(buff))
	if err != nil {
		return body, errors.New("there was an error creating the request")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", flToken))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return body, errors.New("there was an error completing the request")
	}
	switch resp.StatusCode {
	case 401:
		return body, errors.New("authentication failed: you may need to run auth again")
	case 400:
		return body, errors.New("there were some validation errors for your request")
	case 404:
		return body, errors.New("item not found on server")
	case 500:
		return body, errors.New("there was an internal server error")
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, errors.New("there was an error reading the response body")
	}
	return body, nil
}
