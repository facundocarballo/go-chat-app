## Server - Golang
To create a Web Socket, first we need a http server listening for http request.

### Dependencies
#### Web Socket
```bash
    go get github.com/gorilla/websocket
```
#### Handlers - CORS
```bash
    go github.com/gorilla/handlers
```
### Server code
```golang
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.HandleWebSocket(w, r, database)
	})

    corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowCredentials(),
    )

    go websocket.HandleMessages()

    println("Server listening on port" + SERVER_PORT + " ...")
    if err := http.ListenAndServe(SERVER_PORT, corsMiddleware(http.DefaultServeMux)); err != nil {
	    panic(err)
    }
```
Here we have just one endpoint.

> '/ws'
This endpoint will upgrade the connection between the client and the server, to use Web Sockets.
Will be looping, looking for messages of this particular client.

### CORS Middleware
We need this kind of policy to allow or client side to make the Web Socket connection with the server.

### Gorutine - Handle Messages
This will be a function running in the background of our main server, looking for messages that any client connected with our Web Socket could send to us.

And we as a server, will have to handle that. Reading the message and them sending to the corresponded client connected.

### Web Sockets Handlers
To read all the Web Socket methods that we use in this app, you can go to the **websocket** folder and here you will find the implementation of each method.