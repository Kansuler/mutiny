# Mutiny

A package that will attempt mutiny on your application. Serialize request body units of any possible combination.

## Example use

### Define a payload structure for your tests

```go
package test

import (
	"github.com/Kansuler/mutiny"
	"net/http"
	"testing"
)

// Payload is a struct that will be used to generate several test payloads for a request,
// each field will use `Pass` values of the `mutiny.PossibleValues` by default.
type Payload struct {
	Name    mutiny.PossibleValues
	Age     mutiny.PossibleValues
	Address mutiny.PossibleValues
}

// UseName is a method receiver that tells the Name field to use a certain value for a field. e.g `mutiny.Fail`, 
// `mutiny.Nil`, `mutiny.Pass` or `mutiny.Erroneous`
func (p Payload) UseName(av mutiny.AssignedValue) Payload {
	p.Name = mutiny.SelectValue(p.Name, av)
	return p
}

// UseAge is a method receiver that tells the Age field to use a certain value for a field. e.g `mutiny.Fail`, 
// `mutiny.Nil`, `mutiny.Pass` or `mutiny.Erroneous`
func (p Payload) UseAge(av mutiny.AssignedValue) Payload {
	p.Name = mutiny.SelectValue(p.Name, av)
	return p
}

// UseAddress is a method receiver that tells the struct to use a certain value for a field. e.g `mutiny.Fail`, 
// `mutiny.Nil`, `mutiny.Pass` or `mutiny.Erroneous`
func (p Payload) UseAddress(av mutiny.AssignedValue) Payload {
	p.Name = mutiny.SelectValue(p.Name, av)
	return p
}

func TestServer(t *testing.T) {
	// Define the payload structure with possible values for each field of the struct.
	payload := Payload{
		Name: mutiny.PossibleValues{
			Pass:      []any{"John"},
			Fail:      []any{"", nil},
			Erroneous: []any{1, 1.1, true},
		},
		Age: mutiny.PossibleValues{
			Pass:      []any{20},
			Fail:      []any{15, -1, nil},
			Erroneous: []any{"", true},
		},
		Address: mutiny.PossibleValues{
			Pass:      []any{"Street 1"},
			Fail:      []any{"Street Without Number", nil},
			Erroneous: []any{1, -1, true},
		},
	}

	// Test passing values only
	testUnits := mutiny.Riot(payload)

	for _, unit := range testUnits {
		resp, err := http.Post("http://localhost:8080", "application/json", unit.RequestBody)
		// Do something with response/error ...
	}

	// Test all variants of payload with failed names
	testUnits = mutiny.Riot(payload.UseName(mutiny.Fail))

	for _, unit := range testUnits {
		resp, err := http.Post("http://localhost:8080", "application/json", unit.RequestBody)
		// Do something with response/error ...
	}

	// Test all variants of payload with failed name and erroneous address
	testUnits = mutiny.Riot(payload.UseName(mutiny.Fail).UseAddress(mutiny.Erroneous))

	for _, unit := range testUnits {
		resp, err := http.Post("http://localhost:8080", "application/json", unit.RequestBody)
		// Do something with response/error ...
	}
}
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

