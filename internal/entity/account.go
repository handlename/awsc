package entity

type Account struct {
	profile string
}

func NewAccount(profile string) *Account {
	return &Account{profile: profile}
}

func (a Account) Profile() string { return a.profile }
