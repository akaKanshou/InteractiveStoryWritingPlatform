// TODO: Add error results like in create user
const createBtn = document.getElementById("create");
const editBtn = document.getElementById("edit");

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

const publicBtn = document.getElementById("public");
const privateBtn = document.getElementById("private");

publicBtn.addEventListener("click", () => {
    publicBtn.classList.add("active");
    privateBtn.classList.remove("active");
});

privateBtn.addEventListener("click", () => {
    privateBtn.classList.add("active");
    publicBtn.classList.remove("active");
});

document.addEventListener("DOMContentLoaded", (event) => {
    if (createBtn) {
        createBtn.addEventListener("click", async () => {

            const title = document.getElementById("titlewrite").value;
            const discription = document.getElementById("discrip2").value;

            const formData = new FormData();
            formData.append("story_name", title);
            formData.append("description", discription);
            formData.append("visibility", document.querySelector(".active").id);
            const response = await fetch("https://localhost:8080/api/newstory", {
                method: "POST",
                body: formData
            });
            const data = await response.json();

            if (response.status === 201) {
                window.location = `https://localhost:8080/story/${data.story_id}`;
                return
            }

            console.log(data);
        });
    }

    if (editBtn) {
        editBtn.addEventListener("click", async () => {
            const title = document.getElementById("titlewrite").value;
            const discription = document.getElementById("discrip2").value;
            const storyid = document.querySelector("body").dataset.storyid;

            const formData = new FormData();
            formData.append("story_name", title);
            formData.append("description", discription);
            formData.append("visibility", document.querySelector(".active").id);
            formData.append("story_id", storyid);

            const res = await fetch("https://localhost:8080/api/editstory", {
                method: "POST",
                body: formData
            })
            if (res.status === 200) {
                window.location = `https://localhost:8080/story/view/${storyid}`;
            } else {
                const jsonResp = await res.json();
                console.log(jsonResp);
            }
        })
    }

    const visibility = document.querySelector(".visi").dataset.visibility
    if (!visibility) {
        return;
    }
    if (visibility === "1") {
        privateBtn.classList.add("active");
    } else {
        publicBtn.classList.add("active");
    }
})

document.querySelectorAll(".logoAndId").forEach((element, index) => {
    element.addEventListener("click", () => {
        window.location = element.dataset.href;
    })
})