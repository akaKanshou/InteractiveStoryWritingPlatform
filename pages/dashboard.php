<?php

session_start();

if (!isset($_SESSION['user_id'])) {
    header('Location: authenticate.php');
    exit;
}

$username = htmlspecialchars($_SESSION['username'] ?? 'Writer');
$email = htmlspecialchars($_SESSION['email'] ?? '');
$profilePic = $_SESSION['profile_pic'] ?? '';
$userId = $_SESSION['user_id'];
?>
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>ForgeVerse Dashboard</title>
    <meta name="description" content="Your storytelling dashboard. Create, manage, and explore your interactive stories.">
    <link rel="stylesheet" href="../assets/css/global.css">
    <link rel="stylesheet" href="../assets/css/Dash.css">
</head>
<body>

    <div class="toast-container" id="toastContainer"></div>

    <nav class="dashboard-nav glass-strong" id="dashboardNav">
        <div class="nav-left">
            <div class="nav-logo">
                <svg width="32" height="32" viewBox="0 0 48 48" fill="none">
                    <defs>
                        <linearGradient id="navLogoGrad" x1="0%" y1="0%" x2="100%" y2="100%">
                            <stop offset="0%" stop-color="#a78bfa"/>
                            <stop offset="100%" stop-color="#6366f1"/>
                        </linearGradient>
                    </defs>
                    <path d="M24 4L6 14v20l18 10 18-10V14L24 4z" fill="url(#navLogoGrad)" opacity="0.2"/>
                    <path d="M24 4L6 14v20l18 10 18-10V14L24 4z" stroke="url(#navLogoGrad)" stroke-width="2" fill="none"/>
                    <path d="M16 20h16M16 24h12M16 28h8" stroke="#a78bfa" stroke-width="2" stroke-linecap="round"/>
                </svg>
                <span class="nav-brand">ForgeVerse</span>
            </div>
        </div>

        <div class="nav-center">
            <div class="nav-search">
                <span class="search-icon">🔍</span>
                <input type="text" id="searchInput" class="search-input" placeholder="Search your stories...">
            </div>
        </div>

        <div class="nav-right">
            <button class="btn-icon nav-action" id="notificationBtn" title="Notifications">
                <span>🔔</span>
                <span class="notification-dot"></span>
            </button>
            <div class="nav-profile" id="navProfile">
                <div class="profile-avatar" id="profileAvatar">
                    <?php echo strtoupper(substr($username, 0, 1)); ?>
                </div>
                <span class="profile-name"><?php echo $username; ?></span>
                <span class="profile-chevron">▾</span>

                <div class="profile-dropdown glass-strong" id="profileDropdown">
                    <div class="dropdown-header">
                        <div class="dropdown-avatar"><?php echo strtoupper(substr($username, 0, 1)); ?></div>
                        <div>
                            <div class="dropdown-name"><?php echo $username; ?></div>
                            <div class="dropdown-email"><?php echo $email; ?></div>
                        </div>
                    </div>
                    <div class="dropdown-divider"></div>
                    <a href="#" class="dropdown-item" id="profileLink">
                        <span>👤</span> My Profile
                    </a>
                    <a href="#" class="dropdown-item" id="settingsLink">
                        <span>⚙️</span> Settings
                    </a>
                    <div class="dropdown-divider"></div>
                    <a href="#" class="dropdown-item logout" id="logoutBtn">
                        <span>🚪</span> Sign Out
                    </a>
                </div>
            </div>
        </div>
    </nav>

    <main class="dashboard-main" id="dashboardMain">


        <section class="dash-header">
            <div class="dash-header-content">
                <h1 class="dash-title">Welcome back, <span class="text-accent"><?php echo $username; ?></span> ✨</h1>
                <p class="dash-subtitle">Continue your creative journey. Your stories await.</p>
            </div>
            <button class="btn btn-primary create-story-btn" id="createStoryBtn">
                <span>✦</span> Create New Story
            </button>
        </section>


        <section class="dash-stats" id="dashStats">
            <div class="stat-card glass-card">
                <div class="stat-icon" style="background: rgba(139,92,246,0.15); color: var(--accent-purple-light);">📖</div>
                <div class="stat-info">
                    <div class="stat-value" id="totalStories">0</div>
                    <div class="stat-label">Total Stories</div>
                </div>
            </div>
            <div class="stat-card glass-card">
                <div class="stat-icon" style="background: rgba(16,185,129,0.15); color: var(--accent-green);">📝</div>
                <div class="stat-info">
                    <div class="stat-value" id="totalScenes">0</div>
                    <div class="stat-label">Total Scenes</div>
                </div>
            </div>
            <div class="stat-card glass-card">
                <div class="stat-icon" style="background: rgba(59,130,246,0.15); color: var(--accent-blue);">🌐</div>
                <div class="stat-info">
                    <div class="stat-value" id="publicStories">0</div>
                    <div class="stat-label">Published</div>
                </div>
            </div>
            <div class="stat-card glass-card">
                <div class="stat-icon" style="background: rgba(236,72,153,0.15); color: var(--accent-pink);">🔒</div>
                <div class="stat-info">
                    <div class="stat-value" id="privateStories">0</div>
                    <div class="stat-label">Private</div>
                </div>
            </div>
        </section>

 
        <section class="dash-stories-section">
            <div class="section-header">
                <h2>Your Stories</h2>
                <div class="section-filters">
                    <button class="filter-btn active" data-filter="all">All</button>
                    <button class="filter-btn" data-filter="PUBLIC">Public</button>
                    <button class="filter-btn" data-filter="PRIVATE">Private</button>
                </div>
            </div>

  
            <div class="loading-spinner" id="storiesLoading"></div>

            <div class="empty-state hidden" id="emptyState">
                <div class="empty-illustration">📚</div>
                <h3>No stories yet</h3>
                <p>Every great tale starts with a single word. Create your first story and let your imagination run wild!</p>
                <button class="btn btn-primary" id="emptyCreateBtn">
                    <span>✦</span> Create Your First Story
                </button>
            </div>


            <div class="stories-grid hidden" id="storiesGrid">

            </div>
        </section>

    </main>

 
    <div class="modal-backdrop" id="storyModalBackdrop"></div>
    <div class="modal" id="storyModal">
        <div class="modal-header">
            <h3 id="storyModalTitle">Create New Story</h3>
            <button class="modal-close" id="storyModalClose">✕</button>
        </div>

        <form id="storyForm" novalidate>
            <input type="hidden" id="editStoryId" value="">

            <div class="thumbnail-upload" id="thumbnailUpload">
                <div class="thumbnail-preview" id="thumbnailPreview">
                    <span class="thumbnail-placeholder">
                        <span>🖼️</span>
                        <span>Click to upload cover</span>
                    </span>
                    <img src="" alt="" id="thumbnailImg" class="hidden">
                </div>
                <input type="file" id="thumbnailInput" accept="image/*" class="hidden">
            </div>

  
            <div class="input-group">
                <label class="input-label" for="storyTitle">Story Title</label>
                <input type="text" id="storyTitle" class="input-field" placeholder="Enter your story title" maxlength="100" required>
                <div class="input-error" id="storyTitleError"></div>
            </div>

            <div class="input-group">
                <label class="input-label" for="storyDescription">Description</label>
                <textarea id="storyDescription" class="input-field textarea" placeholder="Describe your story in a few words..." maxlength="500" rows="3"></textarea>
                <div class="char-counter"><span id="descCharCount">0</span>/500</div>
                <div class="input-error" id="storyDescriptionError"></div>
            </div>
            <div class="input-group">
                <label class="input-label" for="storyGenre">Genre</label>
                <select id="storyGenre" class="input-field select-field">
                    <option value="">Select a genre</option>
                    <option value="Fantasy">Fantasy</option>
                    <option value="Adventure">Adventure</option>
                    <option value="Romance">Romance</option>
                    <option value="Sci-Fi">Sci-Fi</option>
                    <option value="Mystery">Mystery</option>
                    <option value="Horror">Horror</option>
                    <option value="Comedy">Comedy</option>
                    <option value="Drama">Drama</option>
                </select>
                <div class="input-error" id="storyGenreError"></div>
            </div>

            <div class="input-group">
                <label class="input-label">Visibility</label>
                <div class="visibility-toggle">
                    <button type="button" class="vis-btn active" data-vis="PRIVATE" id="visPrivate">
                        <span>🔒</span> Private
                    </button>
                    <button type="button" class="vis-btn" data-vis="PUBLIC" id="visPublic">
                        <span>🌐</span> Public
                    </button>
                </div>
                <input type="hidden" id="storyVisibility" value="PRIVATE">
            </div>

            <div class="modal-actions">
                <button type="button" class="btn btn-secondary" id="storyModalCancel">Cancel</button>
                <button type="submit" class="btn btn-primary" id="storyModalSave">
                    <span class="btn-text">Create Story</span>
                </button>
            </div>
        </form>
    </div>

    <div class="modal-backdrop" id="deleteModalBackdrop"></div>
    <div class="modal delete-modal" id="deleteModal">
        <div class="modal-header">
            <h3>Delete Story</h3>
            <button class="modal-close" id="deleteModalClose">✕</button>
        </div>
        <p class="delete-message">Are you sure you want to delete "<strong id="deleteStoryName"></strong>"? This action cannot be undone.</p>
        <input type="hidden" id="deleteStoryId" value="">
        <div class="modal-actions">
            <button type="button" class="btn btn-secondary" id="deleteModalCancel">Cancel</button>
            <button type="button" class="btn btn-danger" id="deleteConfirmBtn">
                <span class="btn-text">Delete Story</span>
            </button>
        </div>
    </div>

    <script>
        window.ForgeVerse = {
            userId: <?php echo json_encode($userId); ?>,
            username: <?php echo json_encode($username); ?>
        };
    </script>
    <script src="../assets/js/Dashboard.js?v=<?php echo time(); ?>"></script>
</body>
</html>