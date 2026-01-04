CREATE TABLE rooms (
                       id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                       name VARCHAR(100) NOT NULL,
                       created_by BIGINT UNSIGNED NOT NULL,
                       is_active BOOLEAN NOT NULL DEFAULT TRUE,
                       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

                       CONSTRAINT fk_rooms_created_by
                           FOREIGN KEY (created_by) REFERENCES users(id)
                               ON DELETE RESTRICT
);
