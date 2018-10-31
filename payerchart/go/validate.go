package main
/* Imports
 * 2 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 1 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"fmt"
    	"github.com/hyperledger/fabric/core/chaincode/shim"
 )

/*checks whether the provided hash is valid or not
 arg - hash value
 returns bool value with the message
*/
func CheckHashKey(stub shim.ChaincodeStubInterface, arg string) (bool, string) {

	hash := arg
	
	queryString := fmt.Sprintf("{\"selector\":{\"_rev\":\"%s\"}}", hash)

	getHash, err := stub.GetQueryResult(queryString)

	invalid := fmt.Sprintf("Invalid Hash Key")
	
	if err != nil {
		return false, invalid
	}
	if !getHash.HasNext() {
		return false, invalid
	}

	return true, ""
}
/*checks whether the provided hash is valid or not
 arg - payerid value
 returns bool value with the message
*/
func CheckPayerId(stub shim.ChaincodeStubInterface, arg string) (bool, string) {

	payerId := arg
	
	queryString := fmt.Sprintf("{\"selector\":{\"payerId\":\"%s\"}}", payerId)

	getId, err := stub.GetQueryResult(queryString)
	
	invalid := fmt.Sprintf("Invalid PayerId")
	
	if err != nil {
		return false, invalid
	}
	if !getId.HasNext() {
		return false, invalid
	}

	return true, ""
}
