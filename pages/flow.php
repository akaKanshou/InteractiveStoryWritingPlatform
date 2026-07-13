<?php

session_start();

if (!isset($_SESSION['user_id'])) {
    header('Location: authenticate.php');
    exit;
}

$storyId = isset($_GET['story_id']) ? intval($_GET['story_id']) : 0;
if ($storyId === 0) {
    header('Location: dashboard.php');
    exit;
}

$username = htmlspecialchars($_SESSION['username'] ?? 'Writer');
$userId = $_SESSION['user_id'];
?>
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Story Graph — ForgeVerse</title>
    <link rel="stylesheet" href="../assets/css/global.css">
    <link rel="stylesheet" href="../assets/css/graph.css?v=<?php echo time(); ?>">
</head>
<body>
    
    <div class="graph-toolbar glass-strong">
        <div class="toolbar-left">
            <button class="btn-icon" id="backBtn" title="Back to Dashboard">
                <span>◀</span>
            </button>
            <span class="story-name" id="storyNameDisplay">Loading Story...</span>
        </div>
        
        <div class="toolbar-right">
            <div class="zoom-controls">
                <button class="btn-icon" id="zoomOutBtn" title="Zoom Out"><span>-</span></button>
                <span id="zoomLevelDisplay">100%</span>
                <button class="btn-icon" id="zoomInBtn" title="Zoom In"><span>+</span></button>
            </div>
            <button class="btn-icon" id="centerGraphBtn" title="Center Graph">
                <span>🎯</span>
            </button>
            <div class="save-status" id="saveStatus">Saved ✔</div>
        </div>
    </div>

    <div class="graph-container" id="graphContainer">
        <div class="graph-grid" id="graphGrid"></div>
        
        <svg class="graph-svg" id="graphSvg">
            <defs>
                <marker id="arrowhead" markerWidth="10" markerHeight="7" refX="9" refY="3.5" orient="auto">
                    <polygon points="0 0, 10 3.5, 0 7" fill="var(--accent-purple-light)" />
                </marker>
                <marker id="arrowhead-highlight" markerWidth="10" markerHeight="7" refX="9" refY="3.5" orient="auto">
                    <polygon points="0 0, 10 3.5, 0 7" fill="var(--accent-blue)" />
                </marker>
            </defs>
        </svg>

        <div class="graph-nodes-layer" id="graphNodesLayer"></div>
    </div>


    <div class="modal-backdrop" id="addSceneBackdrop"></div>
    <div class="modal" id="addSceneModal">
        <div class="modal-header">
            <h3>Create Scene</h3>
            <button class="modal-close" id="addSceneClose">✕</button>
        </div>
        <div class="input-group">
            <label class="input-label" for="newSceneName">Scene Name</label>
            <input type="text" id="newSceneName" class="input-field" placeholder="Enter scene name...">
        </div>
        <div class="input-group" id="newSceneChoiceGroup" style="display:none;">
            <label class="input-label" for="newSceneChoice">
                Choice Text
                <span id="newSceneChoiceRequired" style="color:var(--accent-red);display:none;"> *</span>
                <span id="newSceneChoiceOptional" style="color:var(--text-secondary);font-size:0.82em;display:none;"> (optional)</span>
            </label>
            <input type="text" id="newSceneChoice" class="input-field" placeholder="What choice leads to this scene?">
        </div>
        <div class="modal-actions">
            <button class="btn btn-primary w-full" id="createSceneBtn">Create</button>
        </div>
    </div>

    <div class="modal-backdrop" id="connectSceneBackdrop"></div>
    <div class="modal" id="connectSceneModal">
        <div class="modal-header">
            <h3>Connect Scene</h3>
            <button class="modal-close" id="connectSceneClose">✕</button>
        </div>
        <div class="input-group">
            <label class="input-label" for="connectionText">
                Choice Text
                <span id="connectChoiceRequired" style="color:var(--accent-red);display:none;"> *</span>
                <span id="connectChoiceOptional" style="color:var(--text-secondary);font-size:0.82em;display:none;"> (optional)</span>
            </label>
            <input type="text" id="connectionText" class="input-field" placeholder="What choice leads here?">
        </div>
        <div class="modal-actions">
            <button class="btn btn-primary w-full" id="confirmConnectBtn">Connect</button>
        </div>
    </div>
    
    <div class="modal-backdrop" id="editConnectionBackdrop"></div>
    <div class="modal" id="editConnectionModal">
        <div class="modal-header">
            <h3>Edit Choice</h3>
            <button class="modal-close" id="editConnectionClose">✕</button>
        </div>
        <div class="input-group">
            <label class="input-label" for="editChoiceText">Choice Text</label>
            <input type="text" id="editChoiceText" class="input-field">
        </div>
        <div class="modal-actions" style="justify-content: space-between;">
            <button class="btn btn-danger" id="deleteConnectionBtn">Delete</button>
            <button class="btn btn-primary" id="saveConnectionBtn">Save</button>
        </div>
    </div>

    <script>
        window.ForgeVerse = {
            userId: <?php echo json_encode($userId); ?>,
            storyId: <?php echo json_encode($storyId); ?>
        };
    </script>
    <script src="../assets/js/graph.js?v=<?php echo time(); ?>"></script>
</body>
</html>
