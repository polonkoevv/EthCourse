CREATE DATABASE IF NOT EXISTS ipfs;

-- Создание последовательности для автоинкремента music_id
CREATE SEQUENCE IF NOT EXISTS music_music_id_seq;

-- Создание таблицы music
CREATE TABLE IF NOT EXISTS music (
    music_id smallint DEFAULT nextval('music_music_id_seq') PRIMARY KEY,
    title character varying(100),
    cid character varying(100)
);

-- Установка владельца последовательности
ALTER SEQUENCE music_music_id_seq OWNED BY music.music_id;