package usecase

import (
	"context"
	"lucassantoss1701/bank/internal/entity"
	"time"

	"github.com/golang-jwt/jwt"
)

type ILoginUseCase interface {
	Execute(ctx context.Context, input *LoginUseCaseInput) (*LoginUseCaseOutput, error)
}

type LoginUseCase struct {
	repostiory entity.AccountRepository
}

func NewLoginUseCase(repostiory entity.AccountRepository) *LoginUseCase {
	return &LoginUseCase{
		repostiory: repostiory,
	}
}

func (l *LoginUseCase) Execute(ctx context.Context, input *LoginUseCaseInput) (*LoginUseCaseOutput, error) {

	account, err := l.repostiory.FindByCPF(ctx, input.CPF)
	if err != nil {
		return nil, err
	}

	if !account.SecretIsCorrect(input.Secret) {
		return nil, entity.NewErrorHandler(entity.UNAUTHORIZED_ERROR).Add("secret is incorrect")
	}

	return NewLoginUseCaseOutput(&account, input.SecretJWT), nil

}

type LoginUseCaseInput struct {
	CPF       string `json:"cpf"`
	Secret    string `json:"secret"`
	SecretJWT string `json:"-"`
}

func NewLoginUseCaseInput(CPF string, secret string, secretJWT string) *LoginUseCaseInput {
	return &LoginUseCaseInput{
		CPF:       CPF,
		Secret:    secret,
		SecretJWT: secretJWT,
	}
}

type LoginUseCaseOutput struct {
	Token string `json:"token"`
}

func NewLoginUseCaseOutput(account *entity.Account, secretJWT string) *LoginUseCaseOutput {
	claims := jwt.MapClaims{}
	claims["account_id"] = account.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secretJWT))

	return &LoginUseCaseOutput{
		Token: tokenString,
	}
}
