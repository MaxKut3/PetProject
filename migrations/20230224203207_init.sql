-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE Transactions (
    id_transaction      serial primary key,
    type_transaction    int, -- 0 - списание / 1 - пополнение
    userID              int,
    id_service          int,
    amount              int,
    time_oper           date default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE Transactions;
-- +goose StatementEnd
