import Vue from "nativescript-vue";
import Home from "./components/Home";
var WS = require('nativescript-websockets');
 
var mySocket = new WS("ws://echo.websocket.org",{protocols: [/* 'chat', 'video' */], timeout: 6000, allowCellular: true, headers: { 'Authorization': 'Basic ...' }});
mySocket.on('open', function(socket) { console.log("Hey I'm open"); socket.send("Hello"); });
mySocket.on('message', function(socket, message) { console.log("Got a message", message); });
mySocket.on('close', function(socket, code, reason) { console.log("Socket was closed because: ", reason, " code: ", code); });
mySocket.on('error', function(socket, error) { console.log("Socket had an error", error);});
 

new Vue({

    template: `
        <Frame>
            <Home/>
        </Frame>`,

    components: {
        Home
    }
}).$start();
