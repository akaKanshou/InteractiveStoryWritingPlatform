function sendTo(location) {
    window.location.href = location
}

function check(min, max) {
    let input = document.getElementById("usernameInput")
    if (!input) return false

    let str = input.value.toLowerCase()
    if ((str.length < min) || (str.length > max)) {
        console.log("Username should be 3-18 characters long")
        return false
    }

    for (let char in str) {
        if ((char >= '0') && (char <= '9')) {
            continue
        }

        if ((char >= 'a') && (char <= 'z')) {
            continue
        }

        console.log("Username contains invalid characters")
        return false
    }

    return true
}

async function checkAndSubmit() {
    if (!check(3, 18)) return
    if (!check(1, 26)) return

    const formData = new FormData()
    formData.append("username", document.getElementById("usernameInput").value.toLowerCase())
    formData.append("displayname", document.getElementById("displayName").value)

    const req = new Request("https://localhost:8080/user/create", {
        method: "POST",
        body: formData,
    })

    const resp = await fetch(req)
    if (resp.status === 201) {
        console.log("registered")
    } else {
        const errString = await resp.text()
        console.log(errString)
    }
}