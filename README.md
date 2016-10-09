# nathttpd

nathttpd is http request forward service in LAN via rabbitmq message. Then you can invoke the LAN http service from internet.
In some cases, such as file server, we need to invoke delete file from web server, but the web server is hosting in internet. 
we can't invoke the serivce directly, so we need httpnatd.
