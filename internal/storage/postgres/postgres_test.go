package postgres_test

import (
	"fmt"
	"testing"

	sq "github.com/Masterminds/squirrel"
)

func TestQueryBuilder(t *testing.T) {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	users := qb.Select("*").From("users")

	active := users.Where("active = $1", "A")
	active = active.Where("abc = $1", "B")

	sql, args, err := active.ToSql()

	t.Log(fmt.Sprintf(sql, args), err)

}
