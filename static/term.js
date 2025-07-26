const ws = new WebSocket("ws://" + document.location.host + "/ws");

ws.onopen = function (event) {
  console.log("we dem bois");
};

ws.onmessage = function (event) {
  document.getElementById("tty").textContent += "\n" + event.data;
};

ws.onclose = function (event) {
  console.log(event.reason);
};
