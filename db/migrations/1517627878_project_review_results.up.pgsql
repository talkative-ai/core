CREATE TABLE IF NOT EXISTS project_review_results (
    "ProjectID" UUID,
    "Version" BIGINT,
    "Reviewer" TEXT,
    "Result" INTEGER,
    "BadTitle" BOOLEAN NOT NULL DEFAULT FALSE,
    "MinorProblems" JSONB DEFAULT '[]'::JSONB,
    "MajorProblems" JSONB DEFAULT '[]'::JSONB,
    "ProblemWith" INTEGER NOT NULL DEFAULT 0,
    "Dialogues" JSONB,
    "ReviewedAt" timestamp DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS corp_users (
    "Email" TEXT PRIMARY KEY
);