package transaction

import (
	"context"
	"fmt"
	"reporting/schema"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/go-rel/rel"
)

type TransactionRepository interface {
	Report(ctx context.Context, fil ReportFilter) ([]schema.TransactionReport, int, error)
	Reporting(ctx context.Context, fil ReportFilter) ([]schema.TransactionReport, error)
}

type Transaction struct {
	DB rel.Repository
}

type ReportFilter struct {
	MerchantID   uint64
	OutletID     uint64
	MerchantName string
	OutletName   string
	Date         time.Time
	StartDate    time.Time
	EndDate      time.Time
	Limit        int
	Page         int
}

func (f ReportFilter) IsValidMerchantID() bool {
	return f.MerchantID > 0
}

func (f ReportFilter) IsValidOutletID() bool {
	return f.OutletID > 0
}

func (f ReportFilter) IsValidMerchantName() bool {
	return len(f.MerchantName) > 0
}

func (f ReportFilter) IsValidOutletName() bool {
	return len(f.OutletName) > 0
}

func (f ReportFilter) IsValidRangeDate() bool {
	if f.StartDate.IsZero() || f.EndDate.IsZero() {
		return false
	}

	if f.StartDate.After(f.EndDate) {
		return false
	}

	return true
}

func (t *Transaction) Report(ctx context.Context, fil ReportFilter) ([]schema.TransactionReport, int, error) {
	if fil.Page <= 0 {
		fil.Page = 1
	}

	date := fmt.Sprintf("%s-01", fil.Date.Format("2006-01"))
	offset := (fil.Page - 1) * fil.Limit

	dialect := goqu.Dialect("mysql")
	qdl, qdlArgs, err := dialect.From(goqu.L(
		`(
			SELECT LAST_DAY(?) - INTERVAL (? + (10 * ?) + (100 * ?)) DAY AS Date
			FROM (SELECT 0 AS ? UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) AS ?
			CROSS JOIN (SELECT 0 AS ? UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) AS ?
			CROSS JOIN (SELECT 0 AS ? UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) AS ?
		)`,
		date, goqu.I("a.a"),
		goqu.I("b.a"), goqu.I("c.a"),
		goqu.I("a"), goqu.I("a"),
		goqu.I("a"), goqu.I("b"),
		goqu.I("a"), goqu.I("c")).As("a"),
	).Select(
		goqu.C("Date").Table("a"),
	).Where(
		goqu.C("Date").Table("a").Between(goqu.Range(date, goqu.L("LAST_DAY(?)", date))),
	).ToSQL()
	if err != nil {
		return nil, 0, err
	}

	whereMerchant := []exp.Expression{
		goqu.C("id").Table("m").Eq(goqu.C("merchant_id").Table("t")),
	}

	if fil.IsValidMerchantName() {
		whereMerchant = append(whereMerchant, goqu.C("merchant_name").Table("m").Eq(fil.MerchantName))
	}

	wherOutlet := []exp.Expression{
		goqu.C("id").Table("o").Eq(goqu.C("outlet_id").Table("t")),
		goqu.C("merchant_id").Table("o").Eq(fil.MerchantID),
	}

	if fil.IsValidOutletName() {
		wherOutlet = append(wherOutlet, goqu.C("outlet_name").Table("o").Eq(fil.OutletName))
	}

	queryMerchant, _, _ := goqu.Dialect("mysql").From(goqu.T("Merchants").As("m")).Select("id").Where(whereMerchant...).ToSQL()
	queryOutlet, _, _ := goqu.Dialect("mysql").From(goqu.T("Outlets").As("o")).Select("id").Where(wherOutlet...).ToSQL()

	join := []exp.Expression{
		goqu.I("dateList.Date").Eq(goqu.L("DATE(?)", goqu.I("t.created_at"))),
		goqu.L("EXISTS (" + queryMerchant + ")"),
		goqu.L("EXISTS (" + queryOutlet + ")"),
	}

	if fil.IsValidMerchantID() {
		join = append(join, goqu.C("merchant_id").Table("t").Eq(fil.MerchantID))
	}

	if fil.IsValidOutletID() {
		join = append(join, goqu.C("outlet_id").Table("t").Eq(fil.OutletID))
	}

	dataset := dialect.From(
		goqu.L("("+qdl+")", qdlArgs...).As("dateList"),
	).LeftJoin(
		goqu.T("Transactions").As("t"), goqu.On(join...),
	).Select(
		goqu.C("Date").Table("dateList").As("date"),
		goqu.L("IFNULL(SUM(?), 0)", goqu.I("t.bill_total")).As("omzet"),
	).GroupBy(
		goqu.C("Date").Table("dateList"),
	).Order(
		goqu.C("Date").Table("dateList").Asc(),
	).Limit(uint(fil.Limit)).Offset(uint(offset))

	query, args, err := dataset.ToSQL()
	if err != nil {
		return nil, 0, err
	}

	trxReport := []schema.TransactionReport{}
	if err := t.DB.FindAll(ctx, &trxReport, rel.SQL(query, args...)); err != nil {
		return nil, 0, err
	}

	queryCount, argsCount, err := dataset.ClearSelect().GroupBy().ClearOrder().ClearLimit().ClearOffset().
		Select(
			goqu.COUNT(
				goqu.DISTINCT(goqu.I("dateList.Date")),
			).As("count"),
		).ToSQL()
	if err != nil {
		return nil, 0, err
	}

	countTrxReport := schema.CountTransactionReport{}
	if err := t.DB.Find(ctx, &countTrxReport, rel.SQL(queryCount, argsCount...)); err != nil {
		return nil, 0, err
	}

	return trxReport, countTrxReport.Count, nil
}

func (t *Transaction) Reporting(ctx context.Context, fil ReportFilter) ([]schema.TransactionReport, error) {
	dialect := goqu.Dialect("mysql").From("Transactions").Select(
		goqu.L("DATE(?)", goqu.I("created_at")).As("date"),
		goqu.L("SUM(?)", goqu.I("bill_total")).As("omzet"),
	)

	where := []exp.Expression{
		goqu.C("merchant_id").Eq(fil.MerchantID),
	}

	if fil.IsValidOutletID() {
		where = append(where, goqu.C("outlet_id").Eq(fil.OutletID))
	}

	if fil.IsValidRangeDate() {
		where = append(where, goqu.C("created_at").Between(
			goqu.Range(fil.StartDate.Format("2006-01-02 15:04:05"), fil.EndDate.Format("2006-01-02 15:04:05")),
		))
	}

	query, args, err := dialect.Where(where...).GroupBy(goqu.C("date")).Order(goqu.C("date").Asc()).ToSQL()
	if err != nil {
		return nil, err
	}

	trxReport := []schema.TransactionReport{}
	if err := t.DB.FindAll(ctx, &trxReport, rel.SQL(query, args...)); err != nil {
		return nil, err
	}

	return trxReport, nil
}
