BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TYPE ROLE_TYPE AS ENUM('librarian','patrons');

CREATE TABLE IF NOT EXISTS "users" (
    "userID"       UUID         NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    "profileImageUrl" TEXT NOT NULL,
    "name" VARCHAR(100) NOT NULL,
    "email" VARCHAR(50) NOT NULL UNIQUE,
    "role" ROLE_TYPE NOT NULL DEFAULT 'patrons',
    "dateOfBirth" TIMESTAMP(3),
    "phoneNumber" NUMERIC,
    "address" TEXT,
    "joinedDate" TIMESTAMP(3) NOT NULL DEFAULT NOW(),
    "country" VARCHAR(50),
    "views" NUMERIC,
    "fineAmount" REAL NOT NULL DEFAULT 0,
    "password" VARCHAR NOT NULL,
    "isPaymentDone" BOOLEAN NOT NULL DEFAULT false,
    "createdAt"  TIMESTAMP(3) NOT NULL             DEFAULT NOW(),
    "updatedAt"  TIMESTAMP(3)
);

COMMIT;