const localuser = "jayrsawal";
let ws;

$(document).ready(function() {
   ws = new WebSocket("ws://localhost:8888/connect");
   ws.onopen = function (evt) {
      var request = {
         "type": 0,
         "data": {
            "deviceID": localuser
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
      if(response.type == 1) {
        writeMessage(response.data, response.data.username == localuser ? true : false);
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
      });
   });
}

function writeMessage(data, personal) {
   let html = "<span class='message bubble'>\
      <span class='username'>" + data.username + "</span> \
      <span class='text'>" + data.msg + "</span></span>";
   let $msg = $(html);
   if(personal) {
      $msg.addClass("personal");
   }
   $("#board").append($msg);

}

function writePersonalMessage() {
   let msg = $("#input-box input").val();
   if(msg == "") return false;
   let request = {
      "type": 1,
      "data": {
         "msg": msg,
         "deviceID": localuser,
         "username": localuser
      }
   };
   // writeMessage(data, true);
   $("#input-box input").val("");
   ws.send(JSON.stringify(request));
   return false;
}