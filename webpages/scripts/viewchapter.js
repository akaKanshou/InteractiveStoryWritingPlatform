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

