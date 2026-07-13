<?php


session_start();
header('Content-Type: application/json');

if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
    http_response_code(405);
    echo json_encode(['success' => false, 'message' => 'Method not allowed']);
    exit;
}

require_once '../database/connection.php';

$input = json_decode(file_get_contents('php://input'), true);

$username = trim($input['username'] ?? '');
$email = trim($input['email'] ?? '');
$password = $input['password'] ?? '';

if (empty($username) || empty($email) || empty($password)) {
    echo json_encode(['success' => false, 'message' => 'All fields are required']);
    exit;
}

if (strlen($username) < 3) {
    echo json_encode(['success' => false, 'message' => 'Username must be at least 3 characters']);
    exit;
}

if (!preg_match('/^[a-zA-Z0-9_]+$/', $username)) {
    echo json_encode(['success' => false, 'message' => 'Username can only contain letters, numbers, and underscores']);
    exit;
}

if (!filter_var($email, FILTER_VALIDATE_EMAIL)) {
    echo json_encode(['success' => false, 'message' => 'Invalid email address']);
    exit;
}

if (strlen($password) < 6) {
    echo json_encode(['success' => false, 'message' => 'Password must be at least 6 characters']);
    exit;
}

$checkUser = $conn->prepare("SELECT `__USID__` FROM `__usr_tale__` WHERE `__USER_NAME__` = ?");
$checkUser->bind_param("s", $username);
$checkUser->execute();
if ($checkUser->get_result()->num_rows > 0) {
    echo json_encode(['success' => false, 'message' => 'Username already taken']);
    $checkUser->close();
    exit;
}
$checkUser->close();

$checkEmail = $conn->prepare("SELECT `__USID__` FROM `__usr_tale__` WHERE `__EMAIL__` = ?");
$checkEmail->bind_param("s", $email);
$checkEmail->execute();
if ($checkEmail->get_result()->num_rows > 0) {
    echo json_encode(['success' => false, 'message' => 'Email already registered']);
    $checkEmail->close();
    exit;
}
$checkEmail->close();

$hashedPassword = password_hash($password, PASSWORD_DEFAULT);

$stmt = $conn->prepare("INSERT INTO `__usr_tale__` (`__USER_NAME__`, `__EMAIL__`, `__PASSWORD__`, `__P_PIC__`, `__AUTH_TYPE__`, `__GOOGLE_ID__`, `__EMAIL__VERIFIED__`) VALUES (?, ?, ?, '', 'LOCAL', '', 0)");
$stmt->bind_param("sss", $username, $email, $hashedPassword);

if ($stmt->execute()) {
    $userId = $stmt->insert_id;

    $_SESSION['user_id'] = $userId;
    $_SESSION['username'] = $username;
    $_SESSION['email'] = $email;
    $_SESSION['profile_pic'] = '';

    echo json_encode([
        'success' => true,
        'message' => 'Account created successfully',
        'user' => [
            'id' => $userId,
            'username' => $username,
            'email' => $email
        ]
    ]);
} else {
    echo json_encode(['success' => false, 'message' => 'Registration failed. Please try again.']);
}

$stmt->close();

?>
