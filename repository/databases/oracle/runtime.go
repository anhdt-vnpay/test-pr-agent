package oracle

import (
	"context"
	"database/sql"

	"github.com/blcvn/corev3-libs/flogging"
	cmap "github.com/orcaman/concurrent-map/v2"
)

type OracleDB struct {
	Client      *sql.DB
	Schema      string
	StmtManager cmap.ConcurrentMap[string, *sql.Stmt]
	logger      *flogging.FabricLogger
}

func NewOracleDB(client *sql.DB, schema string) *OracleDB {
	return &OracleDB{
		Client:      client,
		Schema:      schema,
		StmtManager: cmap.New[*sql.Stmt](),
		logger:      flogging.MustGetLogger("oracle.repository.onchain_db"),
	}
}

func (odb *OracleDB) InitTables(ctx context.Context) error {

	networkChannelPairs := [][]string{}

	odb.InitTableAccountDelta(ctx)
	odb.InitTableAccountBalance(ctx, networkChannelPairs)
	odb.InitTableAccountTransactions(ctx)
	odb.InitTableOnchainTransaction(ctx)
	odb.InitTableTask(ctx)

	odb.InitTableBlock(ctx, networkChannelPairs)
	odb.InitTableRawTransaction(ctx, networkChannelPairs)
	odb.InitTableAccountInit(ctx)

	return nil

}
