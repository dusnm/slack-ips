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
	findByIDQuery = `
	SELECT u.id, 
	       u.username, 
	       u.name, 
	       u.bank_account_number, 
	       u.city, 
	       u.ips_string,
	       s.qr_fg_color,
	       s.qr_bg_color,
	       s.qr_shape,
	       s.qr_logo,
	       s.qr_show_logo
	FROM users u LEFT JOIN settings s ON u.id = s.user_id
	WHERE u.id = ?
	`
	insertQuery     = `INSERT INTO users (id, username, name, bank_account_number, city, ips_string) VALUES (?, ?, ?, ?, ?, ?)`
	updateQuery     = `UPDATE users SET name = ?, bank_account_number = ?, city = ?, ips_string = ? WHERE id = ?`
	deleteByIDQuery = `DELETE FROM users WHERE id = ?`
)

type (
	Repository struct {
		db             *sql.DB
		findByIDStmt   *sql.Stmt
		insertStmt     *sql.Stmt
		updateStmt     *sql.Stmt
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

	updateStmt, err := db.Prepare(updateQuery)
	if err != nil {
		logger.
			Fatal().
			Err(err).
			Msg("failed to prepare update statement")
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
		updateStmt:     updateStmt,
		deleteByIDStmt: deleteByIDStmt,
		logger:         logger,
	}
}

func (r *Repository) Close() error {
	r.logger.Info().Msg("closing")
	return errors.Join(
		r.findByIDStmt.Close(),
		r.insertStmt.Close(),
		r.updateStmt.Close(),
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
		&result.Settings.QRFGColor,
		&result.Settings.QRBGColor,
		&result.Settings.QRShape,
		&result.Settings.QRLogo,
		&result.Settings.QRShowLogo,
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

func (r *Repository) Update(ctx context.Context, ID string, payload command.Init) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := r.updateStmt.ExecContext(
		ctx,
		payload.Name,
		payload.BankAccountNumber,
		payload.City,
		payload.ToIPSString(),
		ID,
	)

	return err
}

func (r *Repository) DeleteByID(ctx context.Context, ID string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := r.deleteByIDStmt.ExecContext(ctx, ID)
	return err
}
