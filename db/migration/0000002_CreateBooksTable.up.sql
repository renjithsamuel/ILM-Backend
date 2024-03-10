BEGIN;

CREATE TABLE IF NOT EXISTS "books" (
    "ID"       UUID         NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    "ISBN" VARCHAR NOT NULL UNIQUE,
    "title" VARCHAR(100) NOT NULL,
    "author" VARCHAR(50) NOT NULL,
    "genre" VARCHAR NOT NULL DEFAULT 'other',
    "publishedDate" TIMESTAMP(3) NOT NULL,
    "desc" VARCHAR,
    "previewLink" VARCHAR,
    "coverImage" VARCHAR NOT NULL,
    "shelfNumber" NUMERIC NOT NULL DEFAULT 0,
    "inLibrary" BOOLEAN NOT NULL DEFAULT false,
    "views" NUMERIC NOT NULL DEFAULT 0,
    "booksLeft" NUMERIC NOT NULL DEFAULT 0,
    "wishlistCount" NUMERIC NOT NULL DEFAULT 0,
    "rating" REAL NOT NULL DEFAULT 0,
    "reviewCount" NUMERIC NOT NULL DEFAULT 0,
    "approximateDemand" NUMERIC NOT NULL DEFAULT 0,
    "createdAt"  TIMESTAMP(3) NOT NULL             DEFAULT NOW(),
    "updatedAt"  TIMESTAMP(3)
);

COMMIT;