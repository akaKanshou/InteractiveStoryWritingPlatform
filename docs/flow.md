# ForgeVerse — Story Graph (Flow) Documentation

> **Files involved:**
> - [`pages/flow.php`](../pages/flow.php) — HTML shell page for the graph
> - [`assets/js/graph.js`](../assets/js/graph.js) — All graph logic (750+ lines)
> - [`assets/css/graph.css`](../assets/css/graph.css) — Graph-specific styles
> - [`api/load_graph.php`](../api/load_graph.php) — Load scenes + edges for a story
> - [`api/save_graph.php`](../api/save_graph.php) — Save the entire graph state

---

## What Is the Story Graph?

The Story Graph is the **heart of ForgeVerse**. It is a visual, interactive node editor where a writer builds the **structure** of their interactive story.

- Each **node** = one Scene
- Each **arrow (edge)** = one Choice that leads from one scene to another
- The graph is **not** a flowchart or diagram tool — it specifically models story branching

> The graph manages **structure only**. The content of each scene is written in the separate Scene Editor (Sprint 4 — not yet built).

---

## How to Access

From the Dashboard, clicking any story card navigates to:

```
pages/flow.php?story_id=<ID>
```

`flow.php` validates the session and the `story_id` parameter, then injects them into JS:

```html
<script>
    window.ForgeVerse = {
        userId:  <?php echo json_encode($userId); ?>,
        storyId: <?php echo json_encode($storyId); ?>
    };
</script>
```

`graph.js` reads `window.ForgeVerse.storyId` to load the correct story's graph.

---

## Page Layout

```
┌──────────────────────────────────────────────────────────────┐
│ [◀ Back]  [Story Name]         [−] [100%] [+]  [🎯]  Saved ✔ │  ← Toolbar
├──────────────────────────────────────────────────────────────┤
│                                                              │
│         . . . . . . . . . . . . . . . . .                    │
│         .   [Start Node]                .                    │
│         .       ↓                      .                    │
│         .   [Scene 1] ────→ [Scene 2]  .                    │
│         .       ↓                      .                    │
│         .   [End Node]                 .   ← Infinite Canvas │
│         . . . . . . . . . . . . . . . .                    │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

The canvas takes up the **full screen**. The toolbar floats on top using `position: absolute`.

---

## Canvas Features

| Feature | How it works |
|---|---|
| **Pan** | `mousedown` on empty canvas → `mousemove` updates `state.panX/panY` |
| **Drag nodes** | `mousedown` on node → `mousemove` updates `node.x/node.y` |
| **Zoom (mouse wheel)** | `wheel` event → adjusts `state.zoom` and re-positions `panX/panY` to zoom toward cursor |
| **Zoom (toolbar buttons)** | `±0.2` per click, clamped `[0.2 → 3.0]` |
| **Center graph** | Calculates bounding box of all nodes, sets pan to center it |
| **Snap to grid** | Node positions are rounded to nearest 20px on drag |
| **Infinite canvas** | `graphNodesLayer` and `graphSvg` are both absolutely positioned and transformed |

The canvas is made of **3 layers** stacked by z-index:

| Layer | Element | z-index | Purpose |
|---|---|---|---|
| Background grid | `#graphGrid` | 0 | Dot grid pattern (CSS `radial-gradient`) |
| SVG connections | `#graphSvg` | 10 | Arrow paths between nodes |
| Nodes | `#graphNodesLayer` | 20 | Draggable node divs + choice labels |

All three layers receive the **same CSS transform** (`translate + scale`) to stay in sync.

---

## State Object

`graph.js` maintains a single `state` object:

```javascript
const state = {
    zoom: 1,                      // Current zoom level (0.2 to 3.0)
    panX: 0,                      // Canvas X offset in pixels
    panY: 0,                      // Canvas Y offset in pixels
    isPanning: false,             // True while panning the canvas
    startX: 0,                    // Mouse anchor X for pan/drag
    startY: 0,                    // Mouse anchor Y for pan/drag
    nodes: [],                    // Array of node objects { id, name, type, x, y }
    edges: [],                    // Array of edge objects { id, from, to, text }
    isDraggingNode: false,        // True while a node is being dragged
    draggedNode: null,            // Reference to the node being dragged
    connectionMode: false,        // True when user is drawing a new connection
    connectionSourceNodeId: null, // ID of node that initiated connection mode
    connectionTargetNodeId: null, // ID of node selected as connection target
    edgeBeingEdited: null,        // Edge object currently open in edit modal
    saveTimeout: null,            // Debounce timer reference for auto-save
    parentNodeIdForNewScene: null,// Which node's "+" button was clicked
    _choiceRequired: false        // Whether the connect modal's text field is mandatory
};
```

