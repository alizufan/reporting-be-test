package merchant

import (
	"fmt"
	"net/http"
	"reporting/libs/util"
	"reporting/service/transaction"
	"strconv"
)

type TransactionHandler struct {
	TransactionSrv transaction.TransactionService
}

func (m *TransactionHandler) Report(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		q   = r.URL.Query()
		p   = util.GetJWTPayload(ctx)

		outletId, _ = strconv.ParseUint(q.Get("outlet_id"), 10, 64)
		limit, _    = strconv.Atoi(q.Get("limit"))
		page, _     = strconv.Atoi(q.Get("page"))

		req = transaction.TrxRequest{
			MerchantID: p.MerchantID,
			OutletID:   outletId,
			Date:       q.Get("date"),
			OutletName: q.Get("outlet_name"),
			Limit:      limit,
			Page:       page,
		}
	)
	fmt.Println(p.MerchantID)

	res, err := m.TransactionSrv.Report(ctx, req)
	if err != nil {
		util.ErrorHTTPResponse(ctx, rw, err)
		return
	}

	util.HTTPResponse(rw, http.StatusOK, "success get trx report", res, nil)
}
