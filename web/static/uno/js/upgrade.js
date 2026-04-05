// TODO: I have no clue if this is a safe way to turn the uri into websocket
const wsUri = document.location.href.replace("http", "ws") + "/upgrade";
const websocket = new WebSocket(wsUri);

// const websocket = new WebSocket(wsUri);

// var conn;
// //var wsstatus = document.getElementById("ws");
// var msg = document.getElementById("msg");
// var log = document.getElementById("log");
// var state = document.getElementById("state");
// var players = document.getElementById("players");
// const websockettype = "ws";
//
// // Websocket function{{{
// function websocket_connect(ws) {
//         //wslog("Trying to connect");
//         conn = new WebSocket(
//                 ws +
//                         ":" +
//                         document.location.href.replace(
//                                 document.location.protocol,
//                                 "",
//                         ) +
//                         "/ws",
//         );
//
//         conn.onclose = function (ev) {
//                 //wslog("Connection closed", ev);
//         };
//         conn.onmessage = on_message;
//
//         conn.onerror = function (ev) {
//                 //wslog("oopsie dasies ", ev);
//                 conn = websocket_connect(ws);
//         };
//         return conn;
// } //}}}
