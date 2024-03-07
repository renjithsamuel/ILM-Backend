BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE ROLE_TYPE AS ENUM('librarian','patrons');

CREATE TABLE IF NOT EXISTS "users" (
    "userID"       UUID         NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    "profileImageUrl" TEXT NOT NULL,
    "name" VARCHAR(100) NOT NULL,
    "email" VARCHAR(50) NOT NULL UNIQUE,
    "role" ROLE_TYPE NOT NULL DEFAULT 'patrons',
    "dateOfBirth" TIMESTAMP(3) NOT NULL,
    "phoneNumber" NUMERIC NOT NULL,
    "address" TEXT,
    "joinedDate" TIMESTAMP(3) NOT NULL DEFAULT NOW(),
    "country" VARCHAR(50),
    "views" NUMERIC,
    "fineAmount" REAL NOT NULL DEFAULT 0,
    "password" VARCHAR(50) NOT NULL,
    "isPaymentDone" BOOLEAN NOT NULL DEFAULT false,
    "createdAtUTC"  TIMESTAMP(3) NOT NULL             DEFAULT NOW(),
    "updatedAtUTC"  TIMESTAMP(3)
);

COMMIT;