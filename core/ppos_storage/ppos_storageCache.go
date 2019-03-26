package ppos_storage

import (
	"errors"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"math/big"
	"sync"

	//"sync"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"fmt"
)

const (
	PREVIOUS = iota
	CURRENT
	NEXT
	IMMEDIATE
	RESERVE
	REFUND
)

var (
	ParamsIllegalErr 			= errors.New("Params illegal")
	TicketNotFindErr        	= errors.New("The Ticket not find")
)

var ticketCache = sync.Map{}

func PutTicket(txHash common.Hash, ticket *types.Ticket) {
	ticketCache.Store(txHash, ticket)
}

func GetTicket(txHash common.Hash) *types.Ticket {
	if value, ok := ticketCache.Load(txHash); ok {
		return value.(*types.Ticket)
	}
	return nil
}

func RemoveTicket(txHash common.Hash) {
	ticketCache.Delete(txHash)
}

type refundStorage map[discover.NodeID]types.RefundQueue

type candidate_temp struct {
	// previous witness
	pres types.CandidateQueue
	// current witness
	currs types.CandidateQueue
	// next witness
	nexts types.CandidateQueue
	//immediate
	imms types.CandidateQueue
	// reserve
	res types.CandidateQueue
	// refund
	refunds refundStorage
}

type ticketDependency struct {
	// ticket age
	//Age uint64
	// ticket count
	Num uint32
	// ticketIds
	Tinfo []*ticketInfo
}

type ticketInfo struct {
	TxHash			common.Hash
	// The number of remaining tickets
	Remaining		uint32
	Price 			*big.Int
}

func (t *ticketInfo) SubRemaining() {
	if t.Remaining > 0 {
		t.Remaining--
	}
}

//func (td *ticketDependency) AddAge(number uint64) {
//	if number > 0 {
//		td.Age += number
//	}
//}
//
//func (td *ticketDependency) SubAge(number uint64) {
//	if number > 0 && td.Age >= number {
//		td.Age -= number
//	}
//}

func (td *ticketDependency) subNum() {
	if td.Num > 0 {
		td.Num--
	}
}

type ticket_temp struct {
	// total remian  k-v
	Sq int32
	// ticketInfo  map[txHash]ticketInfo
	//Infos map[common.Hash]*types.Ticket
	// ExpireTicket  map[blockNumber]txHash
	//Ets map[string][]common.Hash
	// ticket's attachment  of node
	Dependencys map[discover.NodeID]*ticketDependency
}

type Ppos_storage struct {
	c_storage *candidate_temp
	t_storage *ticket_temp
}

func (ps *Ppos_storage) Copy() *Ppos_storage {

	//if  verifyStorageEmpty(ps) {
	//	return NewPPOS_storage()
	//}

	if nil == ps {
		return NewPPOS_storage()
	}

	ppos_storage := &Ppos_storage{
		c_storage: ps.CopyCandidateStorage(),
		t_storage: ps.CopyTicketStorage(),
	}
	return ppos_storage
}


func NewPPOS_storage () *Ppos_storage {

	cache := &Ppos_storage{
		c_storage: &candidate_temp{
			pres: 	make(types.CandidateQueue, 0),
			currs:  make(types.CandidateQueue, 0),
			nexts: 	make(types.CandidateQueue, 0),
			imms: 	make(types.CandidateQueue, 0),
			res: 	make(types.CandidateQueue, 0),
			refunds: make(refundStorage, 0),
		},

		t_storage: &ticket_temp{
			Sq: 	-1,
			Dependencys: 	make(map[discover.NodeID]*ticketDependency),
		},
	}

	/*cache := new(Ppos_storage)

	c := new(candidate_temp)
	t:= new(ticket_temp)

	c.pres = make(types.CandidateQueue, 0)
	c.currs = make(types.CandidateQueue, 0)
	c.nexts = make(types.CandidateQueue, 0)

	c.imms = make(types.CandidateQueue, 0)
	c.res = make(types.CandidateQueue, 0)

	c.refunds = make(refundStorage, 0)

	t.Sq = -1
	t.Dependencys = make(map[discover.NodeID]*ticketDependency)

	cache.c_storage = c
	cache.t_storage = t*/


	return cache
}

