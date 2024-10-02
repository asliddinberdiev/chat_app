package user

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/asliddinberdiev/chat_app/conf"
	"github.com/asliddinberdiev/chat_app/utils"
	"github.com/golang-jwt/jwt/v4"
)

type service struct {
	Repository
	timeout time.Duration
}

func NewService(repository Repository) Service {
	return &service{
		Repository: repository,
		timeout:    time.Duration(2) * time.Second,
	}
}

func (s *service) Create(ctx context.Context, req *CreateUserReq) (*CreateUserRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	uuid, err := utils.UUID()
	if err != nil {
		return nil, err
	}

	hashPass, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:       uuid,
		Username: req.Username,
		Email:    req.Email,
		Password: hashPass,
	}

	r, err := s.Repository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	res := &CreateUserRes{
		ID:       r.ID,
		Username: r.Username,
		Email:    r.Email,
	}

	return res, nil
}

type MyJWTClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (s *service) Login(ctx context.Context, req *LoginReq) (*LoginRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	user, err := s.Repository.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if ok := utils.CheckPassword(user.Password, req.Password); !ok {
		return nil, errors.New("email or password wrong")
	}

	accessTime, err := strconv.Atoi(conf.Cfg.App.AccessTime)
	if err != nil {
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyJWTClaims{
		ID:       user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    user.ID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(accessTime))),
		},
	})

	ss, err := token.SignedString([]byte(conf.Cfg.App.TokenKey))
	if err != nil {
		return nil, err
	}

	return &LoginRes{AccessToken: ss, ID: user.ID, Username: user.Username}, nil
}
