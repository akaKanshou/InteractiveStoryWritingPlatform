function sendTo(location) {
    window.location.href = `https://localhost:8080${location}`
}

function toggleDropDown(element) {
    const profile = document.getElementById(element)
    if (profile) {
        profile.hidden = !profile.hidden
        profile.classList.toggle("open")
    }
}

function closeAllDropdowns() {
    const profile = document.getElementById("profileDropdown")
    if (profile) {
        profile.hidden = true
        profile.classList.remove("open")
    }
}

let Stories = [], cards = []
const VisValPrivate = 1, VisValPublic = 2

async function showContent(visibility) {
    const errorBox = document.getElementById("errorState")
    const spinner = document.getElementById("loading-spinner-div")
    const latestStoriesGrid = document.getElementById("latestStoriesGrid")
    const emptyState = document.getElementById("emptyState")

    latestStoriesGrid.classList.add("hidden")
    emptyState.classList.add("hidden")
    errorBox.classList.add("hidden")

    spinner.classList.remove("hidden")

    latestStoriesGrid.innerHTML = ""
    const dat = visibility === "all"
        ? cards
        : cards.filter((card) => card.querySelector(".card-visibility-badge").textContent.toLowerCase() === visibility)

    for (const card of dat) {
        latestStoriesGrid.appendChild(card)
        card.classList.remove("dummy")
    }

    if (dat.length === 0) {
        spinner.classList.add("hidden")
        latestStoriesGrid.classList.add("hidden")

        emptyState.classList.remove("hidden")
        return
    }

    spinner.classList.add("hidden")
    emptyState.classList.add("hidden")

    latestStoriesGrid.classList.remove("hidden")
}

async function loadContent(visibility) {
    const latestStoriesGrid = document.getElementById("latestStoriesGrid")
    const emptyState = document.getElementById("emptyState")
    const spinner = document.getElementById("loading-spinner-div")
    const errorBox = document.getElementById("errorState")

    errorBox.classList.add("hidden")
    latestStoriesGrid.classList.add("hidden")
    emptyState.classList.add("hidden")

    spinner.classList.remove("hidden")

    const json = await fetch(`https://localhost:8080/api/getmystories?visibility=${visibility}`)
        .then((res) => res.json())
        .catch((reason) => {
            console.log(reason)
            return {
                message: reason,
                stories: null
            }
        })

    if (!json.stories) {
        document.getElementById("errorMessage").textContent = json.message

        spinner.classList.add("hidden")
        latestStoriesGrid.classList.add("hidden")
        emptyState.classList.add("hidden")

        errorBox.classList.remove("hidden")
        return
    }

    Stories = json.stories

    setStats()
}

async function setStats() {
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

    // const visibility = document.querySelector(".filter-btn.active").dataset.filter
    // loadContent(visibility)
    //     .then(() => createCards())
    //     .then(() => showContent(visibility))
    //     .catch((reason) => {
    //         console.log(reason)
    //     })

    for (let i = 0; i < 2; i++) {
        document.getElementById("editorsPicksStoriesGrid").appendChild(document.querySelector(".story-card").cloneNode(true))
    }

    for (let i = 0; i < 5; i++) {
        document.getElementById("trendingSectionStoriesGrid").appendChild(document.querySelector(".story-card").cloneNode(true))
    }

    for (let i = 0; i < 2; i++) {
        document.getElementById("latestStoriesGrid").appendChild(document.querySelector(".story-card").cloneNode(true))
    }

    for (let i = 0; i < 20; i++) {
        document.getElementById("recentlyUpdatedStoriesGrid").appendChild(document.querySelector(".story-card").cloneNode(true))
    }

    window.scrollTo(0, 0)
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