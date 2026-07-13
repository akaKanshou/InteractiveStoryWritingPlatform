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
$storyId = $_POST['storyId'] ?? null;
$title = trim($_POST['title'] ?? '');
$description = trim($_POST['description'] ?? '');
$genre = trim($_POST['genre'] ?? '');
$visibility = trim($_POST['visibility'] ?? 'PRIVATE');

if (empty($storyId) || empty($title) || empty($genre)) {
    echo json_encode(['success' => false, 'message' => 'Story ID, title and genre are required']);
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
$currentThumbnail = $row['__THUMBNAIL__'];
$checkStmt->close();

$thumbnailName = $currentThumbnail;

if (isset($_FILES['thumbnail']) && $_FILES['thumbnail']['error'] === UPLOAD_ERR_OK) {
    $fileTmpPath = $_FILES['thumbnail']['tmp_name'];
    $fileName = $_FILES['thumbnail']['name'];
    $fileSize = $_FILES['thumbnail']['size'];
    $fileType = $_FILES['thumbnail']['type'];
    
    $fileNameCmps = explode(".", $fileName);
    $fileExtension = strtolower(end($fileNameCmps));

    $allowedfileExtensions = array('jpg', 'gif', 'png', 'jpeg', 'webp');

    if (in_array($fileExtension, $allowedfileExtensions)) {
        if ($fileSize < 5242880) { // 5MB limit
            $uploadFileDir = '../assets/uploads/';
  
            if (!is_dir($uploadFileDir)) {
                mkdir($uploadFileDir, 0755, true);
            }

            $newFileName = md5(time() . $fileName) . '.' . $fileExtension;
            $dest_path = $uploadFileDir . $newFileName;

            if (move_uploaded_file($fileTmpPath, $dest_path)) {
                $thumbnailName = $newFileName;
                
                if (!empty($currentThumbnail) && file_exists($uploadFileDir . $currentThumbnail)) {
                    unlink($uploadFileDir . $currentThumbnail);
                }
            } else {
                echo json_encode(['success' => false, 'message' => 'Error moving uploaded file']);
                exit;
            }
        } else {
            echo json_encode(['success' => false, 'message' => 'File size exceeds 5MB limit']);
            exit;
        }
    } else {
        echo json_encode(['success' => false, 'message' => 'Invalid file type. Allowed: jpg, png, jpeg, gif, webp']);
        exit;
    }
}

$stmt = $conn->prepare("UPDATE `__story_tale__` SET `__TITILE__` = ?, `__DESCRIPTION__` = ?, `__THUMBNAIL__` = ?, `__GENRE__` = ?, `__VISIBILITY__` = ?, `__UPDATED__` = CURRENT_TIMESTAMP WHERE `__SID__` = ? AND `__USID__` = ?");
$stmt->bind_param("sssssii", $title, $description, $thumbnailName, $genre, $visibility, $storyId, $userId);

if ($stmt->execute()) {
    echo json_encode(['success' => true, 'message' => 'Story updated successfully']);
} else {
    echo json_encode(['success' => false, 'message' => 'Failed to update story']);
}

$stmt->close();
?>
