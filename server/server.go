package server

import (
	"os"
	"fmt"
	"github.com/carsonsx/nathttpd/const"
	"github.com/carsonsx/nathttpd/message"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"github.com/carsonsx/nathttpd/conf"
	"github.com/kimiazhu/log4go"
	"io/ioutil"
	"bytes"
	"net/http"
)

// PrintServerAndExit will print our version and exit.
func PrintServerAndExit() {
	fmt.Printf("nats-server version %s\n", constant.VERSION)
	os.Exit(0)
}


func LogError(err error, msg string) {
	if err != nil {
		log4go.Error("%s: %s", msg, err)
	}
}


func getResponseOfError(err error) string {
	return "{code:-1,msg:\"" + err.Error() + "\"}"
}

func getResponse(resp *http.Response, err error) string {
	if err != nil {
		return getResponseOfError(err)
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return getResponseOfError(err)
		} else {
			return string(body)
		}
	}
}

func invoke(request message.HttpRequestMessage) string {
	log.Print(request)
	var bodyReader *bytes.Reader
	bodyType := "application/json;charset=UTF-8"
	if request.JsonData != "" {
		bodyReader = bytes.NewReader([]byte(request.JsonData))
	} else if len(request.FormData) > 0 {
		bodyReader = bytes.NewReader([]byte(request.FormData.Encode()))
		bodyType = "application/x-www-form-urlencoded"
	}
	var req *http.Request
	var response string
	var err error
	if bodyReader == nil {
		req, err = http.NewRequest(request.Method, request.Url, nil)
	} else {
		req, err = http.NewRequest(request.Method, request.Url, bodyReader)
	}
	if err != nil {
		response = getResponseOfError(err)
	} else {
		req.Header.Set("Content-Type", bodyType)
		response = getResponse(http.DefaultClient.Do(req))
	}

	log.Println(response)

	if err != nil && request.Error != nil {
		return invoke(*request.Error)
	}
	if err == nil && request.Next != nil {
		return invoke(*request.Next)
	}

	return response
}

func Run() {

	conn, err := amqp.Dial(conf.MQConf.Url)
	LogError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	LogError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	LogError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		conf.MQConf.ReqQueue, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	LogError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var response string
			var request message.HttpRequestMessage
			err = json.Unmarshal(d.Body, &request)
			if err != nil {
				response = getResponseOfError(err)
			} else {
				response = invoke(request)
			}
			if d.ReplyTo != "" {
				err = ch.Publish(
					"",        // exchange
					d.ReplyTo, // routing key
					false,     // mandatory
					false,     // immediate
					amqp.Publishing{
						ContentType:   "text/plain",
						CorrelationId: d.CorrelationId,
						Body:          []byte(response),
					})
				LogError(err, "Failed to publish a message")
			}
			d.Ack(false)
		}
	}()

	log4go.Info(" [*] Awaiting HTTP requests")
	<-forever
}
