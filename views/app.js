let ws;
let currentRoom = "lobby";
let username = null;
let roomUsers = [];

const usernameModal = document.getElementById("usernameModal");
const usernameInput = document.getElementById("usernameInput");
const usernameConfirm = document.getElementById("usernameConfirm");
const userInfo = document.getElementById("userInfo");
var currentRoomName = document.getElementById("currentRoomName");

/* ------------------------ CONNECT ------------------------ */
function connect(username) {
    ws = new WebSocket("ws://localhost:8080/ws");

    ws.onopen = () => {
        addSystemMessage("Connected to server.");

        ws.send(JSON.stringify({
            command: "/username",
            args: [username]
        }));
    };

    ws.onmessage = (event) => {
        let msg;
        try {
            msg = JSON.parse(event.data);
        } catch (e) {
            console.error("Invalid JSON:", event.data);
            return;
        }
        handleMessage(msg);
    };

    ws.onclose = () => addSystemMessage("Disconnected from server.");
}

/* -------------------- MESSAGE HANDLING --------------------- */
function handleMessage(msg) {
    switch (msg.type) {

        case "username_rejected":
            alert("Username already taken. Pick another.");
            location.reload();
            break;

        case "username_accepted":
            addSystemMessage(`Welcome ${msg.text}!`);
            getRooms();
            break;

        case "message":
            addChatMessage(msg.from, msg.text);
            break;

        case "system":
            addSystemMessage(msg.text);
            break;

        case "error":
            addErrorMessage(msg.text);
            break;

        case "user_joined":
            addSystemMessage(`🔵 ${msg.from} joined ${msg.room}`);
            ws.send(JSON.stringify({
                command: "/list"
            }));
            break;

        case "user_left":
            addSystemMessage(`🔴 ${msg.from} left ${msg.room}`);
            ws.send(JSON.stringify({
                command: "/list"
            }));
            break;

        case "user_renamed":
            addSystemMessage(`📝 ${msg.text}`);
            break;

        case "rooms_list":
            updateRoomList(msg.data.rooms);
            break;

        case "room_created":
            addRoom(msg.room);
            addSystemMessage(`Room created: ${msg.room}`);
            break;

        case "private":
            addPrivateMessage(msg.from, msg.text);
            break;

        case "users_list":
            updateUserList(msg.data.users);
            break;

        default:
            console.log("Unknown message:", msg);
    }
}

/* ------------------------- UI HELPERS ---------------------- */
function addChatMessage(user, text) {
    addMessage(`[${user}] ${text}`);
}

function addSystemMessage(text) {
    addMessage(`🛈 ${text}`, "system");
}

function addErrorMessage(text) {
    addMessage(`⚠️ ${text}`, "error");
}

function addPrivateMessage(from, text) {
    addMessage(`(PM) ${from}: ${text}`, "private");
}

function addMessage(text, type = "normal") {
    let messages = document.getElementById("messages");

    let div = document.createElement("div");
    div.className = `message ${type}`;
    div.innerText = text;

    messages.appendChild(div);
    messages.scrollTop = messages.scrollHeight;
}

/* ------------------------- ROOMS --------------------------- */

function updateRoomList(rooms) {
    const roomList = document.getElementById("roomList");
    roomList.innerHTML = ""; // clear old

    rooms.forEach(roomName => {
        const li = document.createElement("li");
        li.className = "room";
        li.dataset.room = roomName;
        li.textContent = roomName;

        li.onclick = () => joinRoom(roomName);
        roomList.appendChild(li);
    });
}


function addRoom(roomName) {
    const roomList = document.getElementById("roomList");

    let li = document.createElement("li");
    li.className = "room";
    li.dataset.room = roomName;
    li.innerText = roomName;

    li.onclick = () => joinRoom(roomName);

    roomList.appendChild(li);
}

function getRooms() {

    ws.send(JSON.stringify({
        command: "/rooms"
    }));
}

function joinRoom(roomName) {
    currentRoom = roomName;

    ws.send(JSON.stringify({
        command: "/join",
        args: [roomName]
    }));

    ws.send(JSON.stringify({
        command: "/list"
    }));

    currentRoomName.innerText = roomName;
    addSystemMessage(`Switched to room: ${roomName}`);
    highlightActiveRoom();
}

function highlightActiveRoom() {
    document.querySelectorAll(".room").forEach(li => {
        li.classList.toggle("active", li.dataset.room === currentRoom);
    });
}

function updateUserList(users) {
    roomUsers = users;
    userInfo.innerText = `${users.length} user${users.length !== 1 ? 's' : ''} online`;
    // addSystemMessage("Users in room: " + users.join(", "));
}

/* ------------------------- SEND MESSAGE -------------------- */
function sendMessage() {
    let msg = document.getElementById("msg").value.trim();
    if (!msg) return;

    ws.send(JSON.stringify({
        type: "message",
        text: msg,
        room: currentRoom
    }));

    document.getElementById("msg").value = "";
}

/* ------------------------- CREATE ROOM MODAL --------------------------- */

const modalOverlay = document.getElementById("modalOverlay");
const newRoomInput = document.getElementById("newRoomName");

document.getElementById("createRoomBtn").onclick = () => {
    modalOverlay.classList.remove("hidden");
    newRoomInput.value = "";
    newRoomInput.focus();
};

document.getElementById("createRoomCancel").onclick = () => {
    modalOverlay.classList.add("hidden");
};

document.getElementById("createRoomConfirm").onclick = () => {
    const name = newRoomInput.value.trim();
    if (!name) {
        addErrorMessage("Room name cannot be empty.");
        return;
    }

    ws.send(JSON.stringify({
        command: "/create",
        args: [name]
    }));

    modalOverlay.classList.add("hidden");
};

/* ------------------------- USERNAME MODAL --------------------------- */

usernameConfirm.onclick = () => {
    const input = usernameInput.value.trim();
    if (!input) return alert("Please enter a username");

    username = input;

    usernameModal.classList.add("hidden");
    connect(username);
};

usernameInput.addEventListener("keypress", (e) => {
    if (e.key === "Enter") usernameConfirm.onclick();
});

/* ------------------------- INIT ---------------------------- */
window.onload = () => {
    usernameModal.classList.remove("hidden");
    usernameInput.focus();
};
