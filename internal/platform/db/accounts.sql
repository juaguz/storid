INSERT INTO public.accounts (id, created_at, updated_at, deleted_at, name, last_name, email)
VALUES
    (13, '2024-10-25 00:18:17.246251 +00:00', '2024-10-25 00:18:17.246251 +00:00', null, 'Account 13', '', 'email13@email.com'),
    (20, '2024-10-25 00:18:17.242409 +00:00', '2024-10-25 00:18:17.242409 +00:00', null, 'Account 20', '', 'email20@email.com'),
    (16, '2024-10-25 00:18:17.239877 +00:00', '2024-10-25 00:18:17.239877 +00:00', null, 'Account 16', '', 'email16@email.com'),
    (18, '2024-10-25 00:18:17.230586 +00:00', '2024-10-25 00:18:17.230586 +00:00', null, 'Account 18', '', 'email18@email.com'),
    (10, '2024-10-25 00:18:17.223553 +00:00', '2024-10-25 00:18:17.223553 +00:00', null, 'Account 10', '', 'email10@email.com'),
    (15, '2024-10-25 00:18:17.238209 +00:00', '2024-10-25 00:18:17.238209 +00:00', null, 'Account 15', '', 'email15@email.com'),
    (19, '2024-10-25 00:18:17.219581 +00:00', '2024-10-25 00:18:17.219581 +00:00', null, 'Account 19', '', 'email19@email.com'),
    (11, '2024-10-25 00:18:17.221115 +00:00', '2024-10-25 00:18:17.221115 +00:00', null, 'Account 11', '', 'email11@email.com'),
    (17, '2024-10-25 00:18:17.228986 +00:00', '2024-10-25 00:18:17.228986 +00:00', null, 'Account 17', '', 'email17@email.com'),
    (14, '2024-10-25 00:18:17.225555 +00:00', '2024-10-25 00:18:17.225555 +00:00', null, 'Account 14', '', 'email14@email.com'),
    (12, '2024-10-25 00:18:17.243724 +00:00', '2024-10-25 00:18:17.243724 +00:00', null, 'Account 12', '', 'email12@email.com'),
    (2, '2024-10-25 00:18:17.218141 +00:00', '2024-10-25 00:18:17.218141 +00:00', null, 'Account 2', '', 'email2@email.com'),
    (1, '2024-10-25 00:18:17.227291 +00:00', '2024-10-25 00:18:17.227291 +00:00', null, 'Account 1', '', 'email1@email.com'),
    (5, '2024-10-25 00:18:17.232645 +00:00', '2024-10-25 00:18:17.232645 +00:00', null, 'Account 5', '', 'email5@email.com'),
    (8, '2024-10-25 00:18:17.212917 +00:00', '2024-10-25 00:18:17.212917 +00:00', null, 'Account 8', '', 'email8@email.com'),
    (9, '2024-10-25 00:18:17.245009 +00:00', '2024-10-25 00:18:17.245009 +00:00', null, 'Account 9', '', 'email9@email.com'),
    (6, '2024-10-25 00:18:17.236577 +00:00', '2024-10-25 00:18:17.236577 +00:00', null, 'Account 6', '', 'email6@email.com'),
    (4, '2024-10-25 00:18:17.216388 +00:00', '2024-10-25 00:18:17.216388 +00:00', null, 'Account 4', '', 'email4@email.com'),
    (3, '2024-10-25 00:18:17.234306 +00:00', '2024-10-25 00:18:17.234306 +00:00', null, 'Account 3', '', 'email3@email.com'),
    (7, '2024-10-25 00:18:17.222316 +00:00', '2024-10-25 00:18:17.222316 +00:00', null, 'Account 7', '', 'email7@email.com')
ON CONFLICT (id) DO NOTHING;