function autoScrollTextArea(textAreaId) {
  const textArea = document.getElementById(textAreaId);
  if (textArea) {
    textArea.scrollTop = textArea.scrollHeight;
  }
}

const ws = new WebSocket("ws://" + document.location.host + "/ws");

const prompt = document.getElementById("prompt");

prompt.addEventListener("keypress", function (event) {
  if (event.key === "Enter") {
    const command = event.target.value;
    console.log(command);
    event.target.value = "";

    ws.send(command);
  }
});

ws.onopen = function (event) {
  console.log("we dem bois");
};

ws.onmessage = function (event) {
  document.getElementById("tty").textContent += "\n" + event.data;
  autoScrollTextArea("tty");
};

ws.onclose = function (event) {
  console.log(event.reason);
};
