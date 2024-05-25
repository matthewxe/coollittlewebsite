const carousel = document.getElementById("carousel");

function scrollit() {
        carousel.scrollBy({ left: 100, behavior: "smooth" });
        console.log(
                carousel.scrollWidth - carousel.offsetWidth ==
                        carousel.scrollLeft,
        );
        if (
                carousel.scrollWidth - carousel.offsetWidth ==
                carousel.scrollLeft
        ) {
                console.log("scrollout");
                return window.requestAnimationFrame(fadeout);
        }
        window.requestAnimationFrame(scrollit);
}

function fadeout() {
        // console.log(carousel.style.opacity);
        if (carousel.style.opacity <= 0) {
                // console.log("fadeout");
                return window.requestAnimationFrame(fadein);
        }
        carousel.style.opacity -= 0.1;
        window.requestAnimationFrame(fadeout);
}

function fadein() {
        // console.log("startfadein");
        if (carousel.style.opacity >= 1) {
                // console.log("fadein");
                carousel.scrollLeft = 0;
                return window.requestAnimationFrame(scrollit);
        }
        carousel.style.opacity += 0.1;
        window.requestAnimationFrame(fadein);
}

window.requestAnimationFrame(scrollit);
