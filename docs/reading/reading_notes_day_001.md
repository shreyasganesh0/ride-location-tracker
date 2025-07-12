# WebSockets

## Introduction
- Websockets are HTTP connections upgraded via 103 response
- Websocket servers are usually their own server so reverse proxies
  detect the protocol and route them after preprocessing the handshakes

## Websocket Handshake
- Server listens for socket connections on the TCP socket
    - Usually listens on 443 or 80 for HTTP or may be blocked by WAF and proxies
    - Browsers usually expect secure connections for ws
- Handshake is the bridge between the HTTP and WS protocol
- details of connection are negotiated during handshake
- the request-uri is not defined by the spec so it can be called anything 
  like /chat or /gaming for the app to handle multiple ws apps

- Client Handshake Request
    - Starts the websocket handshake
    - must be 1.1 version HTTP or greater
    - contacts endpoint with headers
        - Connection: Upgrade
        - Sec-WebSocket-Key: 
        - Sec-WebSocket-Version: 13
    - client from browsers send an Origin header which can be used for
      validation of the client

- Server Handshake Response
    - send back response
        - 101 Switching Protocols
        - Upgrade: Websocket
        - Connection: Upgrade
        - Sec-WebSocket-Accept
            - Derived from Sec-WebSocket-Key
            - concats with its magic string 
            - encode the SHA1 hashs with base64 
    - if the version header not recoginized it has to send the 
      correct versions

## Gorrila Websockets Chat
- Server has two types Client and Hub
- Client communicate between connection and Hub
- Hub tracks registered clients and broadcasts the messages to them
- 1 goroutine for Hub and 2 for each Client
- goroutine communication happens via channels
- websocket reads messages from the Client buffer and sends messages to the Client
    - the Client then sends those messages to the Hub from broadcasting
- When the Hub unregisters a client
    - deletes client pointer from the map
    - closees the clients send channel to say no more messages sent to that client
    - the Hub loops over all registered clients and sends message to the 
      clients send channel
    - if the send buffer is full it assumes the client is dead and unregisters it
- The Client
    - upgrades the HTTP connection via the handler
    - creates a client and registers it to the Hub
    - messages are transffered from the client send buf to the websocket connection
    - the handler calls read pump which transfers inbound messages from the websocket
      to the hub
- each connection supports 1 one concurrent reader and writer


    

