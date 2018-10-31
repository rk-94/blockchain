package main
/* Imports
 * 2 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 1 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"fmt"
	"bytes"
    	"github.com/hyperledger/fabric/core/chaincode/shim"
 )

/*checks whether the provided hash is valid or not
 args - array of string having hash and payerid as values
 returns array of byte and error
*/
func CheckPayerWithHash(stub shim.ChaincodeStubInterface, args[] string) ([]byte, error){
	
	queryString := fmt.Sprintf("{\"selector\":{\"_rev\":\"%s\", \"payerId\":\"%s\"}}", args[0], args[1])
	
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
