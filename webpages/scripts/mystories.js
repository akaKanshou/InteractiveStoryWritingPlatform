function sendTo(location) {
    window.location.href = `https://localhost:8080${location}`
}

function toggleDropDown(element) {
    document.getElementById(element).hidden=!document.getElementById("profileDropdown").hidden
    document.getElementById(element).classList.toggle("open")
}

function closeAllDropdowns() {
    document.getElementById("profileDropdown").hidden = true
    document.getElementById("profileDropdown").classList.remove("open")
}

let Stories = [], cards = []
const VisValPrivate = 1, VisValPublic = 2

async function showContent(visibility) {
    const spinner = document.getElementById("loading-spinner-div")
    const storiesGrid = document.getElementById("storiesGrid")
    const emptyState = document.getElementById("emptyState")

    storiesGrid.classList.add("hidden")
    emptyState.classList.add("hidden")

    spinner.classList.remove("hidden")

    storiesGrid.innerHTML = ""
    const dat = visibility === "all"
        ? cards
        : cards.filter((card) => card.querySelector(".card-visibility-badge").textContent.toLowerCase() === visibility)

    for (const card of dat) {
        storiesGrid.appendChild(card)
        card.classList.remove("dummy")
    }

    if (dat.length === 0) {
        spinner.classList.add("hidden")
        storiesGrid.classList.add("hidden")

        emptyState.classList.remove("hidden")
        return
    }

    spinner.classList.add("hidden")
    emptyState.classList.add("hidden")

    storiesGrid.classList.remove("hidden")
}

async function loadContent(visibility) {
    const res = await fetch(`https://localhost:8080/api/getmystories?visibility=${visibility}`)
    const json = await res.json()

    if (!json.stories) {
        document.getElementById("errorMessage").textContent=json.message
        document.getElementById("errorState").classList.remove("hidden")
        return
    }

    Stories.push(...json.stories)

    setStats()
}

async function setStats(){
    const totalStories = document.getElementById("totalStories")
    totalStories.textContent = Stories.length.toString()

    const n = Stories.filter((story) => story.visibility === VisValPublic).length

    const publicStories = document.getElementById("publicStories")
    const privateStories = document.getElementById("privateStories")

    publicStories.textContent = n.toString()
    privateStories.textContent = (Stories.length - n).toString()
}

document.addEventListener("DOMContentLoaded", async () => {
    document.addEventListener("click", (e) => {
        closeAllDropdowns()
    })

    document.getElementById("navProfile").addEventListener("click", (e) => {
        toggleDropDown("profileDropdown")
        e.stopPropagation()
    })

    document.addEventListener("click", (e) => {
        if (e.target.classList.contains("filter-btn")) {
            const list = document.getElementsByClassName("filter-btn");
            for (const item of list) {
                item.classList.remove('active')
            }
            e.target.classList.add('active')
        }

        showContent(e.target.dataset.filter)
    })

    const visibility = document.querySelector(".filter-btn.active").dataset.filter
    loadContent(visibility)
        .then(() => createCards())
        .then(() => showContent(visibility))
})

async function createCards() {
    for (const story of Stories) {
        cards.push(newStoryCard(story))
    }
}

function newStoryCard(story) {
    const card = document.getElementById("dummyStoryCard1").cloneNode(true)
    card.id = story.story_id
    card.classList.remove("dummy")

    const title = card.querySelector(`.card-title`)
    title.textContent = story.story_name

    const visibility = card.querySelector(`.card-visibility-badge`)
    visibility.textContent = story.visibility === VisValPublic ? "Public" : "Private";

    const storyDescription = card.querySelector(`.card-description`)
    storyDescription.textContent = story.description

    return card
}