# GroupChat

## Description
this is a chat application written in golang using rabbitMQ as message broker and Gorrila as websocket
## How It Works
open a websocket connection on "localhost:8080/{room id}"
message will broadcast to every body who is registered to the room
you can access to rooms chat history by an app running on port 8050 (localhost:8050/api/history/{room id})
## message sample
```js
 {
    "username" : "sobhankazemi",
    "user_id" : 1,
    "message" : "hello github"
 }
```
## How To Run :
```bash
chmod 777 run.sh
./run.sh
```
