package core

import (
	"context"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/extra/bunbig"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/tvm/cell"

	"github.com/getnimbus/anton/addr"
)

type Transaction struct {
	//ch.CHModel    `ch:"transactions,partition:toYYYYMM(created_at)" json:"-"`
	bun.BaseModel `bun:"table:transactions" json:"-"`

	Address addr.Address `ch:"type:String,pk" bun:"type:bytea" json:"address"`
	Hash    []byte       `bun:"type:bytea,pk,notnull" json:"hash"`
	HashHex string       `ch:"-" bun:"-" json:"hash_hex"`

	CreatedLT uint64        `ch:",pk" bun:",notnull" json:"created_lt"`
	Account   *AccountState `ch:"-" bun:"rel:has-one,join:address=address,join:created_lt=last_tx_lt" json:"account"`

	Workchain  int32  `bun:"type:integer,notnull" json:"workchain"`
	Shard      int64  `bun:"type:bigint,notnull" json:"shard"`
	BlockSeqNo uint32 `bun:"type:integer,notnull" json:"block_seq_no"`

	PrevTxHash    []byte `bun:"type:bytea" json:"prev_tx_hash,omitempty"`
	PrevTxHashHex string `ch:"-" bun:"-" json:"prev_tx_hash_hex,omitempty"`
	PrevTxLT      uint64 `json:"prev_tx_lt,omitempty"`

	InMsgHash    []byte      `json:"in_msg_hash"`
	InMsgHashHex string      `ch:"-" bun:"-" json:"in_msg_hash_hex"`
	InMsg        *Message    `ch:"-" bun:"rel:belongs-to,join:in_msg_hash=hash" json:"in_msg"`
	InAmount     *bunbig.Int `ch:"type:UInt256" bun:"type:numeric,notnull" json:"in_amount,omitempty"`

	OutMsg      []*Message  `ch:"-" bun:"rel:has-many,join:address=src_address,join:created_lt=src_tx_lt" json:"out_msg,omitempty"`
	OutMsgCount uint16      `bun:",notnull" json:"out_msg_count"`
	OutAmount   *bunbig.Int `ch:"type:UInt256" bun:"type:numeric,notnull" json:"out_amount,omitempty"`

	TotalFees *bunbig.Int `ch:"type:UInt256" bun:"type:numeric" json:"total_fees"`

	Description           []byte `bun:"type:bytea,notnull" json:"description_boc,omitempty"`
	DescriptionLoaded     any    `ch:"-" bun:"-" json:"description,omitempty"`
	ComputePhaseExitCode  int32  `ch:"type:Int32" bun:",notnull" json:"compute_phase_exit_code"`
	ActionPhaseResultCode int32  `ch:"type:Int32" bun:",notnull" json:"action_phase_result_code"`

	OrigStatus AccountStatus `ch:",lc" bun:"type:account_status,notnull" json:"orig_status"`
	EndStatus  AccountStatus `ch:",lc" bun:"type:account_status,notnull" json:"end_status"`

	CreatedAt   time.Time `bun:"type:timestamp without time zone,notnull" json:"created_at"`
	DateKey     string    `ch:"-" bun:"-" json:"date_key"`
	TimestampMs string    `ch:"-" bun:"-" json:"timestamp_ms"`
}

func (tx *Transaction) LoadDescription() error { // TODO: optionally load description in API
	var d tlb.TransactionDescription

	c, err := cell.FromBOC(tx.Description)
	if err != nil {
		return errors.Wrap(err, "load description boc")
	}

	if err := tlb.LoadFromCell(&d, c.BeginParse()); err != nil {
		return errors.Wrap(err, "load description from cell")
	}

	tx.DescriptionLoaded = d

	return nil
}

func (tx *Transaction) PartitionKey() string {
	return hex.EncodeToString(tx.Hash)
}

func (tx *Transaction) WithDateKey() *Transaction {
	if tx.DateKey != "" {
		return tx
	}
	tx.DateKey = tx.CreatedAt.Format(time.DateOnly)
	tx.TimestampMs = strconv.FormatInt(tx.CreatedAt.UnixMilli(), 10)
	return tx
}

type TransactionRepository interface {
	AddTransactions(ctx context.Context, tx bun.Tx, transactions []*Transaction) error
}