---

## Node Types

| Type | Icon | Color | Rules |
|---|---|---|---|
| `start` | 🟢 | Green border | 1 per story. Cannot be deleted. Cannot receive connections. No choice text ever. Title not clickable (void node). |
| `scene` | 🎬 | Purple border | The main story scenes. Can be connected to/from freely. Clicking title opens Scene Editor (future). |
| `end` | 🛑 | Pink border | Terminal node. Cannot have outgoing connections. |

---

## Node Anatomy

```
┌─────────────────────────┐
│ 🎬  Scene Name          │  ← Header (click name → opens Scene Editor)
├─────────────────────────┤
│  [+]  [🗑]  [🔗]    [●] │  ← Actions + Status indicator
└─────────────────────────┘
    ●  = port-in  (top centre)
    ●  = port-out (bottom centre)
```

For **Start node:**
```
┌─────────────────────────┐
│ 🟢  Start               │  ← Non-clickable title (void node)
├─────────────────────────┤
│  [+]  [🔗]          [●] │  ← No delete button. Status shows if connected.
└─────────────────────────┘
```

For **End node:**
```
┌─────────────────────────┐
│ 🛑  End                 │
├─────────────────────────┤
│  (no buttons)           │  ← No outgoing connections possible
└─────────────────────────┘
      ●  = port-in only
```

---

## Status Indicator

Every node has a small circular status badge:

| Badge | Meaning | Clicking it does... |
|---|---|---|
| 🟢 Green ✔ | Node is reachable / connected | Nothing (informational) |
| 🔴 Red ✕ (scene) | Scene has no incoming connections — unreachable | Shows alert: use 🔗 on another node |
| 🔴 Red ✕ (start) | Start node has no outgoing connections | Enters Connection Mode from Start |

**Condition for scene red X:**
```javascript
const incomingEdges = state.edges.filter(e => e.to == node.id);
if (incomingEdges.length === 0) → show red X
```

---

## Creating a Scene (+ Button)

1. Click **[+]** on any node
2. A modal opens: **"Create Scene"**
3. Enter scene name
4. Optionally enter choice text (see rules below)
5. Click **Create**

The new scene node is placed to the **right** of the parent node (`x + 250`). Collision avoidance moves it down 120px if another node is already there.

A new edge is automatically created from the parent to the new scene.

---

## Choice Text Rules

| Situation | Choice field | Validation |
|---|---|---|
| Source is Start node | Hidden | Never shown |
| First connection from a scene node | Optional | Shown with *(optional)* label |
| Additional connections from same node | Required | Shown with **\*** label, blocked if empty |

This applies to both the **Create Scene** modal and the **Connect Scene** modal.

---

## Connecting Existing Scenes (🔗 Button / Connection Mode)

1. Click **[🔗]** on any scene node (or the Start node)
2. The cursor changes to a crosshair — **Connection Mode is active**
3. A dashed animated line follows the mouse cursor from the source node
4. Click another node to target it
5. If source is **Start node** → edge created instantly (no modal, no choice text)
6. Otherwise → **"Connect Scene"** modal opens with optional/required choice text
7. Press **Escape** to cancel Connection Mode at any time

**Rules enforced during connection:**
- Cannot connect a node to itself
- Cannot connect to the Start node (no incoming connections on Start)
- Cannot create a duplicate connection (same from → to)

---

## Editing a Connection (Click the Arrow)

Click any arrow or its label to open **"Edit Choice"** modal:

```
┌──────────────────────┐
│  Edit Choice         │
│  Choice Text: [___]  │
│  [Delete]   [Save]   │
└──────────────────────┘
```

- **Save** → updates the edge's `text` in `state.edges`, re-renders, auto-saves
- **Delete** → removes the edge from `state.edges`, re-renders, auto-saves
- Deleting a connection **never** deletes either scene

---

## Auto-Save System

The graph **never has a Save button**. It saves automatically after every meaningful action:

