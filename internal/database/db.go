package datapack

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
)

// DB struct represents the database connection
type DB struct {
    Conn *sql.DB
}

// NewDB creates a new instance of DB with the provided parameters
func NewDB(username, password, dbname, url string, port int) (*DB, error) {
    db := &DB{}

    connStr := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    url,port, username, password, dbname)
	fmt.Println(connStr);
    conn, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }

    db.Conn = conn

    // Try to ping the database to check if the connection is successful
    if err := db.Conn.Ping(); err != nil {
		panic(err)
        db.Conn.Close()
		fmt.Println("not successfull bois")
        return nil, err
    }


    fmt.Println("Successfully connected to the database")
	if err := createNodeMetaTable(db); err != nil {
        db.Conn.Close()
		panic(err);
        fmt.Println("Failed to create node-meta table:", err)
        return nil, err
    }

    fmt.Println("node-meta table created or already exists")

	
    return db, nil
}

// Close closes the database connection
func (db *DB) Close() error {
    return db.Conn.Close()
}

func createNodeMetaTable(db *DB) error {
	var tableExists bool
err := db.Conn.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1)", "nodemeta").Scan(&tableExists)
if err != nil {
    // Handle error
    return err
}

// Truncate the table if it exists
if tableExists {
    _, err := db.Conn.Exec(`TRUNCATE TABLE nodemeta;`)
    if err != nil {
        // Handle error
        return err
    }
}
    _, err = db.Conn.Exec(`
      
      
		CREATE TABLE IF NOT EXISTS nodemeta (
            id SERIAL PRIMARY KEY,
            datasource VARCHAR(255) NOT NULL,
            tablename VARCHAR(255) NOT NULL
        );
    `)

	
    return err
}

func AddNodeMeta(db *DB, datasource, tablename string) error {
    _, err := db.Conn.Exec(`
        INSERT INTO nodemeta (datasource, tablename)
        VALUES ($1, $2)
    `, datasource, tablename)
    return err
}