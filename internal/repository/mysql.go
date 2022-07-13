package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/albertogviana/port-service/internal/port"
)

type mysqlPort struct {
	Name      string
	City      string
	Country   string
	Alias     []byte
	Regions   []byte
	Latitude  float64
	Longitude float64
	Province  string
	Timezone  string
	Unloc     string
	Code      string
	CreateAt  time.Time
	UpdateAt  time.Time
}

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{
		db: db,
	}
}

func (r *MySQLRepository) Create(ctx context.Context, p *port.Port) error {
	now := time.Now()

	mp := mysqlPort{
		Name:     p.Name,
		City:     p.City,
		Country:  p.Country,
		Province: p.Province,
		Timezone: p.Timezone,
		Unloc:    p.Unlocs[0],
		Code:     p.Code,
		CreateAt: now,
		UpdateAt: now,
	}

	var err error

	mp.Regions, err = json.Marshal(p.Regions)
	if err != nil {
		return fmt.Errorf(
			"error during the marshal json on repository.MySQLRepository.Create: %w",
			err,
		)
	}

	mp.Alias, err = json.Marshal(p.Alias)
	if err != nil {
		return fmt.Errorf(
			"error during the marshal json on repository.MySQLRepository.Create: %w",
			err,
		)
	}

	if len(p.Coordinates) == 2 {
		mp.Latitude = p.Coordinates[0]
		mp.Longitude = p.Coordinates[1]
	}

	q := `
		INSERT INTO port(
			unloc,
			name,
			city,
			country,
			alias,
			regions,
			latitude,
			longitude,
			province,
			timezone,
			code,
			created_at,
			updated_at
		)
		VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?)
	`

	stmt, err := r.db.PrepareContext(ctx, q)
	if err != nil {
		return fmt.Errorf(
			"error during the prepare statement on repository.MySQLRepository.Create: %w",
			err,
		)
	}

	defer stmt.Close()

	_, err = stmt.Exec(
		mp.Unloc,
		mp.Name,
		mp.City,
		mp.Country,
		mp.Alias,
		mp.Regions,
		mp.Latitude,
		mp.Longitude,
		mp.Province,
		mp.Timezone,
		mp.Code,
		mp.CreateAt,
		mp.UpdateAt,
	)

	if err != nil {
		return fmt.Errorf(
			"error during the exec statement on repository.MySQLRepository.Create: %w",
			err,
		)
	}

	return nil
}

func (r *MySQLRepository) Update(ctx context.Context, p *port.Port) error {
	now := time.Now()

	mp := mysqlPort{
		Name:     p.Name,
		City:     p.City,
		Country:  p.Country,
		Province: p.Province,
		Timezone: p.Timezone,
		Code:     p.Code,
		Unloc:    p.Unlocs[0],
		UpdateAt: now,
	}

	var err error

	mp.Regions, err = json.Marshal(p.Regions)
	if err != nil {
		return fmt.Errorf(
			"error during the marshal json on repository.MySQLRepository.Create: %w",
			err,
		)
	}

	mp.Alias, err = json.Marshal(p.Alias)
	if err != nil {
		return fmt.Errorf(
			"error during the marshal json on repository.MySQLRepository.Create: %w",
			err,
		)
	}

	if len(p.Coordinates) == 2 {
		mp.Latitude = p.Coordinates[0]
		mp.Longitude = p.Coordinates[1]
	}

	q := `
		UPDATE port SET
			name = ?,
			city = ?,
			country = ?,
			alias = ?,
			regions = ?,
			latitude = ?,
			longitude = ?,
			province = ?,
			timezone = ?,
			code = ?,
			updated_at = ?
		WHERE unloc = ?
	`

	stmt, err := r.db.PrepareContext(ctx, q)
	if err != nil {
		return fmt.Errorf(
			"error during the prepare statement on repository.MySQLRepository.Update: %w",
			err,
		)
	}

	defer stmt.Close()

	result, err := stmt.Exec(
		mp.Name,
		mp.City,
		mp.Country,
		mp.Alias,
		mp.Regions,
		mp.Latitude,
		mp.Longitude,
		mp.Province,
		mp.Timezone,
		mp.Code,
		mp.UpdateAt,
		mp.Unloc,
	)
	if err != nil {
		return fmt.Errorf(
			"error during the exec statement on repository.MySQLRepository.Update: %w",
			err,
		)
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return fmt.Errorf(
			"error during the exec statement on repository.MySQLRepository.Update: %w",
			err,
		)
	}

	return nil
}

func (r *MySQLRepository) FindByUnLoc(ctx context.Context, unloc string) (*port.Port, error) {
	q := `
		SELECT 
			unloc,
			name,
			city,
			country,
			alias,
			regions,
			latitude,
			longitude,
			province,
			timezone,
			code,
			created_at,
			updated_at
		FROM port
		WHERE unloc = ?
	`
	stmt, err := r.db.PrepareContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf(
			"error during the prepare statement on repository.MySQLRepository.FindByUnLoc: %w",
			err,
		)
	}

	defer stmt.Close()

	var mp mysqlPort

	err = stmt.QueryRowContext(ctx, unloc).Scan(
		&mp.Unloc,
		&mp.Name,
		&mp.City,
		&mp.Country,
		&mp.Alias,
		&mp.Regions,
		&mp.Latitude,
		&mp.Longitude,
		&mp.Province,
		&mp.Timezone,
		&mp.Code,
		&mp.CreateAt,
		&mp.UpdateAt,
	)

	switch err {
	case nil:
		p := port.Port{
			Name:    mp.Name,
			City:    mp.City,
			Country: mp.Country,
			Coordinates: []float64{
				mp.Latitude,
				mp.Longitude,
			},
			Province: mp.Province,
			Timezone: mp.Timezone,
			Code:     mp.Code,
			Unlocs:   []string{mp.Unloc},
		}

		err := json.Unmarshal(mp.Regions, &p.Regions)
		if err != nil {
			return nil, fmt.Errorf(
				"error during the unmarshal json on repository.MySQLRepository.FindByUnLoc: %w",
				err,
			)
		}

		err = json.Unmarshal(mp.Alias, &p.Alias)
		if err != nil {
			return nil, fmt.Errorf(
				"error during the unmarshal json on repository.MySQLRepository.FindByUnLoc: %w",
				err,
			)
		}

		return &p, nil
	case sql.ErrNoRows:
		return nil, port.ErrPortNotFound
	default:
		return nil, fmt.Errorf(
			"error during the query execution on repository.MySQLRepository.FindByUnLoc: %w",
			err,
		)
	}
}
