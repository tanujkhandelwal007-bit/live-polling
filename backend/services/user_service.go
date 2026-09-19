package services

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"live-polling-backend/models"
	"live-polling-backend/repositories"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository *repositories.UserRepository
}

func NewUserService(
	repository *repositories.UserRepository,
) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (s *UserService) Register(
	ctx context.Context,
	name string,
	email string,
	password string,
) (*models.User, error) {

	name = strings.TrimSpace(name)
	email = strings.TrimSpace(strings.ToLower(email))

	if name == "" {
		return nil, errors.New("name cannot be empty")
	}

	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	if password == "" {
		return nil, errors.New("password cannot be empty")
	}

	if len(password) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}

	// Check whether email already exists.
	existingUser, err := s.repository.GetUserByEmail(
		ctx,
		email,
	)

	if err == nil && existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password before saving it.
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:        bson.NewObjectID(),
		Name:      name,
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}

	err = s.repository.CreateUser(
		ctx,
		user,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(
	ctx context.Context,
	email string,
	password string,
) (string, *models.User, error) {

	email = strings.TrimSpace(strings.ToLower(email))

	user, err := s.repository.GetUserByEmail(
		ctx,
		email,
	)

	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)

	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	// JWT secret from environment variable.
	jwtSecret := os.Getenv("JWT_SECRET")

	// Local development fallback.
	if jwtSecret == "" {
		jwtSecret = "live-polling-secret-key"
	}

	// Create JWT token.
	claims := jwt.MapClaims{
		"userId": user.ID.Hex(),
		"email":  user.Email,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	tokenString, err := token.SignedString(
		[]byte(jwtSecret),
	)

	if err != nil {
		return "", nil, err
	}

	return tokenString, user, nil
}
