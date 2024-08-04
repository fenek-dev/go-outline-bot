package pg

import "context"

func (p *Postgres) GetBalance(ctx context.Context, userId uint64) (balance uint32, err error) {
	err = p.conn.QueryRow(ctx, "SELECT balance FROM users WHERE id = $1", userId).Scan(&balance)

	if err != nil {
		return 0, err
	}

	return balance, nil
}
