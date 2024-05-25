const carousel = document.getElementById("carousel")

carousel.addEventListener("scroll", (event) => {
        // scrollit();
        console.log("yo");
});
carousel.scrollLeft = 500;

async function scrollit() {
        carousel.scroll(0, 0);
        while (true) {
                carousel.scrollBy(1, 0);
                await new Promise((r) => setTimeout(r, 100));
        }
}
scrollit();
