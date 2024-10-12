package entity

import "github.com/aws/aws-sdk-go-v2/aws"

type Account struct {
	profile string
	region  string
}

func NewAccount(profile string, config aws.Config) *Account {
	return &Account{
		profile: profile,
		region:  config.Region,
	}
}

func NewAccount(profile string) *Account {
	return &Account{profile: profile}
}

func (a Account) Profile() string { return a.profile }
