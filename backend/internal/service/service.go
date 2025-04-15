package service

import (
	"context"
	"os"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/polonkoevv/ethcourse/internal/model"
	"github.com/polonkoevv/ethcourse/internal/storage/postgres"
)

type Service struct {
	sh *shell.Shell
	pg *postgres.Postgres
}

func NewService(sh *shell.Shell, pg *postgres.Postgres) *Service {
	return &Service{sh: sh, pg: pg}
}

func (s *Service) UploadFile(ctx context.Context, name string, file *os.File) (string, error) {
	cid, err := s.sh.Add(file)
	if err != nil {
		return "", err
	}

	err = s.sh.Pin(cid)
	if err != nil {
		return "", err
	}

	music := model.Music{
		Title: name,
		CID:   cid,
	}

	_, err = s.pg.CreateMusic(ctx, music)
	if err != nil {
		return "", err
	}
	return cid, nil
}

func (s *Service) GetAllMusic(ctx context.Context) ([]model.Music, error) {
	return s.pg.GetAllMusic(ctx)
}
