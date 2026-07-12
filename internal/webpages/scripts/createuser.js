function sendTo(location) {
    window.location.href = location
}

function showError(errorString) {
    const errBox = document.getElementById("errorBox")
    if (!errBox) return

    errBox.textContent=errorString
    errBox.hidden=false
}

function checkUsername() {
    let input = document.getElementById("usernameInput")
    if (!input) return false

    let str = input.value.toLowerCase()
    if ((str.length < 3) || (str.length > 18)) {
        showError(`Username should be 3 to 18 characters long`)
        return false
    }

    for (let char in str) {
        if ((char >= '0') && (char <= '9')) {
            continue
        }

        if ((char >= 'a') && (char <= 'z')) {
            continue
        }

        showError("Username contains invalid characters")
        return false
    }

    return true
}

function checkDisplayName() {
    let input = document.getElementById("displayName")
    if (!input) return false

    let str = input.value.toLowerCase()
    if ((str.length < 1) || (str.length > 26)) {
        showError(`Username should be 1 to 26 characters long`)
        return false
    }

    for (let char in str) {
        if (char === ' ') {
            continue
        }

        if ((char >= '0') && (char <= '9')) {
            continue
        }

        if ((char >= 'a') && (char <= 'z')) {
            continue
        }

        showError("Display name contains invalid characters")
        return false
    }

    return true
}

async function checkAndSubmit() {
    if (!checkUsername()) return
    if (!checkDisplayName()) return

    const formData = new FormData()
    formData.append("username", document.getElementById("usernameInput").value.toLowerCase())
    formData.append("displayname", document.getElementById("displayName").value)

    const req = new Request("https://localhost:8080/user/create", {
        method: "POST",
        body: formData,
    })

    const resp = await fetch(req)
    if (resp.status === 201) {
        sendTo("https://localhost:8080/")
    } else {
        const errString = await resp.text()
        showError(errString)
    }
}