package main

func main() {
  db := Connect()
  Migrate(db)
  defer db.Close()
}
