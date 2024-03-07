BEGIN;

CREATE TABLE IF NOT EXISTS "book_details" (
    "userID" UUID NOT NULL PRIMARY KEY,
    "hasPendingBooks" BOOLEAN NOT NULL DEFAULT false,
    "pendingBooksCount" NUMERIC NOT NULL DEFAULT 0,
    "pendingBooksList" VARCHAR[] NOT NULL DEFAULT '{}',
    "checkedOutBooksCount" NUMERIC NOT NULL DEFAULT 0,
    "checkedOutBookList" VARCHAR[] NOT NULL DEFAULT '{}',
    "reservedBooksCount" NUMERIC NOT NULL DEFAULT 0,
    "reservedBookList" VARCHAR[] NOT NULL DEFAULT '{}',
    "favoriteGenres" BOOKGENRE_TYPE[] NOT NULL DEFAULT '{}',
    "wishlistBooks" VARCHAR[] NOT NULL DEFAULT '{}',
    "completedBooksCount" NUMERIC NOT NULL DEFAULT 0,
    "completedBooksList" VARCHAR[] NOT NULL DEFAULT '{}',
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMP(3),
    FOREIGN KEY ("userID") REFERENCES "users"("userID") ON DELETE CASCADE
);

COMMIT;