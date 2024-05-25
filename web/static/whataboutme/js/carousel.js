const carousel = document.getElementById("carousel");

function scrollit() {
        carousel.scrollBy({ left: 50, behavior: "smooth" });
        if (
                carousel.scrollWidth - carousel.offsetWidth ==
                carousel.scrollLeft
        ) {
                return window.requestAnimationFrame(fadeout);
        }
        window.requestAnimationFrame(scrollit);
}

function fadeout() {
        if (carousel.style.opacity <= 0) {
                carousel.scrollLeft = 0;
                return window.requestAnimationFrame(fadein);
        }
        carousel.style.opacity -= 0.1;
        window.requestAnimationFrame(fadeout);
}

function fadein() {
        if (carousel.style.opacity >= 1) {
                return window.requestAnimationFrame(scrollit);
        }
        carousel.style.opacity = +carousel.style.opacity + 0.1;
        window.requestAnimationFrame(fadein);
}

if (carousel.scrollWidth - carousel.offsetWidth != carousel.scrollLeft) {
        window.requestAnimationFrame(scrollit);
}
