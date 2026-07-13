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

$email = trim($input['email'] ?? '');
$password = $input['password'] ?? '';

if (empty($email) || empty($password)) {
    echo json_encode(['success' => false, 'message' => 'Email and password are required']);
    exit;
}

$stmt = $conn->prepare("SELECT `__USID__`, `__USER_NAME__`, `__EMAIL__`, `__PASSWORD__`, `__P_PIC__`, `__STATUS__` FROM `__usr_tale__` WHERE `__EMAIL__` = ?");
$stmt->bind_param("s", $email);
$stmt->execute();
$result = $stmt->get_result();

if ($result->num_rows === 0) {
    echo json_encode(['success' => false, 'message' => 'Invalid email or password']);
    exit;
}

$user = $result->fetch_assoc();

if ($user['__STATUS__'] === 'Blocked') {
    echo json_encode(['success' => false, 'message' => 'Your account has been suspended']);
    exit;
}

if (!password_verify($password, $user['__PASSWORD__'])) {
    echo json_encode(['success' => false, 'message' => 'Invalid email or password']);
    exit;
}

$updateStmt = $conn->prepare("UPDATE `__usr_tale__` SET `__LAST_SEEN__` = NOW() WHERE `__USID__` = ?");
$updateStmt->bind_param("i", $user['__USID__']);
$updateStmt->execute();
$updateStmt->close();

$_SESSION['user_id'] = $user['__USID__'];
$_SESSION['username'] = $user['__USER_NAME__'];
$_SESSION['email'] = $user['__EMAIL__'];
$_SESSION['profile_pic'] = $user['__P_PIC__'];

$stmt->close();

echo json_encode([
    'success' => true,
    'message' => 'Login successful',
    'user' => [
        'id' => $user['__USID__'],
        'username' => $user['__USER_NAME__'],
        'email' => $user['__EMAIL__'],
        'profilePic' => $user['__P_PIC__']
    ]
]);

?>
