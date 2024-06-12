package store

type Store struct {
    UserRepo UserRepo
}

func NewStore(ur UserRepo) *Store {
    return &Store {
        UserRepo: ur,
    }
}

