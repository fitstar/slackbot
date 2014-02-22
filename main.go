package main

import (
	"github.com/gorilla/schema"
	"github.com/trinchan/slackbot/robots"
	"io"
	"log"
	"net/http"
	"strconv"
	"fmt"
)

func main() {
	http.HandleFunc("/slack", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err == nil {
			command := new(robots.SlashCommand)
			decoder := schema.NewDecoder()
			err := decoder.Decode(command, r.PostForm)
			if err != nil {
				log.Println("Couldn't parse post request:", err)
			}
			robot := GetRobot(command)
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			if robot != nil {
				io.WriteString(w, robot.Run(command))
			} else {
				io.WriteString(w, "No robot for that command yet :(")
			}
		}
	})
	
	http.HandleFunc("/slack-webhook", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err == nil {
			message := new(robots.SlashCommand)
			decoder := schema.NewDecoder()
			err := decoder.Decode(message, r.PostForm)
			if err != nil {
				log.Println("Couldn't parse post request:", err)
			}
			//robot := GetRobot(command)
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			//io.WriteString(w, robot.Run(command))
			
			response := new(robots.IncomingWebhook)
			response.Channel = message.Channel_ID
			response.Text = fmt.Sprintf("I heard: @%s", message.Text)
			response.Unfurl_Links = true
			response.Parse = "full"
			robots.MakeIncomingWebhookCall(response)
		}
	})

	StartServer()
}

func StartServer() {
	port := robots.Config.Port
	log.Printf("Starting HTTP server on %d", port)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Fatal("Server start error: ", err)
	}
}

func GetRobot(command *robots.SlashCommand) (robot robots.Robot) {
	if RobotInitFunction, ok := robots.Robots[command.Command]; ok {
		return RobotInitFunction()
	} else {
		return nil
	}
}
