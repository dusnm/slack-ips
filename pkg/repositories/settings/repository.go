package settings

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
	upsertQuery = `
	INSERT INTO settings (
		user_id, 
	    qr_fg_color, 
		qr_bg_color, 
	    qr_shape, 
	    qr_logo,
	    qr_show_logo
	) 
	VALUES (?, ?, ?, ?, ?, ?) 
	ON CONFLICT(user_id) DO 
	UPDATE SET 
		qr_fg_color=excluded.qr_fg_color,
	    qr_bg_color=excluded.qr_bg_color,
	    qr_shape=excluded.qr_shape,
	    qr_logo=
	        CASE 
				WHEN excluded.qr_logo IS NOT NULL AND settings.qr_logo != excluded.qr_logo
				THEN excluded.qr_logo
				ELSE settings.qr_logo
			END,
	    qr_show_logo=excluded.qr_show_logo
	RETURNING 
	    qr_fg_color,
	    qr_bg_color,
	    qr_shape,
	    qr_logo,
	    qr_show_logo
	`
)

type (
	Repository struct {
		db         *sql.DB
		upsertStmt *sql.Stmt
		logger     zerolog.Logger
	}
)

func New(
	db *sql.DB,
	logger zerolog.Logger,
) *Repository {
	stmt, err := db.Prepare(upsertQuery)
	if err != nil {
		logger.
			Fatal().
			Err(err).
			Msg("failed to prepare upsert statement")
	}

	return &Repository{
		db:         db,
		upsertStmt: stmt,
		logger:     logger,
	}
}

func (r *Repository) Close() error {
	r.logger.Info().Msg("closing")
	return r.upsertStmt.Close()
}

func (r *Repository) UpsertByUserID(
	ctx context.Context,
	userID string,
	payload command.Settings,
) (models.Settings, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result := models.Settings{}
	row := r.upsertStmt.QueryRowContext(
		ctx,
		userID,
		payload.QRFGColor,
		payload.QRBGColor,
		payload.QRShape,
		payload.QRLogo,
		payload.QRShowLogo,
	)

	err := row.Scan(
		&result.QRFGColor,
		&result.QRBGColor,
		&result.QRShape,
		&result.QRLogo,
		&result.QRShowLogo,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Settings{}, repositories.ErrNotFound
		}

		return models.Settings{}, err
	}

	return result, nil
}
