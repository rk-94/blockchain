package main 
/* Imports
 * 6 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 4 specific Hyperledger Fabric specific libraries for Smart Contracts
 */ 
 import (
	"fmt"
	"bytes"
	"strconv"
	//"net/http"
	//"io/ioutil"
	"encoding/json"
        "github.com/hyperledger/fabric/core/chaincode/shim"
        pb "github.com/hyperledger/fabric/protos/peer"
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
                stub.PutState("PATIENT"+strconv.Itoa((i*i*i*i*i)+5678912340987345271), patAsBytes)
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
	
	var buffer bytes.Buffer

	payerId := args[1]
	hash := args[0]
			
	checkPayer, errr := CheckPayerId(stub, payerId)
	if errr != ""{
		return shim.Error(errr)
	}
	if !checkPayer{
		return shim.Error("Invalid Payer")
	}
	
	checkHash, errr := CheckHashKey(stub, hash)
	if errr != ""{
		return shim.Error(errr)
	}
	if !checkHash{
		return shim.Error("Invalid Hash")
	}

	checkPayerWithHash, err := CheckPayerWithHash(stub, args)
	if err != nil{
		return shim.Error(err.Error())
	}
	buffer.WriteString(string(checkPayerWithHash[:]))
	if len(checkPayerWithHash) == 0{
		return shim.Error("Invalid Payer with Hash")
	}
	
	return shim.Success(buffer.Bytes())

}

//calls QueryByHash function from query package
func (t *SimpleChaincode) queryByHash(stub shim.ChaincodeStubInterface, args[]string) pb.Response {
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	resultsIterator, err := QueryByHash(stub, args)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(resultsIterator)
}

//calls QueryOnFilter function from query package
func (t *SimpleChaincode) queryCustom(stub shim.ChaincodeStubInterface, args[]string) pb.Response {
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}	
	resultsIterator, err := QueryOnFilter(stub, args)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(resultsIterator)
}

//call insert function to insert new record into the database
func (t *SimpleChaincode) insertData(stub shim.ChaincodeStubInterface, args[]string) pb.Response {
	
	var buffer bytes.Buffer

	if len(args) != 10 {
		return shim.Error("Incorrect number of arguments. Expecting 10")
	}

	url_ref.url = InsertData(stub,args)
	fmt.Println(url_ref.url)

	response := shim.Success(nil)
	buffer.WriteString("Reached here")

	fmt.Println(response)
	
	/*status := response.GetStatus()

	fmt.Println(status)
	if status == 200 {
		url := url_ref.url
		var dat map[string]interface{}
		resp, err := http.Get(url)

		if err != nil {
		panic(err)
		}
		defer resp.Body.Close()
	
		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
            		panic(err.Error())
		}
		err1 := json.Unmarshal(contents, &dat)
		if err1 != nil{
			panic(err.Error())
		}
	
		fmt.Println(dat["_rev"])
    		fmt.Printf("%s\n", string(contents))
	
		hash = fmt.Sprintf("The Hash for the added patient is: %s", dat["_rev"].(string))
	}else {
		return shim.Error("data not inserted")
	}*/
		
	return shim.Success(nil)
}
