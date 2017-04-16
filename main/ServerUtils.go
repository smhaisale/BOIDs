package main

import (
    "net/http"
    "io/ioutil"
    "log"
    "encoding/json"
    "bytes"
)

var client = http.Client{}

type GetRequest struct {
    data string
}

// All contained variable names must begin with a capital letter to be visible by JSONWrapper
func toJsonString(object interface{}) string {
    json, err := json.Marshal(object)
    if err != nil {
        log.Fatal(err)
    }
    return string(json)
}

// Tested for arbitrary objects
func fromJsonString(object interface{}, message string) error {
    err := json.Unmarshal([]byte(message), object)
    if err != nil {
        log.Fatal(err)
        return err
    }
    return nil
}

func getRequestBody(msg interface {}, req *http.Request) interface{} {

    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
        panic(err)
    }
    log.Println(string(body))
    err = json.Unmarshal(body, msg)
    if err != nil {
        log.Printf("error: %v", err)
    }
    defer req.Body.Close()
    return msg
}

func getResponseBody(msg interface {}, resp *http.Response) error {

    log.Println(resp)

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }
    log.Println(string(body))
    err = json.Unmarshal(body, msg)
    if err != nil {
        log.Printf("error: %v", err)
    }
    defer resp.Body.Close()
    return err
}

func getDroneFromServer(droneAddress string) (drone Drone, err error) {
    resp, err := client.Get("http://" + droneAddress + DRONE_GET_INFO_URL)
    if err != nil {
        log.Println("Error! ", err)
        return
    }
    err = getResponseBody(&drone, resp)
    return drone, err
}

func addDroneToServer(droneId string, droneAddress string) error {
    _, err := client.Get("http://" + droneAddress + DRONE_ADD_DRONE_URL + "?id=" + droneId)
    if err != nil {
        log.Println("Error! ", err)
    }
    return err
}

// Takes a URL and does a GET request with request body as the provided data. Returns response as a json string.
func makeGetRequest(url string, data string) (string, error) {
    log.Println("Make GET request to " + url + " with data " + data)
    req, err := http.NewRequest("GET", url, bytes.NewBufferString(data))
    if err != nil {
        return "", err
    }
    response, err := client.Do(req)
    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Println("Error in making GET request! ", err)
        return "", err
    }
    return string(body), err
}