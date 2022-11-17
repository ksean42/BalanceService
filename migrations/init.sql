CREATE TABLE IF NOT EXISTS balance (
	user_id serial,
	balance DECIMAL (19,2),
	CONSTRAINT PK_balance_id PRIMARY KEY(user_id)
);

CREATE TABLE IF NOT EXISTS service (
	service_id serial,
	name varchar(20) NOT NULL,
	price DECIMAL (19,4) NOT NULL,
	CONSTRAINT PK_service_service_id PRIMARY KEY(service_id)
);

CREATE TABLE IF NOT EXISTS reserve_account (
	reserve_id serial,
	order_id int UNIQUE NOT NULL,
	service_id int NOT NULL,
	user_id int NOT NULL,
	amount DECIMAL(19,2) NOT NULL,
	CONSTRAINT PK_reserve_account_reserve_id PRIMARY KEY(reserve_id),
    --CONSTRAINT FK_reserve_account_service_id FOREIGN KEY(service_id) REFERENCES service(service_id),
    CONSTRAINT FK_reserve_account_user_id FOREIGN KEY(user_id) REFERENCES balance(user_id)
);

CREATE TABLE IF NOT EXISTS history (
	transaction_id serial,
	order_id int UNIQUE NOT NULL, --???
	service_id int NOT NULL,
	user_id int NOT NULL,
	amount DECIMAL(19,2) NOT NULL,
	date timestamp NOT NULL,
	comment text,
	CONSTRAINT PK_transaction_id PRIMARY KEY(transaction_id),
	--CONSTRAINT FK_history_service_id FOREIGN KEY(service_id) REFERENCES service(service_id),
	CONSTRAINT FK_history_user_id FOREIGN KEY(user_id) REFERENCES balance(user_id)
);


CREATE INDEX idx_history_date ON history(date);
CREATE INDEX idx_reserve_order_id ON reserve_account(order_id);