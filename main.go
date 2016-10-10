package main

import (
	"github.com/carsonsx/nathttpd/conf"
	"flag"
	"fmt"
	"os"
	"github.com/carsonsx/nathttpd/server"
)

var usageStr = `
Usage: nathttpd [options]

Server Options:
    -u, --url <url>                  Message queue connnection url (default: amqp://0.0.0.0:5672)
    --request_queue <name>           Http request queue name (default: nat_http_request)
    --response_queue <name>          Http response queue name (default: nat_http_response)

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


func main() {

	var showVersion bool

	flag.StringVar(&conf.MQConf.Url, "url", conf.MQConf.Url, "Http message queue url.")
	flag.StringVar(&conf.MQConf.Url, "u", conf.MQConf.Url, "Http message queue url.")
	flag.StringVar(&conf.MQConf.ReqQueue, "request_queue", conf.MQConf.ReqQueue, "Http request message queue name.")
	flag.StringVar(&conf.MQConf.ResQueue, "response_queue", conf.MQConf.ResQueue, "Http respnse message queue name.")
	flag.BoolVar(&showVersion, "version", false, "Print version information.")
	flag.BoolVar(&showVersion, "v", false, "Print version information.")
	flag.Usage = usage
	flag.Parse()

	// Show version and exit
	if showVersion {
		server.PrintServerAndExit()
	}

	server.CreateRabbitQueue(conf.MQConf.Url, conf.MQConf.ReqQueue)
	server.CreateRabbitQueue(conf.MQConf.Url, conf.MQConf.ResQueue)

	server.Run()
}
