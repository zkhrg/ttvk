package repository

import (
	"database/sql"
	"fmt"
	"log"

	"back/domain/entities"
	"back/domain/interfaces"

	sq "github.com/Masterminds/squirrel"
)

type EntityRepo struct {
	DB      *sql.DB
	Builder sq.StatementBuilderType
}

func NewEntityRepo(db *sql.DB) interfaces.EntityRepository {
	return &EntityRepo{
		DB:      db,
		Builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *EntityRepo) GetByID(id string) (*entities.Entity, error) {
	query, args, err := r.Builder.
		Select("*").
		From("ip_logs").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.DB.QueryRow(query, args...)
	entity := &entities.Entity{}
	if err := row.Scan(&entity.ID, &entity.Address, &entity.PingTime, &entity.LastSuccess); err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *EntityRepo) EditByID(newAddress entities.EntityRequest) (*entities.Entity, error) {
	// Проверяем, существует ли запись с таким IP
	query, args, err := r.Builder.
		Select("id").
		From("ip_logs").
		Where(sq.Eq{"ip_address": newAddress.Address}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var existingID string
	err = r.DB.QueryRow(query, args...).Scan(&existingID)

	if err == nil {
		// Запись существует, выполняем обновление
		query, args, err = r.Builder.
			Update("ip_logs").
			Set("ping_time", newAddress.PingTime).
			Set("last_success", newAddress.LastSuccess).
			Where(sq.Eq{"ip_address": newAddress.Address}).
			ToSql()
		if err != nil {
			return nil, err
		}

		_, err = r.DB.Exec(query, args...)
		if err != nil {
			return nil, err
		}
	} else if err == sql.ErrNoRows {
		// Записи нет, создаем новую
		query, args, err = r.Builder.
			Insert("ip_logs").
			Columns("ip_address", "ping_time", "last_success").
			Values(newAddress.Address, newAddress.PingTime, newAddress.LastSuccess).
			ToSql()
		if err != nil {
			return nil, err
		}

		_, err = r.DB.Exec(query, args...)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	// Возвращаем актуальные данные из БД
	query, args, err = r.Builder.
		Select("id", "ip_address", "ping_time", "last_success").
		From("ip_logs").
		Where(sq.Eq{"ip_address": newAddress.Address}).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.DB.QueryRow(query, args...)
	entity := &entities.Entity{}
	if err := row.Scan(&entity.ID, &entity.Address, &entity.PingTime, &entity.LastSuccess); err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *EntityRepo) CreateAddress(addressData entities.EntityRequest) (*entities.Entity, error) {
	// Строим запрос на вставку нового пользователя
	query, args, err := r.Builder.
		Insert("ip_logs").
		Columns("ip_address", "ping_time", "last_success").
		Values(addressData.Address, addressData.PingTime, addressData.LastSuccess).
		Suffix("RETURNING id, ip_address, ping_time, last_success").
		ToSql()
	if err != nil {
		return nil, err
	}

	// Выполняем запрос на вставку
	row := r.DB.QueryRow(query, args...)

	// Создаем структуру для возврата вставленных данных
	entity := &entities.Entity{}
	if err := row.Scan(&entity.ID, &entity.Address, &entity.PingTime, &entity.LastSuccess); err != nil {
		return nil, err
	}

	// Возвращаем вставленный объект
	return entity, nil
}

func (r *EntityRepo) GetFullInfo() ([]*entities.Entity, error) {
	query, args, err := r.Builder.
		Select("*").
		From("ip_logs").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, _ := r.DB.Query(query, args...)

	// Slice to store results
	var addresses []*entities.Entity

	// Iterate over rows
	for rows.Next() {
		var address entities.Entity
		if err := rows.Scan(&address.ID, &address.Address, &address.PingTime, &address.LastSuccess); err != nil {
			log.Fatal(err)
		}
		addresses = append(addresses, &address)
	}

	// Check for errors in iteration
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Print results
	for _, address := range addresses {
		fmt.Printf("Address: %+v\n", *address)
	}

	return addresses, nil
}
