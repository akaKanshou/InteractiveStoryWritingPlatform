# ForgeVerse — Database Documentation

> **Stack:** MySQL (MariaDB 10.4) · PHP 8.0 · XAMPP · Database Name: `usr_d`

---

## Connection Setup

**File:** [`database/connection.php`](../database/connection.php)

```php
$host     = "localhost";
$user     = "root";
$password = "";          // blank for local XAMPP
$database = "usr_d";

$conn = new mysqli($host, $user, $password, $database);
$conn->set_charset("utf8mb4");
```

Every API file includes this with:
```php
require_once '../database/connection.php';
```

The `$conn` variable is then available globally in that file.

---

## Naming Convention

> **Important:** All table and column names use a double-underscore wrapper — `__name__`.
> This is intentional to avoid collisions with MySQL reserved words.

| Pattern | Example |
|---|---|
| Table name | `__story_tale__` |
| Primary key | `__SID__`, `__USID__`, `__SCENE_ID__` |
| Foreign key | `__USID__`, `__STORY_ID__`, `__FROM_SCENE__` |
| Field names | `__TITILE__`, `__EMAIL__`, `__POS_X__` |

---

## Tables Overview

```
usr_d
├── __usr_tale__       ← Users
├── __story_tale__     ← Stories
├── __scene_tale__     ← Scenes (nodes in the graph)
└── __edge_tale__      ← Connections between scenes (graph edges)
```

---

## Table 1: `__usr_tale__` — Users

Stores all registered user accounts.

| Column | Type | Notes |
|---|---|---|
| `__USID__` | `int` AUTO_INCREMENT PK | User ID |
| `__USER_NAME__` | `varchar(50)` UNIQUE | Username |
| `__EMAIL__` | `varchar(100)` UNIQUE | Email address |
| `__PASSWORD__` | `varchar(255)` | Bcrypt hashed password |
| `__P_PIC__` | `varchar(255)` | Profile picture filename (in `assets/uploads/`) |
| `__ORIGIN__` | `timestamp` | Account creation time |
| `__LAST_SEEN__` | `timestamp` | Last login timestamp |
| `__STATUS__` | `enum('Active','Blocked')` | Default: `Active` |
| `__AUTH_TYPE__` | `enum('LOCAL','GOOGLE')` | Login method |
| `__GOOGLE_ID__` | `varchar(255)` | Google OAuth ID (blank for LOCAL) |
| `__EMAIL__VERIFIED__` | `tinyint(1)` | Email verified flag (0/1) |

**Unique Indexes:** `__USER_NAME__`, `__EMAIL__`

**Session fields stored on login:**

```php
$_SESSION['user_id']     = $row['__USID__'];
$_SESSION['username']    = $row['__USER_NAME__'];
$_SESSION['email']       = $row['__EMAIL__'];
$_SESSION['profile_pic'] = $row['__P_PIC__'];
```

---

## Table 2: `__story_tale__` — Stories

Each row is one story belonging to one user.

| Column | Type | Notes |
|---|---|---|
| `__SID__` | `int` AUTO_INCREMENT PK | Story ID |
| `__USID__` | `int` FK → `__usr_tale__.__USID__` | Owner's user ID |
| `__TITILE__` | `varchar(100)` | Story title *(note: intentional typo in DB — do not rename)* |
| `__DESCRIPTION__` | `text` | Story description |
| `__THUMBNAIL__` | `varchar(255)` | Cover image filename in `assets/uploads/` |
| `__GENRE__` | `varchar(50)` | e.g. `Fantasy`, `Horror`, `Sci-Fi` |
| `__VISIBILITY__` | `enum('PUBLIC','PRIVATE')` | Story visibility |
| `__SCENE_COUNT__` | `int` | Auto-synced by `save_graph.php` on every graph save |
| `__CREATED__` | `timestamp` | Creation time |
| `__UPDATED__` | `timestamp` | Last update time |

> **`__SCENE_COUNT__`** is kept in sync automatically every time `save_graph.php` runs:
> ```sql
> UPDATE __story_tale__
> SET __SCENE_COUNT__ = (SELECT COUNT(*) FROM __scene_tale__ WHERE __STORY_ID__ = ?)
> WHERE __SID__ = ?
> ```

---

## Table 3: `__scene_tale__` — Scenes (Graph Nodes)

Each row is one scene (visual node) in a story's graph.

| Column | Type | Notes |
|---|---|---|
| `__SCENE_ID__` | `int` AUTO_INCREMENT PK | Scene ID |
| `__STORY_ID__` | `int` FK → `__story_tale__.__SID__` CASCADE DELETE | Parent story |
| `__SCENE_NAME__` | `varchar(100)` | Scene title shown on graph node |
| `__SCENE_TEXT__` | `longtext` NULL | Scene content **and** node type metadata stored as JSON |
| `__POS_X__` | `int` DEFAULT 0 | Canvas X position (pixels) |
| `__POS_Y__` | `int` DEFAULT 0 | Canvas Y position (pixels) |
| `__CREATED__` | `timestamp` | Row creation time |
| `__UPDATED__` | `timestamp` ON UPDATE | Auto-updates on every change |

