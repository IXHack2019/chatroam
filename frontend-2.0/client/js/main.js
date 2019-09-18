let local_device_id;
let local_username;
let ws;

$(document).ready(function() {
   local_device_id = parseInt(Math.random().toString().replace(".")).toString()+Date.now().toString();
   local_username = "CHATROAM";
   ws = new WebSocket("ws://localhost:8888/connect");
   ws.onopen = function (evt) {
      var request = {
         "type": 0,
         "data": {
            "deviceID": local_device_id
         }
      }
      sendRequestWithPosition(request);
   }
   ws.onclose = function (evt) {
      ws = null;
   }
   ws.onmessage = function (evt) {
      console.log("RESPONSE: " + evt.data);
      response = JSON.parse(evt.data);
      switch(response.type) {
         case 1: 
            writeMessage(response.data, response.data.deviceID == local_device_id ? true : false);
            break;
         case 0:
            local_username = response.username;
            console.log("New username set: " + local_username);
            break;
      } 
   }
   ws.onerror = function (evt) {
      console.log("ERROR: " + evt.data);
   }

   $(document).on('keypress',function(e) {
      if(e.which == 13) {
         writePersonalMessage();
      }
   });
});

function sendRequestWithPosition(request) {
   getPosition()
      .then(function(coords) {
         request.data.lat = coords[0];
         request.data.lng = coords[1];
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
            "deviceID": local_device_id
         }
      }
      sendRequestWithPosition(request);
   }
}

function writeMessage(data, personal) {
   let html = "<span class='message bubble'>\
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
         "deviceID": local_device_id,
         "username": local_username
      }
   };
   // writeMessage(data, true);
   $("#input-box input").val("");
   ws.send(JSON.stringify(request));
   return false;
}