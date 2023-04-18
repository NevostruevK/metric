package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/NevostruevK/metric/internal/storage"
	"github.com/NevostruevK/metric/internal/util/logger"
	"github.com/NevostruevK/metric/internal/util/metrics"
	_ "github.com/lib/pq"
)

const errDuplicateID = `pq: duplicate key value violates unique constraint "metrics_id_key"`

type metricSQL struct {
	id    string
	mtype string
	delta sql.NullInt64
	value sql.NullFloat64
}

type DB struct {
	db             *sql.DB
	stmtInsGauge   *sql.Stmt
	stmtInsCounter *sql.Stmt
	stmtGetMetric  *sql.Stmt
	stmtUpdGauge   *sql.Stmt
	stmtUpdCounter *sql.Stmt
	logger         *log.Logger
	init           bool
}

func NewDB(ctx context.Context, connStr string) (*DB, error) {
	db := &DB{db: nil, logger: logger.NewLogger("postgres : ", log.LstdFlags|log.Lshortfile), init: false}
	if connStr == "" {
		db.logger.Println("DATABASE_DSN is empty, database wasn't initialized")
		return db, fmt.Errorf(" Empty address data base")
	}
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return db, err
	}

	if _, err = conn.ExecContext(ctx, schemaSQL); err != nil {
		return db, err
	}

/*	stmtInsGauge, err := conn.Prepare(insertGaugeSQL)
	if err != nil {
		return db, err
	}
	stmtInsCounter, err := conn.Prepare(insertCounterSQL)
	if err != nil {
		return db, err
	}
*/	stmtGetMetric, err := conn.PrepareContext(ctx, getMetricSQL)
	if err != nil {
		return db, err
	}
/*	stmtUpdGauge, err := conn.Prepare(updateGaugeSQL)
	if err != nil {
		return db, err
	}

	stmtUpdCounter, err := conn.Prepare(updateCounterSQL)
	if err != nil {
		return db, err
	}
*/
	db.db = conn
//	db.stmtInsGauge = stmtInsGauge
//	db.stmtInsCounter = stmtInsCounter
	db.stmtGetMetric = stmtGetMetric
//	db.stmtUpdGauge = stmtUpdGauge
//	db.stmtUpdCounter = stmtUpdCounter
	db.init = true
	return db, nil
}

func (db *DB) ShowMetrics(ctx context.Context) error {
	l := logger.NewLogger("", 0)
	db.logger.Println("Show metrics")
	rows, err := db.db.QueryContext(ctx, "SELECT * FROM metrics")
	if err != nil {
		db.logger.Println(err)
		return err
	}
	defer rows.Close()
	mSQL := metricSQL{}
	for rows.Next() {
		newDelta := sql.NullInt64{}
		newValue := sql.NullFloat64{}
		err = rows.Scan(&mSQL.id, &mSQL.mtype, &newDelta, &newValue)
		if err != nil {
			db.logger.Println(err)
			continue
		}
		m := metrics.Metrics{ID: mSQL.id, MType: mSQL.mtype}
		if mSQL.mtype == metrics.Gauge {
			if newValue.Valid {
				m.Value = &newValue.Float64
				l.Println(m)
			}
			continue
		}
		if mSQL.mtype == metrics.Counter {
			if newDelta.Valid {
				m.Delta = &newDelta.Int64
				l.Println(m)
			}
			continue
		}
	}
	return rows.Err()
}

func (db *DB) GetAllMetrics(ctx context.Context) ([]metrics.Metrics, error) {
	var size int
	err := db.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM metrics").Scan(&size)
	if err != nil {
		db.logger.Println(err)
		return nil, err
	}
	sM := make([]metrics.Metrics, 0, size+10)
	rows, err := db.db.QueryContext(ctx, "SELECT * FROM metrics")
	if err != nil {
		db.logger.Println(err)
		return nil, err
	}
	defer rows.Close()
	mSQL := metricSQL{}
	for rows.Next() {
		newDelta := sql.NullInt64{}
		newValue := sql.NullFloat64{}
		err = rows.Scan(&mSQL.id, &mSQL.mtype, &newDelta, &newValue)
		if err != nil {
			db.logger.Println(err)
			continue
		}
		m := metrics.Metrics{ID: mSQL.id, MType: mSQL.mtype}
		if mSQL.mtype == metrics.Gauge {
			if newValue.Valid {
				m.Value = &newValue.Float64
				sM = append(sM, m)
			}
			db.logger.Printf("ERROR : invalid gauge metric %s\n", mSQL.id)
			continue
		}
		if mSQL.mtype == metrics.Counter {
			if newDelta.Valid {
				m.Delta = &newDelta.Int64
				sM = append(sM, m)
			}
			db.logger.Printf("ERROR : invalid counter metric %s\n", mSQL.id)
			continue
		}
	}
	return sM, rows.Err()
}

