package main

import (
	"bytes"
	"encoding/json"
	"github.com/streadway/amqp"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/carsonsx/nathttpd/conf"
	"github.com/carsonsx/nathttpd/message"
	"flag"
	"fmt"
	"os"
)

var usageStr = `
Usage: httpnat [options]

Server Options:
    -a, --addr <host>                Bind to host address (default: 0.0.0.0)
    -p, --port <port>                Use port for clients (default: 4222)

Logging Options:
    -l, --log <file>                 File to redirect log output
    -T, --logtime                    Timestamp log entries (default: true)
    -s, --syslog                     Enable syslog as log method
    -r, --remote_syslog <addr>       Syslog server addr (udp://localhost:514)
    -D, --debug                      Enable debugging output
    -V, --trace                      Trace the raw protocol
    -DV                              Debug and trace

Authorization Options:
        --user <user>                User required for connections
        --pass <password>            Password required for connections

Common Options:
    -h, --help                       Show this message
    -v, --version                    Show version
`

// usage will print out the flag options for the server.
func usage() {
	fmt.Printf("%s\n", usageStr)
	os.Exit(0)
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
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

	if request.Next != nil {
		return invoke(*request.Next)
	}

	return response
}

func main() {

	flag.StringVar(&conf.MQConf.MQUrl, "mqurl", conf.MQConf.MQUrl, "Http message queue url.")
	flag.StringVar(&conf.MQConf.MQName, "mqname", conf.MQConf.MQName, "Http message queue name.")
	flag.Usage = usage
	flag.Parse()

	conn, err := amqp.Dial(conf.MQConf.MQUrl)
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	FailOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		conf.MQConf.MQName, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	FailOnError(err, "Failed to register a consumer")

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
				FailOnError(err, "Failed to publish a message")
			}
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Awaiting RPC requests")
	<-forever
}
