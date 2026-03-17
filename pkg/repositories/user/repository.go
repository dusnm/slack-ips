package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/dusnm/slack-ips/pkg/dto/command"
	"github.com/dusnm/slack-ips/pkg/models"
	"github.com/dusnm/slack-ips/pkg/repositories"
	"github.com/rs/zerolog"
)

const (
	findByIDQuery   = "SELECT id, username, name, bank_account_number, city, ips_string FROM users WHERE id = ?"
	insertQuery     = `INSERT INTO users (id, username, name, bank_account_number, city, ips_string) VALUES (?, ?, ?, ?, ?, ?)`
	deleteByIDQuery = `DELETE FROM users WHERE id = ?`
)

type (
	Repository struct {
		db             *sql.DB
		findByIDStmt   *sql.Stmt
		insertStmt     *sql.Stmt
		deleteByIDStmt *sql.Stmt
		logger         zerolog.Logger
	}
)

func New(
	db *sql.DB,
	logger zerolog.Logger,
) *Repository {
	findByIDStmt, err := db.Prepare(findByIDQuery)
	if err != nil {
		logger.
			Fatal().
			Err(err).
			Msg("failed to prepare find by id statement")
	}

	insertStmt, err := db.Prepare(insertQuery)
	if err != nil {
		logger.
			Fatal().
			Err(err).
			Msg("failed to prepare insert statement")
	}

	deleteByIDStmt, err := db.Prepare(deleteByIDQuery)
	if err != nil {
		logger.
			Fatal().
			Err(err).
			Msg("failed to prepare delete statement")
	}

	return &Repository{
		db:             db,
		findByIDStmt:   findByIDStmt,
		insertStmt:     insertStmt,
		deleteByIDStmt: deleteByIDStmt,
		logger:         logger,
	}
}

func (r *Repository) Close() error {
	r.logger.Info().Msg("closing")
	return errors.Join(
		r.findByIDStmt.Close(),
		r.insertStmt.Close(),
		r.deleteByIDStmt.Close(),
	)
}

func (r *Repository) FindByID(ctx context.Context, ID string) (models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result := models.User{}
	row := r.findByIDStmt.QueryRowContext(ctx, ID)
	err := row.Scan(
		&result.ID,
		&result.Username,
		&result.Name,
		&result.BankAccountNumber,
		&result.City,
		&result.IPSString,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, repositories.ErrNotFound
		}

		return models.User{}, err
	}

	return result, nil
}

func (r *Repository) Insert(ctx context.Context, payload command.Init) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := r.insertStmt.ExecContext(
		ctx,
		payload.UserID,
		payload.UserName,
		payload.Name,
		payload.BankAccountNumber,
		payload.City,
		payload.ToIPSString(),
	)

	return err
}

func (r *Repository) DeleteByID(ctx context.Context, ID string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := r.deleteByIDStmt.ExecContext(ctx, ID)
	return err
}
