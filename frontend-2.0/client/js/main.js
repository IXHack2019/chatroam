let local_device_id;
let local_username;
let ws;
let connecting;

$(document).ready(function() {
   local_device_id = parseInt(Math.random().toString().replace(".")).toString()+Date.now().toString();
   buildConnection();
   $(document).on('keypress',function(e) {
      if(e.which == 13) {
         writePersonalMessage();
      }
   });
});

function buildConnection() {
   ws = new WebSocket("ws://dumtard.com:80/connect");
   ws.onopen = function (evt) {
      connect();
   }
   ws.onclose = function (evt) {
      console.log("CLOSED");
      reconnect();
   }
   ws.onmessage = function (evt) {
      console.log("RESPONSE: " + evt.data);
      response = JSON.parse(evt.data);
      switch(response.type) {
         case 1: 
            writeMessage(response.data, response.data.deviceId == local_device_id ? true : false);
            break;
         case 0:
            clearInterval(connecting);
            if(local_username === undefined) local_username = response.username;
            console.log("New username set: " + local_username);
            break;
      } 
   }
   ws.onerror = function (evt) {
      console.log("ERROR: " + evt.data);
   }
}

function connect() {
   if(ws.readyState === WebSocket.CLOSED) {
      buildConnection();
   }
   if(ws.readyState === WebSocket.OPEN) {
      clearInterval(connecting);
   }
   var request = {
      "type": 0,
      "data": {
         "deviceId": local_device_id
      }
   }
   sendRequestWithPosition(request);
}

function reconnect() {
   connect();
   clearInterval(connecting);
   connecting = setInterval(connect, 1000);
}

function sendRequestWithPosition(request) {
   getPosition()
      .then(function(coords) {
         request.data.lat = coords[0];
         request.data.lon = coords[1];
         console.log("SEND: " + JSON.stringify(request));
         ws.send(JSON.stringify(request));
      });
}

function getPosition() {
   return new Promise(function(resolve, reject) {
      navigator.geolocation.getCurrentPosition(function (position) {
         resolve([position.coords.latitude, position.coords.longitude]);
      }, throwPositionError);
   });
}

function throwPositionError(error) {
  switch(error.code) {
    case error.PERMISSION_DENIED:
      console.log("User denied the request for Geolocation.");
      break;
    case error.POSITION_UNAVAILABLE:
      console.log("Location information is unavailable.");
      break;
    case error.TIMEOUT:
      console.log("The request to get user location timed out.");
      break;
    case error.UNKNOWN_ERROR:
      console.log("An unknown error occurred.");
      break;
  }
   ws.onopen = function (evt) {
      var request = {
         "type": 0,
         "data": {
            "deviceId": local_device_id
         }
      }
      sendRequestWithPosition(request);
   }
}

function writeMessage(data, personal) {
   let color = idToRGB(data.deviceId, 0.5);
   if(personal) {
      color = "white";
   }
   let html = "<span class='message bubble' style='background-color: " + color + "'>\
      <div class='username'>" + data.username + "</div> \
      <div class='text'>" + data.msg + "</span></div>";
   let $msg = $(html);
   if(personal) {
      $msg.addClass("personal");
   } 
   $("#board").append($msg);
   window.scrollTo(0,document.body.scrollHeight);
}

function writePersonalMessage() {
   let msg = $("#input-box input").val();
   if(msg == "") return false;
   let request = {
      "type": 1,
      "data": {
         "msg": msg,
         "deviceId": local_device_id,
         "username": local_username
      }
   };
   // writeMessage(data, true);
   $("#input-box input").val("");
   ws.send(JSON.stringify(request));
   return false;
}

function idToRGB(id, alpha=1) {
   return hexToRGB(intToRGB(parseInt(id)), alpha);
}

function intToRGB(i){
    var c = (i & 0x00FFFFFF)
        .toString(16)
        .toUpperCase();

    return "#" + "00000".substring(0, 6 - c.length) + c;
}

function hexToRGB(hex, alpha=1.0) { 
  // Expand shorthand form (e.g. "03F") to full form (e.g. "0033FF")
  var shorthandRegex = /^#?([a-f\d])([a-f\d])([a-f\d])$/i;
  hex = hex.replace(shorthandRegex, function(m, r, g, b) {
    return r + r + g + g + b + b;
  });
  var result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
  return result ? "rgba(" + parseInt(result[1],16) + ", " + parseInt(result[2], 16) 
      + ", " + parseInt(result[3], 16) + ", " + alpha + ");" : "#9693cf"
}