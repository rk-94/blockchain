package insertion
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

//Define the structure of the data in database
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

/*will insert the new record into the database
 args- array of string carrying data
 returns the url to fetch the hash key
*/
func InsertData(stub shim.ChaincodeStubInterface, args[]string) (string){
	if len(args) != 10 {
		panic(fmt.Sprintf("incorrect arguments"))
	}
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
	url := fmt.Sprintf("http://rdctstbc001:5984/mychannel_mycc/%s",id)
	return url	
}