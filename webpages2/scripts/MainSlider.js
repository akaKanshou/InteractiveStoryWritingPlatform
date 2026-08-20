
const slides = document.querySelectorAll(".sliderSlide");
const dots = document.querySelectorAll(".sliderDot");
let i = 0;

function show(n){
    slides[i].classList.remove("active");
    dots[i].classList.remove("active");

    i = (n + slides.length) % slides.length;

    slides[i].classList.add("active");
    dots[i].classList.add("active");
}

document.querySelector(".nextArrow").onclick = () => show(i + 1);
document.querySelector(".prevArrow").onclick = () => show(i - 1);

dots.forEach((dot, index) => dot.onclick = () => show(index));

setInterval(() => show(i + 1), 7500);