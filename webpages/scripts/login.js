function sendTo(location) {
    window.location.href = `https://localhost:8080${location}`
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

    // Setting on clicks
    document.getElementById("googleAuthBtn").onclick = () => {
        let checkbox = document.getElementById("rememberMe")
        let query
        if ((checkbox === null) || (!checkbox.checked)) {
            query = "false"
        } else {
            query = "true"
        }
        sendTo("/auth/google?remember=" + query);
    }

    document.getElementById("hero-logoid").onclick = () => {
        sendTo("/");
    }
})