document.addEventListener('DOMContentLoaded', async () => {
    const chapter_id = document.getElementById('editorJS').dataset.id;

    const req = new Request(`https://localhost:8080/api/getcontent/${chapter_id}`, {
        method: "GET",
    });

    const resp = await fetch(req);

    if (resp.status !== 200) {
        console.log("FUCK")
        return
    }

    const editorJS = new EditorJS({
        readOnly: true,
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

        data: await resp.json(),
    });
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