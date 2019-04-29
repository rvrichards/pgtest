package main

// Connect to a postgres table, perform insert, update, and Select.

import (
  "database/sql"
  "fmt"
  _ "github.com/lib/pq"
)

const (
  host     = "localhost"   
  port     = 5432
  user     = "test"
  password = "test"
  dbname   = "realm"
)

// This is for the query call, so we get and have access to the entire record.
type Client struct {
  ID        int
  Name      string
  Age       int
  Knowledge string
  Item     string
}

func main() {
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)
    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
    panic(err)
    }
    defer db.Close()

    // ADD
    sqlStatement := ` INSERT INTO client (name, age, knowledge, item)
    VALUES ($1, $2, $3, $4) 
    RETURNING id`    // This returns the ID. Other SQLs will return this on INSERT.
    id := 0
    //err = db.QueryRow(sqlStatement, "fred", 42, "This is a test", "Blood Axe").Scan(&id)
    //err = db.QueryRow(sqlStatement, "wilma", 42, "This is not a test", "Sword of Truth").Scan(&id)
    //err = db.QueryRow(sqlStatement, "betty", 41, "Yabba Dabba", "Turin Shoud").Scan(&id)
    //err = db.QueryRow(sqlStatement, "bambam", 2, "Yabba Dabba Doo", "A club").Scan(&id)
    err = db.QueryRow(sqlStatement, "Bob2", 2, "Bob2 was here", "BobBob").Scan(&id)
    if err != nil {  panic(err)  }
    fmt.Println("New record ID is:", id)

    // UPDATE
    sqlStatement = ` UPDATE client
    SET name = $2, age = $3
    WHERE id = $1;`
    res, err := db.Exec(sqlStatement, 1, "Barney", 44 )
    if err != nil {
    panic(err)
    }
    count, errUpdate := res.RowsAffected()
    if errUpdate != nil {
    panic(errUpdate)    
    }
    fmt.Println("Records modified is:", count)

    // SELECT
    sqlStatement = `SELECT * FROM client WHERE id=$1;`
    var client Client
    row := db.QueryRow(sqlStatement, 5)   // QueryRow - get one record
    err = row.Scan(&client.ID, &client.Name, &client.Age, &client.Knowledge, &client.Item)
    switch err {
        case sql.ErrNoRows:
            fmt.Println("No rows were returned!")
            return
        case nil:
            fmt.Println(client)
        default:
            panic(err)
        }
    rows, err := db.Query("SELECT id, name FROM client LIMIT $1", 3)
    if err != nil { panic(err) }
    defer rows.Close()
        for rows.Next() {
        var id int
        var name string
        err = rows.Scan(&id, &name)
        if err != nil { panic(err) }
        fmt.Println("Returning: ", id, name)
    }
    // Any more errors?
    err = rows.Err()
    if err != nil {
        panic(err)
    } 

    // Query - select mulptile rows
    rows, err = db.Query("SELECT id, name FROM client WHERE age=$1 LIMIT $2", 42, 3)
    if err != nil { panic(err)  }
    defer rows.Close()
    for rows.Next() {
        var id int
        var firstName string
        err = rows.Scan(&id, &firstName)
        if err != nil { panic(err) }
        fmt.Println(id, firstName)
    }
    // Any more errors?
    err = rows.Err()
    if err != nil { panic(err) }       
}