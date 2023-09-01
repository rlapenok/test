# Telegram bot with gRPC server
Under the hood of this bot, a gRPC server is running, which processes the requests described in the protofile ```./proto/tgbot.proto``` and sends it to everyone who subscribes to its updates
## Dependencies
You must have Docker and make installed
## ENV install
Specify ```ENV TOKEN``` in the docker file (you can get it from BotFather) and ```ENV PORT``` on which the GRPS server will be launched
## Usage
### Start
```
make all
```
### Stop
```
make stop
```
### Sending messages
Send messages to ```localhost:8080``` using gRPC using PostMan, or write your own gRPC client that will send protofiles
## 
Pull requests are welcome
