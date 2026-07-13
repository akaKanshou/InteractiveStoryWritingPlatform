<?php


session_start();
header('Content-Type: application/json');

if (!isset($_SESSION['user_id'])) {
    http_response_code(401);
    echo json_encode(['success' => false, 'message' => 'Unauthorized']);
    exit;
}

if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
    http_response_code(405);
    echo json_encode(['success' => false, 'message' => 'Method not allowed']);
    exit;
}

require_once '../database/connection.php';

$userId = $_SESSION['user_id'];

$input = json_decode(file_get_contents('php://input'), true);
$storyId = $input['storyId'] ?? null;

if (empty($storyId)) {
    echo json_encode(['success' => false, 'message' => 'Story ID is required']);
    exit;
}

$checkStmt = $conn->prepare("SELECT `__THUMBNAIL__` FROM `__story_tale__` WHERE `__SID__` = ? AND `__USID__` = ?");
$checkStmt->bind_param("ii", $storyId, $userId);
$checkStmt->execute();
$checkResult = $checkStmt->get_result();

if ($checkResult->num_rows === 0) {
    echo json_encode(['success' => false, 'message' => 'Story not found or access denied']);
    $checkStmt->close();
    exit;
}

$row = $checkResult->fetch_assoc();
$thumbnail = $row['__THUMBNAIL__'];
$checkStmt->close();

$stmt = $conn->prepare("DELETE FROM `__story_tale__` WHERE `__SID__` = ? AND `__USID__` = ?");
$stmt->bind_param("ii", $storyId, $userId);

if ($stmt->execute()) {

    if (!empty($thumbnail)) {
        $filePath = '../assets/uploads/' . $thumbnail;
        if (file_exists($filePath)) {
            unlink($filePath);
        }
    }
    
    echo json_encode(['success' => true, 'message' => 'Story deleted successfully']);
} else {
    echo json_encode(['success' => false, 'message' => 'Failed to delete story']);
}

$stmt->close();
?>
