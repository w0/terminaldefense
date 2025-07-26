const ws = new WebSocket("ws://" + document.location.host + "/ws");

ws.onopen = function (event) {
  console.log("we dem bois");
};

ws.onmessage = function (event) {
  document.getElementById("wsget").textContent = event.data;
};
