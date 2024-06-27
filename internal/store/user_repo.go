package store

import (
	"context"

	"github.com/Anand-S23/complete-auth/internal/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type UserRepo interface {
    InsertUser(context.Context, *models.User) error

    GetUser(context.Context, string) (*models.User, error)
    GetUserByEmail(context.Context, string) (*models.User, error)
    GetBaseUser(context.Context, string) (*models.User, error)
    GetBaseUserByEmail(context.Context, string) (*models.User, error)
    GetUserProfile(context.Context, string) (*models.User, error)

    UpdateUser(context.Context, *models.User) error
    UpdateBaseUser(context.Context, *models.User) error
    UpdateUserProfile(context.Context, *models.UserProfile) error

    DeleteUser(context.Context, string) error
}

// PostgresUserRepo

type PgUserRepo struct {
    Db *sqlx.DB
}

func NewPgUserRepo(db *sqlx.DB) *PgUserRepo{
    return &PgUserRepo {
        Db: db,
    }
}

func (pg *PgUserRepo) InsertUser(ctx context.Context, user *models.User) error {
    tx, err := pg.Db.BeginTxx(ctx, nil)
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

    insertUserCommand := `
        insert into users (id, email, password_hash, oauth_provider, oauth_id) 
        values (:id, :email, :password_hash, :oauth_provider, :oauth_id)
    `
    _, err = tx.NamedExecContext(ctx, insertUserCommand, user)
    if err != nil {
        return err
    }

    insertProfileCommand := `
        insert into user_profiles (user_id, first_name, last_name, phone_number, pfp_url) 
        values (:user_id, :first_name, :last_name, :phone_number, :pfp_url)
    `
    _, err = tx.NamedExecContext(ctx, insertProfileCommand, user.Profile)
    if err != nil {
        return err
    }

    return nil
}

func (pg *PgUserRepo) GetUser(ctx context.Context, id string) (*models.User, error) {
    tx, err := pg.Db.BeginTxx(ctx, nil)
    if err != nil {
        return nil, err
    }

    defer func() {
        if err != nil {
            tx.Rollback()
        } else {
            err = tx.Commit()
        }
    }()

    var user models.User
    err = tx.GetContext(ctx, &user, `select * from users where id = $1`, id)
    if err != nil {
        return nil, err
    }

    err = tx.GetContext(ctx, &user.Profile, `select * from user_profiles where user_id = $1`, id)
    if err != nil {
        return nil, err
    }

    return &user, nil
}

func (pg *PgUserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
    tx, err := pg.Db.BeginTxx(ctx, nil)
    if err != nil {
        return nil, err
    }

    defer func() {
        if err != nil {
            tx.Rollback()
        } else {
            err = tx.Commit()
        }
    }()

    var user models.User
    err = tx.GetContext(ctx, &user, `select * from users where email = $1`, email)
    if err != nil {
        return nil, err
    }

    err = tx.GetContext(ctx, &user.Profile, `select * from user_profiles where user_id = $1`, user.ID)
    if err != nil {
        return nil, err
    }

    return &user, nil
}

func (pg *PgUserRepo) GetBaseUser(ctx context.Context, id string) (*models.User, error) {
    var user models.User
    err := pg.Db.GetContext(ctx, &user, `select * from users where id = $1`, id)
    return &user, err
}

func (pg *PgUserRepo) GetBaseUserByEmail(ctx context.Context, email string) (*models.User, error) {
    var user models.User
    err := pg.Db.GetContext(ctx, &user, `select * from users where email = $1`, email)
    return &user, err
}

func (pg *PgUserRepo) GetUserProfile(ctx context.Context, id string) (*models.UserProfile, error) {
    var profile models.UserProfile
    err := pg.Db.GetContext(ctx, &profile, `select * from user_profiles where user_id = $1`, id)
    return &profile, err
}

func (pg *PgUserRepo) UpdateUser(ctx context.Context, user *models.User) error {
    tx, err := pg.Db.BeginTxx(ctx, nil)
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

    updateUserCommand := `update users set email = :email, password_hash = :password_hash where id = $1`
    _, err = tx.NamedExecContext(ctx, updateUserCommand, user)
    if err != nil {
        return nil
    }

    updateProfileCommand := `update user_profiles set first_name = :first_name, last_name = :last_name, phone_number = :phone_number, pfp_url = :pfp_url where user_id = $1`
    _, err = tx.NamedExecContext(ctx, updateProfileCommand, user.Profile)
    if err != nil {
        return nil
    }

    return nil
}

func (pg *PgUserRepo) UpdateBaseUser(ctx context.Context, user *models.User) error {
    updateUserCommand := `update users set email = :email, password_hash = :password_hash where id = $1`
    _, err := pg.Db.NamedExecContext(ctx, updateUserCommand, user)
    return err
}

func (pg *PgUserRepo) UpdateUserProfile(ctx context.Context, user *models.User) error {
    updateProfileCommand := `update user_profiles set first_name = :first_name, last_name = :last_name, phone_number = :phone_number, pfp_url = :pfp_url where user_id = $1`
    _, err := pg.Db.NamedExecContext(ctx, updateProfileCommand, user)
    return err
}

func (pg *PgUserRepo) DeleteUser(ctx context.Context, id string) error {
    tx, err := pg.Db.BeginTxx(ctx, nil)
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

    _, err = tx.ExecContext(ctx, `delete from users where id = $1`, id)
    if err != nil {
        return err
    }

    _, err = tx.ExecContext(ctx, `delete from user_profiles where user_id = $1`, id)
    if err != nil {
        return err
    }

    return nil
}

