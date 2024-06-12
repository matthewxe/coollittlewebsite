var conn;
var wsstatus = document.getElementById("ws");
var msg = document.getElementById("msg");
var log = document.getElementById("log");
const websockettype = "wss";

function wslog(message, ev) {
        wsstatus.innerText = message;
        console.log(message, ev);
}

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
                Date: Date.now(),
        };

        conn.send(JSON.stringify(json));
        msg.value = "";
        return false;
};

if (window["WebSocket"]) {
        conn = websocket_connect(websockettype);
} else {
        wslog("Your browser does not support websockets");
}

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
        conn.onmessage = function (ev) {
                wslog("New message", ev);
                var messages = ev.data.split("\n");
                for (var i = 0; i < messages.length; i++) {
                        const msg = JSON.parse(messages[i]);
                        console.log("JDSL", msg);
                        const time = new Date(msg.Date);
                        const timeStr = time.toLocaleTimeString();

                        var item = document.createElement("div");
                        item.innerText =
                                timeStr + " " + msg.Player + ": " + msg.Text;

                        appendLog(item);
                }
        };

        conn.onerror = function (ev) {
                wslog("oopsie dasies ", ev);
                conn = websocket_connect(websockettype);
        };
        return conn;
}

function appendLog(item) {
        var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.appendChild(item);
        if (doScroll) {
                log.scrollTop = log.scrollHeight - log.clientHeight;
        }
}
