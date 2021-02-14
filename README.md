Simple Broker
============

#Getting Started
run simple-broker with docker

docker run hub.tisserv.net/tislib/simple-broker -p 8712:8712

Producer:

send message to topic named queue-1
curl -XPOST "https://127.0.0.1:8712/queue-1"

Consumer:
curl "https://127.0.0.1:8712/queue-1"


<br/>
<br/>
<br/>
<br/>
<br/>


Topic Design notes:[link](design/topic-design.md)
