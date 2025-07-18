# Leaflets and Websockets in Javascript

## Leaflet

### Setup
- Include the leaflet CSS in <head>
```
<link rel="stylesheet" href="https://unpkg.com/leaflet@1.9.4/dis/leaflet.css"
    integrity="sha256-p4NxAoJBhIIN+hmNHrzRCf9tD/miZyoHS5obTRR9BMY="
    crossorigin=""/>
```
- Inlcude leaflet js file after the CSS
```
<script src="https://unpkg.com/leaflet@1.9.4/dist/leaflet.js"
    integrity="sha256-20nQCchB9co0qIjJZRGuk2/Z9VM+kNiyxNV1lvTlZBo="
    crossorigin=""></script>
```

- put a div element with certain id where the map should be
```
<div id="map"></div>
```

- map must have a defined height
set it in the CSS
#map {height: 1080px;}

### Creating a Map
- setting a view to coords and zoom level
```
var map = L.map('map').setView([51.505, -0.09], 13);
```

- all mouse and touch interactions enabled by default
- Set a tile layer
    - set the url template for tile image, attribution text and maximum zoom
    - OpenStreetMap tile layer can be used to do this
```
L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 19,
    attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>'
}).addTo(map);
```
- provider agnostic means we can use any tile provider not just open street maps

### Markers, Circles and Polygons
- these are the other things we can add to the map other than tile layers
```
var marker = L.marker([51.5, -0.89]).addTo(map);
```
- adding a circle (need to specify radius)
```
var circle = L.circle([61.23, 23,01], {
    color: 'red',
    fillColor: '#f03',
    fillOpacity: 0.5,
    radius: 500,
}).addTo(map);
```
- adding a polygon using an arr of arr of coords
```
var polygon = L.polygon([
    [p1_x, p1_y],
    [p2_x, p2_y],
    [p3_x, p3_y]
])addTo(map);
```

### Working with Popups
- .bindPopup("popup text").openPopup
- this works on all three object markers, circle and polygons
- openPopup only works on markers and opens popup immediately
- can also use popUp as layer for standalone
```
var popup = L.popup().setLatLng([34.34, -1.33]).setContent("abcd").openOn(map);
```

### Event Handling
- you can subscribe to events like clicks on a marker or zoom
- a callback function can be registered to handle the event
```
function onMapClick(e) {
    alert("You clicked the map at" + e.LatLng)
}
map.on('click', onMapClick)
```
- each object has its own set of events
- first arg of listener func is event object


## Websockets in Javascript Clients

### Creating a Websocket Object
- creating a websocket object will automatically try to open a connection to the server
    ```
    websocket = new WebSocket(url, protocols);
    ```
    - url: represents the url endpoint of the servers websocket endpoint
        - has  to be using protocol wss:// or ws://
        - recent browser verions also allow https:// or http://
    - protocols: single protocol string or an array of protocol strings.
      indicates sub protocols so a single server can implement 
      different types of interactions
    - will thow a SecurityError if destination doesnt allow access
        - may happen if using an insecure agent

### Connection errors
- first an error event is sent to the Websocket object so handlers can run
- then a close event object is sent with reason for connection closing

```
const exampleWebSocket = new WebSocket(
    "wss://myserver.com/wsendpoint",
    "protocolOne"
)
```
- the state of exampleWebSocket.readyState = CONNECTION becomes OPEN 
  when the connection is ready
- exampleWebSocket.protocol will tell us which protocol is being 
  selected to respond by the server

### Sending data to the server
```
exampleSocket.send("Heres some data")
```
- data can be sent as a string Blob or ArrayBuffer
- calling send immediately after establishing socket doesnt garuntee the send
    - since setup is async
    - to make sure send works we have to register the send to the event handler
      for the onopen event
    ```
    exampleSocket.onopen = (event) => {
        exampleSocket.send("heres the message")
    }

- we can send over JSON
```
function sendText() {

    const msg = {
        type: "message", 
        text: document.getEleemntById("text").value,
        id: clientID,
        date: Date.now(),
    };

    exampleSocket.send(JSON.stringify(msg))
    ```
### Recieving messages from the server
```
exapmleSocket.onmessage = (event) => {
    console.log(event.data)
}
```
- recieve the message from the onmessage event
- JSON data recieve from server 
    - different types of possible messages that can be sent
        - ws handshake
        - message
    ```
    exampleSocket.onmessage = (event) => {

        const f = document.getElementById("chat").contentDocument;
        let text = ""
        const msg = JSON.parse(event.data)
        const time = new Date(msg.date)
        const timeStr = time.toLocaleTimeString();

        switch (msg.type) {

            case "id":
                clientId = msg.id
            case "username":
                text = `user $(msg.name)`
                break;
            case "message":
                text = `$(msg.name) : $(msg.text)`
        }
        if (text.length) {
            f.write(text)
            document.getElementById("chat").contentWindow.scrollByPages(1)
        }
### Closing the connection
```
exampleSocket.close();
```
- close the connetion
- check the attribute exampleSocket.bufferedAmount before closing to see if any
  data is buffered


