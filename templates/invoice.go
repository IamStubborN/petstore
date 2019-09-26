package templates

import (
	"html/template"

	"github.com/pkg/errors"
)

var invoice = `
<img src="https://redocly.github.io/redoc/petstore-logo.png" alt="banner" style = zoom:50% />

# Invoice orders
### from {{.FromDate}} to {{.ToDate}}

| № | User | Pet | Category | Ship Date | Quantity | Price |
| - | ---- | --- | -------- | ----------- | -----  | ----- |
{{range $i := .InvItems}}|{{$i.ID}}|{{$i.User}}|{{$i.Pet}}|{{$i.Category}}|{{$i.ShipDate}}|{{$i.Quantity}}|$ {{$i.Price}}|
{{end}}

<p style="text-align: end; margin-right: 5%">
  <strong>Total price: $ {{.TotalPrice}}</strong>
</p>`

func GetInvoiceTemplate() (*template.Template, error) {
	temp, err := template.New("invoice").Parse(invoice)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse template file")
	}

	return temp, nil
}
