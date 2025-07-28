const term = new Terminal();

term.open(document.getElementById("terminal"));

const ws = new WebSocket("ws://" + document.location.host + "/ws");

let currentLine = "";

ws.onmessage = function (event) {
  try {
    const p = JSON.parse(event.data);

    const d = atob(p);
    console.log(d);
  } catch (e) {
    term.write(event.data);
  }
};

ws.onclose = function (event) {
  console.log(event.reason);
};

term.onKey(({ key, domEvent }) => {
  console.log(key.charCodeAt(0));
  if (key.charCodeAt(0) == 13) {
    ws.send(
      JSON.stringify({
        type: "CMD",
        data: {
          command: currentLine.trimStart(),
        },
      }),
    );
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
