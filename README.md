# Pytricia Go

A lazy port of pytricia from python to golang.

# Usage
``` go
import "github.com/tannerklineintz/pytricia-go"
func main() {
    i := pytricia.NewPyTricia()

    i.Insert("8.8.8.0/24", "test123")

    val := i.Get("8.8.8.10/31")
    val = i.Get("8.8.8.8")

    if i.Contains("8.8.8.8") {
        // ...
    }
    
    ch := i.Children()
    pr := i.Parent()
}
```

# Future implementations
.HasKey()
.ToList()