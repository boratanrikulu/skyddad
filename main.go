package main

func main() {
  db := Connect()
  defer db.Close()
}
