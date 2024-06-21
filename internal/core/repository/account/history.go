package account

import (
	"context"

	"github.com/getnimbus/anton/abi"
	"github.com/getnimbus/anton/abi/known"
	"github.com/getnimbus/anton/internal/core/aggregate/history"
)

func getContractTypes(types []abi.ContractName) (ret []abi.ContractName) {
	for _, t := range types {
		ret = append(ret, t)
		if t == "wallet" {
			ret = append(ret, known.GetAllWalletNames()...)
		}
	}
	return
}

func (r *Repository) AggregateAccountsHistory(ctx context.Context, req *history.AccountsReq) (*history.AccountsRes, error) {
	//var res history.AccountsRes
	//
	//q := r.ch.NewSelect().Model((*core.AccountState)(nil))
	//
	//if len(req.ContractTypes) > 0 {
	//	q = q.Where("hasAny(types, ?)", ch.Array(getContractTypes(req.ContractTypes)))
	//}
	//if req.MinterAddress != nil {
	//	q = q.Where("minter_address = ?", req.MinterAddress)
	//}
	//
	//switch req.Metric {
	//case history.ActiveAddresses:
	//	q = q.ColumnExpr("uniqExact(address) as value")
	//default:
	//	return nil, errors.Wrapf(core.ErrInvalidArg, "invalid account metric %s", req.Metric)
	//}
	//
	//rounding, err := history.GetRoundingFunction(req.Interval)
	//if err != nil {
	//	return nil, err
	//}
	//q = q.ColumnExpr(fmt.Sprintf(rounding, "updated_at") + " as timestamp")
	//q = q.Group("timestamp")
	//
	//if !req.From.IsZero() {
	//	q = q.Where("updated_at > ?", req.From)
	//}
	//if !req.To.IsZero() {
	//	q = q.Where("updated_at < ?", req.To)
	//}
	//
	//q = q.Order("timestamp ASC")
	//
	//if err := q.Scan(ctx, &res.CountRes); err != nil {
	//	return nil, err
	//}
	//
	//return &res, nil
	return nil, nil
}
