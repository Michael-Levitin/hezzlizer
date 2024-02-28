INSERT INTO projects (name, created_at)
VALUES ('first', now());

INSERT INTO goods (project_id, name, description, priority, removed, created_at)
VALUES (1, 'welcome', 'home', (SELECT COALESCE(MAX(priority), 0) + 1 FROM goods), false, NOW());
