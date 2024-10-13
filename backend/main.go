package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Customer struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	Email            string    `json:"email"`
	RegistrationDate time.Time `json:"registration_date"`
}

type Transaction struct {
	ID               int       `json:"id"`
	CustomerID       int       `json:"customer_id"`
	BorrowFee        float64   `json:"borrow_fee"`
	TransactionDate  time.Time `json:"transaction_date"`
	TransactionCount int       `json:"transaction_count"`
}

//TIP To run your code, right-click the code and select <b>Run</b>. Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.

var db *sql.DB

func main() {
	//連接資料庫
	//db, err := sql.Open("mysql", "root:jo2930yen@tcp(127.0.0.1:3306)/exampleDB")
	var err error
	db, err = sql.Open("mysql", "test:test@tcp(mariadb:3306)/pretest?parseTime=true")
	if err != nil {
		log.Fatal("Failed to open database connection:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Successfully connected to database")

	//初始化數據
	err = initializeData()
	if err != nil {
		log.Fatal("Failed to initialize data:", err)
	}

	r := gin.Default()

	// 添加 CORS 中間件
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}
	r.Use(cors.New(config))

	api := r.Group("/api")
	{
		api.GET("/customers", getCustomers)
		api.GET("/customers/:id", getCustomer)
		api.POST("/customers", createCustomer)
		api.PUT("/customers/:id", updateCustomer)
		api.GET("/customers/:id/transactions", getCustomerTransactions)
	}

	if err := r.Run("0.0.0.0:8080"); err != nil {
		log.Fatal("Failed to run server:", err)
	}

}

func initializeData() error {
	err := clearTables()
	if err != nil {
		return fmt.Errorf("清除table資料失敗: %v", err)
	}

	customers, err := generateCustomers(1000)
	if err != nil {
		return fmt.Errorf("客戶資料生成失敗: %v", err)
	}

	err = generateTransactions(customers, 5000)
	if err != nil {
		return fmt.Errorf("交易資料生成失敗: %v", err)
	}

	log.Println("資料生成完畢")
	return nil
}

func clearTables() error {
	_, err := db.Exec("DELETE FROM transactions")
	if err != nil {
		return fmt.Errorf("failed to clear transactions table: %v", err)
	}

	_, err = db.Exec("DELETE FROM customers")
	if err != nil {
		return fmt.Errorf("failed to clear customers table: %v", err)
	}

	_, err = db.Exec("ALTER TABLE transactions AUTO_INCREMENT = 1")
	if err != nil {
		return fmt.Errorf("failed to reset transactions auto increment: %v", err)
	}

	_, err = db.Exec("ALTER TABLE customers AUTO_INCREMENT = 1")
	if err != nil {
		return fmt.Errorf("failed to reset customers auto increment: %v", err)
	}

	log.Println("Tables cleared and auto increment reset")
	return nil
}

func generateCustomers(count int) (map[int]time.Time, error) {
	customers := make(map[int]time.Time)
	for i := 0; i < count; i++ {
		name := fmt.Sprintf("Customer%d", i+1)
		email := fmt.Sprintf("customer%d@gmail.com", i+1)
		registrationDate := randomDate()

		result, err := db.Exec("INSERT INTO customers (name, email, registration_date) VALUES (?, ?, ?)", name, email, registrationDate)
		if err != nil {
			log.Fatal(err)
		}

		id, _ := result.LastInsertId()
		customers[int(id)] = registrationDate
	}
	fmt.Printf("已生成 %d 筆客戶資料\n", count)
	return customers, nil
}

func generateTransactions(customers map[int]time.Time, count int) error {
	transactionCounts := make(map[int]int)
	customerIDs := make([]int, 0, len(customers))
	for id := range customers {
		customerIDs = append(customerIDs, id)
	}

	for i := 0; i < count; i++ {
		customerIndex := rand.Intn(len(customerIDs))
		if rand.Float32() < 0.7 { //70% 的機會選擇前半部分的客戶
			customerIndex = rand.Intn(len(customerIDs) / 2)
		}
		customerID := customerIDs[customerIndex]
		registrationDate := customers[customerID]

		//確保交易日期在註冊後18個月內
		maxTransactionDate := registrationDate.AddDate(0, 18, 0)
		if maxTransactionDate.After(time.Now()) {
			maxTransactionDate = time.Now()
		}
		transactionDate := randomDateBetween(registrationDate, maxTransactionDate)

		//生成隨機金額
		borrowFee := rand.Float64()*9999 + 10

		transactionCounts[customerID]++
		transactionCount := transactionCounts[customerID]

		_, err := db.Exec(`
            INSERT INTO transactions 
            (customer_id, borrow_fee, transaction_date, transaction_count) 
            VALUES (?, ?, ?, ?)`,
			customerID, borrowFee, transactionDate, transactionCount)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Printf("已生成 %d 筆交易資料\n", count)

	return nil
}

func randomDate() time.Time {
	now := time.Now()
	min := now.AddDate(-2, 0, 0).Unix()
	max := now.Unix()
	randTime := rand.Int63n(max-min) + min
	return time.Unix(randTime, 0)
}

func randomDateBetween(start, end time.Time) time.Time {
	delta := end.Sub(start)
	deltaNanos := delta.Nanoseconds()
	randomNanos := rand.Int63n(deltaNanos)
	return start.Add(time.Duration(randomNanos))
}

func getCustomers(c *gin.Context) {
	rows, err := db.Query("SELECT id, name, email, registration_date FROM customers")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var customers []Customer
	for rows.Next() {
		var customer Customer
		if err := rows.Scan(&customer.ID, &customer.Name, &customer.Email, &customer.RegistrationDate); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		customers = append(customers, customer)
	}
	c.JSON(http.StatusOK, customers)
}

func getCustomer(c *gin.Context) {
	id := c.Param("id")
	var cust Customer
	err := db.QueryRow("SELECT id, name, email, registration_date FROM customers WHERE id = ?", id).Scan(&cust.ID, &cust.Name, &cust.Email, &cust.RegistrationDate)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	// 獲取過去一年的交易金額
	var totalAmount float64
	oneYearAgo := time.Now().AddDate(-1, 0, 0)
	err = db.QueryRow("SELECT COALESCE(SUM(borrow_fee), 0) FROM transactions WHERE customer_id = ? AND transaction_date > ?", id, oneYearAgo).Scan(&totalAmount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"customer": cust, "total_amount_last_year": totalAmount})
}

func createCustomer(c *gin.Context) {
	var newCustomer Customer
	if err := c.BindJSON(&newCustomer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO customers (name, email, registration_date) VALUES (?, ?, ?)", newCustomer.Name, newCustomer.Email, time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	newCustomer.ID = int(id)
	c.JSON(http.StatusCreated, newCustomer)
}

func updateCustomer(c *gin.Context) {
	id := c.Param("id")
	var updatedCustomer Customer
	if err := c.BindJSON(&updatedCustomer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("UPDATE customers SET name = ?, email = ? WHERE id = ?", updatedCustomer.Name, updatedCustomer.Email, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedCustomer)
}

func getCustomerTransactions(c *gin.Context) {
	customerID := c.Param("id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	query := `SELECT id, customer_id, borrow_fee, transaction_date, transaction_count 
              FROM transactions 
              WHERE customer_id = ? AND transaction_date BETWEEN ? AND ?`

	rows, err := db.Query(query, customerID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		if err := rows.Scan(&t.ID, &t.CustomerID, &t.BorrowFee, &t.TransactionDate, &t.TransactionCount); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		transactions = append(transactions, t)
	}

	c.JSON(http.StatusOK, transactions)
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
