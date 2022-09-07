package merchant

import (
	"fmt"
	"net/http"
	"reporting/libs/util"
	"reporting/service/transaction"
	"strconv"
	"time"
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

func (m *TransactionHandler) Reporting(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		q   = r.URL.Query()
		p   = util.GetJWTPayload(ctx)

		outletId, _ = strconv.ParseUint(q.Get("outlet_id"), 10, 64)
		limit, _    = strconv.Atoi(q.Get("limit"))
		page, _     = strconv.Atoi(q.Get("page"))

		tz, _        = time.LoadLocation("Asia/Jakarta")
		startDate, _ = time.ParseInLocation("2006-01-02", q.Get("start_date"), tz)
		endDate, _   = time.ParseInLocation("2006-01-02", q.Get("end_date"), tz)

		req = transaction.TrxRequest{
			StartDate:  startDate,
			EndDate:    time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, endDate.Location()),
			MerchantID: p.MerchantID,
			OutletID:   outletId,
			Limit:      limit,
			Page:       page,
		}
	)

	res, err := m.TransactionSrv.Reporting(ctx, &req)
	if err != nil {
		util.ErrorHTTPResponse(ctx, rw, err)
		return
	}

	util.HTTPResponse(rw, http.StatusOK, "success get trx report", res, nil)
}
