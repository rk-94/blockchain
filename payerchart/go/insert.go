package main
/* Imports
 * 1 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 1 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import(
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
	"fmt"
	"bytes"
)

/*will insert the new record into the database
 args- array of string carrying data
 returns the url to fetch the hash key
*/
func insertData(stub shim.ChaincodeStubInterface, args[]string) (string, string){

	var fhirUrls []string
	var data Data

	queryString := fmt.Sprintf("{\"selector\":{\"submitterId\":\"%s\",\"payerId\":\"%s\"}}", args[3], args[2])
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return "error", "There's no such record"
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer

	queryResponse, err := resultsIterator.Next()
	if err != nil {
		return "error", "There's no such record"
	}
	buffer.WriteString(string(queryResponse.Value))
	jsonData := buffer.Bytes()
  	if jsonData == nil{
		return "error", "No Data found"
	}
	json.Unmarshal(jsonData, &data)
	if data.SubmitterIndicator == 0{
		return "error", "Unsubscribed"
	}
	for i:=0 ; i<len(data.Fhir); i++{
		if data.Fhir[i]=="Procedures"{
			fhirUrls = append(fhirUrls, fmt.Sprintf("https://fhir-open.sandboxcerner.com/dstu2/0b8a0111-e8e6-4c26-a91c-5069cbc6b1ca/%s?patient=%s",data.Fhir[i],args[1]))
		}
		if data.Fhir[i]=="Demographics"{
			fhirUrls = append(fhirUrls, fmt.Sprintf("https://fhir-open.sandboxcerner.com/dstu2/0b8a0111-e8e6-4c26-a91c-5069cbc6b1ca/%s?patient=%s",data.Fhir[i],args[1]))
		}
		if data.Fhir[i]=="DocumentRef"{
			fhirUrls = append(fhirUrls, fmt.Sprintf("https://fhir-open.sandboxcerner.com/dstu2/0b8a0111-e8e6-4c26-a91c-5069cbc6b1ca/%s?patient=%s",data.Fhir[i],args[1]))
		}
		if data.Fhir[i]=="DiagnoticReport"{
			fhirUrls = append(fhirUrls, fmt.Sprintf("https://fhir-open.sandboxcerner.com/dstu2/0b8a0111-e8e6-4c26-a91c-5069cbc6b1ca/%s?patient=%s",data.Fhir[i],args[1]))
		}
		if data.Fhir[i]=="Observations"{
			fhirUrls = append(fhirUrls, fmt.Sprintf("https://fhir-open.sandboxcerner.com/dstu2/0b8a0111-e8e6-4c26-a91c-5069cbc6b1ca/%s?patient=%s",data.Fhir[i],args[1]))
		}
		if data.Fhir[i]=="MedicationOrder"{
			fhirUrls = append(fhirUrls, fmt.Sprintf("https://fhir-open.sandboxcerner.com/dstu2/0b8a0111-e8e6-4c26-a91c-5069cbc6b1ca/%s?patient=%s",data.Fhir[i],args[1]))
		}
	} 

	var Payerdata = &Payer{ClaimId: args[0], FhirUrl: fhirUrls,  PatientId: args[1], PayerId: args[2], SubmitterId: args[3],}
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