| Action that triggers save |
|---|
| New scene created |
| Scene deleted |
| Scene dragged (on mouse release) |
| Connection created |
| Connection deleted |
| Connection text edited |

The save is **debounced 1 second** — if multiple actions happen quickly, only one save fires.

**Save status display:**

```
Saving...      ← yellow text, appears immediately
Saved ✔        ← green text, fades out after 2 seconds
Save Failed!   ← red text, stays visible
```

---

## Graph Rendering Pipeline

Every change calls `renderAll()`:

```
renderAll()
  ├── updateTransform()   → applies zoom + pan CSS transform to all layers
  ├── renderNodes()       → clears graphNodesLayer, builds node divs from state.nodes
  └── renderEdges()       → clears graphSvg, draws SVG paths from state.edges
                            + appends choice labels as divs to graphNodesLayer
```

Edges use **Bezier curves** (`getBezierPath`):
```javascript
`M ${x1} ${y1} C ${x1} ${y1+offset}, ${x2} ${y2-offset}, ${x2} ${y2}`
```

Each edge has two SVG elements:
- **`.connection-hitbox`** — transparent, 20px wide for easier clicking
- **`.connection-line`** — visible 3px purple stroke with arrowhead marker

---

## Graph Loading (`api/load_graph.php`)

`GET api/load_graph.php?story_id=<ID>`

1. Verifies session and ownership
2. Fetches all rows from `__scene_tale__` for this story → builds `nodes` array
3. Fetches all rows from `__edge_tale__` for this story → builds `edges` array
4. Returns story name + nodes + edges as JSON

**Node type detection during load:**
- Reads `__SCENE_TEXT__` as JSON → looks for `{ "type": "..." }`
- Fallback: scene named `"Start"` with no type → treated as `start`

**Response shape:**
```json
{
  "status": "success",
  "story_name": "My Epic Story",
  "nodes": [
    { "id": 1, "name": "Start", "type": "start", "x": 100, "y": 300 },
    { "id": 2, "name": "Forest Path", "type": "scene", "x": 350, "y": 300 }
  ],
  "edges": [
    { "id": 1, "from": 1, "to": 2, "text": "" }
  ]
}
```

If the story has **no scenes at all** (brand new story), `graph.js` auto-creates a Start node and auto-saves.

---

## Graph Saving (`api/save_graph.php`)

`POST api/save_graph.php`

**Request body (JSON):**
```json
{
  "story_id": 1,
  "nodes": [ { "id": 1, "name": "Start", "type": "start", "x": 100, "y": 300 }, ... ],
  "edges": [ { "id": 1, "from": 1, "to": 2, "text": "" }, ... ]
}
```

**ID handling:**
- Numeric IDs → existing rows → UPDATE
- String IDs (e.g. `"temp_1720123456_789"`) → new rows → INSERT, returns `idMap`

**Full save logic:**
1. Upsert all nodes (insert new, update existing)
2. Delete any scenes that no longer exist in the payload
3. Delete ALL edges for this story
4. Re-insert all edges (using `idMap` to resolve temp IDs to real DB IDs)
5. Update `__SCENE_COUNT__` in `__story_tale__`
6. Commit transaction

**Response:**
```json
{ "status": "success", "idMap": { "temp_1720123456_789": 5 } }
```

---

## Toolbar Buttons

| Button | ID | Action |
|---|---|---|
| ◀ Back | `#backBtn` | `window.location.href = 'dashboard.php'` |
| − (Zoom Out) | `#zoomOutBtn` | `state.zoom -= 0.2`, clamped to 0.2 |
| + (Zoom In) | `#zoomInBtn` | `state.zoom += 0.2`, clamped to 3.0 |
| 🎯 Center | `#centerGraphBtn` | Calculates bounding box of all nodes, centers view |
| Zoom display | `#zoomLevelDisplay` | Shows current zoom as `"100%"` |
| Save status | `#saveStatus` | Shows "Saving..." / "Saved ✔" / "Save Failed!" |

---

## Modals

Three modals live in `flow.php`:

### 1. `#addSceneModal` — Create Scene
| Field | ID | Notes |
|---|---|---|
| Scene Name | `#newSceneName` | Required |
| Choice Text | `#newSceneChoice` | Hidden for Start, optional/required for scenes |

### 2. `#connectSceneModal` — Connect to Existing Scene
| Field | ID | Notes |
|---|---|---|
| Choice Text | `#connectionText` | Optional or required depending on source's existing connections |

