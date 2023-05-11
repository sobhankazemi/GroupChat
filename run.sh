#!/bin/bash
docker container  run -d --hostname rabbit --name rabbit1 -p 15672:15672 -p 5672:5672 rabbitmq:3-management
cd Chat
go run *.go &
cd ../History
go run *.go &
