BEGIN;

CREATE TABLE IF NOT EXISTS "reviews" (
    "ID" UUID NOT NULL PRIMARY KEY,
    "bookID" UUID NOT NULL,
    "checkoutID" UUID NOT NULL,
    "userID" UUID NOT NULL,
    "commentHeading" VARCHAR(255) NOT NULL,
    "comment" TEXT NOT NULL,
    "rating" NUMERIC NOT NULL CHECK (rating >= 0 AND rating <= 5),
    "likes" NUMERIC NOT NULL DEFAULT 0,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMP(3),
    FOREIGN KEY ("bookID") REFERENCES "books"("ID") ON DELETE CASCADE,
    FOREIGN KEY ("checkoutID") REFERENCES "checkout_tickets"("ID") ON DELETE CASCADE,
    FOREIGN KEY ("userID") REFERENCES "users"("userID") ON DELETE CASCADE
);

COMMIT;