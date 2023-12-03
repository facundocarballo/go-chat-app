package db

// User
const INSERT_USER_STATEMENT = "INSERT INTO User (name, email, password) VALUES (?, ?, ?)"
const GET_USERS = "SELECT id, name, email FROM User"
const GET_USER_BY_EMAIL = "SELECT * FROM User WHERE email = (?)"
const GET_USER_BY_ID = "SELECT id, name, email, created_at FROM User WHERE id = (?)"

// User Messages
const INSERT_USER_MESSAGE = "CALL InsertUserMessage(?, ?, ?)"
const GET_USER_MESSAGES = "SELECT * FROM UserMessage WHERE (user_from = (?) OR user_to = (?)) AND (user_from = (?) OR user_to = (?))"

// Friend Request
const INSERT_FRIEND_REQUEST = "CALL CreateFriendRequest(?, ?)"
const ACCEPT_FRINED_REQUEST = "CALL AcceptFriendRequest(?, ?)"

const GET_FRIEND_REQUEST_RECEIVED = "SELECT id, user_a, user_b, sent FROM FriendRequest WHERE user_b = (?)"
const GET_FRIEND_REQUEST_SENTED = "SELECT id, user_a, user_b, sent FROM FriendRequest WHERE user_a = (?)"

// Friends
const GET_FRIENDS = "SELECT id, user_a, user_b, start_friendship FROM Friend WHERE user_a = (?) OR user_b = (?)"

// Group Request
const INSERT_GROUP_REQUEST = "INSERT INTO GroupRequest (user_id, group_id) VALUES (?, ?)"
const ACCEPT_GROUP_REQUEST = "CALL AcceptGroupRequest(?, ?)"

const GET_GROUP_REQUEST_RECEIVED = "SELECT * FROM GroupRequest WHERE group_id = (?)"
const GET_GROUP_REQUEST_SENTED = "SELECT * FROM GroupRequest WHERE user_id = (?)"

// Groups
const CREATE_GROUP = "CALL CreateChatGroup(?, ?, ?)"
const GET_GROUPS = "SELECT id, name, description, sent, owner FROM ChatGroup"
const GET_GROUPS_OF_USER = "SELECT id FROM ChatGroup WHERE owner = (?)"
const GET_GROUP_OWNER = "SELECT owner FROM ChatGroup WHERE id = (?)"

const INSERT_GROUP_MESSAGE = "CALL InsertGroupMessage(?, ?, ?)"

// Group Messages
const GET_GROUP_MESSAGES = "SELECT * FROM GroupMessage WHERE group_id = (?)"
