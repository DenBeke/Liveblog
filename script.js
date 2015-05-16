var sock = null;
var wsuri = "ws://localhost:1234";

window.onload = function() {

	sock = new WebSocket(wsuri);

	sock.onopen = function() {
		console.log("connected to " + wsuri);
	}

	sock.onclose = function(e) {
		console.log("connection closed (" + e.code + ")");
	}

	sock.onmessage = function(e) {
		var json = JSON.parse(e.data);
		$('<li>' + json.Content + ' <span class="timestamp">' + timeAgo(parseInt(json.Time)) + '</span></li>').hide().prependTo('#messages').fadeIn(1000);
	}
};


function timeAgo(time){
  var units = [
	{ name: "second", limit: 60, in_seconds: 1 },
	{ name: "minute", limit: 3600, in_seconds: 60 },
	{ name: "hour", limit: 86400, in_seconds: 3600  },
	{ name: "day", limit: 604800, in_seconds: 86400 },
	{ name: "week", limit: 2629743, in_seconds: 604800  },
	{ name: "month", limit: 31556926, in_seconds: 2629743 },
	{ name: "year", limit: null, in_seconds: 31556926 }
  ];
  var diff = (new Date() - new Date(time*1000)) / 1000;
  if (diff < 5) return "now";

  var i = 0, unit;
  while (unit = units[i++]) {
	if (diff < unit.limit || !unit.limit){
	  var diff =  Math.floor(diff / unit.in_seconds);
	  return diff + " " + unit.name + (diff>1 ? "s" : "");
	}
  };
}