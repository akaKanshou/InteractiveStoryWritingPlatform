
document.addEventListener('DOMContentLoaded', () => {


    const navProfile = document.getElementById('navProfile');
    const profileDropdown = document.getElementById('profileDropdown');
    const logoutBtn = document.getElementById('logoutBtn');
    const searchInput = document.getElementById('searchInput');
    const notificationBtn = document.getElementById('notificationBtn');

    const storiesGrid = document.getElementById('storiesGrid');
    const storiesLoading = document.getElementById('storiesLoading');
    const emptyState = document.getElementById('emptyState');
    const createStoryBtn = document.getElementById('createStoryBtn');
    const emptyCreateBtn = document.getElementById('emptyCreateBtn');

    const totalStoriesEl = document.getElementById('totalStories');
    const totalScenesEl = document.getElementById('totalScenes');
    const publicStoriesEl = document.getElementById('publicStories');
    const privateStoriesEl = document.getElementById('privateStories');

    const storyModalBackdrop = document.getElementById('storyModalBackdrop');
    const storyModal = document.getElementById('storyModal');
    const storyModalTitle = document.getElementById('storyModalTitle');
    const storyModalClose = document.getElementById('storyModalClose');
    const storyModalCancel = document.getElementById('storyModalCancel');
    const storyModalSave = document.getElementById('storyModalSave');
    const storyForm = document.getElementById('storyForm');
    const editStoryId = document.getElementById('editStoryId');
    const storyTitle = document.getElementById('storyTitle');
    const storyDescription = document.getElementById('storyDescription');
    const storyGenre = document.getElementById('storyGenre');
    const storyVisibility = document.getElementById('storyVisibility');
    const thumbnailUpload = document.getElementById('thumbnailUpload');
    const thumbnailInput = document.getElementById('thumbnailInput');
    const thumbnailPreview = document.getElementById('thumbnailPreview');
    const thumbnailImg = document.getElementById('thumbnailImg');
    const thumbnailPlaceholder = thumbnailPreview.querySelector('.thumbnail-placeholder');
    const descCharCount = document.getElementById('descCharCount');

    const deleteModalBackdrop = document.getElementById('deleteModalBackdrop');
    const deleteModal = document.getElementById('deleteModal');
    const deleteModalClose = document.getElementById('deleteModalClose');
    const deleteModalCancel = document.getElementById('deleteModalCancel');
    const deleteConfirmBtn = document.getElementById('deleteConfirmBtn');
    const deleteStoryId = document.getElementById('deleteStoryId');
    const deleteStoryName = document.getElementById('deleteStoryName');

    const filterBtns = document.querySelectorAll('.filter-btn');

    const toastContainer = document.getElementById('toastContainer');

    let allStories = [];
    let currentFilter = 'all';

    const genreEmojis = {
        'Fantasy': '🏰',
        'Adventure': '⚔️',
        'Romance': '💕',
        'Sci-Fi': '🚀',
        'Mystery': '🔍',
        'Horror': '👻',
        'Comedy': '😄',
        'Drama': '🎭',
        '': '📖'
    };

    navProfile.addEventListener('click', (e) => {
        e.stopPropagation();
        navProfile.classList.toggle('open');
        profileDropdown.classList.toggle('open');
    });

    document.addEventListener('click', (e) => {
        if (!navProfile.contains(e.target)) {
            navProfile.classList.remove('open');
            profileDropdown.classList.remove('open');
        }
    });

    logoutBtn.addEventListener('click', async (e) => {
        e.preventDefault();
        try {
            await fetch('../api/logout.php', { method: 'POST' });
        } catch (err) { }
        window.location.href = 'authenticate.php';
    });

    notificationBtn.addEventListener('click', () => {
        showToast('Notifications coming soon!', 'warning');
    });

    async function loadStories() {
        storiesLoading.classList.remove('hidden');
        storiesGrid.classList.add('hidden');
        emptyState.classList.add('hidden');

        try {
            const response = await fetch('../api/getStories.php');
            const data = await response.json();

            if (data.success) {
                allStories = data.stories || [];
                updateStats();
                renderStories();
            } else {
                showToast(data.message || 'Failed to load stories', 'error');
                showEmptyState();
            }
        } catch (err) {
            showToast('Connection error. Please refresh.', 'error');
            showEmptyState();
        } finally {
            storiesLoading.classList.add('hidden');
        }
    }

    function updateStats() {
        const total = allStories.length;
        const publicCount = allStories.filter(s => s.visibility === 'PUBLIC').length;
        const privateCount = allStories.filter(s => s.visibility === 'PRIVATE').length;
        const scenesCount = allStories.reduce((sum, s) => sum + (parseInt(s.sceneCount) || 0), 0);

        animateCounter(totalStoriesEl, total);
        animateCounter(totalScenesEl, scenesCount);
        animateCounter(publicStoriesEl, publicCount);
        animateCounter(privateStoriesEl, privateCount);
    }

    function animateCounter(el, target) {
        const current = parseInt(el.textContent) || 0;
        if (current === target) return;
        const duration = 400;
        const start = performance.now();
        function tick(now) {
            const progress = Math.min((now - start) / duration, 1);
            el.textContent = Math.round(current + (target - current) * progress);
            if (progress < 1) requestAnimationFrame(tick);
        }
        requestAnimationFrame(tick);
    }

    function renderStories() {
        let stories = [...allStories];

        if (currentFilter !== 'all') {
            stories = stories.filter(s => s.visibility === currentFilter);
        }

        const query = searchInput.value.trim().toLowerCase();
        if (query) {
            stories = stories.filter(s =>
                s.title.toLowerCase().includes(query) ||
                s.description.toLowerCase().includes(query) ||
                s.genre.toLowerCase().includes(query)
            );
        }

        if (stories.length === 0 && allStories.length === 0) {
            showEmptyState();
            return;
        }

        if (stories.length === 0) {
            storiesGrid.innerHTML = '<div class="empty-state"><p class="text-muted">No stories match your search or filter.</p></div>';
            storiesGrid.classList.remove('hidden');
            emptyState.classList.add('hidden');
            return;
        }

        storiesGrid.innerHTML = stories.map((story, i) => createStoryCard(story, i)).join('');
        storiesGrid.classList.remove('hidden');
        emptyState.classList.add('hidden');

        bindCardEvents();
    }

    function createStoryCard(story, index) {
        const emoji = genreEmojis[story.genre] || '📖';
        const genreClass = story.genre ? story.genre.toLowerCase().replace('-', '') : '';
        const thumbnail = story.thumbnail
            ? `<img src="../assets/uploads/${story.thumbnail}" alt="${escapeHtml(story.title)}">`
            : `<div class="card-thumbnail-placeholder">${emoji}</div>`;
        const date = formatDate(story.created);

        return `
            <div class="story-card" data-id="${story.id}" style="animation-delay: ${index * 0.05}s">
                <div class="card-thumbnail">
                    ${thumbnail}
                    ${story.genre ? `<span class="badge ${genreClass} card-genre-badge">${story.genre}</span>` : ''}
                    <span class="card-visibility-badge">${story.visibility === 'PUBLIC' ? '🌐 Public' : '🔒 Private'}</span>
                </div>
                <div class="card-body">
                    <div class="card-title">${escapeHtml(story.title)}</div>
                    <div class="card-description">${escapeHtml(story.description || 'No description yet...')}</div>
                    <div class="card-meta">
                        <div class="card-meta-left">
                            <span class="card-meta-item">📅 ${date}</span>
                            <span class="card-meta-item">🎬 ${story.sceneCount || 0} scenes</span>
                        </div>
                        <div class="card-actions">
                            <button class="card-action-btn edit-btn" data-id="${story.id}" title="Edit">✏️</button>
                            <button class="card-action-btn delete card-delete-btn" data-id="${story.id}" data-title="${escapeHtml(story.title)}" title="Delete">🗑️</button>
                        </div>
                    </div>
                </div>
            </div>
        `;
    }

    function showEmptyState() {
        storiesGrid.classList.add('hidden');
        emptyState.classList.remove('hidden');
    }

    function bindCardEvents() {
        document.querySelectorAll('.story-card').forEach(card => {
            card.addEventListener('click', (e) => {
                if (e.target.closest('.card-action-btn')) return;
                const storyId = card.dataset.id;
                window.location.href = `flow.php?story_id=${storyId}`;
            });
        });

        document.querySelectorAll('.edit-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                e.stopPropagation();
                const id = btn.dataset.id;
                const story = allStories.find(s => s.id == id);
                if (story) openEditModal(story);
            });
        });

        document.querySelectorAll('.card-delete-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                e.stopPropagation();
                const id = btn.dataset.id;
                const title = btn.dataset.title;
                openDeleteModal(id, title);
            });
        });
    }

    searchInput.addEventListener('input', debounce(() => {
        renderStories();
    }, 300));

    filterBtns.forEach(btn => {
        btn.addEventListener('click', () => {
            filterBtns.forEach(b => b.classList.remove('active'));
            btn.classList.add('active');
            currentFilter = btn.dataset.filter;
            renderStories();
        });
    });

    function openCreateModal() {
        storyModalTitle.textContent = 'Create New Story';
        storyModalSave.querySelector('.btn-text').textContent = 'Create Story';
        editStoryId.value = '';
        storyForm.reset();
        storyVisibility.value = 'PRIVATE';
        updateVisibilityUI('PRIVATE');
        thumbnailImg.classList.add('hidden');
        thumbnailImg.src = '';
        thumbnailPlaceholder.classList.remove('hidden');
        descCharCount.textContent = '0';
        clearFormErrors();
        showModal(storyModalBackdrop, storyModal);
    }

    function openEditModal(story) {
        storyModalTitle.textContent = 'Edit Story';
        storyModalSave.querySelector('.btn-text').textContent = 'Save Changes';
        editStoryId.value = story.id;
        storyTitle.value = story.title;
        storyDescription.value = story.description || '';
        storyGenre.value = story.genre || '';
        storyVisibility.value = story.visibility || 'PRIVATE';
        updateVisibilityUI(story.visibility || 'PRIVATE');
        descCharCount.textContent = (story.description || '').length;

        if (story.thumbnail) {
            thumbnailImg.src = '../assets/uploads/' + story.thumbnail;
            thumbnailImg.classList.remove('hidden');
            thumbnailPlaceholder.classList.add('hidden');
        } else {
            thumbnailImg.classList.add('hidden');
            thumbnailImg.src = '';
            thumbnailPlaceholder.classList.remove('hidden');
        }

        clearFormErrors();
        showModal(storyModalBackdrop, storyModal);
    }

    createStoryBtn.addEventListener('click', openCreateModal);
    emptyCreateBtn.addEventListener('click', openCreateModal);

    storyModalClose.addEventListener('click', () => closeModal(storyModalBackdrop, storyModal));
    storyModalCancel.addEventListener('click', () => closeModal(storyModalBackdrop, storyModal));
    storyModalBackdrop.addEventListener('click', () => closeModal(storyModalBackdrop, storyModal));

    document.querySelectorAll('.vis-btn').forEach(btn => {
        btn.addEventListener('click', () => {
            const vis = btn.dataset.vis;
            storyVisibility.value = vis;
            updateVisibilityUI(vis);
        });
    });

    function updateVisibilityUI(vis) {
        document.querySelectorAll('.vis-btn').forEach(b => {
            b.classList.toggle('active', b.dataset.vis === vis);
        });
    }

    thumbnailUpload.addEventListener('click', () => thumbnailInput.click());

    thumbnailInput.addEventListener('change', (e) => {
        const file = e.target.files[0];
        if (file) {
            if (!file.type.startsWith('image/')) {
                showToast('Please select an image file', 'error');
                return;
            }
            if (file.size > 5 * 1024 * 1024) {
                showToast('Image must be under 5MB', 'error');
                return;
            }
            const reader = new FileReader();
            reader.onload = (ev) => {
                thumbnailImg.src = ev.target.result;
                thumbnailImg.classList.remove('hidden');
                thumbnailPlaceholder.classList.add('hidden');
            };
            reader.readAsDataURL(file);
        }
    });

    storyDescription.addEventListener('input', () => {
        descCharCount.textContent = storyDescription.value.length;
    });

    storyForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        clearFormErrors();

        const title = storyTitle.value.trim();
        const description = storyDescription.value.trim();
        const genre = storyGenre.value;
        const visibility = storyVisibility.value;
        const isEdit = !!editStoryId.value;

        let valid = true;
        if (!title) {
            showFormError('storyTitleError', 'Title is required');
            valid = false;
        }
        if (!genre) {
            showFormError('storyGenreError', 'Please select a genre');
            valid = false;
        }
        if (!valid) return;

        setLoading(storyModalSave, true);

        try {
            const formData = new FormData();
            formData.append('title', title);
            formData.append('description', description);
            formData.append('genre', genre);
            formData.append('visibility', visibility);

            if (thumbnailInput.files[0]) {
                formData.append('thumbnail', thumbnailInput.files[0]);
            }

            let url = '../api/createStory.php';
            if (isEdit) {
                url = '../api/updateStory.php';
                formData.append('storyId', editStoryId.value);
            }

            const response = await fetch(url, {
                method: 'POST',
                body: formData
            });

            const data = await response.json();

            if (data.success) {
                showToast(isEdit ? 'Story updated successfully!' : 'Story created successfully!', 'success');
                closeModal(storyModalBackdrop, storyModal);
                loadStories();
            } else {
                showToast(data.message || 'Operation failed', 'error');
            }
        } catch (err) {
            showToast('Connection error. Please try again.', 'error');
        } finally {
            setLoading(storyModalSave, false);
        }
    });

    function openDeleteModal(id, title) {
        deleteStoryId.value = id;
        deleteStoryName.textContent = title;
        showModal(deleteModalBackdrop, deleteModal);
    }

    deleteModalClose.addEventListener('click', () => closeModal(deleteModalBackdrop, deleteModal));
    deleteModalCancel.addEventListener('click', () => closeModal(deleteModalBackdrop, deleteModal));
    deleteModalBackdrop.addEventListener('click', () => closeModal(deleteModalBackdrop, deleteModal));

    deleteConfirmBtn.addEventListener('click', async () => {
        const id = deleteStoryId.value;
        setLoading(deleteConfirmBtn, true);

        try {
            const response = await fetch('../api/deleteStory.php', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ storyId: id })
            });

            const data = await response.json();

            if (data.success) {
                showToast('Story deleted successfully', 'success');
                closeModal(deleteModalBackdrop, deleteModal);
                loadStories();
            } else {
                showToast(data.message || 'Delete failed', 'error');
            }
        } catch (err) {
            showToast('Connection error. Please try again.', 'error');
        } finally {
            setLoading(deleteConfirmBtn, false);
        }
    });

    function showModal(backdrop, modal) {
        backdrop.classList.add('active');
        modal.classList.add('active');
        document.body.style.overflow = 'hidden';
    }

    function closeModal(backdrop, modal) {
        backdrop.classList.remove('active');
        modal.classList.remove('active');
        document.body.style.overflow = '';
    }

    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape') {
            closeModal(storyModalBackdrop, storyModal);
            closeModal(deleteModalBackdrop, deleteModal);
        }
    });

    function showFormError(errorId, message) {
        const el = document.getElementById(errorId);
        if (el) {
            el.textContent = message;
            el.classList.add('visible');
        }
    }

    function clearFormErrors() {
        document.querySelectorAll('.modal .input-error').forEach(el => {
            el.textContent = '';
            el.classList.remove('visible');
        });
        document.querySelectorAll('.modal .input-field').forEach(el => {
            el.classList.remove('error');
        });
    }

    function setLoading(btn, loading) {
        if (loading) {
            btn.disabled = true;
            btn.classList.add('loading');
            if (!btn.querySelector('.spinner')) {
                const spinner = document.createElement('div');
                spinner.classList.add('spinner');
                btn.appendChild(spinner);
            }
        } else {
            btn.disabled = false;
            btn.classList.remove('loading');
            const spinner = btn.querySelector('.spinner');
            if (spinner) spinner.remove();
        }
    }

    function showToast(message, type = 'success') {
        const icons = { success: '✓', error: '✕', warning: '⚠' };
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

    function escapeHtml(str) {
        const div = document.createElement('div');
        div.textContent = str;
        return div.innerHTML;
    }

    function formatDate(dateStr) {
        if (!dateStr) return '';
        const date = new Date(dateStr);
        const now = new Date();
        const diff = now - date;
        const days = Math.floor(diff / (1000 * 60 * 60 * 24));

        if (days === 0) return 'Today';
        if (days === 1) return 'Yesterday';
        if (days < 7) return days + 'd ago';
        if (days < 30) return Math.floor(days / 7) + 'w ago';

        return date.toLocaleDateString('en-US', {
            month: 'short',
            day: 'numeric',
            year: date.getFullYear() !== now.getFullYear() ? 'numeric' : undefined
        });
    }

    function debounce(fn, delay) {
        let timer;
        return (...args) => {
            clearTimeout(timer);
            timer = setTimeout(() => fn(...args), delay);
        };
    }

    loadStories();

});