### Node Type Detection via `__SCENE_TEXT__`

When a new special node (start/end) is created on the frontend, `save_graph.php` stores the type as JSON:

```json
{ "type": "start" }
{ "type": "end" }
```

`load_graph.php` reads this JSON back to restore the correct node type:

```php
$data = json_decode($row['__SCENE_TEXT__'], true);
if (isset($data['type'])) {
    $type = $data['type'];  // 'start', 'end', or 'scene'
}
```

**Fallback:** If a scene is named `"Start"` and no type JSON is stored, it is treated as `start`. Same for `"End"`.

### Position Storage

Every time a node is dragged or created, the graph auto-saves X/Y to the DB. This means **reopening the graph always restores the exact previous layout**.

---

## Table 4: `__edge_tale__` — Connections (Graph Edges)

Each row is one directed connection between two scenes, representing a story choice.

| Column | Type | Notes |
|---|---|---|
| `__EDGE_ID__` | `int` AUTO_INCREMENT PK | Edge ID |
| `__STORY_ID__` | `int` FK → `__story_tale__.__SID__` CASCADE DELETE | Parent story |
| `__FROM_SCENE__` | `int` FK → `__scene_tale__.__SCENE_ID__` CASCADE DELETE | Source scene (node the arrow leaves from) |
| `__TO_SCENE__` | `int` FK → `__scene_tale__.__SCENE_ID__` CASCADE DELETE | Target scene (node the arrow points to) |
| `__CHOICE_TEXT__` | `varchar(255)` | Label shown on the connection arrow |
| `__CREATED__` | `timestamp` | Creation time |

**Key rules for `__CHOICE_TEXT__`:**

| Situation | Value stored |
|---|---|
| Connection from Start node | `""` (empty string — no choice ever needed) |
| Single connection from scene node | `""` (optional — user may leave it blank) |
| Multiple connections from same scene | Non-empty string required (enforced by frontend) |

---

## Foreign Key Relationships & Cascade Rules

```
__usr_tale__
    │
    └─── __story_tale__  (no FK defined — deletion handled by API logic)
              │
              ├─── __scene_tale__   FK: __STORY_ID__ → __SID__  ON DELETE CASCADE
              │
              └─── __edge_tale__    FK: __STORY_ID__ → __SID__  ON DELETE CASCADE
                        │
                        ├─── FK: __FROM_SCENE__ → __SCENE_ID__  ON DELETE CASCADE
                        └─── FK: __TO_SCENE__   → __SCENE_ID__  ON DELETE CASCADE
```

> **Cascade Delete means:**
> - Delete a **story** → all its scenes and all edges are deleted automatically by MySQL
> - Delete a **scene** → all edges that reference it (from or to) are deleted automatically

---

## Save Graph Logic (`api/save_graph.php`)

The save works as a **full sync** approach — edges are always replaced, scenes are diffed:

1. **Verify ownership** — confirm `story_id` belongs to `$_SESSION['user_id']`
2. **Begin transaction**
3. **Scenes:**
   - Numeric ID → UPDATE existing scene (name, x, y)
   - Temp string ID (e.g. `temp_1234_56`) → INSERT new scene, record new DB ID in `$idMapping`
   - Any scene in DB but NOT in the payload → DELETE it
4. **Edges:** DELETE all edges for this story, then re-INSERT all from payload (using `$idMapping` to resolve temp IDs)
5. **Update `__SCENE_COUNT__`** in `__story_tale__`
6. **Commit** — return `{ status: 'success', idMap: {...} }`

---

## SQL Dump File

**File:** [`database/usr_d.sql`](../database/usr_d.sql)

Import this to recreate the entire schema from scratch:

1. Open `http://localhost/phpmyadmin`
2. Create database named `usr_d`
3. Click **Import** → select `usr_d.sql` → Execute

---

## Planned Future Tables

| Table | Purpose |
|---|---|
| `__character_tale__` | Character profiles with personality & voice |
| `__object_tale__` | Reusable object library |
| `__choice_tale__` | Extended choice conditions and variables |
| `__asset_tale__` | Uploaded media assets per scene |
| `__history_tale__` | Version history and undo log |

---

## Developer Rules

> [!IMPORTANT]
> - **Never** create one table per story. All stories share the same tables, filtered by `__STORY_ID__` / `__USID__`.
> - **Always** use prepared statements with `bind_param()` — never concatenate user input into SQL.
> - **Always** verify ownership before any write by checking `__USID__ = $_SESSION['user_id']`.
> - The column `__TITILE__` is a typo in the original schema — **do not rename it** without a migration script.
