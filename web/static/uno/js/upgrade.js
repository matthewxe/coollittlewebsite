// TODO: I have no clue if this is a safe way to turn the uri into websocket
const wsUri = document.location.href.replace("http", "ws") + "/upgrade";
const websocket = new WebSocket(wsUri);

var username;

websocket.addEventListener("open", () => {
        setContent("status", "Connected");
});

websocket.addEventListener("error", (e) => {
        setContent("status", "Error" + e);
});

websocket.addEventListener("message", (e) => {
        const msg = JSON.parse(e.data);
        console.log(msg);

        switch (msg.type) {
                case "name":
                        username = msg.data.name;
                        setContent("name", username);
                        break;
                case "count":
                        setContent("count", msg.data.count);
                        break;
                case "opp":
                        console.log("opp");
                        break;
        }
});

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

function increment() {
        const json = { type: "increment", data: { name: username } };
        websocket.send(JSON.stringify(json));
}

// function createOpponent(name)
