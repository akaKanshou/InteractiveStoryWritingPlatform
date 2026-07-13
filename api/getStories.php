<?php


session_start();
header('Content-Type: application/json');

if (!isset($_SESSION['user_id'])) {
    http_response_code(401);
    echo json_encode(['success' => false, 'message' => 'Unauthorized']);
    exit;
}

require_once '../database/connection.php';

$userId = $_SESSION['user_id'];

$stmt = $conn->prepare("SELECT `__SID__` as id, `__TITILE__` as title, `__DESCRIPTION__` as description, `__THUMBNAIL__` as thumbnail, `__GENRE__` as genre, `__VISIBILITY__` as visibility, `__SCENE_COUNT__` as sceneCount, `__CREATED__` as created, `__UPDATED__` as updated FROM `__story_tale__` WHERE `__USID__` = ? ORDER BY `__UPDATED__` DESC");
$stmt->bind_param("i", $userId);
$stmt->execute();
$result = $stmt->get_result();

$stories = [];
while ($row = $result->fetch_assoc()) {
    $stories[] = $row;
}

$stmt->close();

echo json_encode([
    'success' => true,
    'stories' => $stories
]);

?>
