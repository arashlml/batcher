package repository

import (
	"batcher/batcher"
	"batcher/entity"
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync/atomic"
)

type Repo struct {
	Counter         int
	dbInsertCounter int64
	batcher         *batcher.Batcher[entity.Request]
	db              *sql.DB
}

func NewRepository(db *sql.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) BatchInsert(ctx context.Context, items []entity.Request) error {
	if len(items) == 0 {
		log.Println("REPOSITORY : ITEM IS EMPTY!")
		return nil
	}
	atomic.AddInt64(&r.dbInsertCounter, int64(len(items)))
	log.Printf("DB INSERT ITEMS = %d\n", atomic.LoadInt64(&r.dbInsertCounter))

	query := `INSERT INTO requests (field1, field2, field3, field4, field5) VALUES `
	args := []any{}
	log.Println("BatchInsert:", len(items))
	for i, item := range items {
		if i > 0 {
			query += ", "
		}
		placeholderStart := i*5 + 1
		query += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d)", placeholderStart, placeholderStart+1, placeholderStart+2, placeholderStart+3, placeholderStart+4)
		args = append(args, item.Field1, item.Field2, item.Field3, item.Field4, item.Field5)
	}
	result, err := r.db.ExecContext(ctx, query, args...)
	log.Printf("Sql result: --> %v", result)
	if err != nil {
		log.Printf("error executing the query --> %v ", err)
		return err
	}
	log.Println("Successfully inserted the batch")
	return nil
}
