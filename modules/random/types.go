package random

const (
	ModuleName = "random"
)

func (m Random) Convert() interface{} {
	return QueryRandomResp{
		RequestTxHash: m.RequestTxHash.String(),
		Height:        m.Height,
		Value:         m.Value,
	}
}

type Requests []Request

func (m Requests) Convert() interface{} {
	var res []QueryRandomRequestQueueResp

	for _, request := range m {
		q := QueryRandomRequestQueueResp{
			Height:           request.Height,
			Consumer:         request.Consumer.String(),
			TxHash:           request.TxHash.String(),
			Oracle:           request.Oracle,
			ServiceFeeCap:    request.ServiceFeeCap,
			ServiceContextID: request.ServiceContextID.String(),
		}
		res = append(res, q)
	}
	return res
}
