# camera-go
Camera and microphone via websockets.

I'm using the Golang [Gorilla](https://github.com/gorilla/websocket/) websocket framework for the server,
until replaced by http/2.

DISCLAIMER: This is currently a personal experiment in progress.
I'm very open for comments. Also, this is one of my first Golang experiences.

Usage, assuming you have your Go environment prepared:

```
git clone https://github.com/nkoster/camera-go
cd camera-go
go get github.com/gorilla/websocket
go build
./camera-go -html ./www/
````

or

```
go run *.go -html ./www/
```

Open http://localhost:8086 on, for example, two laptops with cam and mic, and exchange IDs.
