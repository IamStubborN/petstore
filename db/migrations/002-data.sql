-- +migrate Up
-- +migrate StatementBegin
INSERT INTO category (name, price) VALUES ('Cat', '35.00');
INSERT INTO category (name, price) VALUES ('Dog', '49.99');
INSERT INTO category (name, price) VALUES ('Horse', '77.00');
INSERT INTO category (name, price) VALUES ('Monkey', '129.99');
INSERT INTO category (name, price) VALUES ('Tarzan', '999.99');

INSERT INTO order_status (name) VALUES ('placed');
INSERT INTO order_status (name) VALUES ('approved');
INSERT INTO order_status (name) VALUES ('delivered');

INSERT INTO pet_status (name) VALUES ('available');
INSERT INTO pet_status (name) VALUES ('pending');
INSERT INTO pet_status (name) VALUES ('sold');

INSERT INTO user_status (name, allowed_methods) VALUES ('admin', '{GET,POST,PUT,DELETE}');
INSERT INTO user_status (name, allowed_methods) VALUES ('guest', '{GET}');
INSERT INTO user_status (name, allowed_methods) VALUES ('developer', '{GET,POST,PUT,DELETE}');
INSERT INTO user_status (name, allowed_methods) VALUES ('tester', '{GET,POST,PUT}');

INSERT INTO tag (name) VALUES ('small');
INSERT INTO tag (name) VALUES ('average');
INSERT INTO tag (name) VALUES ('large');
INSERT INTO tag (name) VALUES ('best');
INSERT INTO tag (name) VALUES ('cool');

-- +migrate StatementEnd