const term = new Terminal();

term.open(document.getElementById("terminal"));

const ws = new WebSocket("ws://" + document.location.host + "/ws");

let currentLine = "";

ws.onmessage = function (event) {
  term.write(event.data);
};

ws.onclose = function (event) {
  console.log(event.reason);
};

term.onKey(({ key, domEvent }) => {
  console.log(key.charCodeAt(0));
  if (key.charCodeAt(0) == 13) {
    term.write("\n");
    ws.send(currentLine);
    currentLine = "";
  }

  if (key.charCodeAt(0) == 127) {
    if (currentLine.length > 0) {
      currentLine = currentLine.slice(0, -1);
      term.write("\x1b[D\x1b[K"); //move left and delete
    }
    return;
  }
  currentLine += key;
  term.write(key);
});
