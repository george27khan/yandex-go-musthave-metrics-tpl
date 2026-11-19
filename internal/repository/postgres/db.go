package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"yandex-go-musthave-metrics-tpl/internal/model"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBStorage struct {
	db *sql.DB
}

func New(dsn string) (*DBStorage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	if err = initDB(db); err != nil {
		return nil, fmt.Errorf("failed to init database: %w", err)
	}
	return &DBStorage{db: db}, nil
}

func (s *DBStorage) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s *DBStorage) Add(ctx context.Context, m model.Metrics) error {

	res, err := s.db.ExecContext(
		ctx,
		`update metrics 
				  set delta = $1,
					  value = $2
			    where name = $3`,
		m.Delta, m.Value, m.ID)
	if err != nil {
		return err
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if cnt == 0 {
		res, err = s.db.ExecContext(
			ctx,
			`insert into metrics(name,type,delta, value,hash)
				   values ($1,$2,$3,$4,$5)`,
			m.ID, m.MType, m.Delta, m.Value, m.Hash)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *DBStorage) GetAll(ctx context.Context) ([]model.Metrics, error) {
	var metrics []model.Metrics
	rows, err := s.db.QueryContext(ctx,
		`select t.name, t.type, t.delta, t.value, t.hash 
				from metrics t order by t.name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var m model.Metrics
		if err = rows.Scan(&m.ID, &m.MType, &m.Delta, &m.Value, &m.Hash); err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}
	return metrics, nil
}

func (s *DBStorage) Get(ctx context.Context, name string) (model.Metrics, error) {
	row := s.db.QueryRowContext(ctx,
		`select t.name, t.type, t.delta, t.value, t.hash 
				from metrics t where t.name = $1`, name)
	var m model.Metrics
	if err := row.Scan(&m.ID, &m.MType, &m.Delta, &m.Value, &m.Hash); err != nil {
		return model.Metrics{}, err
	}
	return m, nil
}

func initDB(db *sql.DB) error {
	_, err := db.Exec(`create table if not exists metrics (
    	name varchar PRIMARY KEY,
    	type varchar,
    	delta int,
        value double precision,
        hash varchar)`)
	return err
}
