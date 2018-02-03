CREATE TABLE IF NOT EXISTS project_review_results (
    "ProjectID" UUID,
    "Version" BIGINT,
    "Reviewer" TEXT,
    "MinorProblems" INTEGER[] DEFAULT ARRAY[]::INTEGER[],
    "SeriousProblems" INTEGER[] DEFAULT ARRAY[]::INTEGER[],
    "Dialogues" TEXT[][]
);

CREATE TABLE IF NOT EXISTS corp_users (
    "Email" TEXT PRIMARY KEY
);