/** candidate related func */

//func (p *Ppos_storage) CopyCandidateStorage ()  *candidate_temp {
//	start := common.NewTimer()
//	start.Begin()
//
//	temp := new(candidate_temp)
//
//	type result struct {
//		Status int
//		Data   interface{}
//	}
//	var wg sync.WaitGroup
//	wg.Add(6)
//	resCh := make(chan *result, 6)
//
//	loadQueueFunc := func(flag int) {
//		res := new(result)
//		switch flag {
//		case PREVIOUS:
//			res.Status = PREVIOUS
//			res.Data = p.c_storage.pres.DeepCopy()
//			resCh <- res
//		case CURRENT:
//			res.Status = CURRENT
//			res.Data = p.c_storage.currs.DeepCopy()
//			resCh <- res
//		case NEXT:
//			res.Status = NEXT
//			res.Data = p.c_storage.nexts.DeepCopy()
//			resCh <- res
//		case IMMEDIATE:
//			res.Status = IMMEDIATE
//			res.Data = p.c_storage.imms.DeepCopy()
//			resCh <- res
//		case RESERVE:
//			res.Status = RESERVE
//			res.Data = p.c_storage.res.DeepCopy()
//			resCh <- res
//		}
//		wg.Done()
//	}
//
//	go loadQueueFunc(PREVIOUS)
//	go loadQueueFunc(CURRENT)
//	go loadQueueFunc(NEXT)
//	go loadQueueFunc(IMMEDIATE)
//	go loadQueueFunc(RESERVE)
//
//	go func() {
//		res := new(result)
//		cache := make(refundStorage, len(p.c_storage.refunds))
//		for nodeId, queue := range p.c_storage.refunds {
//			cache[nodeId] = queue.DeepCopy()
//		}
//		res.Status = REFUND
//		res.Data = cache
//		resCh <- res
//		wg.Done()
//	}()
//	wg.Wait()
//	close(resCh)
//	for res := range resCh {
//		switch res.Status {
//		case PREVIOUS:
//			temp.pres = res.Data.(types.CandidateQueue)
//		case CURRENT:
//			temp.currs = res.Data.(types.CandidateQueue)
//		case NEXT:
//			temp.nexts = res.Data.(types.CandidateQueue)
//		case IMMEDIATE:
//			temp.imms = res.Data.(types.CandidateQueue)
//		case RESERVE:
//			temp.res = res.Data.(types.CandidateQueue)
//		case REFUND:
//			temp.refunds = res.Data.(refundStorage)
//		}
//	}
//	log.Debug("CopyCandidateStorage", "Time spent", fmt.Sprintf("%v ms", start.End()))
//	return temp
//}

func (p *Ppos_storage) CopyCandidateStorage ()  *candidate_temp {
	start := common.NewTimer()
	start.Begin()

	cache := make(refundStorage, len(p.c_storage.refunds))
	for nodeId, queue := range p.c_storage.refunds {
		cache[nodeId] = queue.DeepCopy()
	}

	temp := &candidate_temp{
		pres: 		p.c_storage.pres.DeepCopy(),
		currs: 		p.c_storage.currs.DeepCopy(),
		nexts: 		p.c_storage.nexts.DeepCopy(),

		imms: 		p.c_storage.imms.DeepCopy(),
		res: 		p.c_storage.res.DeepCopy(),

		refunds: 	cache,
	}

	/*temp := new(candidate_temp)

	temp.pres = p.c_storage.pres.DeepCopy()
	temp.currs = p.c_storage.currs.DeepCopy()
	temp.nexts = p.c_storage.nexts.DeepCopy()

	temp.imms = p.c_storage.imms.DeepCopy()
	temp.res = p.c_storage.res.DeepCopy()

	temp.refunds = cache*/



	//PrintObject("CopyCandidateStorage前: pres", p.c_storage.pres)
	//PrintObject("CopyCandidateStorage前: currs", p.c_storage.currs)
	//PrintObject("CopyCandidateStorage前: nexts", p.c_storage.nexts)
	//PrintObject("CopyCandidateStorage前: imms", p.c_storage.imms)
	//PrintObject("CopyCandidateStorage前: res", p.c_storage.res)
	//PrintObject("CopyCandidateStorage前: refunds", p.c_storage.refunds)





	//PrintObject("CopyCandidateStorage后: pres", temp.pres)
	//PrintObject("CopyCandidateStorage后: currs", temp.currs)
	//PrintObject("CopyCandidateStorage后: nexts", temp.nexts)
	//PrintObject("CopyCandidateStorage后: imms", temp.imms)
	//PrintObject("CopyCandidateStorage后: res", temp.res)
	//PrintObject("CopyCandidateStorage后: refunds", temp.refunds)


	log.Debug("CopyCandidateStorage", "Time spent", fmt.Sprintf("%v ms", start.End()))
	return temp
}

