package otelsql

import (
	"database/sql/driver"

	"go.opentelemetry.io/otel/trace"
)

type countableRows struct {
	rows driver.Rows
	span trace.Span

	count int
}

func newCountableRows(
	rows driver.Rows,
	span trace.Span,
) *countableRows {
	return &countableRows{rows: rows, span: span}
}

func (c *countableRows) Columns() []string {
	return c.rows.Columns()
}

func (c *countableRows) Close() error {
	defer func() {
		c.span.SetAttributes(dbRowsUnmarshalled.Int(c.count))
		c.span.End()
	}()

	return c.rows.Close()
}

func (c *countableRows) Next(dest []driver.Value) error {
	err := c.rows.Next(dest)
	if err == nil {
		c.count++
	}
	return err
}
