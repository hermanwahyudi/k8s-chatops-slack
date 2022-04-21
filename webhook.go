/*
 * @author: Herman Wahyudi, DevOps
 *
 */
 
package main

import (
    "fmt"
    "net/http"
    "log"
    "encoding/json"
    "html"
    "strings"
    "bytes"
    "io/ioutil"
    "github.com/tkanos/gonfig"
    "strconv"
)

var configuration Configuration
var eventSubs EventSubscriptions
var greetings string

type Configuration struct {
    Port string
    Name string
    Token string
    SlackURL string
}

type EventSubscriptions struct {
    Challenge string `json:"challenge"`

    Event struct {
        Type string `json:"type"`
        Text string `json:"text"`
        User string `json:"user"`
        Channel string `json:"channel"`
    }
}

func index(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Welcome ", configuration.Name)

    log.Println("Welcome ", configuration.Name)
}

func eventSubscriptions(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        json.NewDecoder(r.Body).Decode(&eventSubs)
        json.NewEncoder(w).Encode(eventSubs)

        fmt.Println(eventSubs.Challenge)

        log.Println(eventSubs.Challenge)

        if eventSubs.Event.Type == "app_mention" {
            text := strings.ToLower(eventSubs.Event.Text)
            channel := eventSubs.Event.Channel
            user:= eventSubs.Event.User

            log.Println(text, channel, user)

            botResponse(text, channel, user)
        }

    } else {
        fmt.Println(w, "Execution %q %s", r.Method, html.EscapeString(r.URL.Path))
        fmt.Fprintf(w, "Execution %q %s", r.Method, html.EscapeString(r.URL.Path))

        log.Println("Execution ", r.Method, html.EscapeString(r.URL.Path))
    }
}

func botResponse(text, channel, user string) {
    statusResponse, responseDict := false, [] string{"hi", "hai", "hey", "hallo", "halo", "hello", "helo"}

    for _, value := range responseDict { 
        if strings.Contains(text, value) { 
            statusResponse, greetings = true, value
        }
    }
    
    if statusResponse { 
        responseSlack(greetings + " <@" + user + ">", channel)
    } else if(strings.Contains(text, "kabar")) {
        responseSlack("Kabar baik.", channel)
    } else {
        responseSlack("Ada apa ya mention saya <@" + user + ">", channel)
    }
}

func responseSlack(text, channel_id string){
    payload := map[string] interface{} {
                "text": text,
                "channel": channel_id,
                "username": configuration.Name,
                "icon_emoji": ":hammer_and_pick:",
            }

    bytesRep, err := json.Marshal(payload)
    if err != nil {
        log.Println("There are some errors: ", err)
    }

    req, _ := http.NewRequest("POST", configuration.SlackURL, bytes.NewBuffer(bytesRep))

    req.Header.Add("Authorization", configuration.Token)
    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("cache-control", "no-cache")

    res, _ := http.DefaultClient.Do(req)

    defer res.Body.Close()
    body, _ := ioutil.ReadAll(res.Body)

    log.Println(res)
    log.Println(body)
}

func processNumber(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        command := r.FormValue("command")
        
        text, username, command := r.FormValue("text"), r.FormValue("user_name"), r.FormValue("command")
        response_url, channelname, channelid := r.FormValue("response_url"), r.FormValue("channel_name"), r.FormValue("channel_id")
        
        fmt.Println(text, username, command, channelname, channelid, response_url)
        
        split := strings.Split(text, " ")

        num1, err := strconv.Atoi(split[0])
        num2, err := strconv.Atoi(split[1])
        var result int

        if command == "/add-number" {
            result = num1 + num2
        } else if(command == "/subs-number") {
            result = num1 - num2
        } else if(command == "/multiple-number") {
            result = num1 * num2
        } else {
            result = num1 / num2
        }
        
        fmt.Println(result, err)

        responseSlack("Result: " + strconv.Itoa(result), channelid)
    } else {
        fmt.Println(w, "Execution %q %s", r.Method, html.EscapeString(r.URL.Path))
        fmt.Fprintf(w, "Execution %q %s", r.Method, html.EscapeString(r.URL.Path))

        log.Println("Execution ", r.Method, html.EscapeString(r.URL.Path))
    }
}

/**
  * main function
  */

func main() {
    mux := http.NewServeMux()

    errConfig := gonfig.GetConf("config/config.json", &configuration)
    if errConfig != nil {
        log.Println(errConfig)
    }

    log.Println("Rest API - Listening on port:", configuration.Port)

    mux.HandleFunc("/", index)
    mux.HandleFunc("/event-subscriptions-fqrt3gddk", eventSubscriptions)
    mux.HandleFunc("/add-number", processNumber)
    mux.HandleFunc("/subs-number", processNumber)
    mux.HandleFunc("/multiple-number", processNumber)
    mux.HandleFunc("/divide-number", processNumber)

    log.Println(http.ListenAndServe(":" + configuration.Port, mux))
}