func (p *Ppos_storage) CopyTicketStorage() *ticket_temp {

	start := common.NewTimer()
	start.Begin()

	cache := make(map[discover.NodeID]*ticketDependency, len(p.t_storage.Dependencys))

	for key := range p.t_storage.Dependencys {
		temp := p.t_storage.Dependencys[key]
		tinfos := make([]*ticketInfo, len(temp.Tinfo))

		for j, tin := range temp.Tinfo {

			t := &ticketInfo{
				TxHash: 	tin.TxHash,
				Remaining: 	tin.Remaining,
				Price: 		tin.Price,
			}
			tinfos[j] = t

		}

		cache[key] = &ticketDependency{
			temp.Num,
			tinfos,
		}
	}

	ticket_cache := &ticket_temp{
		Sq: 	p.t_storage.Sq,
		Dependencys: 	cache,
	}


	/*ticket_cache := new(ticket_temp)

	ticket_cache.Sq = p.t_storage.Sq
	//ticket_cache.Infos = make(map[common.Hash]*types.Ticket)
	//ticket_cache.Ets = make(map[string][]common.Hash)
	ticket_cache.Dependencys = make(map[discover.NodeID]*ticketDependency)

	//for key := range p.t_storage.Infos {
	//	ticket := p.t_storage.Infos[key]
	//	ticket_cache.Infos[key] = ticket.DeepCopy()
	//}
	//for key := range p.t_storage.Ets {
	//	value := p.t_storage.Ets[key]
	//	list := make([]common.Hash, len(value))
	//	copy(list, value)
	//	ticket_cache.Ets[key] = list
	//}
	//for key := range p.t_storage.Dependencys {
	//	temp := p.t_storage.Dependencys[key]
	//	tids := make([]common.Hash, len(temp.Tids))
	//	copy(tids, temp.Tids)
	//	ticket_cache.Dependencys[key] = &ticketDependency{
	//		//temp.Age,
	//		temp.Num,
	//		tids,
	//	}
	//}
	for key := range p.t_storage.Dependencys {
		temp := p.t_storage.Dependencys[key]

		tinfos := make([]*ticketInfo, len(temp.Tinfo))

		for j, tin := range temp.Tinfo {

			t := &ticketInfo{
				TxHash: 	tin.TxHash,
				Remaining: 	tin.Remaining,
				Price: 		tin.Price,
			}
			tinfos[j] = t

		}
		ticket_cache.Dependencys[key] = &ticketDependency{
			//temp.Age,
			temp.Num,
			tinfos,
		}
	}*/


	//PrintObject("CopyTicketStorage前:", p.t_storage)

	//PrintObject("CopyTicketStorage后:", ticket_cache)

	log.Debug("CopyTicketStorage", "Time spent", fmt.Sprintf("%v ms", start.End()))
	return ticket_cache
}

// Get CandidateQueue
// flag:
// 0: previous witness
// 1: current witness
// 2: next witness
// 3: immediate
// 4: reserve
func (p *Ppos_storage) GetCandidateQueue(flag int) types.CandidateQueue {
	switch flag {
	case PREVIOUS:
		return p.c_storage.pres
	case CURRENT:
		return p.c_storage.currs
	case NEXT:
		return p.c_storage.nexts
	case IMMEDIATE:
		return p.c_storage.imms
	case RESERVE:
		return p.c_storage.res
	default:
		return nil
	}
}

