let local_device_id;
let local_username;
let ws;
let connecting;
let shifted;

$.fn.textWidth = function(){
    var calc = '<span style="display:none">' + $(this).text() + '</span>';
    $('body').append(calc);
    var width = $('body').find('span:last').width();
    $('body').find('span:last').remove();
   return width;
};

$.fn.marquee = function(args) {
   var that = $(this);
   var textWidth = that.textWidth(),
       offset = that.width(),
       width = offset,
       css = {
           'text-indent' : that.css('text-indent'),
           'overflow' : that.css('overflow'),
           'white-space' : that.css('white-space')
       },
       marqueeCss = {
           'text-indent' : width,
           'overflow' : 'hidden',
           'white-space' : 'nowrap'
       },
       args = $.extend(true, { count: -1, speed: 1e1, leftToRight: false }, args),
       i = 0,
       stop = textWidth*-1,
       dfd = $.Deferred();

   function go() {
       if(!that.length) return dfd.reject();
       if(width == stop) {
           i++;
           if(i == args.count) {
               that.css(css);
               return dfd.resolve();
           }
           if(args.leftToRight) {
               width = textWidth*-1;
           } else {
               width = offset;
           }
       }
       that.css('text-indent', width + 'px');
       if(args.leftToRight) {
           width++;
       } else {
           width--;
       }
       setTimeout(go, args.speed);
   };
   if(args.leftToRight) {
       width = textWidth*-1;
       width++;
       stop = offset;
   } else {
       width--;            
   }
   that.css(marqueeCss);
   go();
   return dfd.promise();
};

$(document).ready(function() {
   local_device_id = parseInt(Math.random().toString().replace(".")).toString()+Date.now().toString();
   buildConnection();
   $(document).on('keyup keydown',function(e) {
      if(e.which == 13) {
         if(!e.shiftKey) {
            writePersonalMessage();
         } else {
            advertise();
         }
      }
   });
});

function buildConnection() {
   ws = new WebSocket("wss://dumtard.com/connect");
   ws.onopen = function (evt) {
      connect();
   }
   ws.onclose = function (evt) {
      console.log("CLOSED");
      // reconnect();
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
   // if(ws.readyState === WebSocket.CLOSED) {
   //    buildConnection();
   // }
   if(ws.readyState === WebSocket.OPEN) {
      clearInterval(connecting);
   }
   console.log("send");
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
      if(navigator.geolocation) {
         resolve([43.6425829792536, -79.38512721871393]);
      } else {
         navigator.geolocation.getCurrentPosition(function (position) {
            resolve([position.coords.latitude, position.coords.longitude]);
         }, throwPositionError);
      }
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
   let html = "<div class='";
   if(personal) {
      html += "personalbox";
      color = "#D6EAF8";
   } else {
      html += "senderbox"
   }
   html += "'><div class='username'>" + data.username + "</div> \
      <span class='message bubble' style='background-color: " + color + "'>\
      <div class='text'>" + data.msg + "</span></div></div>";
   let $msg = $(html);
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
   if(id === undefined) return "rgba(100,100,100,0.2)"
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

function advertise(duration=5000) {
   let text = $("#text-box").val();
   if(text == "") return;
   $("#megaphone span").html(text);
   $("#megaphone").addClass("active");
   setTimeout(function() {
      $("#megaphone").removeClass("active");
   }, duration)
}