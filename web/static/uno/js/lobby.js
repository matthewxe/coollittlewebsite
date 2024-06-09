var conn;
var id = document.getElementById("ws");
var msg = document.getElementById("msg");
var log = document.getElementById("log");

/*if (window["WebSocket"]) {
        conn = websoccket_connect("wss");
} else {
        console.log("your browser does not support websockets");
}

function websoccket_connect(ws) {
        console.log("trying to connect...");
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
                console.log("onclose", ev);
        };
        conn.onmessage = function (ev) {
                console.log("onmessage", ev);
        };
        conn.onerror = function (ev) {
                console.log("oopsie dasies ", ev);
                conn = websoccket_connect("ws");
        };
        return conn;
}*/

function appendLog(item) {
        var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.appendChild(item);
        if (doScroll) {
                log.scrollTop = log.scrollHeight - log.clientHeight;
        }
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
        conn = new WebSocket(
                "ws:" +
                        document.location.href.replace(
                                document.location.protocol,
                                "",
                        ) +
                        "/ws",
        );
        conn.onclose = function (evt) {
                var item = document.createElement("div");
                item.innerHTML = "<b>Connection closed.</b>";
                appendLog(item);
        };
        conn.onmessage = function (evt) {
                var messages = evt.data.split("\n");
                for (var i = 0; i < messages.length; i++) {
                        var item = document.createElement("div");
                        item.innerText = messages[i];
                        appendLog(item);
                }
        };
} else {
        var item = document.createElement("div");
        item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
        appendLog(item);
}
