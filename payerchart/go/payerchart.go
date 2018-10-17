package main 
/* Imports
 * 6 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 4 specific Hyperledger Fabric specific libraries for Smart Contracts
 */ 
 import (
	"fmt"
	"bytes"
	"strconv"
	"net/http"
	"io/ioutil"
	"encoding/json"
        "github.com/hyperledger/fabric/core/chaincode/shim"
        pb "github.com/hyperledger/fabric/protos/peer"
        pv "github.com/chaincode/payerchart/go/validate"
	iv "github.com/chaincode/payerchart/go/insertion"
) 

// Define the Smart Contract structure 
type SimpleChaincode struct {
}

//Define the structure for url
type Url struct {
	url string
}

var url_ref Url

//Define the data structure
type Payer struct {
	ClientId        string `json:"clientId"`
	DateAsserted    string `json:"dateAsserted"`
	PatDisplay      string `json:"patDisplay"`
	PatReference    string `json:"patReference"`
	PayerId         string `json:"payerId"`
	PeriodStart     string `json:"periodStart"`
	SourceDisplay   string `json:"sourceDisplay"`
	SourceReference string `json:"sourceReference"`
	Status          string `json:"status"`
	WasTaken        string `json:"wasTaken"`
}
// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init initializes chaincode
// ===========================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}
// Invoke - Our entry point for Invocations
// ========================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function,args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "initLedger" { //read the data from the json file
		return t.initLedger(stub)
	}else if function == "isValid"{
		return t.isValid(stub,args)
	}else if function == "queryByHash"{
		return t.queryByHash(stub,args)
	}else if function == "queryCustom"{
		return t.queryCustom(stub,args)
	}else if function == "insertData"{
		return t.insertData(stub,args)
	}else if function == "retHash"{
		return t.retHash(stub)
	}
	
	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

//initLedger - populate the database with the data
func (t *SimpleChaincode) initLedger(stub shim.ChaincodeStubInterface) pb.Response {
	patients := []Payer{
		Payer{PayerId: "CI08128",PatReference: "Patient/4342010",PatDisplay: "Smart, Joe", SourceReference: "Practitioner/1912007",SourceDisplay: "who, Doctor",DateAsserted: "2016-06-27T09:57:32.000-05:00",Status: "active",WasTaken: "false",PeriodStart: "2016-06-27T09:00:00.000-07:00",ClientId:"CLM0098"},
		Payer{PayerId: "AT09562",PatReference: "Patient/4342011",PatDisplay: "Himilton, Jack",SourceReference: "Practitioner/1912007",SourceDisplay: "Who, Doctor",DateAsserted: "2017-06-27T09:57:32.000-05:00",Status: "active",WasTaken: "true",PeriodStart: "2017-06-27T09:00:00.000-07:00",ClientId:"CLM0097"},
		Payer{PayerId: "CI02129",PatReference: "Patient/4343019",PatDisplay: "Hilton, Mariya",SourceReference: "Practitioner/1912007",SourceDisplay: "Who, Doctor",DateAsserted: "2015-06-27T09:17:32.000-05:00",Status: "active",WasTaken: "true",PeriodStart: "2015-06-27T09:17:00.000-07:00",ClientId:"CLM0045"}}
	i := 0
	for i < len(patients) {
		fmt.Println("i is ", i)
                patAsBytes, _ := json.Marshal(patients[i])
                stub.PutState("PATIENT"+strconv.Itoa(i+1), patAsBytes)
                fmt.Println("Added", patients[i])
                i = i + 1
	}
	return shim.Success(nil)	
}

//will do the validation of hash and payerId
func (t *SimpleChaincode) isValid(stub shim.ChaincodeStubInterface, args[]string) pb.Response {
	
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	hash := args[0]
	payerId := args[1]
		
	checkHash, errr := pv.CheckHashKey(stub, hash)
	if errr != ""{
		return shim.Error(errr)
	}
	if !checkHash{
		return shim.Error("Invalid Hash")
	} 

	checkPayer, errr := pv.CheckPayerId(stub, payerId)
	if errr != ""{
		return shim.Error(errr)
	}
	if !checkPayer{
		return shim.Error("Invalid Payer")
	}

	checkPayerWithHash, err := pv.CheckPayerWithHash(stub, args)
	if err != nil{
		return shim.Error(err.Error())
	}
	if checkPayerWithHash == nil{
		return shim.Error("Invalid Payer with Hash")
	}
	
	return shim.Success(nil)

}

//to fetch the data from database by hash
func (t *SimpleChaincode) queryByHash(stub shim.ChaincodeStubInterface, args[]string) pb.Response {
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	
	queryString := fmt.Sprintf("{\"selector\":{\"_rev\":\"%s\"}}", args[0])
	resultsIterator, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(resultsIterator)
}

//will query the database by giving custom conditions
func (t *SimpleChaincode) queryCustom(stub shim.ChaincodeStubInterface, args[]string) pb.Response {
	
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
		
	queryString := args[0]
	resultsIterator, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(resultsIterator)
}

//responsible for contacting the database and return the data
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}		
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", resultsIterator)

	return buffer.Bytes(), nil
}

//call insert function to insert new record into the database
func (t *SimpleChaincode) insertData(stub shim.ChaincodeStubInterface, args[]string) pb.Response{

	url_ref.url = iv.InsertData(stub,args)
	fmt.Println(url_ref.url)

	return shim.Success(nil)	
}

//will return the hash value of the recently added record
func (t *SimpleChaincode) retHash(stub shim.ChaincodeStubInterface) pb.Response{

	url := url_ref.url
	var dat map[string]interface{}
	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
            return shim.Error(err.Error())
	}
	err1 := json.Unmarshal(contents, &dat)
	if err1 != nil{
		return shim.Error(err.Error())
	}
	
	fmt.Println(dat["_rev"])
    	fmt.Printf("%s\n", string(contents))
	
	response := fmt.Sprintf("The Hash for the added patient is: %s", dat["_rev"].(string))

	return shim.Success([]byte(response))
}