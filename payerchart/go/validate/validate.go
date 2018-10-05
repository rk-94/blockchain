package validate

import (
	"fmt"
	"bytes"
    	"github.com/hyperledger/fabric/core/chaincode/shim"
    //pb "github.com/hyperledger/fabric/protos/peer"
     
	 )

func CheckHashKey(stub shim.ChaincodeStubInterface, arg string) (bool, string) {

	hash := arg
	getHash, err := stub.GetState(hash)
	
	invalid := fmt.Sprintf("Invalid Hash Key %s",string(getHash))
	
	if err != nil {
		return false, invalid
	}else if getHash == nil {
		return false, invalid
	}

	return true, ""
	/*resultsIterator, err := getQueryResultForQueryString(stub, arg)
	if err != nil {
		return nil, err
	}
	return resultsIterator, err*/
}

func CheckPayerId(stub shim.ChaincodeStubInterface, arg string) (bool, string) {

	payerId := arg
	getId, err := stub.GetState(payerId)

	invalid := fmt.Sprintf("Invalid PayerId %s",string(getId))
	
	if err != nil {
		return false, invalid
	}else if getId == nil {
		return false, invalid
	}

	return true, ""
}

func CheckPayerWithHash(stub shim.ChaincodeStubInterface, args[] string) (bool, error){
	if len(args) < 2 {
		return false, nil
	}
	payerId := args[0]
	hash := args[1]

	//invalid := fmt.Sprintf("Invalid PayerId with hash")
	
	queryString := fmt.Sprintf("{\"selector\":{\"payerid\":\"%s\", \"pathash\":\"%s\"}}", payerId, hash)

	getPayerWithHash, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return false, err
	}
	if len(getPayerWithHash) <= 0 {
		return false, err
	}	
	
	return true, nil
	
}

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