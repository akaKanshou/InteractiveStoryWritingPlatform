# ForgeVerse — Dashboard Documentation

> **Files involved:**
> - [`pages/dashboard.php`](../pages/dashboard.php) — HTML page
> - [`assets/js/Dashboard.js`](../assets/js/Dashboard.js) — All frontend logic
> - [`assets/css/Dash.css`](../assets/css/Dash.css) — Dashboard styles
> - [`api/getStories.php`](../api/getStories.php) — Fetch stories
> - [`api/createStory.php`](../api/createStory.php) — Create new story
> - [`api/updateStory.php`](../api/updateStory.php) — Edit story metadata
> - [`api/deleteStory.php`](../api/deleteStory.php) — Delete story

---

## Overview

The Dashboard is the first screen after login. It lets the user:

- View all their stories as cards
- Create a new story (with cover, title, genre, description, visibility)
- Edit story metadata (same modal, different mode)
- Delete a story (with confirmation modal)
- Search and filter stories
- See stats (total stories, scenes, public, private)
- Navigate to the Story Graph by clicking any story card

---

## Authentication Guard

`dashboard.php` starts with a session check:

```php
session_start();
if (!isset($_SESSION['user_id'])) {
    header('Location: authenticate.php');
    exit;
}
```

If the user is not logged in, they are redirected to the login page immediately. No data is rendered.

---

## PHP Variables Injected Into JS

At the bottom of `dashboard.php`:

```html
<script>
    window.ForgeVerse = {
        userId:   <?php echo json_encode($userId); ?>,
        username: <?php echo json_encode($username); ?>
    };
</script>
```

`Dashboard.js` reads `window.ForgeVerse` for any user-specific logic.

---

## Page Layout

```
┌──────────────────────────────────────────────────┐
│  [Logo] ForgeVerse    [🔍 Search]    [🔔] [User▾] │  ← Navbar
├──────────────────────────────────────────────────┤
│  Welcome back, Username ✨        [✦ Create Story] │  ← Header
├──────────────────────────────────────────────────┤
│  📖 Total │ 📝 Scenes │ 🌐 Public │ 🔒 Private    │  ← Stats Row
├──────────────────────────────────────────────────┤
│  Your Stories   [All] [Public] [Private]          │  ← Filter Bar
│                                                  │
│  ┌────────┐  ┌────────┐  ┌────────┐              │  ← Story Cards Grid
│  │  Card  │  │  Card  │  │  Card  │              │
│  └────────┘  └────────┘  └────────┘              │
└──────────────────────────────────────────────────┘
```

---

## Story Card Structure

Each story is rendered as a card:

```
┌──────────────────────────┐
│  [Cover Image / Emoji]   │  ← Genre badge + Visibility badge overlay
│  [Genre] [Public/Private] │
├──────────────────────────┤
│  Story Title             │
│  Description text...     │
│  📅 Created   🎬 Scenes  │
│              [✏️]  [🗑️] │
└──────────────────────────┘
```

- **Clicking the card** (not the buttons) → opens the Story Graph: `flow.php?story_id={id}`
- **✏️ Edit button** → opens the Edit Story modal, pre-filled with story data
- **🗑️ Delete button** → opens the Delete confirmation modal

Cards are generated in `createStoryCard(story, index)` inside `Dashboard.js`.

---

## Stats Row

Stats are calculated from the already-loaded `allStories` array — **no extra API call**.

| Stat | Source |
|---|---|
| Total Stories | `allStories.length` |
| Total Scenes | Sum of `story.sceneCount` across all stories |
| Published | Count where `visibility === 'PUBLIC'` |
| Private | Count where `visibility === 'PRIVATE'` |

All numbers animate from 0 to target using `animateCounter()`.

---

## Search & Filter

**Search** (`#searchInput`): Filters by title, description, and genre — debounced 300ms.

**Filter buttons** (`[All]`, `[Public]`, `[Private]`): Filters `allStories` by visibility.

Both filters apply simultaneously — search and filter work together.

No API call is made on filter/search. Everything works on the in-memory `allStories` array loaded once at startup.

---

## Create / Edit Story Modal

The same modal (`#storyModal`) handles both **create** and **edit** — the title and button text change.

### Create Mode

Triggered by: `#createStoryBtn` or `#emptyCreateBtn`

Form fields:
| Field | ID | Notes |
|---|---|---|
| Cover image | `#thumbnailInput` | Optional. `jpg/png/gif/webp`, max 5MB |
| Story title | `#storyTitle` | Required |
| Description | `#storyDescription` | Optional, max 500 chars (live counter) |
| Genre | `#storyGenre` | Required (Fantasy, Adventure, etc.) |
| Visibility | `#storyVisibility` | Toggle: Private (default) / Public |

On submit → `POST ../api/createStory.php` (multipart/form-data)

