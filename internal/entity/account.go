package entity

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/morikuni/failure/v2"
	"github.com/rs/zerolog/log"
)

type Account struct {
	profile string
	region  string

	// Additional info originated from AWS STS
	// Refer documentation of aws-sdk-go-v2 for more information
	// https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/sts#GetCallerIdentityOutput

	id     string
	userID string
	arn    string
}

func NewAccount(ctx context.Context, profile string, config aws.Config, options ...AccountOption) (*Account, error) {
	a := &Account{
		profile: profile,
		region:  config.Region,
	}

	for _, opt := range options {
		if err := opt(ctx, a, config); err != nil {
			return nil, failure.Wrap(err)
		}
	}

	return a, nil
}

func (a *Account) Profile() string { return a.profile }
func (a *Account) Region() string  { return a.region }
func (a *Account) ID() string      { return a.id }
func (a *Account) UserID() string  { return a.userID }
func (a *Account) Arn() string     { return a.arn }

type AccountOption func(ctx context.Context, account *Account, config aws.Config) error

func AccountOptionWithAdditionalInfo(ctx context.Context, account *Account, config aws.Config) error {
	client := sts.NewFromConfig(config)
	res, err := client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		log.Warn().Err(err).Msg("failed to get caller identity for additional info")
		return nil
	}

	account.id = *res.Account
	account.userID = *res.UserId
	account.arn = *res.Arn

	return nil
}
