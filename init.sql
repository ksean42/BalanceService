CREATE TABLE IF NOT EXISTS balance (
	user_id serial,
	balance DECIMAL (19,4),
	CONSTRAINT PK_balance_id PRIMARY KEY(user_id)
);

CREATE TABLE IF NOT EXISTS history (
	transaction_id serial,
	order_id serial, --???
	service_id int NOT NULL,
	user_id int NOT NULL,
	amount DECIMAL(19,4) NOT NULL,
	date timestamp NOT NULL,
	approved boolean DEFAULT 'false' NOT NULL,
	comment text,
	CONSTRAINT PK_transaction_id PRIMARY KEY(transaction_id),
	CONSTRAINT FK_history_order_id FOREIGN KEY(order_id) REFERENCES reserve_account(reserve_id),
	CONSTRAINT FK_history_service_id FOREIGN KEY(service_id) REFERENCES service(service_id),
	CONSTRAINT FK_history_user_id FOREIGN KEY(user_id) REFERENCES balance(user_id)
);

CREATE TABLE IF NOT EXISTS reserve_account (
	reserve_id serial,
	order_id int UNIQUE NOT NULL,
	service_id int NOT NULL,
	user_id int NOT NULL,
	amount DECIMAL(19,4) NOT NULL,
	CONSTRAINT PK_reserve_account_reserve_id PRIMARY KEY(reserve_id),
        CONSTRAINT FK_reserve_account_service_id FOREIGN KEY(service_id) REFERENCES service(service_id),
        CONSTRAINT FK_reserve_account_user_id FOREIGN KEY(user_id) REFERENCES balance(user_id)
);

CREATE INDEX idx_history_date ON history(date);

INSERT INTO service(name, price) VALUES
('service_1', 300),
('service_2', 560),
('service_3', 300);

drop table history ; drop table reserve_account ; drop table service ; drop table balance ;