func (db *DB) AddGroupOfMetrics(ctx context.Context, sM []metrics.Metrics) error {
	
	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		db.logger.Println(err)
		return err
	}
	defer func(){
		err := tx.Rollback()
		if !errors.Is(err, sql.ErrTxDone){
			db.logger.Println(err)
		}
	}()
	
	for _, m := range sM{
		err := db.AddMetric(ctx, &m)
		if err != nil{
			db.logger.Printf("ERROR : AddGroupOfMetrics:AddMetric : %v\n",err)			
			return err
		}
	}
	return tx.Commit()
}


func (db *DB) AddMetric(ctx context.Context, rt storage.RepositoryData) error {
	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		db.logger.Println(err)
		return err
	}
	defer func(){
		err := tx.Rollback()
		if !errors.Is(err, sql.ErrTxDone){
			db.logger.Println(err)
		}
	}()

	switch rt.Type(){
	case metrics.Gauge:
		_, err := db.db.ExecContext(ctx, insertGaugeSQL, rt.Name(),  rt.GaugeValue())
		if err != nil{
			if !strings.Contains(err.Error(), errDuplicateID){
				db.logger.Printf("ERROR : Inserted %s %v\n", rt, err)
				return err
			}
			_, err = db.db.ExecContext(ctx, updateGaugeSQL, rt.Name(), rt.GaugeValue())
			if err != nil{
				db.logger.Printf("ERROR : Updated %s %v\n", rt, err)
				return err
			}
			db.logger.Printf("Updated %s \n", rt)
			return tx.Commit()
		} 
	case metrics.Counter:
		_, err := db.db.ExecContext(ctx, insertCounterSQL, rt.Name(), rt.CounterValue())
		if err != nil{
			if !strings.Contains(err.Error(), errDuplicateID){
				db.logger.Printf("ERROR : Inserted %s %v\n", rt, err)
				return err
			}
			mSQL := metricSQL{}
			if err = db.db.QueryRowContext(ctx, getCounterMetricSQL, rt.Name()).Scan(&mSQL.mtype, &mSQL.delta); err != nil{
				db.logger.Printf("ERROR : getCounterMetricSQL %s %v\n", rt, err)
				return err
			}
			if mSQL.mtype == metrics.Counter && mSQL.delta.Valid {
				if err = rt.AddCounterValue(mSQL.delta.Int64); err != nil {
					db.logger.Printf("ERROR : AddCounterValue %s %v\n", rt, err)
					return err
				}
			}				
			_, err = db.db.ExecContext(ctx, updateCounterSQL, rt.Name(), rt.CounterValue())
			if err != nil{
				db.logger.Printf("ERROR : Updated %s %v\n", rt, err)
				return err
			}
			db.logger.Printf("Updated %s \n", rt)
			return tx.Commit()
		}
	default:
		msg := fmt.Sprintf("ERROR : AddMetric is not implemented for type %s\n",rt.Type())
		db.logger.Println(msg)
		return fmt.Errorf(msg)
	}
	db.logger.Printf("Inserted %s \n", rt)
	return tx.Commit()
}

func (db *DB) GetMetric(ctx context.Context, reqType, name string) (storage.RepositoryData, error) {
	if validType := metrics.IsMetricType(reqType); !validType {
		return nil, fmt.Errorf("type %s is not valid metric type", reqType)
	}
	mSQL := metricSQL{}
	err := db.stmtGetMetric.QueryRowContext(ctx, name).Scan(&mSQL.id, &mSQL.mtype, &mSQL.delta, &mSQL.value)
	if err != nil {
		return nil, fmt.Errorf("metric name %s is not found", name)
	}
	if mSQL.mtype != reqType {
		return nil, fmt.Errorf("metric type %s name %s is not found", reqType, name)
	}
	m := metrics.Metrics{ID: mSQL.id, MType: mSQL.mtype}
	if mSQL.mtype == metrics.Gauge {
		if mSQL.value.Valid {
			m.Value = &mSQL.value.Float64
			return &m, nil
		}
		return nil, fmt.Errorf("metric type %s name %s is value is nil", reqType, name)
	}
	if mSQL.mtype == metrics.Counter {
		if mSQL.delta.Valid {
			m.Delta = &mSQL.delta.Int64
			return &m, nil
		}
		return nil, fmt.Errorf("metric type %s name %s is delta is nil", reqType, name)
	}
	return nil, fmt.Errorf("type %s is not valid metric type", reqType)
}

func (db *DB) Close() error {

	if !db.init {
		return fmt.Errorf(" Can't close DB : DataBase wasn't initiated")
	}
	db.stmtInsGauge.Close()
	db.stmtInsCounter.Close()
	db.stmtGetMetric.Close()
	db.stmtUpdGauge.Close()
	db.stmtUpdCounter.Close()
	if err := db.db.Close(); err != nil {
		return fmt.Errorf(" Can't close DB %w", err)
	}
	return nil
}

func (db DB) Ping() error {
	if !db.init {
		return fmt.Errorf(" Can't ping DB : DataBase wasn't initiated")
	}
	if err := db.db.Ping(); err != nil {
		return fmt.Errorf(" Can't ping DB %w", err)
	}
	return nil
}
