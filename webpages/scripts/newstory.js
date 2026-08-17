let visibility = "";

function setButtonOffAll() {
    for (let i = 0; i < document.querySelectorAll("button").length; i++) {
    document.querySelectorAll("button")[i].classList.remove("active")
    }
}

document.getElementById("public").addEventListener("click", (e) => {
    setButtonOffAll();
    visibility = "public";
    // document.getElementById("public").classList.add("active");
    e.target.classList.add("active");
});

document.getElementById("private").addEventListener("click", () => {
    setButtonOffAll();
    visibility = "private";
     e.target.classList.add("active");
});

document.getElementById("create").addEventListener("click", async () => {

    const title = document.getElementById("titlewrite").value;
    const discription = document.getElementById("discrip2").value;

  /*  const info = {
        story_name: title,
        discription: discription,
        visibility: visibility
    };

    const response = await fetch("http://localhost:8080/api/editstory", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(info)
        */
       const formData = new FormData();
       formData.append("story_name", title);
       formData.append("description", discription);
       formData.append("visibility", visibility);
       const response = await fetch("https://localhost:8080/api/newstory", {
        method: "POST",
        body:formData
       });
    const data = await response.json();

    console.log(data);

});