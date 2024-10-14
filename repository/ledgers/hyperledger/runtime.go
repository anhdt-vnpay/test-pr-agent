package hyperledger

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/blcvn/corev4-explorer/entities"

	"github.com/blcvn/corev3-libs/flogging"
	"github.com/blcvn/corev3-libs/ledger/hyperledger/config"
	"github.com/blcvn/corev3-libs/ledger/hyperledger/flow"
	"github.com/blcvn/corev3-libs/ledger/hyperledger/flow/client"
	"github.com/blcvn/corev3-libs/ledger/hyperledger/helper"
	"github.com/blcvn/corev3-libs/ledger/hyperledger/listener"
	"github.com/hyperledger/fabric-protos-go/common"
	cmap "github.com/orcaman/concurrent-map/v2"
)

type ledgerRepo struct {
	listenerOption listener.ListenerOption
	shardingMap    cmap.ConcurrentMap[string, ledgerClient]
	logger         *flogging.FabricLogger
	dbRepo         dbRepo
	// processBlock   func(key string, blockEvent *listener.CommonBlockWrapper) (*listener.CommonBlockWrapper, error)
}

func NewLedgerRepository(fCfg config.FabricConfig) iLedgerRepo {
	lr := &ledgerRepo{}
	lr.logger = flogging.MustGetLogger("corev4-explorer.repository.ledgers.hyperledger")

	listenerOption := listener.ListenerOption{
		RootFolder:                 "./tmp",
		MaxQueryGoroutine:          1,
		BlockMainQueueSize:         100,
		BlockExceptionQueueSize:    100,
		NumberBlockParsingWorker:   100,
		BlockLoadExceptionInterval: 500 * time.Second,
		BlockMinExceptionLoadTime:  500 * time.Second,
		UpdateMissing:              true,
		BlockSyncTime:              10 * time.Second,
	}
	// lr.listenerOption = listenerOption

	flowClient := flow.NewFlow(
		flow.WithFlowClientOption(client.NewClient(
			client.WithLoggerOption(lr.logger),
		)),
		flow.WithLoggerOption(lr.logger),
		flow.WithPeerGRPCTimeoutOption(3*time.Second),
	)

	st := &stateHelper{}

	for _, channelConfig := range fCfg.Network.Channels {
		lc := listener.NewLedgerClient[*Tx](channelConfig,
			listener.WithClientFlowClientOption[*Tx](flowClient),
			listener.WithClientLoggerOption[*Tx](lr.logger),
			listener.WithClientStateHelperOption[*Tx](st),
			listener.WithClientListenerOption[*Tx](listenerOption),
			listener.WithClientTransformerOption[*Tx](&transform{}),
		)
		lc.SetupProcessBlockJob(func(key string, blockEvent *listener.CommonBlockWrapper) (*listener.CommonBlockWrapper, error) {
			lr.processBlockWithConfig(key, blockEvent, channelConfig.ID, channelConfig.ShardId)
			return blockEvent, nil
		})
		lr.shardingMap.Set(channelConfig.ShardId, lc)
	}
	return &ledgerRepo{listenerOption: listenerOption}
}

func (l *ledgerRepo) Start(ctx context.Context) {
	for _, lc := range l.shardingMap.Items() {
		lc.Start(ctx)
	}
}

func (l *ledgerRepo) processBlockWithConfig(key string, blockEvent *listener.CommonBlockWrapper, channelName, shardId string) (*listener.CommonBlockWrapper, error) {
	blockHeader := blockEvent.Block.Header
	blockNumber := blockHeader.Number
	blockData := blockEvent.Block.Data.GetData()
	blockBytes, err := blockEvent.MarshalBytes()
	if err != nil {
		l.logger.Errorf("Failed to marshal block number %d reason %s", blockNumber, err.Error())
		return nil, err
	}
	blockSize := len(blockBytes)
	txCount := len(blockData)
	/**
	 * Read block from block event
	 */
	validationCodes := blockEvent.Block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER]

	blockInformation := &entities.Block{
		Blocknum:      blockHeader.Number,
		Datahash:      hex.EncodeToString(blockHeader.DataHash),
		Prehash:       hex.EncodeToString(blockHeader.PreviousHash),
		Txcount:       int64(txCount),
		BlockTime:     blockEvent.ReceiveTime,
		Blockhash:     hex.EncodeToString(blockHeader.PreviousHash),
		ChannelName:   channelName,
		Blksize:       int64(blockSize),
		PrevBlockhash: blockHeader.String(),
		NetworkName:   shardId,
	}

	/**
	 * Read raw transaction from block data
	 */

	var rawTransactions []*entities.RawTransaction
	for txIndex, ebytes := range blockData {
		rawTx, err := helper.ParseTransation(ebytes, txIndex, blockNumber, validationCodes)
		if err != nil {
			l.logger.Errorf("Shard ID %s Channel %s TxIndex %d", shardId, channelName, txIndex)
			return nil, err
		}
		rawTransactions = append(rawTransactions, &entities.RawTransaction{
			Txhash:         rawTx.TxId,
			ChaincodeName:  rawTx.Chaincode,
			ValidationCode: string(validationCodes),
			BlockNum:       blockNumber,
			NetworkName:    shardId,
			ChannelName:    channelName,
			BlockTime:      blockEvent.ReceiveTime,
			TxData:         string(rawTx.Data),
			FunctionName:   rawTx.FuncName,
		})
	}
	err = l.dbRepo.SaveBlockAndRawTxs(blockInformation, rawTransactions)
	if err != nil {
		return nil, err
	}
	return blockEvent, nil
}
