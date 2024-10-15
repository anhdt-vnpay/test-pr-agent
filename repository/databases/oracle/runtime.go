package oracle

import (
	"context"
	"database/sql"

	"github.com/blcvn/corev3-libs/flogging"
	"github.com/blcvn/corev4-explorer/appconfig"
	tf "github.com/blcvn/corev4-explorer/transform"
	cmap "github.com/orcaman/concurrent-map/v2"
)

type OracleDB struct {
	Client      *sql.DB
	Schema      string
	StmtManager cmap.ConcurrentMap[string, *sql.Stmt]
	ctx         context.Context
	logger      *flogging.FabricLogger
	transform   transform
}

func NewOracleDB(client *sql.DB, schema string) *OracleDB {
	return &OracleDB{
		Client:      client,
		Schema:      schema,
		StmtManager: cmap.New[*sql.Stmt](),
		ctx:         context.Background(),
		logger:      flogging.MustGetLogger("oracle.repository.onchain_db"),
		transform:   tf.NewTransform(),
	}
}

func (odb *OracleDB) InitTables(ctx context.Context) error {

	networkChannelPairs := [][]string{}

	for _, channel := range appconfig.FabricCfg.Network.Channels {
		networkChannelPairs = append(networkChannelPairs, []string{channel.ShardId, channel.ID})
	}

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
