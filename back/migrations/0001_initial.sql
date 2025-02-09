-- Активируем расширение для работы с UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Создаём таблицу для хранения информации об IP-адресах
CREATE TABLE IF NOT EXISTS ip_logs (
   id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
   ip_address INET NOT NULL,
   ping_time INT NOT NULL, -- Время пинга в миллисекундах
   last_success TIMESTAMPTZ NOT NULL
);