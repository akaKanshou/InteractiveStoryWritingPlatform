const hamBurger = document.getElementById("hamburger");
const nav = document.querySelector("nav");

const Ham1 = document.getElementById("hamburgerRow1");
const Ham2 = document.getElementById("hamburgerRow2");
const Ham3 = document.getElementById("hamburgerRow3");

let open = false;

hamBurger.addEventListener("click", () => {

    if (!open) {
        nav.style.left = "0";
        nav.classList.add("open");

        Ham1.style.transform =
            "translateY(8px) rotate(45deg)";

        Ham2.style.opacity = "0";

        Ham3.style.transform =
            "translateY(-8px) rotate(-45deg)";

    } else {
        nav.style.left = "-250px";
        nav.classList.remove("open");

        Ham1.style.transform = "rotate(0deg)";
        Ham2.style.opacity = "1";
        Ham3.style.transform = "rotate(0deg)";
    }

    open = !open;
});

document.querySelectorAll(".logoAndId").forEach((element, index) => {
    element.addEventListener("click", () => {
        window.location = element.dataset.href;
    })
})

const editBtn = document.getElementById("edit");
const deleteBtn = document.getElementById("delete");

editBtn.addEventListener("click", () => {
    const id = document.querySelector("body").dataset.storyid;
    window.location = `https://localhost:8080/story/edit/${id}`;
})

deleteBtn.addEventListener("click", async () => {
    const id = document.querySelector("body").dataset.storyid;
    const res = await fetch(`https://localhost:8080/api/story/${id}`, {method: "DELETE"})
    if (res.status === 200) {
        window.location= `https://localhost:8080/user/mystories`;
    } else {
        const resJson = await res.json()
        console.log(resJson)
    }
})



document.querySelectorAll(".chapterRow").forEach((element, index) => {
    element.addEventListener("click", () => {
        window.location = `https://localhost:8080/chapter/view/${element.dataset.chapterid}`
    })
})