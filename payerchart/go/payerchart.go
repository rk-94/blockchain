package main 
/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 3 specific Hyperledger Fabric specific libraries for Smart Contracts
 */ 
 import (
	"fmt"
	"bytes"
	"strconv"
	"net/http"
	"io/ioutil"
	"encoding/json"
        "github.com/hyperledger/fabric/core/chaincode/shim"
        pb "github.com/hyperledger/fabric/protos/peer"
        pv "github.com/chaincode/payerchart/go/validate"
	iv "github.com/chaincode/payerchart/go/insertion"
	 ) // Define the Smart Contract structure 

type SimpleChaincode struct {
}

type Payer struct {
	Clientid        string `json:"clientid"`
	DateAsserted    string `json:"dateAsserted"`
	Patdisplay      string `json:"patdisplay"`
	Pathash         string `json:"pathash"`
	Patreference    string `json:"patreference"`
	Payerid         string `json:"payerid"`
	Periodstart     string `json:"periodstart"`
	Sourcedisplay   string `json:"sourcedisplay"`
	Sourcereference string `json:"sourcereference"`
	Status          string `json:"status"`
	WasTaken        string `json:"wasTaken"`
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
		return t.initLedger(stub)
	}else if function == "isValid"{
		return t.isValid(stub)
	//}else if function == "validPayer"{
	//	return t.validPayer(stub,arg)
	}else if function == "queryByHash"{
		return t.queryByHash(stub,args)
	}else if function == "queryCustom"{
		return t.queryCustom(stub,args)
	}else if function == "updatePatdisplay"{
		return t.updatePatdisplay(stub,args)
	}else if function == "updateSourceRef"{
		return t.updateSourceRef(stub,args)
	}else if function == "updateSourcedisplay"{
		return t.updateSourcedisplay(stub,args)
	}else if function == "updateStatus"{
		return t.updateStatus(stub,args)
	}else if function == "insertData"{
		return t.insertData(stub,args)
	}
	
	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

func (t *SimpleChaincode) initLedger(stub shim.ChaincodeStubInterface) pb.Response {
	patients := []Payer{
		Payer{Pathash: "1-1c2fae390fa5475d9b809301bbf3f25e",Payerid: "CI08128",Patreference: "Patient/4342010",Patdisplay: "Smart, Joe", Sourcereference: "Practitioner/1912007",Sourcedisplay: "who, Doctor",DateAsserted: "2016-06-27T09:57:32.000-05:00",Status: "active",WasTaken: "false",Periodstart: "2016-06-27T09:00:00.000-07:00",Clientid:"CLM0098"},
		Payer{Pathash: "2-04d8eac1680d237ca25b68b36b8899d3",Payerid: "AT09562",Patreference: "Patient/4342011",Patdisplay: "Himilton, Jack",Sourcereference: "Practitioner/1912007",Sourcedisplay: "Who, Doctor",DateAsserted: "2017-06-27T09:57:32.000-05:00",Status: "active",WasTaken: "true",Periodstart: "2017-06-27T09:00:00.000-07:00",Clientid:"CLM0097"}}
	i := 0
	for i < len(patients) {
		fmt.Println("i is ", i)
                patAsBytes, _ := json.Marshal(patients[i])
                stub.PutState("PATIENT"+strconv.Itoa(i), patAsBytes)
                fmt.Println("Added", patients[i])
                i = i + 1
	}
	/*validhash, err := pv.CheckHashKey(stub,args[0])
	if err != ""{
		return shim.Error(err)
	}
	
	validPayer, err := pv.CheckPayerId(stub,args[1])
	if err != ""{
		return shim.Error(err)
	}

	if !validhash{
		return shim.Error("Invalid Hash")
	}
	
	if !validPayer{
		return shim.Error("Invalid Payer")
	}*/
	return shim.Success(nil)	
}

func (t *SimpleChaincode) queryByHash(stub shim.ChaincodeStubInterface, args[]string) pb.Response {
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	
	queryString := fmt.Sprintf("{\"selector\":{\"_rev\":\"%s\"}}", args[0])
	resultsIterator, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(resultsIterator)
}

func (t *SimpleChaincode) isValid(stub shim.ChaincodeStubInterface) pb.Response {
	
	hash := "1-1c2fae390fa5475d9b809301bbf3f25e"
	Payerid := "CI08128"
	
	var str[]string
	str[0] = Payerid
	str[1] = hash
	
	/*checkHash, err := pv.CheckHashKey(stub, str[1])
	if errr != nil{
		return shim.Error(err)
	}
	if !checkHash{
		return shim.Error("Invalid Hash")
	} 

	checkPayer, errr := pv.CheckPayerId(stub, str[0])
	if errr != ""{
		return shim.Error(errr)
	}
	if !checkPayer{
		return shim.Error("Invalid Payer")
	}*/

	checkPayerWithHash, err := pv.CheckPayerWithHash(stub, str)
	if err != nil{
		return shim.Error(err.Error())
	}
	if !checkPayerWithHash{
		return shim.Error("Invalid Payer with Hash")
	}
	return shim.Success(nil)
}


