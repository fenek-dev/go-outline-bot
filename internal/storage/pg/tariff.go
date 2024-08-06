package pg

import (
	"context"

	"github.com/fenek-dev/go-outline-bot/internal/models"
	"github.com/georgysavva/scany/v2/pgxscan"
)

func (p *Postgres) GetTariff(ctx context.Context, tariffID uint64) (tariff models.Tariff, err error) {

	err = pgxscan.Get(ctx, p.conn, &tariff, "SELECT * FROM tariffs WHERE id = $1", tariffID)
	if err != nil {
		return tariff, err
	}

	return tariff, err
}

func (p *Postgres) GetTariffsByServer(ctx context.Context, serverId uint64) (tariffs []models.Tariff, err error) {
	rows, err := p.conn.Query(ctx, "SELECT * FROM tariffs WHERE server_id = $1", serverId)
	if err != nil {
		return nil, err
	}

	err = pgxscan.ScanAll(&tariffs, rows)
	if err != nil {
		return nil, err
	}

	return tariffs, err
}
