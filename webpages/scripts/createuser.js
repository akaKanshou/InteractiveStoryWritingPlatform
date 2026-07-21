function sendTo(location) {
    window.location.href = `https://localhost:8080${location}`
}

function showError(errorString) {
    const errBox = document.getElementById("usernameErrors")
    if (!errBox) return

    for (let i = 0; i < errorString.length; i++) {
        const li = document.createElement("li");
        li.textContent = errorString[i];
        errBox.appendChild(li);
    }

    errBox.classList.add('visible');
}

function clearError() {
    const errBox = document.getElementById("usernameErrors")
    if (!errBox) return

    errBox.innerHTML = ""

    errBox.classList.remove('visible');
}

function checkUsername() {
    const username = document.getElementById("regUsername").value.toLowerCase()
    if (/^[a-z0-9_]{3,18}$/.test(username)) {
        return true
    }

    let errString = []

    if (!/^[a-z0-9_]$/.test(username)) {
        errString.push("Username must only contain numbers, letters and underscores.")
    }

    if (!/^.{3,18}$/.test(username)) {
        errString.push("Username must be 3 to 18 characters long.")
    }

    showError(errString)

    return false
}

async function checkAndSubmit() {
    clearError()
    if (!checkUsername()) return

    const formData = new FormData()
    formData.append("username", document.getElementById("regUsername").value.toLowerCase())

    const req = new Request("https://localhost:8080/auth/register", {
        method: "POST",
        body: formData,
    })

    const resp = await fetch(req)
    if (resp.status === 201) {
        sendTo("/")
    } else {
        const errString = await resp.json()
        showError([errString.message])
    }
}


addEventListener('DOMContentLoaded', () => {
    // Creating particles
    const count = 20;
    const heroParticles = document.getElementById('heroParticles');
    for (let i = 0; i < count; i++) {
        const particle = document.createElement('div');
        particle.classList.add('particle');
        particle.style.left = Math.random() * 100 + '%';
        particle.style.top = (80 + Math.random() * 30) + '%';
        particle.style.animationDuration = (8 + Math.random() * 15) + 's';
        particle.style.animationDelay = Math.random() * 10 + 's';
        heroParticles.appendChild(particle);
    }

    // Setting on click
    const submitBtn = document.getElementById("registerSubmit")
    if (submitBtn) {
        submitBtn.addEventListener('click', checkAndSubmit)
    }
})