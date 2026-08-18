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

        data: await resp.json(),
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
        formData.append("content", content);
        formData.append("chapter_id", chapter_id);
        // TODO: edge editing from this page
        // formData.append("edgesInc", document.getElementById("edgesIncoming").value)
        // formData.append("edgesOut", document.getElementById("edgesOutgoing").value)

        const req = new Request(`https://localhost:8080/api/editchapter`, {
            method: "POST",
            body: formData,
        });

        const resp = await fetch(req);
        const jsonRes = await resp.json()

        if (resp.status === 200) {
            window.location.href = `https://localhost:8080/chapter/view/${chapter_id}`;
            return
        }

        console.log(resp, jsonRes);
    })

    document.getElementById("delete").addEventListener("click", async function (event) {
        event.preventDefault();

        const req = new Request(`https://localhost:8080/api/deletechapter/${chapter_id}`, {method: "DELETE"});
        const resp = await fetch(req);

        if (resp.status === 200) {
            window.location.href = `https://localhost:8080/index`;
            return
        }

        console.log(resp.body)
    })
})

