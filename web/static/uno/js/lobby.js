var conn;

if (window["WebSocket"]) {
        conn = websoccket_connect("wss");
} else {
        console.log("your browser does not support websockets");
}

function websoccket_connect(ws) {
        console.log("trying to connect...");
        conn = new WebSocket(ws + "://" + document.location.host + "/uno/game");

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
}
