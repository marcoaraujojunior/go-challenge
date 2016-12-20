package model

import (
	"errors"
	"time"
	"services/database"
	"strings"
	"strconv"
	"math"
)

var (
	month = map[int]string{
		1: "January",
		2: "February",
		3: "March",
		4: "April",
		5: "May",
		6: "June",
		7: "July",
		8: "August",
		9: "September",
		10: "October",
		11: "November",
		12: "December"}
)
type Invoice struct {
	ID             uint      `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	ReferenceMonth int
	ReferenceYear  int
	Document       string    `sql:"size:14"`
	Description    string    `sql:"size:255"`
	Amount         float32   `sql:"type:decimal(10,2)"`
	IsActive       bool      `sql:"default:1"`
	CreatedAt      time.Time `sql:"type:datetime;default:current_timestamp"`
	DeactiveAt     time.Time `sql:"type:datetime"`
}

type Invoices struct {
	Records  []Invoice
	Page     int
	PerPage  int
	Total    int
}

func buildOrder(p []string) string {
	var orderBy []string
	fields := strings.SplitN(p[0], ",", 3)
	for _, field := range fields {
		pair := strings.SplitN(string(field), "-", 2)
		order := "ASC"
		if len(pair) > 1 {
			field = pair[1]
			order = "DESC"
		}
		orderBy = append(orderBy, field + " " + order)
	}
	return strings.Join(orderBy, ",")
}

func buildQuery(q map[string][]string) map[string]interface{} {
	var query map[string]interface{}
	var qf Invoice
	var order, mention string
	page := 1
	per_page := 50
	query = make(map[string]interface{})
	if val, ok := q["sort"]; ok {
		order = buildOrder(val)
	}
	if val, ok := q["month"]; ok {
		qf.ReferenceMonth, _ = strconv.Atoi(val[0])
	}
	if val, ok := q["year"]; ok {
		qf.ReferenceYear, _ = strconv.Atoi(val[0])
	}
	if val, ok := q["document"]; ok {
		qf.Document = val[0]
	}
	if val, ok := q["per_page"]; ok {
		per_page, _ = strconv.Atoi(val[0])
	}
	if val, ok := q["page"]; ok {
		page, _ = strconv.Atoi(val[0])
	}
	if val, ok := q["q"]; ok {
		mention = "%" + val[0] + "%"
	}
	query["condition"] = qf
	query["mention"] = mention
	query["order"] = order
	query["limit"] = per_page
	query["offset"] = ( ( page - 1 ) * per_page )
	query["page"] = page
	return query
}

func GetAll(q map[string][]string) Invoices {
	var invoices []Invoice
	var total, count int
	query := buildQuery(q)
	orm := database.GetDb()
	if (query["mention"] != "") {
		orm = orm.Where("document LIKE ?", query["mention"])
	}
	orm = orm.Where(query["condition"])
    if query["order"] != "" {
		orm = orm.Order(query["order"])
    }
	orm = orm.Limit(query["limit"]).
		Offset(query["offset"]).
		Find(&invoices)

	count = toCount(query)

	total = int(math.Ceil(float64(count) / float64(query["limit"].(int))))
	return Invoices{Records: invoices,
		Page: query["page"].(int),
		PerPage: query["limit"].(int),
		Total: total}
}

func toCount(query map[string]interface{}) int {
	var invoices []Invoice
	var total int
	orm := database.GetDb()
	if (query["mention"] != "") {
		orm = orm.Where("document LIKE ?", query["mention"])
	}
	orm = orm.Where(query["condition"]).
		Find(&invoices).
		Count(&total)
    return total
}

func Delete(invoiceNumber string) error {
	var invoice Invoice
	var err error
	invoice, err = Get(invoiceNumber)
	if (err != nil) {
		return err
	}
	invoice.IsActive = false
	invoice.DeactiveAt = time.Now()
	return Save(&invoice)
}

func Update(attributes Invoice) error {
	invoice, err := Get(attributes.Document)
	if (err != nil) {
		return err
	}
	if attributes.ReferenceMonth > 0 {
		invoice.ReferenceMonth = attributes.ReferenceMonth
	}
	if attributes.ReferenceYear > 0 {
		invoice.ReferenceYear = attributes.ReferenceYear
	}
	if len(attributes.Document) > 0 {
		invoice.Document = attributes.Document
	}
	if len(attributes.Description) > 0 {
		invoice.Description = attributes.Description
	}
	if attributes.Amount > 0 {
		invoice.Amount = attributes.Amount
	}

	return Save(&invoice)
}

func Create(attributes *Invoice) error {
	_, err := Get(attributes.Document)
	if (attributes.Document == "") {
		err = errors.New("Attribute Document containing invoice number is required")
		return err
	}
	if (err == nil) {
		err = errors.New("Invoice [" + attributes.Document + "] already exist")
		return err
	}
	return Save(attributes)
}

func inMonth(m int) bool {
	if _, ok := month[m]; ok {
		return true
	}
	return false
}

func Save(i *Invoice) error {
	if (!inMonth(i.ReferenceMonth)) {
		err := errors.New("Invalid month")
		return err
	}
	return database.GetDb().Debug().Save(&i).Error
}

func Get(invoiceNumber string) (Invoice, error) {
	var invoice Invoice
	var err error
	if (invoiceNumber == "") {
		err = errors.New("Attribute invoice number is required")
		return invoice, err
	}
	if database.GetDb().Where("document = ?", invoiceNumber).Find(&invoice).RecordNotFound() {
		err = errors.New("Invoice [" + invoiceNumber + "] not found")
		return invoice, err
	}
	return invoice, nil
}

