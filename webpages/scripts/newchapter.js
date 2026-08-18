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

