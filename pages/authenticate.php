<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome to ForgeVerse</title>
    <meta name="description" content="Sign in or create an account to access ForgeVerse AI — your interactive storytelling platform.">
    <link rel="stylesheet" href="../assets/css/global.css">
    <link rel="stylesheet" href="../assets/css/auth.css">
</head>
<body>

    <div class="toast-container" id="toastContainer"></div>

    <div class="auth-page">

        <div class="auth-hero" id="authHero">
            <div class="hero-particles" id="heroParticles"></div>
            <div class="hero-overlay"></div>
            <div class="hero-content">
                <div class="hero-logo">
                    <div class="logo-icon">
                        <svg width="48" height="48" viewBox="0 0 48 48" fill="none">
                            <defs>
                                <linearGradient id="logoGrad" x1="0%" y1="0%" x2="100%" y2="100%">
                                    <stop offset="0%" stop-color="#a78bfa"/>
                                    <stop offset="100%" stop-color="#6366f1"/>
                                </linearGradient>
                            </defs>
                            <path d="M24 4L6 14v20l18 10 18-10V14L24 4z" fill="url(#logoGrad)" opacity="0.2"/>
                            <path d="M24 4L6 14v20l18 10 18-10V14L24 4z" stroke="url(#logoGrad)" stroke-width="2" fill="none"/>
                            <path d="M16 20h16M16 24h12M16 28h8" stroke="#a78bfa" stroke-width="2" stroke-linecap="round"/>
                        </svg>
                    </div>
                    <div class="logo-text">
                        <h1>ForgeVerse</h1>
                    </div>
                </div>
                <p class="hero-tagline">Where Imagination Becomes Reality</p>
                <p class="hero-description">Create, explore, and experience stories like never before with the power of AI-driven interactive storytelling.</p>
                <div class="hero-features">

                    <div class="hero-feature">
                        <span class="feature-icon">⚡</span>
                        <span>Interactive Scene Editor</span>
                    </div>
                    <div class="hero-feature">
                        <span class="feature-icon">🌟</span>
                        <span>Visual Story Graphs</span>
                    </div>
                </div>
            </div>
        </div>

        <div class="auth-panel">
            <div class="auth-card glass-strong">

                <div class="auth-tabs" id="authTabs">
                    <button class="auth-tab active" data-tab="login" id="loginTabBtn">Sign In</button>
                    <button class="auth-tab" data-tab="register" id="registerTabBtn">Sign Up</button>
                    <div class="tab-indicator" id="tabIndicator"></div>
                </div>

                <div class="auth-welcome" id="authWelcome">
                    <h2 id="welcomeTitle">Welcome Back!</h2>
                    <p id="welcomeSubtitle">Sign in to continue your creative journey</p>
                </div>

                <form class="auth-form active" id="loginForm" novalidate>

                    <div class="input-group">
                        <label class="input-label" for="loginEmail">Email Address</label>
                        <div class="input-icon-wrapper">
                            <input type="email" id="loginEmail" class="input-field" placeholder="Enter your email" autocomplete="email" required>
                            <span class="icon-left">✉</span>
                        </div>
                        <div class="input-error" id="loginEmailError"></div>
                    </div>

                    <div class="input-group">
                        <label class="input-label" for="loginPassword">Password</label>
                        <div class="input-icon-wrapper">
                            <input type="password" id="loginPassword" class="input-field" placeholder="Enter your password" autocomplete="current-password" required>
                            <span class="icon-left">🔒</span>
                            <button type="button" class="password-toggle" data-target="loginPassword" aria-label="Toggle password visibility">👁</button>
                        </div>
                        <div class="input-error" id="loginPasswordError"></div>
                    </div>

                    <div class="auth-extras">
                        <label class="remember-me">
                            <input type="checkbox" id="rememberMe">
                            <span class="checkmark"></span>
                            <span>Remember me</span>
                        </label>
                        <a href="#" class="forgot-link">Forgot Password?</a>
                    </div>

                    <button type="submit" class="btn btn-primary auth-submit" id="loginSubmit">
                        <span class="btn-text">Sign In</span>
                    </button>

                </form>

                <form class="auth-form" id="registerForm" novalidate>

                    <div class="input-group">
                        <label class="input-label" for="regUsername">Username</label>
                        <div class="input-icon-wrapper">
                            <input type="text" id="regUsername" class="input-field" placeholder="Choose a username" autocomplete="username" required>
                            <span class="icon-left">👤</span>
                        </div>
                        <div class="input-error" id="regUsernameError"></div>
                    </div>

                    <div class="input-group">
                        <label class="input-label" for="regEmail">Email Address</label>
                        <div class="input-icon-wrapper">
                            <input type="email" id="regEmail" class="input-field" placeholder="Enter your email" autocomplete="email" required>
                            <span class="icon-left">✉</span>
                        </div>
                        <div class="input-error" id="regEmailError"></div>
                    </div>

                    <div class="input-group">
                        <label class="input-label" for="regPassword">Password</label>
                        <div class="input-icon-wrapper">
                            <input type="password" id="regPassword" class="input-field" placeholder="Create a password" autocomplete="new-password" required>
                            <span class="icon-left">🔒</span>
                            <button type="button" class="password-toggle" data-target="regPassword" aria-label="Toggle password visibility">👁</button>
                        </div>
                        <div class="input-error" id="regPasswordError"></div>

                        <div class="password-strength" id="passwordStrength">
                            <div class="strength-bar">
                                <div class="strength-fill" id="strengthFill"></div>
                            </div>
                            <span class="strength-text" id="strengthText"></span>
                        </div>
                    </div>

                    <div class="input-group">
                        <label class="input-label" for="regConfirmPassword">Confirm Password</label>
                        <div class="input-icon-wrapper">
                            <input type="password" id="regConfirmPassword" class="input-field" placeholder="Confirm your password" autocomplete="new-password" required>
                            <span class="icon-left">🔒</span>
                            <button type="button" class="password-toggle" data-target="regConfirmPassword" aria-label="Toggle password visibility">👁</button>
                        </div>
                        <div class="input-error" id="regConfirmPasswordError"></div>
                    </div>

                    <button type="submit" class="btn btn-primary auth-submit" id="registerSubmit">
                        <span class="btn-text">Create Account</span>
                    </button>

                </form>

                <div class="auth-divider">
                    <span>or continue with</span>
                </div>

                <button type="button" class="btn-google" id="googleAuthBtn">
                    <svg width="20" height="20" viewBox="0 0 24 24">
                        <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 0 1-2.2 3.32v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.1z" fill="#4285F4"/>
                        <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853"/>
                        <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05"/>
                        <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335"/>
                    </svg>
                    <span>Continue with Google</span>
                </button>

                <p class="auth-footer-text" id="authFooterText">
                    Don't have an account? <a href="#" id="switchToRegister">Sign Up</a>
                </p>

            </div>
        </div>

    </div>

    <script src="../assets/js/auth.js"></script>
</body>
</html>