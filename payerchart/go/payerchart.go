package main 
/* Imports
 * 6 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 4 specific Hyperledger Fabric specific libraries for Smart Contracts
 */ 
 import (
	"fmt"
	"bytes"
	"net/http"
	"io/ioutil"
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
	PatientId    	string `json:"patientId"`
	PayerId		string `json:"payerId"`
	FhirUrl    	string `json:"fhirUrl"`
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
	if function == "isValid"{
		return t.isValid(stub,args)
	}else if function == "queryByHash"{
		return t.queryByHash(stub,args)
	}else if function == "queryCustom"{
		return t.queryCustom(stub,args)
	}else if function == "insert"{
		return t.insert(stub,args)
	}else if function == "retHash"{
		return t.retHash(stub)
	}
	
	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
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
func (t *SimpleChaincode) insert(stub shim.ChaincodeStubInterface, args[]string) pb.Response {
	
	var buffer bytes.Buffer

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	url_ref.url = insertData(stub,args)
	buffer.WriteString(url_ref.url)
		
	return shim.Success(buffer.Bytes())
}

func (t *SimpleChaincode) retHash(stub shim.ChaincodeStubInterface) pb.Response {
	var buffer bytes.Buffer
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
    	buffer.WriteString(string(contents))
	
	hash := fmt.Sprintf("The Hash for the added patient is: %s", dat["_rev"].(string))
	buffer.WriteString(hash)
	
	return shim.Success(buffer.Bytes())
}
