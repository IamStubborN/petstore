package workers

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/IamStubborN/petstore/config"
	"github.com/IamStubborN/petstore/db"
	"github.com/IamStubborN/petstore/db/models"
	"github.com/IamStubborN/petstore/fileserver"
	"github.com/IamStubborN/petstore/templates"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type InvoiceWorker struct {
	freq    time.Duration
	genTime time.Time
}

func newInvoiceWorker(cfg *config.Config) Worker {
	genTime, err := time.Parse("15:04", cfg.Invoice.GenerateTime)
	if err != nil {
		zap.L().Error("can't parse duration", zap.Error(err))
	}

	return InvoiceWorker{
		freq:    cfg.Invoice.Frequency,
		genTime: genTime.UTC(),
	}
}

func (iw InvoiceWorker) Run(ctx context.Context) {
	for {
		select {
		case <-time.After(duration(iw.genTime)):
			go invoiceCycle(iw.freq)
		case <-ctx.Done():
			zap.L().Info("invoice worker closed")
			return
		}
	}
}

func invoiceCycle(freq time.Duration) {
	generateInvoice(freq)
	for range time.Tick(freq) {
		generateInvoice(freq)
	}
}

func generateInvoice(freq time.Duration) {
	to := time.Now().UTC()
	from := to.Add(-freq)

	toFormatted := to.Format(time.RFC3339)
	fromFormatted := from.Format(time.RFC3339)

	storeI := db.GetStoreDI()
	invoices, err := storeI.CreateInvoiceByDates(context.Background(), fromFormatted, toFormatted)
	if err != nil {
		zap.L().Error("invoice", zap.Error(err))
		return
	}

	filePath, err := genInvoiceFile(fromFormatted, toFormatted, invoices)
	if err != nil {
		zap.L().Error("can't generate template in filePath", zap.Error(err))
		return
	}

	fm := fileserver.GetFM()
	if err = fm.PutFile("invoices", "text/markdown", filePath); err != nil {
		zap.L().Error("can't put filePath on filePath server", zap.Error(err))
		return
	}

	if err = os.Remove(filePath); err != nil {
		zap.L().Error("can't delete temp filePath", zap.Error(err))
	}
}

func genInvoiceFile(from, to string, invoices []*models.InvoiceItem) (string, error) {
	fromDate := strings.Split(from, "T")[0]
	toDate := strings.Split(to, "T")[0]

	invoicesRefactored, err := refactorShipDate(invoices)
	if err != nil {
		return "", err
	}

	var data = struct {
		FromDate   string
		ToDate     string
		InvItems   []*models.InvoiceItem
		TotalPrice string
	}{fromDate, toDate, invoicesRefactored, calcTotalPrice(invoices)}

	template, err := templates.GetInvoiceTemplate()
	if err != nil {
		return "", err
	}

	generatedTime := time.Now().UTC().Format("15:04:05")
	filePath := path.Join(os.TempDir(), "from_"+fromDate+"_to_"+toDate+"_generated_"+generatedTime+".md")

	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	if err := template.Execute(file, data); err != nil {
		return "", err
	}

	return filePath, nil
}

func calcTotalPrice(invoices []*models.InvoiceItem) string {
	var price float64
	for _, item := range invoices {
		price += float64(item.Quantity) * item.Price
	}

	return fmt.Sprintf("%.2f", price)
}

func refactorShipDate(invoices []*models.InvoiceItem) ([]*models.InvoiceItem, error) {
	for idx, item := range invoices {
		shipDate, err := time.Parse(time.RFC3339, item.ShipDate)
		if err != nil {
			return nil, errors.Wrap(err, "can't parse time")
		}

		invoices[idx].ShipDate = shipDate.Format("2006-01-02 15:04:05")
	}

	return invoices, nil
}

func duration(generateTime time.Time) time.Duration {
	t := time.Now().UTC()
	n := time.Date(t.Year(), t.Month(), t.Day(), generateTime.Hour(), generateTime.Minute(), 0, 0, t.Location()).UTC()

	if t.After(n) {
		n = n.Add(24 * time.Hour)
	}

	d := n.Sub(t)

	return d
}