func (t *SimpleChaincode) queryCustom(stub shim.ChaincodeStubInterface, args[]string) pb.Response {
	
	/*hash := "1-1c2fae390fa5475d9b809301bbf3f25e"
	Payerid := "CI08128"
	
	var str[]string
	str[0] = Payerid
	str[1] = hash
	
	checkHash, errr := pv.CheckHashKey(stub, str[1])
	if errr != ""{
		return shim.Error(errr)
	}
	if !checkHash{
		return shim.Error("Invalid Hash")
	} 

	checkPayer, errr := pv.CheckPayerId(stub, str[0])
	if errr != ""{
		return shim.Error(errr)
	}
	if !checkPayer{
		return shim.Error("Invalid Payer")
	}

	checkPayerWithHash, errr := pv.CheckPayerWithHash(stub, str)
	if errr != ""{
		return shim.Error(errr)
	}
	if !checkPayerWithHash{
		return shim.Error("Invalid Payer with Hash")
	}*/

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
		
	queryString := args[0]
	resultsIterator, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(resultsIterator)
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

func (t *SimpleChaincode) insertData(stub shim.ChaincodeStubInterface, args[]string) pb.Response{
	/*if len(args) != 12 {
		return shim.Error("Incorrect number of arguments. Expecting 12")
	}
	
	var buffer bytes.Buffer
        buffer.WriteString("")

	var Payerdata = Payer{Pathash:args[1], Payerid:args[2], Patreference: args[3], Patdisplay: args[4], Sourcereference: args[5], Sourcedisplay: args[6], DateAsserted: args[7], Status: args[8], WasTaken: args[9], Periodstart: args[10], Clientid: args[11]}

	PayerAsBytes, _ := json.Marshal(Payerdata)
	insertErr := stub.PutState(args[0], PayerAsBytes)
	if insertErr!= nil {
		return shim.Error(insertErr.Error())
	}
	buffer.WriteString("The Hash for the added patient is: ")
	buffer.WriteString(args[1])
	return shim.Success(buffer.Bytes())*/

	var buffer bytes.Buffer
        buffer.WriteString("")
	
	response, err := iv.InsertData(stub,args)
	if err != nil {
		return shim.Error("Insertion failed" +err.Error())
	}
	if (response == 0)&&(err == nil){
		return shim.Error("Incorrect number of arguments. Expecting 12")
	}

	url := fmt.Sprintf("http://rdctstbc001:5984/mychannel_mycc/%s",args[0])
	var dat map[string]interface{}
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            return shim.Error(err.Error())
	}
	err1 := json.Unmarshal(contents, &dat)
	if err1 != nil{
		return shim.Error(err1.Error())
	}
	fmt.Println(dat["_rev"])
    	fmt.Printf("%s\n", string(contents))
	
	buffer.WriteString("The Hash for the added patient is: ")
	buffer.WriteString(dat["_rev"].(string))	

	return shim.Success(buffer.Bytes())
}

func (t *SimpleChaincode) updatePatdisplay(stub shim.ChaincodeStubInterface, args[]string) pb.Response{
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	patientId:= args[0]
	newPatdisplay:= args[1]
	
	patientAsBytes, err := stub.GetState(patientId)
	if err != nil {
		return shim.Error("Failed to get Patient:" + err.Error())
	} else if patientAsBytes == nil {
		return shim.Error("Patient does not exist")
	}
	
	patNameUpdate := Payer{}
	err = json.Unmarshal(patientAsBytes, &patNameUpdate) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	patNameUpdate.Patdisplay = newPatdisplay//change the Patient Name

	patientJSONasBytes, _ := json.Marshal(patNameUpdate)
	err = stub.PutState(patientId, patientJSONasBytes) //rewrite the patient
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- change Patient Name (success)")
	return shim.Success(nil)
}

func (t *SimpleChaincode) updateSourceRef(stub shim.ChaincodeStubInterface, args[]string) pb.Response{
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	patientId := args[0]
	newSourceRef := args[1]
	
	patientAsBytes, err := stub.GetState(patientId)
	if err != nil {
		return shim.Error("Failed to get Patient:" + err.Error())
	} else if patientAsBytes == nil {
		return shim.Error("Patient does not exist")
	}

	sourceRefUpdate := Payer{}
	err = json.Unmarshal(patientAsBytes, &sourceRefUpdate) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	sourceRefUpdate.Sourcereference = newSourceRef//change the Source Reference

	SourceRefJSONasBytes, _ := json.Marshal(sourceRefUpdate)
	err = stub.PutState(patientId, SourceRefJSONasBytes) //rewrite the patient
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- change Source Reference (success)")
	return shim.Success(nil)

}

func (t *SimpleChaincode) updateSourcedisplay(stub shim.ChaincodeStubInterface, args[]string) pb.Response{
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	patientId := args[0]
	newSourcedisplay := args[1]
	
	patientAsBytes, err := stub.GetState(patientId)
	if err != nil {
		return shim.Error("Failed to get Patient:" + err.Error())
	} else if patientAsBytes == nil {
		return shim.Error("Patient does not exist")
	}

	sourceDisplayUpdate := Payer{}
	err = json.Unmarshal(patientAsBytes, &sourceDisplayUpdate) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	sourceDisplayUpdate.Sourcedisplay = newSourcedisplay//change the Source Name

	SourceDispJSONasBytes, _ := json.Marshal(sourceDisplayUpdate)
	err = stub.PutState(patientId, SourceDispJSONasBytes) //rewrite the patient
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- change Source Name (success)")
	return shim.Success(nil)
}

func (t *SimpleChaincode) updateStatus(stub shim.ChaincodeStubInterface, args[]string) pb.Response{
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	patientId := args[0]
	newStatus := args[1]
	
	patientAsBytes, err := stub.GetState(patientId)
	if err != nil {
		return shim.Error("Failed to get Patient:" + err.Error())
	} else if patientAsBytes == nil {
		return shim.Error("Patient does not exist")
	}

	statusUpdate := Payer{}
	err = json.Unmarshal(patientAsBytes, &statusUpdate) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	statusUpdate.Status = newStatus//change the Status

	StatusJSONasBytes, _ := json.Marshal(statusUpdate)
	err = stub.PutState(patientId, StatusJSONasBytes) //rewrite the patient
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- change Status (success)")
	return shim.Success(nil)
}
