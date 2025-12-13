let ws;
let currentRoom = "lobby";
let username = null;
let roomUsers = [];

const usernameModal = document.getElementById("usernameModal");
const usernameInput = document.getElementById("usernameInput");
const usernameConfirm = document.getElementById("usernameConfirm");
const changeNameBtn = document.getElementById("changeNameBtn");
const changeNameModal = document.getElementById("changeNameModal");
const changeNameInput = document.getElementById("changeNameInput");
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
            // Ensure UI reflects the current room (hide/show Leave button)
            updateUIForRoom(currentRoom);
            break;

        case "message":
            addChatMessage(msg.from, msg.text);
            break;

        case "system":
            if (msg.text.startsWith("Username changed to")) {
                const parts = msg.text.split(" ");
                const newName = parts[parts.length - 1];
                username = newName;
            }
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

        case "left_room":
            switchToLobby();
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
    let messages = document.getElementById("messages");

    const time = new Date().toLocaleTimeString([], {hour: '2-digit', minute: '2-digit'});

    let div = document.createElement("div");
    div.className = (user === username) ? "message mine" : "message";

    // avatar + body
    const avatar = document.createElement("div");
    avatar.className = "avatar";
    avatar.innerText = user ? user.charAt(0).toUpperCase() : "?";

    const body = document.createElement("div");
    body.className = "message-body";

    const meta = document.createElement("div");
    meta.className = "meta";

    const userSpan = document.createElement("span");
    userSpan.className = "user";
    userSpan.innerText = user;

    const timeSpan = document.createElement("span");
    timeSpan.className = "time";
    timeSpan.innerText = time;

    meta.appendChild(userSpan);
    meta.appendChild(timeSpan);

    const content = document.createElement("div");
    content.className = "content";
    content.innerText = text;

    body.appendChild(meta);
    body.appendChild(content);

    div.appendChild(avatar);
    div.appendChild(body);

    messages.appendChild(div);
    messages.scrollTop = messages.scrollHeight;
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

    // system / error / private use simpler layout
    if (type === "system" || type === "error") {
        const content = document.createElement("div");
        content.className = "content";
        content.innerText = text;
        div.appendChild(content);
    } else {
        const avatar = document.createElement("div");
        avatar.className = "avatar";
        avatar.innerText = "#";

        const body = document.createElement("div");
        body.className = "message-body";

        const content = document.createElement("div");
        content.className = "content";
        content.innerText = text;

        body.appendChild(content);
        div.appendChild(avatar);
        div.appendChild(body);
    }

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
        if (roomName === currentRoom) li.classList.add("active");

        li.onclick = () => joinRoom(roomName);
        roomList.appendChild(li);
    });

    // Ensure the active room is highlighted after rebuilding the list
    highlightActiveRoom();
}

function highlightRoom(roomName) {
    const items = document.querySelectorAll("#roomList .room");

    items.forEach(item => {
        if (item.dataset.room === roomName) {
            item.classList.add("active");
        } else {
            item.classList.remove("active");
        }
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

    addSystemMessage(`Switched to room: ${roomName}`);
    // Update UI for the newly joined room (show/hide Leave button)
    updateUIForRoom(roomName);
    highlightActiveRoom();
}

function highlightActiveRoom() {
    document.querySelectorAll(".room").forEach(li => {
        li.classList.toggle("active", li.dataset.room === currentRoom);
    });
}

function updateUserList(users) {
    roomUsers = users;
    // show an icon with the user count instead of a text string
    if (userInfo) {
        userInfo.innerHTML = `
            <span class="user-icon"><img src="icons/users.svg" alt="users"></span>
            <span class="presence-dot" aria-hidden="true"></span>
            <span class="user-count">${users.length}</span>`;
    }
    // addSystemMessage("Users in room: " + users.join(", "));
}

document.getElementById("leaveBtn").onclick = () => {
    ws.send(JSON.stringify({
        command: "/leave",
        args: []
    }));
};

function updateUIForRoom(roomName) {
    const leaveBtn = document.getElementById("leaveBtn");
    const info = document.getElementById("userInfo");
    if (roomName === "lobby") {
        leaveBtn.classList.add("hidden");
        if (info) info.classList.add("hidden");
    } else {
        leaveBtn.classList.remove("hidden");
        if (info) info.classList.remove("hidden");
    }
    // Display 'Lobby' with capital first letter when server sends 'lobby'
    const displayName = (roomName === "lobby") ? "Lobby" : roomName;
    document.getElementById("currentRoomName").textContent = displayName;
}

function switchToLobby() {
    highlightRoom("lobby");
    updateUIForRoom("lobby");
    clearMessages();
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
    addChatMessage(username, msg);

    document.getElementById("msg").value = "";
}

function clearMessages() {
    const messages = document.getElementById("messages");
    messages.innerHTML = "";
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

changeNameBtn.onclick = () => {
    changeNameModal.classList.remove("hidden");
    changeNameInput.value = "";
    changeNameInput.focus();
};

document.getElementById("changeNameCancel").onclick = () => {
    changeNameModal.classList.add("hidden");
};

document.getElementById("changeNameConfirm").onclick = () => {
    const newName = changeNameInput.value.trim();
    if (!newName) {
        addErrorMessage("Enter a valid username.");
        return;
    }

    ws.send(JSON.stringify({
        command: "/name",
        args: [newName]
    }));

    changeNameModal.classList.add("hidden");
};

/* ------------------------- INIT ---------------------------- */
window.onload = () => {
    usernameModal.classList.remove("hidden");
    usernameInput.focus();
    // ensure UI elements reflect initial room (hide userInfo in lobby)
    updateUIForRoom(currentRoom);
};