// Set CandidateQueue
func (p *Ppos_storage) SetCandidateQueue(queue types.CandidateQueue, flag int) {
	switch flag {
	case PREVIOUS:
		p.c_storage.pres = queue
	case CURRENT:
		p.c_storage.currs = queue
	case NEXT:
		p.c_storage.nexts = queue
	case IMMEDIATE:
		p.c_storage.imms = queue
	case RESERVE:
		p.c_storage.res = queue
	}
}

// Delete CandidateQueue
func (p *Ppos_storage) DelCandidateQueue(flag int)  {
	switch flag {
	case PREVIOUS:
		p.c_storage.pres = nil
	case CURRENT:
		p.c_storage.currs = nil
	case NEXT:
		p.c_storage.nexts = nil
	case IMMEDIATE:
		p.c_storage.imms = nil
	case RESERVE:
		p.c_storage.res = nil
	}
}

// Get Refund
func (p *Ppos_storage) GetRefunds(nodeId discover.NodeID) types.RefundQueue {
	if queue, ok := p.c_storage.refunds[nodeId]; ok {
		return queue
	} else {
		return make(types.RefundQueue, 0)
	}
}



// Set Refund
func (p *Ppos_storage) SetRefund(nodeId discover.NodeID, refund *types.CandidateRefund) {

	if queue, ok := p.c_storage.refunds[nodeId]; ok {
		queue = append(queue, refund)
		p.c_storage.refunds[nodeId] = queue
	} else {
		queue = make(types.RefundQueue, 1)
		queue[0] = refund
		p.c_storage.refunds[nodeId] = queue
	}
}

func (p *Ppos_storage) SetRefunds(nodeId discover.NodeID, refundArr types.RefundQueue) {
	p.c_storage.refunds[nodeId] = refundArr
}

func (p *Ppos_storage) AppendRefunds(nodeId discover.NodeID, refundArr types.RefundQueue) {
	if queue, ok := p.c_storage.refunds[nodeId]; ok {
		queue = append(queue, refundArr ...)
		p.c_storage.refunds[nodeId] = queue
	} else {
		p.c_storage.refunds[nodeId] = refundArr
	}
}

// Delete RefundArr
func (p *Ppos_storage) DelRefunds(nodeId discover.NodeID) {
	delete(p.c_storage.refunds, nodeId)
}

/** ticket related func */

// Get total remian
func (p *Ppos_storage) GetTotalRemian() int32 {
	return p.t_storage.Sq
}

// Set total remain
func (p *Ppos_storage) SetTotalRemain(count int32) {
	p.t_storage.Sq = count
}

// Get TicketInfo
/*func (p *Ppos_storage) GetTicketInfo(txHash common.Hash) *types.Ticket {
	ticket, ok := p.t_storage.Infos[txHash]
	if ok {
		return ticket
	}
	return nil
}*/

func (p *Ppos_storage) GetTicketInfo(txHash common.Hash) *ticketInfo {
	for _, obj := range p.t_storage.Dependencys {
		for index := range obj.Tinfo {
			 tinfo := obj.Tinfo[index]
			if tinfo.TxHash == txHash {
				return tinfo
			}
		}
	}
	return nil
}

//Set TicketInfo
//func (p *Ppos_storage) SetTicketInfo(txHash common.Hash, ) {
//	p.t_storage.Infos[txHash] = ticket
//}
//
//func (p *Ppos_storage) removeTicketInfo(txHash common.Hash) {
//	delete(p.t_storage.Infos, txHash)
//}

//GetTiketArr
//func (p *Ppos_storage) GetTicketArr(txHashs ...common.Hash) []*types.Ticket {
//	tickets := make([]*types.Ticket, 0)
//	if len(txHashs) > 0 {
//		for index := range txHashs {
//			if ticket := p.GetTicketInfo(txHashs[index]); ticket != nil {
//				newTicket := *ticket
//				tickets = append(tickets, &newTicket)
//			}
//		}
//	}
//	return tickets
//}

