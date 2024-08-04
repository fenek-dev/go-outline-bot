package pg

import (
	"context"
	"github.com/fenek-dev/go-outline-bot/internal/models"
	"github.com/georgysavva/scany/v2/pgxscan"
)

func (p *Postgres) GetServer(ctx context.Context, serverID uint64) (server models.Server, err error) {
	//TODO implement me
	return server, err
}
func (p *Postgres) GetAllServers(ctx context.Context) (servers []*models.Server, err error) {

	rows, err := p.conn.Query(ctx, "SELECT * FROM servers")
	if err != nil {
		return nil, err
	}

	err = pgxscan.ScanAll(&servers, rows)
	if err != nil {
		return nil, err
	}

	return servers, err
}
