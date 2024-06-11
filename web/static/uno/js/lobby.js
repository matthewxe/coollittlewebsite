var conn;
var wsstatus = document.getElementById("ws");
var msg = document.getElementById("msg");
var log = document.getElementById("log");

function wslog(message, ev) {
        wsstatus.innerHTML = message;
        console.log(message, ev);
}

document.getElementById("form").onsubmit = function () {
        if (!conn) {
                return false;
        }
        if (!msg.value) {
                return false;
        }
        conn.send(msg.value);
        msg.value = "";
        return false;
};

if (window["WebSocket"]) {
        conn = websocket_connect("wss");
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
                        var item = document.createElement("div");
                        item.innerText = messages[i];
                        appendLog(item);
                }
        };

        conn.onerror = function (ev) {
                wslog("oopsie dasies ", ev);
                conn = websocket_connect("wss");
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
