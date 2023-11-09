// This is the implementation of a test of the create method
// Date: 8 November 2023
//------------------------------------------------------------------

package main

import (
	"bytes"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Constants
const profilekey = "595edb38-4b39-42b7-81ee-be458a198024"
const secretkey = "76198df7-296f-4939-82b6-6ce3c6699bc5"
const post_url = "https://fumopay.dev/transaction/pay"

// Variables, structs

// Struct used to hold transaction info for processing
type transaction_info struct {
	profile_key   string
	currency      string // assuming string type, EMV and others use BCD
	description   string
	ext_reference string
	ext_customer  string
	amount        int // assuming integer type for this field... not sure
	age           int
	redirect_url  string
	signature     string // base64 in here

}

// Struct used for HTTPS POST
type Post struct {
	MerchantProfile MerchantProfile `json:"merchant_profile"`
	Currency        string          `json:"currency"`
	Description     string          `json:"description"`
	Ext_reference   string          `json:"ext_reference"`
	Ext_customer    string          `json:"ext_customer"`
	Amount          int             `json:"amount"`
	Age             int             `json:"age"`
	Redirect_url    string          `json:"redirect_url"`
	Signature       string          `json:"signature"`
}

type MerchantProfile struct {
	ProfileKey string `json:"profile_key"`
}

var transaction_value float64 = 0
var TransactionInfo = transaction_info{}

// -------------------------------------------------------------------
// Name: SignTxn
// Function: Compute the signature for this transaction
// Parameters: struct with transaction info, secret key
// Returns: base64 encoded signature string
// --------------------------------------------------------------------
func SignTxn(txninfo transaction_info, secretkey string) string {

	stringToHash := fmt.Sprintf("%s%s%s%s%s",
		txninfo.currency, txninfo.description, txninfo.profile_key,
		txninfo.ext_reference, secretkey)

	txnhash := sha512.New()
	txnhash.Write([]byte(stringToHash))
	signature := base64.StdEncoding.EncodeToString(txnhash.Sum(nil))
	return signature
}

// ---------------------------------------------------------------------
//
//	Main Function
//
// ---------------------------------------------------------------------
func main() {

	// Let's populate a transaction value
	transaction_value = 5.00 // £5.00

	// Populate the struct
	TransactionInfo.amount = int(transaction_value * 100)
	TransactionInfo.description = "cash sale"
	TransactionInfo.ext_reference = "ORD00012"
	TransactionInfo.currency = "GBP"
	TransactionInfo.ext_customer = "L. Jones"
	TransactionInfo.age = 1
	TransactionInfo.redirect_url = ""
	TransactionInfo.profile_key = profilekey
	TransactionInfo.signature = SignTxn(TransactionInfo, secretkey)

	// Create JSON payload for HTTPS request

	Payload := Post{
		MerchantProfile: MerchantProfile{
			ProfileKey: profilekey,
		},
		Currency:      TransactionInfo.currency,
		Description:   TransactionInfo.description,
		Ext_reference: TransactionInfo.ext_reference,
		Ext_customer:  TransactionInfo.ext_customer,
		Amount:        TransactionInfo.amount,
		Age:           TransactionInfo.age,
		Redirect_url:  TransactionInfo.redirect_url,
		Signature:     TransactionInfo.signature,
	}

	payload, err := json.Marshal(Payload)
	if err != nil {
		panic(err)
	}
	//fmt.Println("JSON: " + string(payload))

	resp, err := http.Post(post_url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}

	// Okay so far so good, now we need to handle the response
	defer resp.Body.Close()

	// Check the response from the host, according to the API documentation if
	// the transaction is successful, the status is "1" but more importantly a
	// transaction ID is returned along with a fumopay URL

	fmt.Println(fmt.Sprintf("HTTP Response Code: %d", resp.StatusCode))
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Unable to read body of response")
	}
	rawJSON := string(body)
	fmt.Println("Response Body: ", rawJSON)
}
