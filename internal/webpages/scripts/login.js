
document.getElementById("sign-in-btn").onclick = () => {
    let checkbox = document.getElementById("checkRememberMe")
    let query
    if ((checkbox == null) || (!checkbox.checked)) {
        query = "false"
    } else {
        query = "true"
    }
    window.location.href = "https://localhost:8080/auth/google?remember=" + query;
}
