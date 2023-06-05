package mutiny_test

import (
	"encoding/json"
	"fmt"
	"github.com/Kansuler/mutiny"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

// Example payload
type Payload struct {
	CountryCode mutiny.PossibleValues `json:"country_code"`
	Currency    mutiny.PossibleValues `json:"currency"`
	BankAccount mutiny.PossibleValues `json:"bank_account"`
}

func (p Payload) SetCountryCode(value mutiny.AssignedValue) Payload {
	p.CountryCode = mutiny.SelectValue(p.CountryCode, value)
	return p
}

func (p Payload) SetCurrency(value mutiny.AssignedValue) Payload {
	p.Currency = mutiny.SelectValue(p.Currency, value)
	return p
}

func (p Payload) SetBankAccount(value mutiny.AssignedValue) Payload {
	p.BankAccount = mutiny.SelectValue(p.BankAccount, value)
	return p
}

type MutinyTestSuite struct {
	suite.Suite
}

func JSONEq(suite *MutinyTestSuite, expected, actual string) bool {
	var expectedJSONAsInterface, actualJSONAsInterface interface{}

	if err := json.Unmarshal([]byte(expected), &expectedJSONAsInterface); err != nil {
		suite.T().Error(fmt.Sprintf("Expected value ('%s') is not valid json.\nJSON parsing error: '%s'", expected, err.Error()))
		return false
	}

	if err := json.Unmarshal([]byte(actual), &actualJSONAsInterface); err != nil {
		suite.T().Error(fmt.Sprintf("Input ('%s') needs to be valid json.\nJSON parsing error: '%s'", actual, err.Error()))
		return false
	}

	return assert.ObjectsAreEqual(expectedJSONAsInterface, actualJSONAsInterface)
}

func (suite *MutinyTestSuite) TestPayload() {
	payload := Payload{
		CountryCode: mutiny.PossibleValues{
			Pass:      []any{"SE"},
			Fail:      []any{"FEK"},
			Erroneous: []any{"&nbsp;", true, 123},
		},
		Currency: mutiny.PossibleValues{
			Pass:      []any{"EUR", "SEK"},
			Erroneous: []any{"123", "&nbsp;"},
			Fail:      []any{"DKK"},
		},
		BankAccount: mutiny.PossibleValues{
			Pass: []any{
				json.RawMessage(`{"type":"iban","account_number":"SE4550000000058398257466"}`),
			},
			Fail: []any{
				json.RawMessage(`{"type":"national-DK","account_number":"4200874917","bank_code":"0040"}`),
			},
			Erroneous: []any{},
		},
	}

	expectedResults := []string{
		"{\"bank_account\":{\"account_number\":\"SE4550000000058398257466\", \"type\":\"iban\"}, \"country_code\":\"SE\", \"currency\":\"EUR\"}",
		"{\"bank_account\":{\"account_number\":\"SE4550000000058398257466\", \"type\":\"iban\"}, \"country_code\":\"SE\", \"currency\":\"SEK\"}",
	}
	units, err := mutiny.Riot(payload)
	suite.NoError(err)
	suite.Len(units, 2)
	for _, unit := range units {
		var match bool
		for index, expected := range expectedResults {
			if JSONEq(suite, expected, string(unit.RequestBody)) {
				match = true
				expectedResults = append(expectedResults[:index], expectedResults[index+1:]...)
				break
			}
		}

		suite.Truef(match, "Unexpected request body: %s", string(unit.RequestBody))
	}
	suite.Zerof(len(expectedResults), "expected test to consume all expected results")

	expectedResults = []string{
		"{\"bank_account\":{\"account_number\":\"SE4550000000058398257466\", \"type\":\"iban\"}, \"country_code\":\"&nbsp;\", \"currency\":\"EUR\"}",
		"{\"bank_account\":{\"account_number\":\"SE4550000000058398257466\", \"type\":\"iban\"}, \"country_code\":true, \"currency\":\"EUR\"}",
		"{\"bank_account\":{\"account_number\":\"SE4550000000058398257466\", \"type\":\"iban\"}, \"country_code\":123, \"currency\":\"EUR\"}",
		"{\"bank_account\":{\"account_number\":\"SE4550000000058398257466\", \"type\":\"iban\"}, \"country_code\":\"&nbsp;\", \"currency\":\"SEK\"}",
		"{\"bank_account\":{\"account_number\":\"SE4550000000058398257466\", \"type\":\"iban\"}, \"country_code\":true, \"currency\":\"SEK\"}",
		"{\"bank_account\":{\"account_number\":\"SE4550000000058398257466\", \"type\":\"iban\"}, \"country_code\":123, \"currency\":\"SEK\"}",
	}

	units, err = mutiny.Riot(payload.SetCountryCode(mutiny.Erroneous))
	suite.NoError(err)
	suite.Len(units, 6)
	for _, unit := range units {
		var match bool
		for index, expected := range expectedResults {
			if JSONEq(suite, expected, string(unit.RequestBody)) {
				match = true
				expectedResults = append(expectedResults[:index], expectedResults[index+1:]...)
				break
			}
		}

		suite.Truef(match, "Unexpected request body: %s", string(unit.RequestBody))
	}
	suite.Zerof(len(expectedResults), "expected test to consume all expected results")
}

func TestMutinyTestSuite(t *testing.T) {
	suite.Run(t, new(MutinyTestSuite))
}
