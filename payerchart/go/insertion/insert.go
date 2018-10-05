package insertion

import(
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
	"bytes"
)

type Payer struct {
	Clientid        string `json:"clientid"`
	DateAsserted    string `json:"dateAsserted"`
	Patdisplay      string `json:"patdisplay"`
	Pathash         string `json:"pathash"`
	Patreference    string `json:"patreference"`
	Payerid         string `json:"Payerid"`
	Periodstart     string `json:"periodstart"`
	Sourcedisplay   string `json:"sourcedisplay"`
	Sourcereference string `json:"sourcereference"`
	Status          string `json:"status"`
	WasTaken        string `json:"wasTaken"`
}

func InsertData(stub shim.ChaincodeStubInterface, args[]string) (int, error) {
	if len(args) != 12 {
		return 0, nil
	}
	
	var buffer bytes.Buffer
        buffer.WriteString("")

	var Payerdata = Payer{Pathash:args[1], Payerid:args[2], Patreference: args[3], Patdisplay: args[4], Sourcereference: args[5], Sourcedisplay: args[6], DateAsserted: args[7], Status: args[8], WasTaken: args[9], Periodstart: args[10], Clientid: args[11]}

	PayerAsBytes, _ := json.Marshal(Payerdata)
	insertErr := stub.PutState(args[0], PayerAsBytes)
	if insertErr!= nil {
		return 0, insertErr
	}
	return 1, nil
}