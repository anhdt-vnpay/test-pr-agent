package hyperledger

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/blcvn/corev4-explorer/appconfig"
	"github.com/blcvn/corev4-explorer/entities"
	"github.com/blcvn/corev4-explorer/helper/statehelper"
	"github.com/blcvn/corev4-explorer/transform"

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

func NewLedgerRepository(fCfg *config.FabricConfig, dbRepo dbRepo) iLedgerRepo {

	listenerOption := listener.ListenerOption{
		RootFolder:                 appconfig.KvDBRootPath,
		MaxQueryGoroutine:          appconfig.BlockMainQueueSize,
		BlockMainQueueSize:         int64(appconfig.BlockMainQueueSize),
		BlockExceptionQueueSize:    int64(appconfig.BlockExceptionQueueSize),
		NumberBlockParsingWorker:   appconfig.NumberBlockParsingWorker,
		BlockLoadExceptionInterval: appconfig.BlockLoadExceptionInterval,
		BlockMinExceptionLoadTime:  appconfig.BlockMinExceptionLoadTime,
		UpdateMissing:              appconfig.UpdateMissing,
		BlockSyncTime:              appconfig.BlockSyncTime,
	}
	st, err := statehelper.NewStateHelper()
	if err != nil {
		panic(fmt.Sprintf("Can't create state helper: %s", err.Error()))
	}

	tf := transform.NewTransform()

	lr := &ledgerRepo{
		listenerOption: listenerOption,
		shardingMap:    cmap.New[ledgerClient](),
		logger:         flogging.MustGetLogger("corev4-explorer.repository.ledgers.hyperledger"),
		dbRepo:         dbRepo,
	}

	flowClient := flow.NewFlow(
		flow.WithFlowClientOption(client.NewClient(
			client.WithLoggerOption(lr.logger),
		)),
		flow.WithLoggerOption(lr.logger),
		flow.WithPeerGRPCTimeoutOption(appconfig.PeerGRPCTimeout),
	)

	for _, channelConfig := range fCfg.Network.Channels {
		lc := listener.NewLedgerClient[*entities.RawTransaction](channelConfig,
			listener.WithClientFlowClientOption[*entities.RawTransaction](flowClient),
			listener.WithClientLoggerOption[*entities.RawTransaction](lr.logger),
			listener.WithClientStateHelperOption[*entities.RawTransaction](st),
			listener.WithClientListenerOption[*entities.RawTransaction](listenerOption),
			listener.WithClientTransformerOption[*entities.RawTransaction](tf),
			listener.WithClientKvDbTypeOption[*entities.RawTransaction](appconfig.KvDBType),
		)
		lc.SetupProcessBlockJob(func(key string, blockEvent *listener.CommonBlockWrapper) (*listener.CommonBlockWrapper, error) {
			return lr.processBlockWithConfig(key, blockEvent, channelConfig.ID, channelConfig.ShardId)
		})
		lr.shardingMap.Set(channelConfig.ShardId, lc)
	}
	return lr
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
			Txhash:                 rawTx.TxId,
			ChaincodeName:          rawTx.Chaincode,
			ValidationCode:         rawTx.ValidationCode.String(),
			BlockNum:               blockNumber,
			NetworkName:            shardId,
			ChannelName:            channelName,
			BlockTime:              blockEvent.ReceiveTime,
			ChaincodeProposalInput: rawTx.FuncName,
			Payload:                rawTx.Data,
		})
	}
	err = l.dbRepo.SaveBlockAndRawTxs(blockInformation, rawTransactions)
	if err != nil {
		return nil, err
	}
	return blockEvent, nil
}
