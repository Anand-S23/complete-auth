package store

import (
	"context"
	"database/sql"

	"github.com/Anand-S23/complete-auth/internal/models"
)

type UserRepo interface {
    InsertUser(context.Context, *models.User) error
    InsertOAuthUser(ctx context.Context, user *models.User) error
}

// PostgresUserRepo

type PgUserRepo struct {
    Db *sql.DB
}

func NewPgUserRepo(db *sql.DB) *PgUserRepo{
    return &PgUserRepo {
        Db: db,
    }
}

func (pg *PgUserRepo) InsertUser(ctx context.Context, user *models.User) error {
    stmt, err := pg.Db.PrepareContext(ctx, `insert into users (id, email, password_hash) values ($1, $2, $3)`)
    if err != nil {
        return err
    }
    defer stmt.Close()

    _, err = stmt.ExecContext(ctx, user.ID, user.Email, user.Password)
    return err
}

func (pg *PgUserRepo) InsertOAuthUser(ctx context.Context, user *models.User) error {
    stmt, err := pg.Db.PrepareContext(ctx, `insert into users (id, email, oauth_provider, oauth_id) values ($1, $2, $3, $4)`)
    if err != nil {
        return err
    }
    defer stmt.Close()

    _, err = stmt.ExecContext(ctx, user.ID, user.Email, user.OAuthProvider, user.OAuthID)
    return err
}

