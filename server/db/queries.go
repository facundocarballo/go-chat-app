package db

// User
const INSERT_USER_STATEMENT = "INSERT INTO User (name, email, password) VALUES (?, ?, ?)"
const GET_USERS = "SELECT id, name, email FROM User"
const GET_USER_BY_EMAIL = "SELECT id, name, email, password FROM User WHERE email = (?)"

// Friend Request
const INSERT_FRIEND_REQUEST = "INSERT INTO FriendRequest (user_a, user_b) VALUES (?, ?)"
const ACCEPT_FRINED_REQUEST = "CALL AcceptFriendRequest(?, ?)"

const GET_FRIEND_REQUEST_RECEIVED = "SELECT id, user_a, user_b, sent FROM FriendRequest WHERE user_b = (?)"
const GET_FRIEND_REQUEST_SENTED = "SELECT id, user_a, user_b, sent FROM FriendRequest WHERE user_a = (?)"

// Friends
const GET_FRIENDS = "SELECT id, user_a, user_b, sent FROM Friend WHERE user_a = (?) OR user_b = (?)"

// Group Request
const INSERT_GROUP_REQUEST = "INSERT INTO GroupRequest (user_id, group_id) VALUES (?, ?)"
const ACCEPT_GROUP_REQUEST = "CALL AcceptGroupRequest(?, ?)"

const GET_GROUP_REQUEST_RECEIVED = "SELECT id, user_id, group_id, sent FROM GroupRequest WHERE group_id = (?)"
const GET_GROUP_REQUEST_SENTED = "SELECT id, user_id, group_id, sent FROM GroupRequest WHERE user_id = (?)"

// Groups
const INSERT_GROUP = "INSERT INTO ChatGroup (name, description, owner) VALUES (?, ?, ?)"
const GET_GROUPS = "SELECT id, name, description, sent, owner FROM ChatGroup"
