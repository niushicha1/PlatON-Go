package core

import (
	"math/big"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/cbfttypes"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/event"
)

func TestBlockChainReactor_Close(t *testing.T) {
	//log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(6), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))
	//t.Run("close before commit", func(t *testing.T) {
	//	eventmux := new(event.TypeMux)
	//	reacter := NewBlockChainReactor(nil, eventmux)
	//	reacter.bftResultSub = eventmux.Subscribe(cbfttypes.CbftResult{})
	//	go func() { reacter.loop() }()
	//	var parenthash common.Hash
	//	cbftress := make(chan cbfttypes.CbftResult, 5)
	//	go func() {
	//		for i := 1; i < 11; i++ {
	//			header := new(types.Header)
	//			header.Number = big.NewInt(int64(i))
	//			header.Time = big.NewInt(int64(i))
	//			header.ParentHash = parenthash
	//			block := types.NewBlock(header, nil, nil)
	//			snapshotdb.Instance().NewBlock(header.Number, header.ParentHash, block.Hash())
	//			parenthash = block.Hash()
	//			cbftress <- cbfttypes.CbftResult{Block: block}
	//		}
	//		close(cbftress)
	//	}()
	//
	//	for value := range cbftress {
	//		eventmux.Post(value)
	//		if value.Block.Number().Cmp(big.NewInt(8)) == 0 {
	//			snapshotdb.Instance().Close()
	//		}
	//		if value.Block.Number().Cmp(big.NewInt(9)) == 0 {
	//			close(reacter.exitCh)
	//		}
	//	}
	//	time.Sleep(time.Second)
	//	snapshotdb.Instance().Clear()
	//})
	t.Run("close after commit", func(t *testing.T) {
		eventmux := new(event.TypeMux)
		reacter := NewBlockChainReactor(nil, eventmux)
		reacter.bftResultSub = eventmux.Subscribe(cbfttypes.CbftResult{})
		go func() { reacter.loop() }()
		var parenthash common.Hash
		cbftress := make(chan cbfttypes.CbftResult, 5)
		go func() {
			for i := 1; i < 11; i++ {
				header := new(types.Header)
				header.Number = big.NewInt(int64(i))
				header.Time = big.NewInt(int64(i))
				header.ParentHash = parenthash
				block := types.NewBlock(header, nil, nil)
				snapshotdb.Instance().NewBlock(header.Number, header.ParentHash, block.Hash())
				parenthash = block.Hash()
				cbftress <- cbfttypes.CbftResult{Block: block}
			}
			close(cbftress)
		}()

		for value := range cbftress {
			eventmux.Post(value)
		}

		close(reacter.exitCh)
		time.Sleep(time.Second)
		snapshotdb.Instance().Clear()
	})
}