//Get ExpireTicket
//func (p *Ppos_storage) GetExpireTicket(blockNumber *big.Int) []common.Hash {
//	ids, ok := p.t_storage.Ets[blockNumber.String()]
//	if ok {
//		return ids
//	}
//	return nil
//}
//
//// Set ExpireTicket
//func (p *Ppos_storage) SetExpireTicket(blockNumber *big.Int, txHash common.Hash) {
//	ids, ok := p.t_storage.Ets[blockNumber.String()]
//	if !ok {
//		ids = make([]common.Hash, 0)
//	}
//	ids = append(ids, txHash)
//	p.t_storage.Ets[blockNumber.String()] = ids
//}
//
//func (p *Ppos_storage) RemoveExpireTicket(blockNumber *big.Int, txHash common.Hash) {
//	ids, ok := p.t_storage.Ets[blockNumber.String()]
//	if ok {
//		ids = removeTicketId(txHash, ids)
//		if len(ids) == 0 {
//			delete(p.t_storage.Ets, blockNumber.String())
//		} else {
//			p.t_storage.Ets[blockNumber.String()] = ids
//		}
//	}
//}

//Get ticket dependency
func (p *Ppos_storage) GetTicketDependency(nodeId discover.NodeID) *ticketDependency {
	value, ok := p.t_storage.Dependencys[nodeId]
	if ok {
		return value
	}
	return nil
}

// Set ticket dependency
func (p *Ppos_storage) SetTicketDependency(nodeId discover.NodeID, ependency *ticketDependency) {
	p.t_storage.Dependencys[nodeId] = ependency
}

func (p *Ppos_storage) RemoveTicketDependency(nodeId discover.NodeID) {
	delete(p.t_storage.Dependencys, nodeId)
}

/*func (p *Ppos_storage) GetCandidateTxHashs(nodeId discover.NodeID) []common.Hash {
	value, ok := p.t_storage.Dependencys[nodeId]
	if ok {
		return value.Tids
	}
	return nil
}*/

func (p *Ppos_storage) GetCandidateTxHashs(nodeId discover.NodeID) []common.Hash {
	value, ok := p.t_storage.Dependencys[nodeId]
	if ok {
		tids := make([]common.Hash, 0)
		for index := range value.Tinfo {
			tids = append(tids, value.Tinfo[index].TxHash)
		}
		return tids
	}
	return nil
}

/*func (p *Ppos_storage) AppendTicket(nodeId discover.NodeID, txHash common.Hash, ticket *types.Ticket) error {
	p.SetTicketInfo(txHash, ticket)
	value := p.GetTicketDependency(nodeId)
	if nil == value {
		value = new(ticketDependency)
		value.Tids = make([]common.Hash, 0)
	}
	value.Num += ticket.Remaining
	value.Tids = append(value.Tids, txHash)
	p.SetTicketDependency(nodeId, value)
	return nil
}*/

func (p *Ppos_storage) AppendTicket(nodeId discover.NodeID, txHash common.Hash, count uint32, price *big.Int) error {
	value := p.GetTicketDependency(nodeId)
	if nil == value {
		value = new(ticketDependency)
		value.Tinfo = make([]*ticketInfo, 0)
	}
	value.Num += count
	tinfo := &ticketInfo{
		txHash,
		count,
		price,
	}
	value.Tinfo = append(value.Tinfo, tinfo)
	p.SetTicketDependency(nodeId, value)
	return nil
}


/*func (p *Ppos_storage) SubTicket(nodeId discover.NodeID, txHash common.Hash) error {
	value := p.GetTicketDependency(nodeId)
	if nil != value {
		ticket := p.GetTicketInfo(txHash)
		if ticket == nil {
			return TicketNotFindErr
		}
		ticket.SubRemaining()
		value.subNum()
		if ticket.Remaining == 0 {
			p.removeTicketInfo(txHash)
			if list := removeTicketId(txHash, value.Tids); len(list) > 0 {
				value.Tids = list
			} else {
				value.Tids = make([]common.Hash, 0)
			}
		} else {
			p.SetTicketInfo(txHash, ticket)
		}
	}
	return nil
}*/

