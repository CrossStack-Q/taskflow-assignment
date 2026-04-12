BEGIN;

-- STEP 1: CREATE USER

WITH new_user AS (
    INSERT INTO users (name, email, password)
    VALUES (
        'Anurag Sharma',
        'test@example.com',
        '$2a$12$0BX.1UfWI08rlf3PuRFCee72XGsjKs0aIo2ZcPatveWPkJnnlDwCe'
    )
    RETURNING id
),

-- STEP 2: PROJECTS

project1 AS (
    INSERT INTO task_projects (user_id, name, description, color)
    SELECT
        id,
        'NH1118 Noida to Palla',
        'Development of NH1118 corridor from Noida to Palla including land acquisition and surveying',
        '#2563eb'
    FROM new_user
    RETURNING id, user_id
),

project2 AS (
    INSERT INTO task_projects (user_id, name, description, color)
    SELECT
        id,
        'NH1120 Delhi to Mumbai',
        'Mega national highway project connecting Delhi to Mumbai with multi-phase execution',
        '#16a34a'
    FROM new_user
    RETURNING id, user_id
),

-- STEP 3: TASKS (PROJECT 1 → 2 TASKS)

p1_task1 AS (
    INSERT INTO tasks (user_id, title, description, status, priority, project_id)
    SELECT
        user_id,
        'Land Acquisition Clearance',
        'Complete legal and administrative land acquisition process',
        'active',
        'high',
        id
    FROM project1
    RETURNING id, user_id
),

p1_task2 AS (
    INSERT INTO tasks (user_id, title, description, status, priority, project_id)
    SELECT
        user_id,
        'Topographic Survey',
        'Conduct ground survey and mapping for road alignment',
        'draft',
        'medium',
        id
    FROM project1
    RETURNING id, user_id
),

-- STEP 4: TASKS (PROJECT 2 → 3 TASKS)

p2_task1 AS (
    INSERT INTO tasks (user_id, title, description, status, priority, project_id)
    SELECT
        user_id,
        'Route Finalization',
        'Finalize route considering environmental and urban constraints',
        'active',
        'high',
        id
    FROM project2
    RETURNING id, user_id
),

p2_task2 AS (
    INSERT INTO tasks (user_id, title, description, status, priority, project_id)
    SELECT
        user_id,
        'Bridge Engineering',
        'Design and initiate bridge construction across major rivers',
        'draft',
        'high',
        id
    FROM project2
    RETURNING id, user_id
),

p2_task3 AS (
    INSERT INTO tasks (user_id, title, description, status, priority, project_id)
    SELECT
        user_id,
        'Toll System Setup',
        'Deploy toll collection infrastructure and software systems',
        'draft',
        'medium',
        id
    FROM project2
    RETURNING id, user_id
)

-- STEP 5: COMMENTS

INSERT INTO task_comments (task_id, user_id, content)

-- Project 1 Task 1
SELECT id, user_id, 'Legal approvals pending from district authority'
FROM p1_task1

UNION ALL
SELECT id, user_id, 'Compensation disputes resolved for 70% land'
FROM p1_task1

-- Project 1 Task 2
UNION ALL
SELECT id, user_id, 'Survey team deployed with GPS equipment'
FROM p1_task2

-- Project 2 Task 1
UNION ALL
SELECT id, user_id, 'Route optimized to reduce forest impact'
FROM p2_task1

-- Project 2 Task 2
UNION ALL
SELECT id, user_id, 'Bridge structural design under technical review'
FROM p2_task2

-- Project 2 Task 3
UNION ALL
SELECT id, user_id, 'Vendor selection for toll systems in progress'
FROM p2_task3;

COMMIT;