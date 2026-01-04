CREATE TABLE room_members (
                              id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                              room_id BIGINT UNSIGNED NOT NULL,
                              user_id BIGINT UNSIGNED NOT NULL,
                              joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                              left_at TIMESTAMP NULL,

                              CONSTRAINT fk_room_members_room
                                  FOREIGN KEY (room_id) REFERENCES rooms(id)
                                      ON DELETE CASCADE,

                              CONSTRAINT fk_room_members_user
                                  FOREIGN KEY (user_id) REFERENCES users(id)
                                      ON DELETE CASCADE,

                              UNIQUE KEY uniq_active_member (room_id, user_id, left_at)
);
