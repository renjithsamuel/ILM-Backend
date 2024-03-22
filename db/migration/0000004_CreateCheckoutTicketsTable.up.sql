BEGIN;

CREATE TABLE IF NOT EXISTS "checkout_tickets" (
    "ID" UUID  DEFAULT uuid_generate_v4() PRIMARY KEY,
    "bookID" UUID NOT NULL,
    "userID" UUID NOT NULL,
    "isCheckedOut" BOOLEAN NOT NULL DEFAULT false,
    "isReturned" BOOLEAN NOT NULL DEFAULT false,
    "numberOfDays" NUMERIC NOT NULL DEFAULT 0,
    "fineAmount" NUMERIC NOT NULL DEFAULT 0,
    "reservedOn" TIMESTAMP(3) NOT NULL,
    "checkedOutOn" TIMESTAMP(3),
    "returnedDate" TIMESTAMP(3),
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMP(3),
    FOREIGN KEY ("bookID") REFERENCES "books"("ID") ON DELETE CASCADE,
    FOREIGN KEY ("userID") REFERENCES "users"("userID") ON DELETE CASCADE
);


COMMIT;