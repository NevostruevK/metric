package db

import (
	"database/sql"
	"fmt"

	"github.com/NevostruevK/metric/internal/storage"
	"github.com/NevostruevK/metric/internal/util/metrics"
	_ "github.com/lib/pq"
)

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
	init           bool
}

func NewDB(connStr string) (*DB, error) {
	db := &DB{db: nil, init: false}
//	connStr = "user=postgres sslmode=disable"
	if connStr == "" {
		fmt.Println("Empty address data base")
		return db, fmt.Errorf(" Empty address data base")
	}
	//	connStr := "user=postgres sslmode=disable"
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return db, fmt.Errorf(" Can't open DB %w", err)
	}

	if _, err = conn.Exec(schemaSQL); err != nil {
		return db, err
	}
	fmt.Println("Create metrics table")

	stmtInsGauge, err := conn.Prepare(insertGaugeSQL)
	if err != nil {
		return db, err
	}

	stmtInsCounter, err := conn.Prepare(insertCounterSQL)
	if err != nil {
		return db, err
	}

	stmtGetMetric, err := conn.Prepare(getMetricSQL)
	if err != nil {
		return db, err
	}
	stmtUpdGauge, err := conn.Prepare(updateGaugeSQL)
	if err != nil {
		return db, err
	}

	stmtUpdCounter, err := conn.Prepare(updateCounterSQL)
	if err != nil {
		return db, err
	}

	//	fmt.Println("Prepare statement")

	/*	_, err = stmtUpdGauge.Exec("gauge123", 3.123456)
		if err != nil {
			return db, err
		}
	*/
	/*		_, err = stmtInsGauge.Exec("NewGauge", "gauge", 1.1234)
			if err != nil {
				return db, err
			}

			_, err = stmtInsGauge.Exec("NewGauge1", "gauge", 1.12345)
			if err != nil {
				return db, err
			}
			_, err = stmtInsGauge.Exec("NewGauge3", "gauge", 1.123456)
			//	_, err = stmtInsGauge.Exec("NewGauge")
			//	_, err = db.stmtInsGauge.Exec("NewGauge")
			if err != nil {
				return db, err
			}
	*/
	/*	fmt.Println("Try to insert prepared statement")

		var id string
		var mtype string
		var delta sql.NullInt64
		var value sql.NullFloat64

		stmtGetMetric.QueryRow("NewGauge2").Scan(&id, &mtype, &delta, &value)
		// в запросе использован контейнер "?",
		// которому будет передано значение SomeID
		// так можно строить запрос с параметрами
		fmt.Println(id)
		fmt.Println(mtype)
		fmt.Println(delta)
		fmt.Println(value)
		err = stmtGetMetric.QueryRow("NewGauge3").Scan(&id, &mtype, &delta, &value)
		if err != nil {
			return db, err
		}

		//	row, err := stmtGetMetric.QueryRow("NewGauge3")
		fmt.Println("Try to read NewGauge3")
		fmt.Println(id)
		fmt.Println(mtype)
		fmt.Println(delta)
		fmt.Println(value)
	*/
	db.db = conn
	db.stmtInsGauge = stmtInsGauge
	db.stmtInsCounter = stmtInsCounter
	db.stmtGetMetric = stmtGetMetric
	db.stmtUpdGauge = stmtUpdGauge
	db.stmtUpdCounter = stmtUpdCounter
	db.init = true
	return db, nil
}
func (db *DB) ShowMetrics() error {
	rows, err := db.db.Query("SELECT * FROM metrics")
	if err != nil {
		return err
	}
	defer rows.Close()
	mSQL := metricSQL{}
	for rows.Next() {
		newDelta := sql.NullInt64{}
		newValue := sql.NullFloat64{}
		err = rows.Scan(&mSQL.id, &mSQL.mtype, &newDelta, &newValue)
		if err != nil {
			continue
		}
		m := metrics.Metrics{ID: mSQL.id, MType: mSQL.mtype}
		if mSQL.mtype == metrics.Gauge {
			if newValue.Valid {
				m.Value = &newValue.Float64
				fmt.Println(m)
			}
			continue
		}
		if mSQL.mtype == metrics.Counter {
			if newDelta.Valid {
				m.Delta = &newDelta.Int64
				fmt.Println(m)
			}
			continue
		}
	}
	return rows.Err()
}

func (db *DB) GetAllMetrics() ([]metrics.Metrics, error) {
	var size int
	err := db.db.QueryRow("SELECT COUNT(*) FROM metrics").Scan(&size)
	if err != nil {
		return nil, err
	}
	sM := make([]metrics.Metrics, 0, size+10)
	rows, err := db.db.Query("SELECT * FROM metrics")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	mSQL := metricSQL{}
	for rows.Next() {
		newDelta := sql.NullInt64{}
		newValue := sql.NullFloat64{}
		err = rows.Scan(&mSQL.id, &mSQL.mtype, &newDelta, &newValue)
		if err != nil {
			continue
		}
		m := metrics.Metrics{ID: mSQL.id, MType: mSQL.mtype}
		if mSQL.mtype == metrics.Gauge {
			if newValue.Valid {
				m.Value = &newValue.Float64
				sM = append(sM, m)
			}
			continue
		}
		if mSQL.mtype == metrics.Counter {
			if newDelta.Valid {
				m.Delta = &newDelta.Int64
				sM = append(sM, m)
			}
			continue
		}
	}
//	fmt.Println(sM)
	//	return sM, nil
	return sM, rows.Err()
}


