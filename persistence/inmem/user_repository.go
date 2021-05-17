package inmem

import (
	"echo-grpc/user"
	"fmt"
	"sync"
)


type userRepository struct {
	mu *sync.RWMutex
	users []user.User
}

func NewUserRepository() user.UserRepository {
	return &userRepository{}
}

func (r *userRepository) FindUser(username string) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.Unlock()
	for _, user := range r.users {
		if user.Username == username {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("username '%s' not found", username)
}

func (r *userRepository) ListUsers() ([]user.User, error) {
	return r.users, nil
}

func (r *userRepository) DeleteUser(username string) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.Unlock()
	for index, user := range r.users {
		if user.Username == username {
			r.mu.Lock()
			r.users = append(r.users[:index], r.users[index+1:]...)
			return &user, nil
		}
	}
	return nil, fmt.Errorf("username '%s' not found", username)
}
