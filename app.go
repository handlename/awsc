package myapp

import (
	"context"

	errcode "github.com/handlename/my-golang-template/internal/errorcode"
	"github.com/morikuni/failure/v2"
)

type App struct{}

func NewApp() *App {
	return &App{}
}

func (a *App) Run(ctx context.Context) error {
	return failure.New(errcode.ErrNotImplemented, failure.Message("not implemented yet"))
}
