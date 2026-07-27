//TODO: Scroll to top on page reload

let stories

function closeAllDropdowns() {
    const dropdowns = document.querySelectorAll(".dropdown")
    for (let dropdown of dropdowns) {
        dropdown.classList.remove("open")
    }
}

document.addEventListener("DOMContentLoaded", async () => {
    document.addEventListener("click", (e) => {
        if (e.target.classList.contains("filter-btn")) {
            const list = document.getElementsByClassName("filter-btn");
            for (const item of list) {
                item.classList.remove('active')
            }
            e.target.classList.add('active')

            showContent(e.target.dataset.filter)
        }

        closeAllDropdowns()
    })

    document.querySelector(".nav-profile").addEventListener("click", (e)=> {
        const dropdown = document.querySelector(".nav-profile .dropdown")
        dropdown.classList.toggle("open")
        e.stopImmediatePropagation()
    })

    stories = document.querySelectorAll(".story-card")

    showContent("all")
})

function showContent(filter) {
    for  (const story of stories) {
        if ((filter === "all") || (story.querySelector(".card-visibility-badge").dataset.filter === filter)) {
            story.classList.remove("hidden")
        } else {
            story.classList.add("hidden")
        }
    }
}