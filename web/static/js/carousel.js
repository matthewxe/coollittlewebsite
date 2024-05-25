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
                window.requestAnimationFrame(fadeit);
        }
        window.requestAnimationFrame(scrollit);
}

function fadeit() {}

window.requestAnimationFrame(scrollit);
