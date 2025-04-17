package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/polonkoevv/ethcourse/internal/model"
)

func NewPostgres(host, port, user, password, dbname string) (*Postgres, error) {
	link := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)
	conn, err := pgx.Connect(context.Background(), link)
	if err != nil {
		return nil, err
	}

	return &Postgres{conn: conn}, nil
}

type Postgres struct {
	conn *pgx.Conn
}

func (p *Postgres) GetMusicById(ctx context.Context, id int) (*model.Music, error) {
	var music model.Music
	err := p.conn.QueryRow(ctx, "SELECT * FROM music WHERE music_id = $1", id).Scan(&music.ID, &music.Title, &music.CID, &music.OwnerAddr, &music.Signature, &music.UploadedAt)
	if err != nil {
		return nil, err
	}
	return &music, nil
}

func (p *Postgres) GetMusicByCID(ctx context.Context, cid string) (*model.Music, error) {
	var music model.Music
	err := p.conn.QueryRow(ctx, "SELECT * FROM music WHERE cid = $1", cid).Scan(&music.ID, &music.Title, &music.CID, &music.OwnerAddr, &music.Signature, &music.UploadedAt)
	if err != nil {
		return nil, err
	}
	return &music, nil
}

func (p *Postgres) GetAllMusic(ctx context.Context) ([]model.Music, error) {
	rows, err := p.conn.Query(ctx, "SELECT * FROM music")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var music []model.Music
	for rows.Next() {
		var m model.Music
		err := rows.Scan(&m.ID, &m.Title, &m.CID, &m.OwnerAddr, &m.Signature, &m.UploadedAt)
		if err != nil {
			return nil, err
		}
		m.Link = fmt.Sprintf("http://localhost:8080/ipfs/%s", m.CID)
		music = append(music, m)
	}
	return music, nil
}

func (p *Postgres) CreateMusic(ctx context.Context, music model.Music) (int, error) {
	var id int
	err := p.conn.QueryRow(ctx, "INSERT INTO music (title, cid, owner_addr, signature, uploaded_at) VALUES ($1, $2, $3, $4, $5) RETURNING music_id", music.Title, music.CID, music.OwnerAddr, music.Signature, music.UploadedAt).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (p *Postgres) UpdateMusic(ctx context.Context, music model.Music) error {
	_, err := p.conn.Exec(ctx, "UPDATE music SET title = $1, cid = $2, owner_addr = $3, signature = $4, uploaded_at = $5 WHERE music_id = $6", music.Title, music.CID, music.OwnerAddr, music.Signature, music.UploadedAt, music.ID)
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) DeleteMusic(ctx context.Context, id int) error {
	_, err := p.conn.Exec(ctx, "DELETE FROM music WHERE music_id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
