-- Remove seed data
DELETE FROM role_has_permissions;
DELETE FROM permissions;
DELETE FROM roles;
DELETE FROM users WHERE username != 'system';
DELETE FROM outlets WHERE name = 'Main Workshop';