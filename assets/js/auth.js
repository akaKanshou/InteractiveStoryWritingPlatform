

document.addEventListener('DOMContentLoaded', () => {

    const loginTabBtn = document.getElementById('loginTabBtn');
    const registerTabBtn = document.getElementById('registerTabBtn');
    const tabIndicator = document.getElementById('tabIndicator');
    const loginForm = document.getElementById('loginForm');
    const registerForm = document.getElementById('registerForm');
    const welcomeTitle = document.getElementById('welcomeTitle');
    const welcomeSubtitle = document.getElementById('welcomeSubtitle');
    const authFooterText = document.getElementById('authFooterText');
    const switchToRegister = document.getElementById('switchToRegister');
    const googleAuthBtn = document.getElementById('googleAuthBtn');
    const toastContainer = document.getElementById('toastContainer');

    const regPasswordInput = document.getElementById('regPassword');
    const passwordStrength = document.getElementById('passwordStrength');
    const strengthFill = document.getElementById('strengthFill');
    const strengthText = document.getElementById('strengthText');

    const heroParticles = document.getElementById('heroParticles');

    let currentTab = 'login';

    function createParticles() {
        const count = 20;
        for (let i = 0; i < count; i++) {
            const particle = document.createElement('div');
            particle.classList.add('particle');
            particle.style.left = Math.random() * 100 + '%';
            particle.style.top = (80 + Math.random() * 30) + '%';
            particle.style.animationDuration = (8 + Math.random() * 15) + 's';
            particle.style.animationDelay = Math.random() * 10 + 's';
            heroParticles.appendChild(particle);
        }
    }
    createParticles();

    function switchTab(tab) {
        currentTab = tab;

        loginTabBtn.classList.toggle('active', tab === 'login');
        registerTabBtn.classList.toggle('active', tab === 'register');

        tabIndicator.classList.toggle('right', tab === 'register');

        if (tab === 'login') {
            loginForm.classList.add('active');
            registerForm.classList.remove('active');
            welcomeTitle.textContent = 'Welcome Back!';
            welcomeSubtitle.textContent = 'Sign in to continue your creative journey';
            authFooterText.innerHTML = 'Don\'t have an account? <a href="#" id="switchToRegister">Sign Up</a>';
        } else {
            loginForm.classList.remove('active');
            registerForm.classList.add('active');
            welcomeTitle.textContent = 'Join ForgeVerse';
            welcomeSubtitle.textContent = 'Create your account and start writing';
            authFooterText.innerHTML = 'Already have an account? <a href="#" id="switchToLogin">Sign In</a>';
        }

        bindFooterLink();

        clearAllErrors();
    }

    loginTabBtn.addEventListener('click', () => switchTab('login'));
    registerTabBtn.addEventListener('click', () => switchTab('register'));

    function bindFooterLink() {
        const switchLink = document.getElementById('switchToRegister') || document.getElementById('switchToLogin');
        if (switchLink) {
            switchLink.addEventListener('click', (e) => {
                e.preventDefault();
                switchTab(currentTab === 'login' ? 'register' : 'login');
            });
        }
    }
    bindFooterLink();

    document.querySelectorAll('.password-toggle').forEach(btn => {
        btn.addEventListener('click', () => {
            const targetId = btn.dataset.target;
            const input = document.getElementById(targetId);
            if (input.type === 'password') {
                input.type = 'text';
                btn.textContent = '🙈';
            } else {
                input.type = 'password';
                btn.textContent = '👁';
            }
        });
    });

    if (regPasswordInput) {
        regPasswordInput.addEventListener('input', () => {
            const val = regPasswordInput.value;
            if (val.length === 0) {
                passwordStrength.classList.remove('visible');
                return;
            }
            passwordStrength.classList.add('visible');

            let score = 0;
            if (val.length >= 6) score++;
            if (val.length >= 10) score++;
            if (/[A-Z]/.test(val)) score++;
            if (/[0-9]/.test(val)) score++;
            if (/[^A-Za-z0-9]/.test(val)) score++;

            strengthFill.className = 'strength-fill';
            if (score <= 1) {
                strengthFill.classList.add('weak');
                strengthText.textContent = 'Weak';
                strengthText.style.color = 'var(--accent-red)';
            } else if (score === 2) {
                strengthFill.classList.add('fair');
                strengthText.textContent = 'Fair';
                strengthText.style.color = 'var(--accent-amber)';
            } else if (score === 3) {
                strengthFill.classList.add('good');
                strengthText.textContent = 'Good';
                strengthText.style.color = 'var(--accent-blue)';
            } else {
                strengthFill.classList.add('strong');
                strengthText.textContent = 'Strong';
                strengthText.style.color = 'var(--accent-green)';
            }
        });
    }

    function showError(inputId, errorId, message) {
        const input = document.getElementById(inputId);
        const error = document.getElementById(errorId);
        if (input) input.classList.add('error');
        if (error) {
            error.textContent = message;
            error.classList.add('visible');
        }
    }

    function clearError(inputId, errorId) {
        const input = document.getElementById(inputId);
        const error = document.getElementById(errorId);
        if (input) input.classList.remove('error');
        if (error) {
            error.textContent = '';
            error.classList.remove('visible');
        }
    }

    function clearAllErrors() {
        document.querySelectorAll('.input-field').forEach(f => f.classList.remove('error'));
        document.querySelectorAll('.input-error').forEach(e => {
            e.textContent = '';
            e.classList.remove('visible');
        });
    }

    function isValidEmail(email) {
        return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
    }

    document.querySelectorAll('.input-field').forEach(input => {
        input.addEventListener('focus', () => {
            input.classList.remove('error');
            const errorEl = input.closest('.input-group')?.querySelector('.input-error');
            if (errorEl) {
                errorEl.textContent = '';
                errorEl.classList.remove('visible');
            }
        });
    });

    loginForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        clearAllErrors();

        const email = document.getElementById('loginEmail').value.trim();
        const password = document.getElementById('loginPassword').value;

        let valid = true;
        if (!email) {
            showError('loginEmail', 'loginEmailError', 'Email is required');
            valid = false;
        } else if (!isValidEmail(email)) {
            showError('loginEmail', 'loginEmailError', 'Please enter a valid email');
            valid = false;
        }
        if (!password) {
            showError('loginPassword', 'loginPasswordError', 'Password is required');
            valid = false;
        }
        if (!valid) return;

        const submitBtn = document.getElementById('loginSubmit');
        setLoading(submitBtn, true);

        try {
            const response = await fetch('../api/login.php', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email, password })
            });

            const data = await response.json();

            if (data.success) {
                showToast('Welcome back! Redirecting...', 'success');
                setTimeout(() => {
                    window.location.href = 'dashboard.php';
                }, 800);
            } else {
                showToast(data.message || 'Login failed. Please try again.', 'error');
            }
        } catch (err) {
            showToast('Connection error. Please try again.', 'error');
        } finally {
            setLoading(submitBtn, false);
        }
    });


    registerForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        clearAllErrors();

        const username = document.getElementById('regUsername').value.trim();
        const email = document.getElementById('regEmail').value.trim();
        const password = document.getElementById('regPassword').value;
        const confirmPassword = document.getElementById('regConfirmPassword').value;

    
        let valid = true;
        if (!username) {
            showError('regUsername', 'regUsernameError', 'Username is required');
            valid = false;
        } else if (username.length < 3) {
            showError('regUsername', 'regUsernameError', 'Username must be at least 3 characters');
            valid = false;
        } else if (!/^[a-zA-Z0-9_]+$/.test(username)) {
            showError('regUsername', 'regUsernameError', 'Only letters, numbers, and underscores allowed');
            valid = false;
        }

        if (!email) {
            showError('regEmail', 'regEmailError', 'Email is required');
            valid = false;
        } else if (!isValidEmail(email)) {
            showError('regEmail', 'regEmailError', 'Please enter a valid email');
            valid = false;
        }

        if (!password) {
            showError('regPassword', 'regPasswordError', 'Password is required');
            valid = false;
        } else if (password.length < 6) {
            showError('regPassword', 'regPasswordError', 'Password must be at least 6 characters');
            valid = false;
        }

        if (!confirmPassword) {
            showError('regConfirmPassword', 'regConfirmPasswordError', 'Please confirm your password');
            valid = false;
        } else if (password !== confirmPassword) {
            showError('regConfirmPassword', 'regConfirmPasswordError', 'Passwords do not match');
            valid = false;
        }

        if (!valid) return;

      
        const submitBtn = document.getElementById('registerSubmit');
        setLoading(submitBtn, true);

        try {
            const response = await fetch('../api/register.php', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ username, email, password })
            });

            const data = await response.json();

            if (data.success) {
                showToast('Account created! Redirecting...', 'success');
                setTimeout(() => {
                    window.location.href = 'dashboard.php';
                }, 800);
            } else {
                showToast(data.message || 'Registration failed. Please try again.', 'error');
            }
        } catch (err) {
            showToast('Connection error. Please try again.', 'error');
        } finally {
            setLoading(submitBtn, false);
        }
    });

    googleAuthBtn.addEventListener('click', () => {
        showToast('Google Sign-In will be available soon!', 'warning');
    });

    function setLoading(btn, loading) {
        if (loading) {
            btn.disabled = true;
            btn.classList.add('loading');
            const spinner = document.createElement('div');
            spinner.classList.add('spinner');
            btn.appendChild(spinner);
        } else {
            btn.disabled = false;
            btn.classList.remove('loading');
            const spinner = btn.querySelector('.spinner');
            if (spinner) spinner.remove();
        }
    }

    function showToast(message, type = 'success') {
        const icons = {
            success: '✓',
            error: '✕',
            warning: '⚠'
        };

        const toast = document.createElement('div');
        toast.classList.add('toast', type);
        toast.innerHTML = `
            <span class="toast-icon">${icons[type] || '●'}</span>
            <span>${message}</span>
        `;
        toastContainer.appendChild(toast);

        setTimeout(() => {
            toast.classList.add('removing');
            setTimeout(() => toast.remove(), 300);
        }, 3500);
    }

});
