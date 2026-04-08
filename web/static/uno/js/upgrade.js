// TODO: I have no clue if this is a safe way to turn the uri into websocket

var websocket;
var username;

connect();

async function connect() {
        // await fetch("/uno/cookie");
        //
        // if (document.cookie.indexOf("uno=") === -1) {
        //         alert("Ayo cookies arent setting cuh, refresh or sum");
        //         return;
        // }

        if (!window["WebSocket"]) {
                alert("Your browser does not support websockets");
                return;
        }

        const auth = await fetch("/uno/auth").then((response) =>
                response.text(),
        );
        const wsUri = `${document.location.href.replace(
                "http",
                "ws",
        )}/upgrade?room=uno&auth=${auth}`;
        console.log(wsUri);

        websocket = new WebSocket(wsUri);

        websocket.addEventListener("open", () => {
                setContent("status", "Connected");
        });

        websocket.addEventListener("error", (e) => {
                setContent("status", `Error (${e})`);
        });

        websocket.addEventListener("close", (e) => {
                setContent("status", `Closed (${e.code}: ${e.reason})`);
        });

        websocket.addEventListener("message", (e) => handler(e));
}

function handler(e) {
        {
                const msg = JSON.parse(e.data);
                console.log(msg);

                switch (msg.type) {
                        case "name":
                                username = msg.data.name;
                                setContent("name", username);
                                break;
                        case "state":
                                setContent("opponent", "");
                                for (const key in msg.data.Players) {
                                        var player = msg.data.Players[key];
                                        if (key === username) {
                                                setContent(
                                                        "count",
                                                        player.Count,
                                                );
                                                continue;
                                        }
                                        addContent(
                                                "opponent",
                                                `
<hr>
<div>
        <h2>Name: ${key}</h2>
        <h2>Count: ${player.Count}</h2>
</div>
                         `,
                                        );
                                }

                                break;
                }
        }
}

/**
 * @param {string} id - id of the thing you want to change
 * @param {string} content - the thing you want to set it as
 */
function setContent(id, content) {
        var elem = document.getElementById(id);
        if (elem !== null) {
                elem.innerHTML = content;
        }
}

/**
 * @param {string} id - id of the thing you want to change
 * @param {string} content - the thing you want to set it as
 */
function addContent(id, content) {
        var elem = document.getElementById(id);
        if (elem !== null) {
                elem.innerHTML += content;
        }
}

function increment() {
        const json = { type: "increment", data: { name: username } };
        websocket.send(JSON.stringify(json));
}

// function createOpponent(name)
