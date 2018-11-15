package main
/* Imports
 * 3 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 1 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import(
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
	"fmt"
)

/*will insert the new record into the database
 args- array of string carrying data
 returns the url to fetch the hash key
*/
func insertData(stub shim.ChaincodeStubInterface, args[]string) (string){
	
	fhir := args[3]+args[1]
	var Payerdata = &Payer{PatientId: args[1], PayerId: args[2], FhirUrl: fhir}

	PayerAsBytes, err := json.Marshal(Payerdata)
	if err != nil {
		return "oops there is a problem in marshalling"
	}
	insertErr := stub.PutState(args[0], PayerAsBytes)
	if insertErr!= nil {
		panic(insertErr.Error())
	}
	url := fmt.Sprintf("http://rdctstbc001:5984/mychannel_mycc/%s",args[0])
	
	return url	
}