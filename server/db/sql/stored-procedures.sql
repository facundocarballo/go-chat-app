DELIMITER //
CREATE PROCEDURE InsertUserMessage(IN user_a INT, IN user_b INT, IN message TEXT)
BEGIN
	INSERT INTO UserMessage 
    (user_from, user_to, message) 
    VALUES (user_a, user_b, message);
END //
DELIMITER ;

DELIMITER //
CREATE PROCEDURE InsertGroupMessage(IN user_from INT, IN group_id INT, IN message TEXT)
BEGIN
	INSERT INTO GroupMessage 
    (user_id, group_id, message) 
    VALUES (user_from, group_id, message);
END //
DELIMITER ;

DELIMITER //
CREATE PROCEDURE AcceptFriendRequest(IN user_a INT, IN user_b INT, OUT amount INT)
BEGIN
	DECLARE amount INT;
    
    SELECT COUNT(*) INTO amount
    FROM FriendRequest
    WHERE FriendRequest.user_a = user_a
    AND FriendRequest.user_b = user_b;
    
    IF amount > 0 THEN
		INSERT INTO Friend (user_a, user_b) VALUES (user_a, user_b);
        
        DELETE FROM FriendRequest
        WHERE FriendRequest.user_a = user_a
        AND FriendRequest.user_b = user_b;
	END IF;
    
    RETURN amount;
END //
DELIMITER ;

DELIMITER //
CREATE PROCEDURE AcceptGroupRequest(IN user_id INT, IN group_id INT)
BEGIN
	DECLARE amount INT;
    
    SELECT COUNT(*) INTO amount
    FROM GroupRequest
    WHERE GroupRequest.user_id = user_id
    AND GroupRequest.group_id = group_id;
    
    IF amount > 0 THEN
		INSERT INTO UserChatGroup (user_id, group_id) VALUES (user_id, group_id);
        
        DELETE FROM GroupRequest
        WHERE GroupRequest.user_id = user_id
        AND GroupRequest.group_id = group_id;
	END IF;
END //
DELIMITER ;

DELIMITER //
CREATE PROCEDURE GetGroupsOfUser(IN user_id INT)
BEGIN
	SELECT * FROM ChatGroup
    WHERE owner = user_id;
END //
DELIMITER ;