-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
 CREATE TABLE Reserved ( --Можно же через таблицу с заморозкой средств реализовать последовательность операций для отдельного user
    userID      int primary key,
    balance     int
 );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE Reserved;
-- +goose StatementEnd
