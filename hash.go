package main; import "golang.org/x/crypto/bcrypt"; import "fmt"; func main() { hash, _ := bcrypt.GenerateFromPassword([]byte("AmG150571"), bcrypt.DefaultCost); fmt.Println(string(hash)) }
