package otelsql

import (
	"database/sql/driver"

	"go.opentelemetry.io/otel/trace"
)

type countableRows struct {
	rows      driver.Rows
	closeSpan trace.Span

	count int
}

func newCountableRows(
	rows driver.Rows,
	closeSpan trace.Span,
) *countableRows {
	return &countableRows{rows: rows, closeSpan: closeSpan}
}

func (c *countableRows) Columns() []string {
	return c.rows.Columns()
}

func (c *countableRows) Close() error {
	defer func() {
		c.closeSpan.SetAttributes(dbRowsUnmarshalled.Int(c.count))
		c.closeSpan.End()
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
