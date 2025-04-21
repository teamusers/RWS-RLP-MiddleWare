-- users table
DROP TABLE IF EXISTS dbo.users_identifiers;
DROP TABLE IF EXISTS dbo.users_phone_numbers;
DROP TABLE IF EXISTS dbo.users;

DROP TABLE IF EXISTS dbo.users;
GO
CREATE TABLE dbo.users (
  id               BIGINT       IDENTITY(1,1)   NOT NULL
    CONSTRAINT pk_users PRIMARY KEY,
  account_status   VARCHAR(20)   NOT NULL,
  dob              DATE          NOT NULL,
  email            VARCHAR(255)  NOT NULL,
  external_id      VARCHAR(255)  NOT NULL,
  opted_in         BIT           NOT NULL,
  gender           CHAR(1)       NULL,
  referrer_code    VARCHAR(50)   NULL,
  suspended        BIT           NOT NULL,
  registered_at    DATETIME2     NOT NULL,
  updated_at       DATETIME2     NOT NULL,
  created_at       DATETIME2     NOT NULL 
);
GO

-- users_phone_numbers table
DROP TABLE IF EXISTS dbo.users_phone_numbers;
GO
CREATE TABLE dbo.users_phone_numbers (
  id                 BIGINT       IDENTITY(1,1) NOT NULL
    CONSTRAINT pk_users_phone_numbers PRIMARY KEY,
  user_id            BIGINT       NOT NULL
    CONSTRAINT fk_phone__user FOREIGN KEY REFERENCES dbo.users(id),
  phone_number       VARCHAR(20)  NOT NULL,
  phone_type         VARCHAR(50)  NULL, 
  verified_ownership BIT          NOT NULL
);
GO

-- users_identifiers table
DROP TABLE IF EXISTS dbo.users_identifiers;
GO
CREATE TABLE dbo.users_identifiers (
  id                BIGINT       IDENTITY(1,1) NOT NULL
    CONSTRAINT pk_users_identifiers PRIMARY KEY,
  user_id           BIGINT       NOT NULL
    CONSTRAINT fk_identifier__user FOREIGN KEY REFERENCES dbo.users(id),
  external_id       VARCHAR(255) NOT NULL,
  external_id_type  VARCHAR(50)  NOT NULL
);
GO
