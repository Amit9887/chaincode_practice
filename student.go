import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

type studentSmartcontract struct{
}

type student struct{
	Name      string  `json:"Name"`
	Age       int     `json:"Age"`
	Gender    string  `json:"Gender"`
	Address   string  `json:"Address"`
	Course    string  `json:"Course"`
}


func (s *studentSmartcontract) Init(APIstub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}


func (s *studentSmartcontract) Invoke(APIstub shim.ChaincodeStubInterface) pb.Response {
	function, args := APIstub.GetFunctionAndParameters()
	if function == "queryStudent" {
		return s.QueryStudent(APIstub, args)
	} else if function == "initLedger" {
		return s.InitLedger(APIstub)
	} else if function == "createStudent" {
		return s.CreateStudent(APIstub, args)
	} else if function == "queryAllStudent" {
		return s.QueryAllStudent(APIstub)
	} else if function == "changeStudentAddress" {
		return s.ChangeStudentAddress(APIstub, args)
	}else if function == "changeStudentCourse" {
		return s.ChangeStudentCourse(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *studentSmartcontract) QueryStudent(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	studentAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(studentAsBytes)
}

func (s *studentSmartcontract) InitLedger(APIstub shim.ChaincodeStubInterface) pb.Response {
	students := []Student{
		Student{Name: "Amit", Age: "18", Gender: "M", Address: "Rewa", Course: "Cse"},
		Student{Name: "Rahul", Age: "19", Gender: "M", Address: "Delhi", Course: "Eee"},
		Student{Name: "Sneha", Age: "17", Gender: "F", Address: "Bhopal", Course: "Ece"},
		Student{Name: "Ritika", Age: "18", Gender: "F", Address: "Kolkata", Course: "Mbbs"},
		Student{Name: "Abhyudya", Age: "20", Gender: "M", Address: "Kanpur", Course: "Cse"},
		Student{Name: "Vivek", Age: "21", Gender: "M", Address: "Patna", Course: "Ce"},
		Student{Name: "Rajat", Age: "19",Gender: "M", Address: "Gurugram", Course: "Cse"},
		Student{Name: "Chandramohan", Age: "22", Gender: "M", Address: "Nasik", Course: "Mca"},
		Student{Name: "Anjali", Age: "19", Gender: "F", Address: "Ghaziabad", Course: "Cse"},
		Student{Name: "Kunal", Age: "20", Gender: "M", Address: "Jabalpur", Course: "Cse"},
	}

	i := 0
	for i < len(students) {
		fmt.Println("i is ", i)
		studentAsBytes, _ := json.Marshal(students[i])
		APIstub.PutState("STUDENT"+strconv.Itoa(i), studentAsBytes)
		fmt.Println("Added", students[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *studentSmartcontract) CreateStudent(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var car = Car{Name: args[1], Age: args[2], Gender: args[3], Address: args[4], Course: args[5]}

	studentAsBytes, _ := json.Marshal(car)
	APIstub.PutState(args[0], studentAsBytes)

	return shim.Success(nil)
}

func (s *studentSmartcontract) QueryAllStudent(APIstub shim.ChaincodeStubInterface) pb.Response {

	startKey := "101"
	endKey := "1999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
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

	fmt.Printf("- queryAllStudent:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *studentSmartcontract) ChangeStudentAddress(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	studentAsBytes, _ := APIstub.GetState(args[0])
	student := Students{}

	json.Unmarshal(carAsBytes, &student)
	student.Address = args[1]

	studentAsBytes, _ = json.Marshal(student)
	APIstub.PutState(args[0], studentAsBytes)

	return shim.Success(nil)
}

func (s *studentSmartcontract) ChangeStudentCourse(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	studentAsBytes, _ := APIstub.GetState(args[0])
	student := Students{}

	json.Unmarshal(studentAsBytes, &student)
	student.Address = args[1]

	studentAsBytes, _ = json.Marshal(student)
	APIstub.PutState(args[0], studentAsBytes)

	return shim.Success(nil)
}

func main() {

	// Create a new Smart Contract
	err := shim.Start(new(studentSmartcontract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
