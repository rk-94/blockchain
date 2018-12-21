package main
/* Imports
 * 1 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
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
func insertData(stub shim.ChaincodeStubInterface, args[]string) (string, string){
	
	fhir := "https://fhir-open.sandboxcerner.com/dstu2/0b8a0111-e8e6-4c26-a91c-5069cbc6b1ca/Procedure?patient="+args[1]
	var Payerdata = &Payer{ClaimId: args[0], FhirUrl: fhir, PatientId: args[1], PayerId: args[2], SubmitterId: args[3],}
	var response string
	payerData, err := stub.GetState(args[0])
	if err != nil {
		response = fmt.Sprintf("Failed to get payer data: " + err.Error())
		return "error", response
	} else if payerData != nil {
		response = fmt.Sprintf("This claim id already exists: " + args[0])
		return "error", response
	}

	PayerAsBytes, err := json.Marshal(Payerdata)
	if err != nil {
		return "error", "oops there is a problem in marshalling"
	}
	insertErr := stub.PutState(args[0], PayerAsBytes)
	if insertErr!= nil {
		panic(insertErr.Error())
	}
	
	return "success", args[0]	
}