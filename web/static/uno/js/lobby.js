var conn;
var wsstatus = document.getElementById("ws");
var msg = document.getElementById("msg");
var log = document.getElementById("log");
var state = document.getElementById("state");
var players = document.getElementById("players");
const websockettype = "wss";

function wslog(message, ev) {
        wsstatus.innerText = message;
        console.log(message, ev);
}

// Log //{{{
document.getElementById("form").onsubmit = function () {
        if (!conn) {
                return false;
        }
        if (!msg.value) {
                return false;
        }
        const json = {
                Type: "message",
                Text: msg.value,
        };

        conn.send(JSON.stringify(json));
        msg.value = "";
        return false;
}; //}}}

// Check if websockets work {{{
if (window["WebSocket"]) {
        conn = websocket_connect(websockettype);
} else {
        wslog("Your browser does not support websockets");
} //}}}

// Websocket function{{{
function websocket_connect(ws) {
        wslog("Trying to connect");
        conn = new WebSocket(
                ws +
                        ":" +
                        document.location.href.replace(
                                document.location.protocol,
                                "",
                        ) +
                        "/ws",
        );

        conn.onclose = function (ev) {
                wslog("Connection closed", ev);
        };
        conn.onmessage = on_message;

        conn.onerror = function (ev) {
                wslog("oopsie dasies ", ev);
                conn = websocket_connect(websockettype);
        };
        return conn;
} //}}}

// Appendlog{{{
function appendLog(item) {
        var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.appendChild(item);
        if (doScroll) {
                log.scrollTop = log.scrollHeight - log.clientHeight;
        }
} //}}}

//Handle incoming messages{{{
function on_message(ev) {
        wslog("New message", ev);
        const json = JSON.parse(ev.data);
        console.log("OnMessage", json);

        switch (json.Type) {
                case "status":
                        return handleLobbyUpdate(json);
                case "message":
                        return handleMessages(json);
                case "messagelog":
                        return handleMessagesLog(json);
        }
} //}}}

function handleLobbyUpdate(json) {
        console.log("LobbyUpdate");
        //if (json.Id != document.getElementById("id").innerHTML) {
        //        return;
        //}
        state.innerHTML = json.State;
        switch (json.State) {
                case 0:
                        state.innerHTML = "Ready";
                        break;
                case 1:
                        state.innerHTML = "Ongoing";
                        break;
                case 2:
                        state.innerHTML = "Done";
        }
        players.innerHTML = "";
        players.innerHTML += `<li>${handlePlayer(json.Leader, true)}</li>`;

        if (json.Players == null) {
                return;
        }
        for (let index = 0; index < json.Players.length; index++) {
                const element = handlePlayer(json.Players[index], false);
                players.innerHTML += `<li>${element}</li>`;
        }
}

// Does not include leader
function handlePlayer(player, isleader) {
        var out = player.Name;
        if (isleader) {
                out += " (leader)";
        }
        if (player.Current) {
                out += " (you)";
        }
        if (player.Active) {
                out += " /";
        } else {
                out += " X";
        }
        return out;
}

function handleMessages(json) {
        console.log("Message");
        const time = new Date(json.Date);
        const timeStr = time.toLocaleTimeString();

        var item = document.createElement("div");
        /* dont worry its sanitized (i think) */
        item.innerHTML = timeStr + " <b>" + json.Player + "</b>: " + json.Text;

        appendLog(item);
        //var messages = ev.data.split("\n");
        //for (var i = 0; i < messages.length; i++) {
        //        console.log("JDSL", json);
        //}
}

function handleMessagesLog(json) {
        for (let index = 0; index < json.Log.length; index++) {
                const log = json.Log[index];
                handleMessages(log);
        }
}

function startLobby() {
        const json = {
                Type: "start",
        };
        conn.send(JSON.stringify(json));
        conn.send;
}
