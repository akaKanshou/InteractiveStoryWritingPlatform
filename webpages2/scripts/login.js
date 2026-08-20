document.getElementById("googleAuthBtn").addEventListener("click", function() {
    let checkbox = document.getElementById("rememberMe")
    let query
    if ((checkbox === null) || (!checkbox.checked)) {
        query = "false"
    } else {
        query = "true"
    }
    window.location = `https://localhost:8080/auth/google?remember=${query}`;
});