func (p *Ppos_storage) SubTicket(nodeId discover.NodeID, txHash common.Hash) (*ticketInfo, error) {
	value := p.GetTicketDependency(nodeId)
	if nil != value {
		ticket := p.GetTicketInfo(txHash)
		if ticket == nil {
			return nil, TicketNotFindErr
		}
		ticket.SubRemaining()
		value.subNum()
		if ticket.Remaining == 0 {
			if list := removeTinfo(txHash, value.Tinfo); len(list) > 0 {
				value.Tinfo = list
			} else {
				value.Tinfo = make([]*ticketInfo, 0)
			}
		}
		return ticket, nil
	}
	return nil, nil
}

/*func (p *Ppos_storage) RemoveTicket(nodeId discover.NodeID, txHash common.Hash) error {
	ticket := p.GetTicketInfo(txHash)
	if ticket == nil {
		return TicketNotFindErr
	}
	value := p.GetTicketDependency(nodeId)
	if nil != value {
		value.Num -= ticket.Remaining
		if list := removeTicketId(txHash, value.Tids); len(list) > 0 {
			value.Tids = list
		} else {
			value.Tids = make([]common.Hash, 0)
		}
	}
	p.removeTicketInfo(txHash)
	return nil
}*/

func (p *Ppos_storage) RemoveTicket(nodeId discover.NodeID, txHash common.Hash) (*ticketInfo, error) {
	ticket := p.GetTicketInfo(txHash)
	if ticket == nil {
		return nil, TicketNotFindErr
	}
	value := p.GetTicketDependency(nodeId)
	if nil != value {
		value.Num -= ticket.Remaining
		if list := removeTinfo(txHash, value.Tinfo); len(list) > 0 {
			value.Tinfo = list
		} else {
			value.Tinfo = make([]*ticketInfo, 0)
		}
		if value.Num == 0 {
			p.RemoveTicketDependency(nodeId)
		}
	}
	return ticket, nil
}

func (p *Ppos_storage) GetCandidateTicketCount(nodeId discover.NodeID) uint32 {
	if value := p.GetTicketDependency(nodeId); value != nil {
		log.Debug("获取当前node的得票数", "nodeId", nodeId.String(), "tcount", value.Num)
		return value.Num
	}
	log.Debug("获取当前node的得票数", "nodeId", nodeId.String(), "tcount", 0)
	return 0
}

func (p *Ppos_storage) GetCandidateTicketAge(nodeId discover.NodeID) uint64 {
	/*if value := p.GetTicketDependency(nodeId); value != nil {
		return value.Age
	}*/
	return 0
}

func (p *Ppos_storage) SetCandidateTicketAge(nodeId discover.NodeID, age uint64) {
	/*if value := p.GetTicketDependency(nodeId); value != nil {
		value.Age = age
	}*/
}

func (p *Ppos_storage) GetTicketRemainByTxHash(txHash common.Hash) uint32 {
	//PrintObject("Call GetTicketRemainByTxHash", p.t_storage.Dependencys)
	//log.Debug("Call GetTicketRemainByTxHash", "ticketId", txHash.Hex())
	for _, depen := range p.t_storage.Dependencys {
		for _, field := range depen.Tinfo {
			if txHash == field.TxHash {
				return field.Remaining
			}
		}
	}
	return 0
}


func removeTicketId(hash common.Hash, hashs []common.Hash) []common.Hash {
	for index, tempHash := range hashs {
		if tempHash == hash {
			if len(hashs) == 1 {
				return nil
			}
			start := hashs[:index]
			end := hashs[index+1:]
			return append(start, end...)
		}
	}
	return hashs
}

func removeTinfo(hash common.Hash, tinfos []*ticketInfo) []*ticketInfo {
	for index, info := range tinfos {
		if info.TxHash == hash {
			if len(tinfos) == 1 {
				return nil
			}
			start := tinfos[:index]
			end := tinfos[index+1:]
			return append(start, end...)
		}
	}
	return tinfos
}
