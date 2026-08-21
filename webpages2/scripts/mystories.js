document.querySelectorAll(".logoAndId").forEach((element, index) => {
    element.addEventListener("click", () => {
        window.location = element.dataset.href;
    })
})


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
        Ham1.style.transform = "translateY(8px) rotate(45deg)";
        Ham2.style.opacity = "0";
        Ham3.style.transform = "translateY(-8px) rotate(-45deg)";
    } else {
        nav.style.left = "-250px";
        nav.classList.remove("open");

        Ham1.style.transform = "rotate(0deg)";
        Ham1.style.top = "0";

        Ham2.style.opacity = "1";

        Ham3.style.transform = "rotate(0deg)";
        Ham3.style.top = "0";
    }

    open = !open;

});


const storyList = document.getElementById("storyList");
const stories = Array.from(
    storyList.querySelectorAll(".slide")
);

const tabs = document.querySelectorAll(".Hb");
const searchBar = document.getElementById("searchBar");
const storyCount = document.getElementById("storyCount");

let currentFilter = "all";

function updateStories() {

    const search =
        searchBar.value.trim().toLowerCase();

    stories.forEach(story => {

        const type = story.dataset.type;

        const title =
            story.querySelector(".titleOfSlide")
                .textContent
                .toLowerCase();

        const description =
            story.querySelector(".DesOfSlide")
                .textContent
                .toLowerCase();

        const matchesFilter =
            currentFilter === "all" ||
            type === currentFilter;

        const matchesSearch =
            !search ||
            title.includes(search) ||
            description.includes(search);

        story.style.display =
            matchesFilter && matchesSearch
                ? "flex"
                : "none";
    });

    updateCount();
}


function updateCount() {

    const visibleStories =
        stories.filter(
            story => story.style.display !== "none"
        );

    storyCount.textContent =
        visibleStories.length +
        (visibleStories.length === 1
            ? " story"
            : " stories");

}


tabs.forEach(tab => {

    tab.addEventListener("click", () => {
        tabs.forEach(item =>
            item.classList.remove("active")
        );

        tab.classList.add("active");

        currentFilter =
            tab.dataset.filter;

        updateStories();

    });

});


searchBar.addEventListener(
    "input",
    updateStories
);

const gridButton =
    document.querySelector(".b1");

const listButton =
    document.querySelector(".b2");


gridButton.addEventListener("click", () => {

    storyList.classList.remove("list-view");

    gridButton.classList.add("selected");
    listButton.classList.remove("selected");

});


listButton.addEventListener("click", () => {

    storyList.classList.add("list-view");

    listButton.classList.add("selected");
    gridButton.classList.remove("selected");

});

const createBtn = document.querySelector(".Create")

if (createBtn) {
    createBtn.addEventListener("click", () => {
        window.location = "https://localhost:8080/story/new";
    });
}

document.querySelectorAll(".slide").forEach((element, index) => {
    element.addEventListener("click", () => {
        window.location = `https://localhost:8080/story/view/${element.dataset.storyid}` ;
    })
})

updateStories();

//
// let stories
//
// function closeAllDropdowns() {
//     const dropdowns = document.querySelectorAll(".dropdown")
//     for (let dropdown of dropdowns) {
//         dropdown.classList.remove("open")
//     }
// }
//
// document.addEventListener("DOMContentLoaded", async () => {
//     document.addEventListener("click", (e) => {
//         if (e.target.classList.contains("filter-btn")) {
//             const list = document.getElementsByClassName("filter-btn");
//             for (const item of list) {
//                 item.classList.remove('active')
//             }
//             e.target.classList.add('active')
//
//             showContent(e.target.dataset.filter)
//         }
//
//         closeAllDropdowns()
//     })
//
//     document.querySelector(".nav-profile").addEventListener("click", (e)=> {
//         const dropdown = document.querySelector(".nav-profile .dropdown")
//         dropdown.classList.toggle("open")
//         e.stopImmediatePropagation()
//     })
//
//     stories = document.querySelectorAll(".story-card")
//
//     showContent("all")
// })
//
// function showContent(filter) {
//     for  (const story of stories) {
//         if ((filter === "all") || (story.querySelector(".card-visibility-badge").dataset.filter === filter)) {
//             story.classList.remove("hidden")
//         } else {
//             story.classList.add("hidden")
//         }
//     }
// }