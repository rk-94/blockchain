package main 
/* Imports
 * 3 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */ 
 import (
	"fmt"
	"encoding/json"
	"bytes"
	//"net/http"
	"math/rand"
    "github.com/hyperledger/fabric/core/chaincode/shim"
     pb "github.com/hyperledger/fabric/protos/peer"
    ) 

// Define the Smart Contract structure 
type SimpleChaincode struct {
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
		return t.initLedger(stub,args)
	}else if function == "insertData"{
		return t.insertData(stub,args)
	}else if function == "queryContext"{
		return t.queryContext(stub,args)
	}else if function == "queryCustom"{
		return t.queryCustom(stub,args)
	}
	
	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

//initLedger - Adds the first block in the blockchain
func (t *SimpleChaincode) initLedger(stub shim.ChaincodeStubInterface, args[]string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	i := rand.Int()
	var buffer bytes.Buffer
	id := fmt.Sprintf("%s%d","DATA",i)
	var data interface{}
	err := json.Unmarshal([]byte(args[0]), &data)
	if err != nil {
   		return shim.Error(err.Error())
	} 
	dataInBytes, err := json.Marshal(data)
	if err != nil {
		return shim.Error(err.Error())
	 }
	insertErr := stub.PutState(id, dataInBytes)
	if insertErr!= nil {
		return shim.Error(insertErr.Error())
	}
	buffer.WriteString(id)
	return shim.Success(buffer.Bytes())	
}

func validUser(stub shim.ChaincodeStubInterface, username string, password string, context string) (bool, string){
	queryString := fmt.Sprintf("{\"selector\":{\"username\":\"%s\",\"password\":\"%s\",\"context\":\"%s\"}}", username, password, context)
	resultsIterator, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return false, "Error while sending query"
	}
	if len(resultsIterator) == 0 {
		return false, "Invalid User"
	}
	return true, ""

}

//insertData - inserts the data in the database
func (t *SimpleChaincode) insertData(stub shim.ChaincodeStubInterface, args[]string) pb.Response {

	if args[0] == "check" {
		if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4- {\"check\",\"username\",\"password\",\"context\"}")
		}
		resp, message := validUser(stub,args[1],args[2],args[3])
		if !resp {
			return shim.Error(message)
		}else {
			return shim.Success(nil)
		}
	}
	
	i := rand.Int()
	id := fmt.Sprintf("%s%d","DATA",i)
	var data interface{}
	err := json.Unmarshal([]byte(args[0]), &data)
	if err != nil {
   		return shim.Error(err.Error())
	} 
	dataInBytes, err := json.Marshal(data)
	if err != nil {
		return shim.Error(err.Error())
	 }
	insertErr := stub.PutState(id, dataInBytes)
	if insertErr!= nil {
		return shim.Error(insertErr.Error())
	}
	return shim.Success(nil)	
}

//to fetch the data from database by hash
func (t *SimpleChaincode) queryContext(stub shim.ChaincodeStubInterface, args[]string) pb.Response {
	
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3(context, username,password)")
	}	
	
	queryString := fmt.Sprintf("{\"selector\":{\"context\":\"%s\",\"username\":\"%s\",\"password\":\"%s\"}}", args[0], args[1], args[2])
	resultsIterator, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(resultsIterator)
}

func (t *SimpleChaincode) queryCustom(stub shim.ChaincodeStubInterface, args[]string) pb.Response {
	
	if len(args) != 1 {
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

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", resultsIterator)

	return buffer.Bytes(), nil
}
