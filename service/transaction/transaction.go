package transaction

import (
	"context"
	"fmt"
	"math"
	"net/url"
	"reporting/repository/transaction"
	"reporting/schema"
	"time"
)

type TransactionService interface {
	Report(ctx context.Context, req TrxRequest) (*TrxResponse, error)
}

type Transaction struct {
	TransactionRepo transaction.TransactionRepository
}

type TrxRequest struct {
	MerchantID   uint64
	MerchantName string
	OutletID     uint64
	OutletName   string
	Date         string
	Limit        int
	Page         int
}

type Pagination struct {
	Limit     int `json:"limit"`
	Page      int `json:"page"`
	TotalPage int `json:"total_page"`
}
type Link struct {
	Current  string `json:"current"`
	NextPage string `json:"next"`
	PervPage string `json:"prev"`
}

type TrxResponse struct {
	Pagination Pagination                 `json:"pagination"`
	Link       Link                       `json:"link"`
	Data       []schema.TransactionReport `json:"data"`
}

func (t *Transaction) Report(ctx context.Context, req TrxRequest) (*TrxResponse, error) {
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	tDate, _ := time.Parse("2006-01", req.Date)
	if tDate.IsZero() {
		tDate = time.Now()
	}

	res, count, err := t.TransactionRepo.Report(ctx, transaction.ReportFilter{
		MerchantID: req.MerchantID,
		OutletID:   req.OutletID,
		Date:       tDate,
		Limit:      req.Limit,
		Page:       req.Page,
	})
	if err != nil {
		return nil, err
	}

	totalPage := math.Floor(float64(count) / float64(req.Limit))
	if totalPage <= 0 {
		totalPage = 1
	}

	uri := "http://localhost:3000/report"
	URL, _ := url.Parse(uri)
	q := URL.Query()
	if req.OutletID > 0 {
		q.Add("outlet_id", fmt.Sprint(req.OutletID))
	}
	q.Add("limit", fmt.Sprint(req.Limit))
	q.Add("page", fmt.Sprint(req.Page))
	URL.RawQuery = q.Encode()

	link := Link{
		Current: URL.String(),
	}

	if totalPage > 1 {
		URL, _ := url.Parse(uri)
		q := URL.Query()
		if req.OutletID > 0 {
			q.Add("outlet_id", fmt.Sprint(req.OutletID))
		}
		q.Add("limit", fmt.Sprint(req.Limit))
		q.Add("page", fmt.Sprint(req.Page+1))
		URL.RawQuery = q.Encode()
		link.NextPage = URL.String()
	}

	if req.Page > 1 && totalPage > 1 {
		URL, _ := url.Parse(uri)
		q := URL.Query()
		if req.OutletID > 0 {
			q.Add("outlet_id", fmt.Sprint(req.OutletID))
		}
		q.Add("limit", fmt.Sprint(req.Limit))
		q.Add("page", fmt.Sprint(req.Page-1))
		URL.RawQuery = q.Encode()
		link.PervPage = URL.String()
	}

	pagin := Pagination{
		Limit:     req.Limit,
		Page:      req.Page,
		TotalPage: int(totalPage),
	}

	return &TrxResponse{
		Pagination: pagin,
		Link:       link,
		Data:       res,
	}, nil
}
