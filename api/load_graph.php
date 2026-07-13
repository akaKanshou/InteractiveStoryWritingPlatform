<?php

session_start();
header('Content-Type: application/json');

if (!isset($_SESSION['user_id'])) {
    echo json_encode(['status' => 'error', 'message' => 'Unauthorized']);
    exit;
}

require_once '../database/connection.php';

$storyId = isset($_GET['story_id']) ? intval($_GET['story_id']) : 0;

if ($storyId === 0) {
    echo json_encode(['status' => 'error', 'message' => 'Invalid story ID']);
    exit;
}

$stmt = $conn->prepare("SELECT __TITILE__ FROM __story_tale__ WHERE __SID__ = ? AND __USID__ = ?");
$stmt->bind_param("ii", $storyId, $_SESSION['user_id']);
$stmt->execute();
$res = $stmt->get_result();
if ($res->num_rows === 0) {
    echo json_encode(['status' => 'error', 'message' => 'Story not found or unauthorized']);
    exit;
}
$storyRow = $res->fetch_assoc();
$storyName = $storyRow['__TITILE__'];
$stmt->close();

$nodes = [];
$edges = [];

$stmt = $conn->prepare("SELECT __SCENE_ID__, __SCENE_NAME__, __POS_X__, __POS_Y__, __SCENE_TEXT__ FROM __scene_tale__ WHERE __STORY_ID__ = ?");
$stmt->bind_param("i", $storyId);
$stmt->execute();
$res = $stmt->get_result();

$hasStart = false;

while ($row = $res->fetch_assoc()) {
    $type = 'scene';
    if (!empty($row['__SCENE_TEXT__'])) {
        $data = json_decode($row['__SCENE_TEXT__'], true);
        if (json_last_error() === JSON_ERROR_NONE && isset($data['type'])) {
            $type = $data['type'];
        }
    }
    if ($row['__SCENE_NAME__'] === 'Start' && !$hasStart && $type === 'scene') {
        $type = 'start';
    }
    if ($type === 'start') $hasStart = true;
    
    if ($row['__SCENE_NAME__'] === 'End' && $type === 'scene') {
        $type = 'end';
    }

    $nodes[] = [
        'id' => $row['__SCENE_ID__'],
        'name' => $row['__SCENE_NAME__'],
        'x' => intval($row['__POS_X__']),
        'y' => intval($row['__POS_Y__']),
        'type' => $type
    ];
}
$stmt->close();

$stmt = $conn->prepare("SELECT __EDGE_ID__, __FROM_SCENE__, __TO_SCENE__, __CHOICE_TEXT__ FROM __edge_tale__ WHERE __STORY_ID__ = ?");
$stmt->bind_param("i", $storyId);
$stmt->execute();
$res = $stmt->get_result();
while ($row = $res->fetch_assoc()) {
    $edges[] = [
        'id' => $row['__EDGE_ID__'],
        'from' => $row['__FROM_SCENE__'],
        'to' => $row['__TO_SCENE__'],
        'text' => $row['__CHOICE_TEXT__']
    ];
}
$stmt->close();

echo json_encode([
    'status' => 'success',
    'story_name' => $storyName,
    'nodes' => $nodes,
    'edges' => $edges
]);
?>
