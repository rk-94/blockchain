package main 
/* Imports
 * 2 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */ 
 import (
	"fmt"
	"bytes"
        "github.com/hyperledger/fabric/core/chaincode/shim"
        pb "github.com/hyperledger/fabric/protos/peer"
       ) 

// Define the Smart Contract structure 
type SimpleChaincode struct {
}

//Define the data structure
type Payer struct {
	ClaimId		string `json:"claimId"`
	FhirUrl    	string `json:"fhirUrl"`
	PatientId   string `json:"patientId"`
	PayerId		string `json:"payerId"`
	SubmitterId	string `json:"submitterId"`
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
	}
	
	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

//will do the validation of hash and payerId
func (t *SimpleChaincode) isValid(stub shim.ChaincodeStubInterface, args[]string) pb.Response {
	
	var buffer bytes.Buffer
	if args[0] == "payerhash" {
		checkPayerWithHash, err := CheckPayerWithHash(stub, args[1], args[2])
		if err != nil{
			return shim.Error(err.Error())
		}
		buffer.WriteString(string(checkPayerWithHash[:]))
		if len(checkPayerWithHash) == 0{
			return shim.Error("Invalid Payer with Hash")
		}
	}else if args[0] == "hash" {
		hash := args[1]
		checkHash, errr := CheckHashKey(stub, hash)
		if errr != ""{
			return shim.Error(errr)
		}
		if !checkHash{
			return shim.Error("Invalid Hash")
		}
	
	}else if args[0] == "payer" {
		payerId := args[1]
		checkPayer, errr := CheckPayerId(stub, payerId)
		if errr != ""{
			return shim.Error(errr)
		}
		if !checkPayer{
			return shim.Error("Invalid Payer")
		}
	}else {
		return shim.Error("Invalid arguments")
	}
	return shim.Success(nil)		
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
	encounter_id := insertData(stub,args)
	buffer.WriteString(encounter_id)
		
	return shim.Success(buffer.Bytes())
}