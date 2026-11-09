package database

import (
	"database/sql"
	"testing"

	"github.com/danmaciel/ca_golang/internal/entity"
	"github.com/stretchr/testify/suite"

	// sqlite3
	_ "github.com/mattn/go-sqlite3"
)

type OrderRepositoryTestSuite struct {
	suite.Suite
	Db *sql.DB
}

func (suite *OrderRepositoryTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	suite.NoError(err)
	db.Exec("CREATE TABLE orders (id varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, PRIMARY KEY (id))")
	suite.Db = db
}

func (suite *OrderRepositoryTestSuite) TearDownSuite() {
	suite.Db.Close()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuite))
}

func (suite *OrderRepositoryTestSuite) TestGivenAnOrder_WhenSave_ThenShouldSaveOrder() {
	order, err := entity.NewOrder("123", 10.0, 2.0)
	suite.NoError(err)
	suite.NoError(order.CalculateFinalPrice())
	repo := NewOrderRepository(suite.Db)
	err = repo.Save(order)
	suite.NoError(err)

	var orderResult entity.Order
	err = suite.Db.QueryRow("Select id, price, tax, final_price from orders where id = ?", order.ID).
		Scan(&orderResult.ID, &orderResult.Price, &orderResult.Tax, &orderResult.FinalPrice)

	suite.NoError(err)
	suite.Equal(order.ID, orderResult.ID)
	suite.Equal(order.Price, orderResult.Price)
	suite.Equal(order.Tax, orderResult.Tax)
	suite.Equal(order.FinalPrice, orderResult.FinalPrice)
}

func TestGetAllReturnsRowsAndTotal(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec("CREATE TABLE orders (id TEXT PRIMARY KEY, price REAL NOT NULL, tax REAL NOT NULL, final_price REAL NOT NULL)")
	if err != nil {
		t.Fatal(err)
	}

	repo := NewOrderRepository(db)
	order, err := entity.NewOrder("test-order-1", 10, 2)
	if err != nil {
		t.Fatal(err)
	}
	if err := order.CalculateFinalPrice(); err != nil {
		t.Fatal(err)
	}
	if err := repo.Save(order); err != nil {
		t.Fatal(err)
	}
	result, err := repo.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 || result[0].ID != order.ID {
		t.Fatalf("unexpected orders: %+v", result)
	}
	total, err := repo.GetTotal()
	if err != nil || total != 1 {
		t.Fatalf("expected total 1; total=%d err=%v", total, err)
	}
}

func TestGetAllPropagatesScanErrors(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	// A nullable field models a corrupted or older schema. It must not
	// silently produce an incomplete list when scanning into float fields.
	_, err = db.Exec("CREATE TABLE orders (id TEXT PRIMARY KEY, price REAL NOT NULL, tax REAL, final_price REAL NOT NULL)")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, NULL, ?)", "invalid-order", 10, 10)
	if err != nil {
		t.Fatal(err)
	}
	orders, err := NewOrderRepository(db).GetAll()
	if err == nil {
		t.Fatal("expected scan error for NULL tax")
	}
	if orders != nil {
		t.Fatalf("must not return a partial list after scan failure: %+v", orders)
	}
}

func TestGetAllPropagatesQueryErrors(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	orders, err := NewOrderRepository(db).GetAll()
	if err == nil || orders != nil {
		t.Fatalf("missing table must return a query error, got orders=%+v err=%v", orders, err)
	}
}
