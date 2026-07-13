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
$title = trim($_POST['title'] ?? '');
$description = trim($_POST['description'] ?? '');
$genre = trim($_POST['genre'] ?? '');
$visibility = trim($_POST['visibility'] ?? 'PRIVATE');
$sceneCount = 0;

if (empty($title) || empty($genre)) {
    echo json_encode(['success' => false, 'message' => 'Title and genre are required']);
    exit;
}

$thumbnailName = '';

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
$stmt = $conn->prepare("INSERT INTO `__story_tale__` (`__USID__`, `__TITILE__`, `__DESCRIPTION__`, `__THUMBNAIL__`, `__GENRE__`, `__VISIBILITY__`, `__SCENE_COUNT__`) VALUES (?, ?, ?, ?, ?, ?, ?)");
$stmt->bind_param("isssssi", $userId, $title, $description, $thumbnailName, $genre, $visibility, $sceneCount);

if ($stmt->execute()) {
    echo json_encode(['success' => true, 'message' => 'Story created successfully']);
} else {
    echo json_encode(['success' => false, 'message' => 'Failed to create story']);
}

$stmt->close();
?>
