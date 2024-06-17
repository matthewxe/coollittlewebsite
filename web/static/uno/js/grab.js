let ongoing = document.getElementById("ongoing");

function refresh(ongoing) {
        const url = new URL(
                window.location.protocol +
                        "//" +
                        window.location.host +
                        "/uno/list",
        );

        fetch(url)
                .then((response) => response.text())
                .then((text) => {
                        ongoing.innerHTML = text;
                });
}

refresh(ongoing);
