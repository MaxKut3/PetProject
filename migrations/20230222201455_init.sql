-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE Balance (
    userID      int primary key,
    balance     int CHECK ( balance >= 0 )
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE Balance;
-- +goose StatementEnd