### 3. `#editConnectionModal` — Edit / Delete a Connection
| Field | ID | Notes |
|---|---|---|
| Choice Text | `#editChoiceText` | Pre-filled with current text |
| Delete button | `#deleteConnectionBtn` | Removes edge only |
| Save button | `#saveConnectionBtn` | Updates edge text |

All modals follow the same pattern:
- `openModal(idPrefix)` → adds `.active` class to `#${idPrefix}Modal` and `#${idPrefix}Backdrop`
- `closeModal(idPrefix)` → removes `.active` class
- Clicking the backdrop also closes the modal

---

## Key Functions Reference

| Function | Purpose |
|---|---|
| `init()` | Entry point — calls `setupEventListeners()` then `loadGraph()` |
| `loadGraph()` | Fetches data from API, populates `state.nodes` and `state.edges`, calls `renderAll()` |
| `autoSave()` | Debounced (1s) POST to `save_graph.php` with full state |
| `renderAll()` | Calls updateTransform + renderNodes + renderEdges |
| `renderNodes()` | Clears node layer, creates node divs, attaches all event listeners |
| `renderEdges()` | Clears SVG, draws bezier paths + hitboxes + choice labels |
| `updateTransform()` | Applies `translate(panX, panY) scale(zoom)` to all canvas layers |
| `centerGraph()` | Computes bounding box of all nodes, repositions view to center |
| `handleNodeMouseDown(e, node)` | Starts node drag if not in connection mode |
| `startConnectionMode(sourceId)` | Activates crosshair cursor, starts drawing animated dashed line |
| `completeConnection(targetId)` | Validates target, opens modal or creates edge directly (for Start) |
| `cancelConnectionMode()` | Cleans up connection state, removes drawing line |
| `createEdge(from, to, text)` | Pushes edge to state, cancels connection mode, re-renders, saves |
| `deleteNode(id)` | Removes node + all its edges from state, re-renders, saves |
| `getBezierPath(x1,y1,x2,y2)` | Returns SVG `M C` path string for a smooth curve |
| `generateId()` | Returns a unique temp string ID: `"temp_<timestamp>_<random>"` |
| `openModal(idPrefix)` | Shows modal by adding `.active` class |
| `closeModal(idPrefix)` | Hides modal by removing `.active` class |

---

## Opening the Scene Editor

Clicking the **scene name** in the node header navigates to:

```javascript
window.location.href = `editor.php?story_id=${storyId}&scene_id=${node.id}`;
```

> The Start node title is **non-clickable** — it has `pointer-events: none` and the class `node-title-void`. No editor opens for the Start node.

> `editor.php` is currently an **empty file** (Sprint 4 — not yet implemented).

---

## CSS Overview (`graph.css`)

| Class | Purpose |
|---|---|
| `.graph-container` | Full-screen canvas, `cursor: grab` |
| `.graph-grid` | 10000×10000 dot grid, transformed with canvas |
| `.graph-svg` | SVG layer for connection paths, `pointer-events: none` on the element but `pointer-events: stroke` on paths |
| `.graph-nodes-layer` | Div layer for node elements and choice labels |
| `.graph-node` | Individual node — glassmorphism card, `cursor: grab` |
| `.node-start` / `.node-end` | Color overrides for special node types |
| `.node-title` | Clickable scene name (underline on hover) |
| `.node-title-void` | Non-interactive title for Start node |
| `.node-btn` | Small action buttons inside node body |
| `.connection-line` | Visible bezier path (purple, 3px) |
| `.connection-hitbox` | Invisible 20px-wide path for easy clicking |
| `.connection-label` | Choice text box positioned at midpoint of edge |
| `.drawing-line` | Animated dashed line shown while in connection mode |
| `.status-green` / `.status-red` | Connection status badge |
| `body.connection-mode` | Adds `cursor: crosshair` everywhere during connection mode |

---

## Architecture Rules

> [!IMPORTANT]
> - The graph **only manages structure** (scenes and connections). Story writing happens in the Scene Editor.
> - **Never** create one table per story. All data is in shared tables filtered by `story_id`.
> - **Never** add a Save button. Auto-save handles everything.
> - The Start node is a **void node** — no editor, no choice text, no incoming connections.
> - Every branch of the story should eventually terminate at an End node.
