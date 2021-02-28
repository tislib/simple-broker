General Idea

All operations are compatible to rest standards

Producer sends http request to topic (path = topic)

request is forwarded to consumer

broker has predefined cache

broker uses go channels to create communication between consumer and producer

sub paths is used to create fan in and fon out techniques, etc.

example:

POST /users {...id:1}

POST /users {...id:2}

POST /users {...id:3}

POST /users {...id:4}

POST /users/1/orders {...id:5}

POST /users/1/orders {...id:6}

POST /users/1/orders {...id:7}

POST /users/1/orders {...id:8}

// call these endpoints in parallel

GET -H "Accept-group: group1" /users            (receives record {id:1}, {id:3}, {id:5})

GET -H "Accept-group: group1" /users            (receives record {id:2}, {id:4}, {id:6})

GET -H "Accept-group: group1" /users/1/orders   (receives record {id:7})

GET -H "Accept-group: group1" /users/1/orders   (receives record {id:8})

GET -H "Accept-group: group2" /users/1/orders   (receives record {id:5}, {id:7})

GET -H "Accept-group: group2" /users/1/orders   (receives record {id:6}, {id:8})

#internal design

producer sends request to http server (broker)

broker accepts packet and checks if the corresponding queue is already exists if not initializes queue producer channel and its operator

consumer sends request to http server (broker)

broker accepts packet and check if the corresponding queue is already exists if not initializes queue producer channel and its operator

producer sends message to queue specific channel

consumer reads message from queue specific channel

we have one operator channel which is responsible to redirect packets from producer to consumer

each consumer has an index, and consumer channel is opened with specific index

operator part will keep data in array and channel, channel is good for performance and array is good for storage



