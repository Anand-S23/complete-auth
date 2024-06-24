package store

import (
	"context"
	"database/sql"

	"github.com/Anand-S23/complete-auth/internal/models"
)

type UserRepo interface {
    InsertUserWithProfile(context.Context, *models.User) error
    // TODO: Merge into one
    InsertUser(context.Context, *models.User) error
    InsertOAuthUser(context.Context, *models.User) error
    // 
    UpdateUser(context.Context, *models.User) (*models.User, error)
    GetUserByID(context.Context, string) (*models.User, error)
    GetUserByEmail(context.Context, string) (*models.User, error)
    DeleteUser()

    GetUserProfileByID(context.Context, string) (*models.User, error)
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
    tx, err := pg.Db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }

    defer func() {
        if err != nil {
            tx.Rollback()
        } else {
            err = tx.Commit()
        }
    }()

    userInsertStatment, err := tx.PrepareContext(ctx, `insert into users (id, email, password_hash) values ($1, $2, $3)`)
    if err != nil {
        return err
    }
    defer userInsertStatment.Close()

    _, err = userInsertStatment.ExecContext(ctx, user.ID, user.Email, user.Password)
    if err != nil {
        return err
    }

    profileInsertStatment, err := tx.PrepareContext(ctx, `insert into user_profiles (user_id, first_name, last_name, phone_number, pfp_url) values ($1, $2, $3, $4, $5)`)
    if err != nil {
        return err
    }
    defer profileInsertStatment.Close()

    _, err = profileInsertStatment.ExecContext(ctx, user.ID, user.Email, user.Password)
    if err != nil {
        return err
    }

    return nil
}

func (pg *PgUserRepo) InsertOAuthUser(ctx context.Context, user *models.User) error {
    tx, err := pg.Db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }

    defer func() {
        if err != nil {
            tx.Rollback()
        } else {
            err = tx.Commit()
        }
    }()

    userInsertStatment, err := tx.PrepareContext(ctx, `insert into users (id, email, oauth_provider, oauth_id) values ($1, $2, $3, $4)`)
    if err != nil {
        return err
    }
    defer userInsertStatment.Close()

    _, err = userInsertStatment.ExecContext(ctx, user.ID, user.Email, user.OAuthProvider, user.OAuthID)
    if err != nil {
        return err
    }

    profileInsertStatment, err := tx.PrepareContext(ctx, `insert into user_profiles (user_id, first_name, last_name, phone_number, pfp_url) values ($1, $2, $3, $4, $5)`)
    if err != nil {
        return err
    }
    defer profileInsertStatment.Close()

    _, err = profileInsertStatment.ExecContext(ctx, user.ID, user.Email, user.Password)
    if err != nil {
        return err
    }

    return nil
}