func (db *DB) AddGroupOfMetrics(sM []metrics.Metrics) error {
	tx, err := db.db.Begin()
	if err != nil{
		return err
	}
	defer tx.Rollback()

	txStmtGetMetric := tx.Stmt(db.stmtGetMetric)
	txStmtInsGauge := tx.Stmt(db.stmtInsGauge)
	txStmtInsCounter := tx.Stmt(db.stmtInsCounter)
	txStmtUpdGauge := tx.Stmt(db.stmtUpdGauge)
	txStmtUpdCounter := tx.Stmt(db.stmtUpdCounter)

	mSQL := metricSQL{}
	for _, m := range sM{
		err = txStmtGetMetric.QueryRow(m.Name()).Scan(&mSQL.id, &mSQL.mtype, &mSQL.delta, &mSQL.value)
		if err != nil {
			if m.Type() == metrics.Gauge {
				_, err = txStmtInsGauge.Exec(m.Name(), metrics.Gauge, m.GaugeValue())
				if err != nil {
					return err
				}
				continue
			}
			if m.Type() == metrics.Counter {
				_, err = txStmtInsCounter.Exec(m.Name(), metrics.Counter, m.CounterValue())
				if err != nil {
					return err
				}
				continue
			}
			return fmt.Errorf("wrong metric type ")
		}
		if m.Type() == metrics.Counter && mSQL.mtype == metrics.Counter {
			if mSQL.delta.Valid {
				m.AddCounterValue(mSQL.delta.Int64)
			}
		}
		if m.Type() == metrics.Gauge {
//			fmt.Println("Update Gauge ", m.Name(), "  ", m.GaugeValue())
			_, err = txStmtUpdGauge.Exec(m.Name(), m.GaugeValue())
			if err != nil {
				return err
			}
			continue
		}
		if m.Type() == metrics.Counter {
//			fmt.Println("Update Counter ", m.Name(), "  ", m.CounterValue())
			_, err = txStmtUpdCounter.Exec(m.Name(), m.CounterValue())
			if err != nil {
				return err
			}
			continue
		}
		return fmt.Errorf("wrong metric type ")	
	}
	return tx.Commit()
}




func (db *DB) AddMetric(rt storage.RepositoryData) error {

	mSQL := metricSQL{}
	err := db.stmtGetMetric.QueryRow(rt.Name()).Scan(&mSQL.id, &mSQL.mtype, &mSQL.delta, &mSQL.value)
	if err != nil {
		if rt.Type() == metrics.Gauge {
			_, err = db.stmtInsGauge.Exec(rt.Name(), metrics.Gauge, rt.GaugeValue())
			return err
//			if err != nil {
//				return err
//			}
		}
		if rt.Type() == metrics.Counter {
			_, err = db.stmtInsCounter.Exec(rt.Name(), metrics.Counter, rt.CounterValue())
			return err
//			if err != nil {
//				return err
//			}
		}
		return fmt.Errorf("wrong metric type ")
	}
	if rt.Type() == metrics.Counter && mSQL.mtype == metrics.Counter {
		if mSQL.delta.Valid {
			rt.AddCounterValue(mSQL.delta.Int64)
		}
	}
	if rt.Type() == metrics.Gauge {
//		fmt.Println("Updata Gauge ", rt.Name(), "  ", rt.GaugeValue())
		_, err = db.stmtUpdGauge.Exec(rt.Name(), rt.GaugeValue())
		return err

//		if err != nil {
//			return err
//		}
//		return nil
	}
	if rt.Type() == metrics.Counter {
//		fmt.Println("Updata Counter ", rt.Name(), "  ", rt.CounterValue())
		_, err = db.stmtUpdCounter.Exec(rt.Name(), rt.CounterValue())
		return err
//		if err != nil {
//			return err
//		}
//		return nil
	}
	return fmt.Errorf("wrong metric type ")
}

func (db *DB) GetMetric(reqType, name string) (storage.RepositoryData, error) {
	if validType := metrics.IsMetricType(reqType); !validType {
		return nil, fmt.Errorf("type %s is not valid metric type", reqType)
	}
	mSQL := metricSQL{}
	err := db.stmtGetMetric.QueryRow(name).Scan(&mSQL.id, &mSQL.mtype, &mSQL.delta, &mSQL.value)
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
	fmt.Println("CLOSE : ", db)

	if !db.init {
		return fmt.Errorf(" Can't close DB : DataBase wasn't inited")
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
	fmt.Println("PING : ", db)
	if !db.init {
		fmt.Println(" Can't ping DB : DataBase wasn't inited")
		return fmt.Errorf(" Can't ping DB : DataBase wasn't inited")
	}
	if err := db.db.Ping(); err != nil {
		return fmt.Errorf(" Can't ping DB %w", err)
	}
	return nil
}
