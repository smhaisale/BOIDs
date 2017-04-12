package main

import (
    "net/http"
    "io/ioutil"
    "log"
    "encoding/json"
)

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