### Edit Mode

Triggered by clicking the ✏️ button on a story card.

`openEditModal(story)` pre-fills all fields. The hidden `#editStoryId` is set to the story's ID.

On submit → `POST ../api/updateStory.php` (multipart/form-data with `storyId`)

### Thumbnail Handling

- Clicking the upload area triggers the hidden `<input type="file">`.
- Preview is shown via `FileReader.readAsDataURL()` before upload.
- Server renames the file to `md5(time + original_name).ext` and stores it in `assets/uploads/`.
- On edit, the old thumbnail file is deleted from disk when a new one is uploaded.

---

## Delete Story Modal

Triggered by the 🗑️ button on a story card.

1. Shows the story name in a confirmation message
2. On confirm → `POST ../api/deleteStory.php` with `{ storyId: id }` as JSON
3. The server also deletes the thumbnail image from `assets/uploads/` if it exists
4. All scenes and edges are deleted automatically via MySQL CASCADE

---

## Navbar — Profile Dropdown

The avatar/username in the top-right opens a dropdown:
- **My Profile** — placeholder link
- **Settings** — placeholder link
- **Sign Out** → calls `POST ../api/logout.php` then redirects to `authenticate.php`

Closes on outside click via a `document.addEventListener('click', ...)` listener.

---

## API Reference

### `GET api/getStories.php`
Returns all stories for the logged-in user, ordered by `__UPDATED__` DESC.

**Response:**
```json
{
  "success": true,
  "stories": [
    {
      "id": 1,
      "title": "My Story",
      "description": "...",
      "thumbnail": "abc123.png",
      "genre": "Fantasy",
      "visibility": "PRIVATE",
      "sceneCount": 5,
      "created": "2026-07-11 15:49:48",
      "updated": "2026-07-11 16:00:00"
    }
  ]
}
```

### `POST api/createStory.php`
**Method:** POST (multipart/form-data)

**Fields:** `title`, `description`, `genre`, `visibility`, `thumbnail` (file, optional)

**Response:** `{ "success": true, "message": "Story created successfully" }`

### `POST api/updateStory.php`
**Method:** POST (multipart/form-data)

**Fields:** `storyId`, `title`, `description`, `genre`, `visibility`, `thumbnail` (file, optional)

**Response:** `{ "success": true, "message": "Story updated successfully" }`

### `POST api/deleteStory.php`
**Method:** POST (application/json)

**Body:** `{ "storyId": 1 }`

**Response:** `{ "success": true, "message": "Story deleted successfully" }`

> All APIs verify `$_SESSION['user_id']` is set and that the story belongs to that user before doing anything.

---

## Toast Notification System

`showToast(message, type)` — displays a non-blocking popup at the bottom of the screen.

| Type | Color | Icon |
|---|---|---|
| `'success'` | Green | ✓ |
| `'error'` | Red | ✕ |
| `'warning'` | Amber | ⚠ |

Toasts disappear after 3.5 seconds with a fade-out animation. They stack vertically.

---

## Utility Functions in Dashboard.js

| Function | Purpose |
|---|---|
| `loadStories()` | Fetches stories from API, updates stats, renders cards |
| `renderStories()` | Applies current filter + search, builds card HTML |
| `createStoryCard(story, i)` | Returns HTML string for one story card |
| `updateStats()` | Recalculates and animates the 4 stat counters |
| `animateCounter(el, target)` | Smooth count-up animation using `requestAnimationFrame` |
| `openCreateModal()` | Resets and opens the story modal in create mode |
| `openEditModal(story)` | Fills and opens the story modal in edit mode |
| `openDeleteModal(id, title)` | Opens delete confirmation modal |
| `showToast(message, type)` | Creates and auto-removes a toast notification |
| `formatDate(dateStr)` | Returns "Today", "Yesterday", "3d ago", or `Jul 12` |
| `escapeHtml(str)` | Safely escapes HTML to prevent XSS in card rendering |
| `debounce(fn, delay)` | Wraps a function to delay execution (used on search) |
| `setLoading(btn, bool)` | Adds/removes spinner on submit buttons |

---

## Flow: Create a Story

```
User clicks "Create New Story"
    ↓
openCreateModal() — resets form, opens modal
    ↓
User fills in title, genre, optional thumbnail/description/visibility
    ↓
storyForm submit → validation (title + genre required)
    ↓
POST ../api/createStory.php (FormData)
    ↓
Server saves to __story_tale__, renames thumbnail → assets/uploads/
    ↓
{ success: true }
    ↓
showToast("Story created!") + closeModal() + loadStories()
    ↓
New card appears in the grid
```

---

## Flow: Open a Story

```
User clicks a story card (not the edit/delete buttons)
    ↓
window.location.href = `flow.php?story_id=${id}`
    ↓
Story Graph page loads
```
