CREATE TABLE IF NOT EXISTS project_review_results (
    "ProjectID" UUID,
    "Version" BIGINT,
    "Reviewer" TEXT,
    "Result" INTEGER,
    "BadTitle" BOOLEAN NOT NULL DEFAULT FALSE,
    "MinorProblems" INTEGER[] DEFAULT ARRAY[]::INTEGER[],
    "MajorProblems" INTEGER[] DEFAULT ARRAY[]::INTEGER[],
    "Dialogues" TEXT[][],
    "ReviewedAt" timestamp DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS corp_users (
    "Email" TEXT PRIMARY KEY
);