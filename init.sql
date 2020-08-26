CREATE DATABASE soccermanager;
USE soccermanager;

CREATE TABLE account (
   id INTEGER NOT NULL AUTO_INCREMENT,
    confirmed BIT NOT NULL,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    locked BIT NOT NULL,
    login_attempts INTEGER,
    password VARCHAR(255),
    profile VARCHAR(255),
    username VARCHAR(255),
    verification_token VARCHAR(255),
    PRIMARY KEY (id)
) engine=InnoDB;

CREATE TABLE player (
   id INTEGER NOT NULL AUTO_INCREMENT,
    age INTEGER,
    country VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    market_value INTEGER,
    position VARCHAR(255),
    team_id INTEGER,
    PRIMARY KEY (id)
) engine=InnoDB;

CREATE TABLE team (
   id INTEGER NOT NULL AUTO_INCREMENT,
    available_cash INTEGER,
    country VARCHAR(255),
    name VARCHAR(255),
    account_id INTEGER,
    PRIMARY KEY (id)
) engine=InnoDB;

CREATE TABLE transfer_list (
   id INTEGER NOT NULL AUTO_INCREMENT,
    asked_price INTEGER,
    market_value INTEGER,
    transferred BIT NOT NULL,
    player_id INTEGER,
    transferred_from INTEGER,
    transferred_to INTEGER,
    PRIMARY KEY (id)
) engine=InnoDB;

ALTER TABLE account
   ADD CONSTRAINT UK_gex1lmaqpg0ir5g1f5eftyaa1 UNIQUE (username);

ALTER TABLE player
   ADD CONSTRAINT FKdvd6ljes11r44igawmpm1mc5s
   FOREIGN KEY (team_id)
   REFERENCES team (id);

ALTER TABLE team
   ADD CONSTRAINT FK3p8m6jcr5un0gcq49y0ento02
   FOREIGN KEY (account_id)
   REFERENCES account (id);

ALTER TABLE transfer_list
   ADD CONSTRAINT FK5ta1ls744ss66fvgwagecsros
   FOREIGN KEY (player_id)
   REFERENCES player (id);

ALTER TABLE transfer_list
   ADD CONSTRAINT FKgoyf6slgja6unsb9g0xkbgfjm
   FOREIGN KEY (transferred_from)
   REFERENCES team (id);

ALTER TABLE transfer_list
   ADD CONSTRAINT FK5bpxwc2m8q5w6r1oy3a8x1t5g
   FOREIGN KEY (transferred_to)
   REFERENCES team (id);