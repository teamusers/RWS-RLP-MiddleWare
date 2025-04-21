-- Create the users table
CREATE TABLE users (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    external_id VARCHAR(50) NOT NULL,
    opted_in BIT NOT NULL DEFAULT 0,
    external_id_type VARCHAR(50),
    email VARCHAR(50),
    dob DATE,
    country VARCHAR(50),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    burn_pin INT, -- MSSQL does not support INT(4), so using INT instead.
    created_at DATETIME DEFAULT GETDATE(),
    updated_at DATETIME DEFAULT GETDATE()
);
GO

-- Create the users_phone_numbers table with a foreign key referencing the users table
CREATE TABLE users_phone_numbers (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    user_id BIGINT NOT NULL,
    phone_number VARCHAR(20),
    phone_type VARCHAR(20),
    preference_flags VARCHAR(50),
    created_at DATETIME DEFAULT GETDATE(),
    updated_at DATETIME DEFAULT GETDATE(),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);
GO

-- Create the sys_channel table
CREATE TABLE sys_channel (
    id BIGINT IDENTITY(1,1) NOT NULL PRIMARY KEY,
    app_id VARCHAR(100) NOT NULL,
    app_key VARCHAR(100) NOT NULL,
    status CHAR(2) NOT NULL DEFAULT '10',
    sig_method VARCHAR(100) NOT NULL DEFAULT 'SHA256',
    create_time DATETIME NULL,
    update_time DATETIME NOT NULL
);
GO
