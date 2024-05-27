let ongoing = document.getElementById("ongoing");

function refresh(ongoing) {
        console.log("yo");
        const url = new URL(
                window.location.protocol +
                        "//" +
                        window.location.host +
                        "/uno/list",
        );

        fetch(url, {
                headers: { getnames: true },
        })
                .then((response) => response.text())
                .then((text) => {
                        ongoing.innerText = text;
                });
}

refresh(ongoing);
