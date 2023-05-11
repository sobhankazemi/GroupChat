from golang:1.20.4-alpine3.17

MAINTAINER sobhan_kazemi@comp.iust.ac.ir

RUN mkdir app

WORKDIR /app

COPY . . 

RUN go mod tidy
EXPOSE 8080
CMD ["go" , "run" , "."]


