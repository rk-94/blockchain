package main
/* Imports
 * 3 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 1 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import(
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
	"fmt"
	"math/rand"
)

/*will insert the new record into the database
 args- array of string carrying data
 returns the url to fetch the hash key
*/
func InsertData(stub shim.ChaincodeStubInterface, args[]string) (string){
	i := rand.Int()
	id := fmt.Sprintf("%s%d","PATIENT",i)
	
	var Payerdata = Payer{PayerId: args[0], PatReference: args[1], PatDisplay: args[2], SourceReference: args[3], SourceDisplay: args[4], DateAsserted: args[5], Status: args[6], WasTaken: args[7], PeriodStart: args[8], ClientId: args[9]}

	PayerAsBytes, _ := json.Marshal(Payerdata)
	insertErr := stub.PutState(id, PayerAsBytes)
	if insertErr!= nil {
		panic(insertErr.Error())
	}else{
		fmt.Println("Inserted")
	}
	url := fmt.Sprintf("http://rdctstbc001:5984/mychannel_myc/%s",id)
	return url	
}