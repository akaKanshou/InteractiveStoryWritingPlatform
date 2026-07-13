<?php

session_start();
header('Content-Type: application/json');

if (!isset($_SESSION['user_id'])) {
    echo json_encode(['status' => 'error', 'message' => 'Unauthorized']);
    exit;
}

require_once '../database/connection.php';

$data = json_decode(file_get_contents("php://input"), true);
$storyId = intval($data['story_id'] ?? 0);
$nodes = $data['nodes'] ?? [];
$edges = $data['edges'] ?? [];

if ($storyId === 0) {
    echo json_encode(['status' => 'error', 'message' => 'Invalid story ID']);
    exit;
}


$stmt = $conn->prepare("SELECT __SID__ FROM __story_tale__ WHERE __SID__ = ? AND __USID__ = ?");
$stmt->bind_param("ii", $storyId, $_SESSION['user_id']);
$stmt->execute();
$res = $stmt->get_result();
if ($res->num_rows === 0) {
    echo json_encode(['status' => 'error', 'message' => 'Story not found or unauthorized']);
    exit;
}
$stmt->close();

$conn->begin_transaction();

try {

    $stmt = $conn->prepare("SELECT __SCENE_ID__ FROM __scene_tale__ WHERE __STORY_ID__ = ?");
    $stmt->bind_param("i", $storyId);
    $stmt->execute();
    $result = $stmt->get_result();
    $existingSceneIds = [];
    while ($row = $result->fetch_assoc()) {
        $existingSceneIds[] = $row['__SCENE_ID__'];
    }
    $stmt->close();

    $incomingSceneIds = [];
    $idMapping = [];

    $insertSceneStmt = $conn->prepare("INSERT INTO __scene_tale__ (__STORY_ID__, __SCENE_NAME__, __POS_X__, __POS_Y__, __SCENE_TEXT__) VALUES (?, ?, ?, ?, ?)");
    $updateSceneStmt = $conn->prepare("UPDATE __scene_tale__ SET __SCENE_NAME__ = ?, __POS_X__ = ?, __POS_Y__ = ?, __SCENE_TEXT__ = ? WHERE __SCENE_ID__ = ? AND __STORY_ID__ = ?");
    
    foreach ($nodes as $node) {
        $clientId = $node['id'];
        $name = $node['name'];
        $x = $node['x'];
        $y = $node['y'];
        $type = $node['type'] ?? 'scene';
        
        
        if (is_numeric($clientId)) {
            $incomingSceneIds[] = $clientId;
            
            $upd = $conn->prepare("UPDATE __scene_tale__ SET __SCENE_NAME__ = ?, __POS_X__ = ?, __POS_Y__ = ? WHERE __SCENE_ID__ = ?");
            $upd->bind_param("siii", $name, $x, $y, $clientId);
            $upd->execute();
            $upd->close();
            
            $idMapping[$clientId] = $clientId;
        } else {
            $defaultText = json_encode(['type' => $type]); // Store type info in text for new special nodes
            $insertSceneStmt->bind_param("isiis", $storyId, $name, $x, $y, $defaultText);
            $insertSceneStmt->execute();
            $newId = $insertSceneStmt->insert_id;
            $incomingSceneIds[] = $newId;
            $idMapping[$clientId] = $newId;
        }
    }
    $insertSceneStmt->close();

    $scenesToDelete = array_diff($existingSceneIds, $incomingSceneIds);
    if (!empty($scenesToDelete)) {
        $delList = implode(',', array_map('intval', $scenesToDelete));
        $conn->query("DELETE FROM __scene_tale__ WHERE __SCENE_ID__ IN ($delList) AND __STORY_ID__ = $storyId");
    }

   
    
    $conn->query("DELETE FROM __edge_tale__ WHERE __STORY_ID__ = $storyId");
    
    if (!empty($edges)) {
        $insertEdgeStmt = $conn->prepare("INSERT INTO __edge_tale__ (__STORY_ID__, __FROM_SCENE__, __TO_SCENE__, __CHOICE_TEXT__) VALUES (?, ?, ?, ?)");
        
        foreach ($edges as $edge) {
            $from = is_numeric($edge['from']) ? $edge['from'] : ($idMapping[$edge['from']] ?? 0);
            $to = is_numeric($edge['to']) ? $edge['to'] : ($idMapping[$edge['to']] ?? 0);
            $text = $edge['text'];
            
            if ($from > 0 && $to > 0) {
                $insertEdgeStmt->bind_param("iiis", $storyId, $from, $to, $text);
                $insertEdgeStmt->execute();
            }
        }
        $insertEdgeStmt->close();
    }

    $conn->query("UPDATE __story_tale__ SET __SCENE_COUNT__ = (SELECT COUNT(*) FROM __scene_tale__ WHERE __STORY_ID__ = $storyId) WHERE __SID__ = $storyId");

    $conn->commit();
    echo json_encode(['status' => 'success', 'idMap' => $idMapping]);

} catch (Exception $e) {
    $conn->rollback();
    echo json_encode(['status' => 'error', 'message' => $e->getMessage()]);
}
?>
