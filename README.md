# Pytricia Go

A port of [pytricia](https://github.com/jsommers/pytricia) to golang.


# Usage
``` go
import "github.com/tannerklineintz/pytricia-go"
func main() {
    i := pytricia.NewPyTricia()

    i.Insert("8.8.8.0/24", "test123")

    val := i.Get("8.8.8.10/31")
    val = i.Get("8.8.8.8")

    if i.HasKey("8.8.8.0/24") {
        // ...
	}

    if i.Contains("8.8.8.8") {
        // ...
    }

    ch := i.Children()
    pr := i.Parent()
}
```
__More example code available in [test code](./pytricia_test.go)__

# TO DO
 - Keys()
 - Values()
 - Delete()
 - Clear()