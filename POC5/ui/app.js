// ui.js

const socket = new WebSocket("ws://localhost:8080/ws");

socket.addEventListener("message", (event) => {
    console.log(event.data)
    const messageContainer = document.getElementById("message-container");
    const messageElement = document.createElement("div");
    messageElement.innerText = event.data;
    messageContainer.appendChild(messageElement);
});

