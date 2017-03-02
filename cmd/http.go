package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	if total < 1 {
		return nil
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

func getReader(path string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", flServerURL, path), nil)
	if err != nil {
		return nil, errors.New("there was an error creating the request")
	}
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", flToken))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.New("there was an error completing the request")
	}
	switch resp.StatusCode {
	case 401:
		return nil, errors.New("authentication failed: you may need to run auth again")
	case 400:
		return nil, errors.New("there were some validation errors for your request")
	case 404:
		return nil, errors.New("item not found on server")
	case 500:
		return nil, errors.New("there was an internal server error")
	}
	return resp.Body, nil
}

func getJSON(path string, data interface{}) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", flServerURL, path), nil)
	if err != nil {
		return errors.New("there was an error creating the request")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", flToken))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New("there was an error completing the request")
	}
	switch resp.StatusCode {
	case 401:
		return errors.New("authentication failed: you may need to run auth again")
	case 403:
		return errors.New("you do not have access to that resource")
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
		debug(string(body))
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
		msg, _ := ioutil.ReadAll(resp.Body)
		debug(string(msg))
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

func patchJSON(path string, data interface{}) ([]byte, error) {
	var body []byte
	buff, err := json.Marshal(data)
	if err != nil {
		return body, errors.New("there was an error encoding JSON")
	}
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s%s", flServerURL, path), bytes.NewBuffer(buff))
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

func postMultipart(path, contentType string, reader io.Reader) ([]byte, error) {
	var body []byte
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", flServerURL, path), reader)
	if err != nil {
		return body, errors.New("there was an error creating the request")
	}
	req.Header.Set("Content-Type", contentType)
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

func deleteHTTP(path string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s%s", flServerURL, path), nil)
	if err != nil {
		return errors.New("there was an error creating the request")
	}
	req.Header.Set("Content-Type", "application/json")
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
	return nil
}
