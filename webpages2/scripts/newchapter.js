document.addEventListener('DOMContentLoaded', () => {
    const editorJS = new EditorJS({
        readOnly: false,
        holder: "editorJS",
        minHeight: 30,
        tools: {
            header: {
                class: Header,
                inlineToolbar: ["italic", "bold"],
                config: {
                    placeholder: "Heading",
                    levels: [2, 3, 4],
                    defaultLevel: 3,
                },
                shortcut: "CMD+SHIFT+H",
            },

            quote: {
                class: Quote,
                inlineToolbar: true,
                config: {
                    quotePlaceholder: "Enter a quote",
                    captionPlaceholder: "Quote's author",
                },
                shortcut: "CMD+SHIFT+O",
            },

            delimiter: Delimiter,
        },

        defaultBlock: "paragraph",

        onReady: function () {
            console.log("EditorJS is ready");
        },

        onChange: function (api, event) {},

        data: {
            blocks: [
                {
                    type: "header",
                    data: {
                        text: "Write your content here",
                        level: 2,
                    },
                },

                {
                    type: "paragraph",
                    data: {
                        text: "Write your content here.",
                    },
                },
            ],
        },
    });

    document.getElementById("submit").addEventListener("click", async function (event) {
        event.preventDefault();

        const content = await editorJS.save()
            .then((savedData) => {
                return JSON.stringify(savedData)
            })
            .catch((error) => {
                return ""
            });

        const formData = new FormData();
        formData.append("chapter_name", document.getElementById("chapterNameTextInput").value);
        formData.append("story_id", document.getElementById("storyIDTextInput").value);
        formData.append("content", content);
        formData.append("index", document.getElementById("chapterNumNumberInput").value);
        formData.append("edgesInc", document.getElementById("edgesIncoming").value)
        formData.append("edgesOut", document.getElementById("edgesOutgoing").value)

        const req = new Request("https://localhost:8080/api/newchapter", {
            method: "POST",
            body: formData,
        });

        const resp = await fetch(req);
        const jsonRes = await resp.json()

        if (resp.status === 201) {
            window.location.href = `https://localhost:8080/chapter/view/${jsonRes.chapter_id}`;
            return
        }

        console.log(resp, jsonRes);
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

document.addEventListener("DOMContentLoaded", () => {
    const searchParams = new URLSearchParams(window.location.search);
    const storyIDField = document.getElementById("storyIDTextInput");
    if (searchParams.has("story")) {
          storyIDField.value = searchParams.get("story");
    } else {
        storyIDField.value = "";
        storyIDField.disabled = false;
    }
})