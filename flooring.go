package main
import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"errors"
)

var filePath = "/Users/rahulkumar/go/src/flooring/"

type order struct {
	OrderDate time.Time
	OrderNumber int
	CustomerName string
	State string
	taxRate float64
	ProductType string
	Area float64
	CostPerSquareFoot float64
	LaborCostPerSquareFoot float64
	MaterialCost  float64
	LaborCost float64
	Tax float64
	Total float64
}
func validateName(customerName string) bool {
	exp ,_ := regexp.Compile("^[a-zA-Z0-9., ]+$")
	return exp.MatchString(customerName)
}
func validateDate(month string, day string) (time.Time,bool) {
	month1,_ := strconv.Atoi(month)
	if month1 > 12 || month1 <= 0 {
		return time.Time{}, false
	}
	day1,_ := strconv.Atoi(day)
	if month1 == 2 {
		if day1 > 28 && day1 <= 0 {
			return time.Time{}, false
		}
	} else if month1 == 4 || month1 == 6 || month1 == 9 || month1 == 11 {
		if day1 > 30 && day1 <= 0 {
			return time.Time{}, false
		}
	} else {
		if day1 > 31 && day1 <= 0 {
			return time.Time{}, false
		}
	}
	if month1 < 10 {
		month += "0" + month
	}
	if day1 < 10 {
		day = "0" + day
	}
	dateTimeSrting := "2020" + "-" + month + "-" + day + "T" + "00" + ":" + "00" + ":" + "00" + "+00:00"
	time1, err := time.Parse(time.RFC3339, dateTimeSrting)
	fmt.Println(err)
	if err != nil {
		return time.Time{},false
	} else {
		if time.Now().Before(time1) {
			return time1,true
		}
	}
	return time.Time{},false
}
func parseState(line string, state string) (*taxinfo,bool) {
	lines := strings.Split(line,",")
	if strings.ToLower(state) == strings.ToLower(lines[0]) || strings.ToLower(state) == strings.ToLower(lines[1]) {
		tax ,_ := strconv.ParseFloat(lines[2],64)
		return newTaxinfo(lines[0],lines[1],tax),true
	} else {
		return nil, false
	}
}
func validateState(state string) (*taxinfo,bool){
	file1,err := os.Open(filePath + "/taxes.txt")
	if err != nil {
		file1,err = os.Create(filePath + "/taxes.txt")
		if err != nil {
			panic("File Not Found")
		}
	}
	// create a scanner to read from file and split text based on lines
	scanner := bufio.NewScanner(file1)
	scanner.Split(bufio.ScanLines)
	exists := false
	var tax *taxinfo
	// use Scan to iterate through the file
	for scanner.Scan() {
		// append the current line to the slice lines
		tax,exists = parseState(scanner.Text(),state)
		if exists == true {
			break
		}
	}
	return tax,exists
}
func parseProduct(line string, productType string) (*product,bool) {
	lines := strings.Split(line,",")
	if strings.ToLower(productType) == strings.ToLower(lines[0]) {
		costPerSquareFoot ,_ := strconv.ParseFloat(lines[1],64)
		laborCostPerSquareFoot,_ := strconv.ParseFloat(lines[2],64)
		return newProduct(lines[0],costPerSquareFoot,laborCostPerSquareFoot),true
	} else {
		return nil, false
	}
}
func validateProduct(productstr string) (*product,bool){
	file1,err := os.Open(filePath + "/Products.txt")
	if err != nil {
		file1,err = os.Create(filePath + "/Products.txt")
		if err != nil {
			panic("File Not Found")
		}
	}
	// create a scanner to read from file and split text based on lines
	scanner := bufio.NewScanner(file1)
	scanner.Split(bufio.ScanLines)
	exists := false
	var product1 *product
	// use Scan to iterate through the file
	for scanner.Scan() {
		// append the current line to the slice lines
		product1,exists = parseProduct(scanner.Text(),productstr)
		if exists == true {
			break
		}
	}
	return product1,exists
}
func validateArea(area string) (float64,bool){
	area1,_ := strconv.ParseFloat(area,64)
	if area1 >=100 {
		return area1,true
	} else {
		return area1,false
	}
}
func validateOrder(month string, day string, customerName string, state string, producttype string, area string) (*order,bool) {
	orderDate,val1 := validateDate(month,day)
	order, err :=validateOrderWithoutDate(customerName, state,producttype, area)
	fmt.Println(order)
	fmt.Println(val1)
	fmt.Println(err)
	if err == true && val1 == true {
		(*order).OrderDate = orderDate
		return order, true
	}else{
		return nil, false
	}

}

func validateOrderWithoutDate(customerName string, state string, producttype string, area string) (*order,bool) {
	val2 := validateName(customerName)
	taxinfo1, val3 := validateState(state)
	product1,val4 := validateProduct(producttype)
	area1,val5 := validateArea(area)
	if val2 && val3 && val4 && val5 {
		return newOrder(time.Time{},customerName,*taxinfo1,*product1,area1),true
	} else {
		return nil, false
	}
}
func newOrder(orderDate time.Time, customerName string, taxinfo1 taxinfo, product1 product, area float64) *order {
	order1 := &order{OrderDate: orderDate, CustomerName: customerName, State: taxinfo1.StateAbbreviation, ProductType: product1.ProductType, Area: area}
	order1.taxRate = taxinfo1.TaxRate
	order1.CostPerSquareFoot = product1.CostPerSquareFoot
	order1.LaborCostPerSquareFoot = product1.LaborCostPerSquareFoot
	order1.MaterialCost = area * product1.CostPerSquareFoot
	order1.LaborCost = area * product1.LaborCostPerSquareFoot
	order1.Tax = math.Round((order1.MaterialCost + order1.LaborCost) * (taxinfo1.TaxRate/100) * 100)/100
	order1.Total = math.Round((order1.MaterialCost + order1.LaborCost + order1.Tax) *100) /100
	return order1
}
type taxinfo struct {
	StateAbbreviation string
	StateName string
	TaxRate float64
}
func newTaxinfo(stateAbbreviation string, stateName string, taxRate float64) *taxinfo {
	return &taxinfo{StateAbbreviation: stateAbbreviation, StateName: stateName, TaxRate: taxRate}
}
type product struct {
	ProductType string
	CostPerSquareFoot float64
	LaborCostPerSquareFoot float64
}
func newProduct(productType string, costPerSquareFoot float64, laborCostPerSquareFoot float64) *product {
	return &product{ProductType: productType, CostPerSquareFoot: costPerSquareFoot, LaborCostPerSquareFoot: laborCostPerSquareFoot}
}
func (ord1 *order) writeOrder() {
	month := strconv.Itoa(int((ord1.OrderDate.Month())))
	day := strconv.Itoa(int((ord1.OrderDate.Day())))
	if ord1.OrderDate.Month() < 10 {
		month += "0" + month
	}
	if ord1.OrderDate.Day() < 10 {
		day = "0" + day
	}
	file1,err := os.OpenFile(filePath + "/Orders_"+month +
		day + strconv.Itoa(int((ord1.OrderDate.Year()))), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		file1,err = os.Create(filePath + "/Orders_"+month+
			day + strconv.Itoa(int((ord1.OrderDate.Year()))))
		if err != nil {
			panic("File Not Found")
		}
	}
	datewriter := bufio.NewWriter(file1)
	stringToWrite := ord1.toString()
	datewriter.WriteString(stringToWrite + "\n")
	datewriter.Flush()
	file1.Close()
}
func (ord1 *order) toString() string{
	return strconv.Itoa(ord1.OrderNumber) + "," + ord1.CustomerName + "," +
		ord1.State + "," + fmt.Sprintf("%f",ord1.taxRate) + "," + ord1.ProductType + "," +
		fmt.Sprintf("%f",ord1.Area) + "," +
		fmt.Sprintf("%f",ord1.CostPerSquareFoot) + "," +
		fmt.Sprintf("%f",ord1.LaborCostPerSquareFoot) + "," +
		fmt.Sprintf("%f",ord1.MaterialCost) + "," +
		fmt.Sprintf("%f",ord1.LaborCost) + "," +
		fmt.Sprintf("%f",ord1.Tax) + "," +
		fmt.Sprintf("%f",ord1.Total)
}

func display(date string) string{
	file1, err := os.Open(filePath + "/Orders_"+date)
	if err != nil {
		return "File Not Found"
	}
	scanner := bufio.NewScanner(file1)
	scanner.Split(bufio.ScanLines)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	defer file1.Close()
	var newline string

	for _, line := range lines {
		newline += OrderLine(line, date)
	}
	newline += "********************************"
	return newline
}

func OrderLine(line string, date string) string {
	lines := strings.Split(line,",")
	var newline string
	newline += "******************************** \n"
	newline += lines[0] + "|" + date + "\n"
	newline += lines[1] + "\n"
	newline += "Product: " + lines[4] + "\n"
	newline += "Material Cost: " + lines[8] + "\n"
	newline += "Labor Cost: " + lines[9] + "\n"
	newline += "Tax: " + lines[10] + "\n"
	newline += "Total: " + lines[11] + "\n"

	return newline
}

func updateOrder(lines []string, date string, orderNumber string) {

	file1, _ := os.Create(filePath + "/Orders_"+date)
	defer file1.Close()

	datewriter := bufio.NewWriter(file1)

	for _ , line := range lines {
		updateLine, err := regexp.Compile("^" + orderNumber)
		if err != nil {
			fmt.Print(err)
		}
		var stringToWrite string = line
		var er error
		if updateLine.MatchString(line) {
			stringToWrite, er = lineToBeUpdated(line)
			if er != nil{
				stringToWrite = line
			}
		}
		datewriter.WriteString(stringToWrite + "\n")
	}
	datewriter.Flush()

}

func commonInput() (string,string,string,string) {
	customerName := userInput("Customer name: ")
	state := userInput("State: ")
	producttype := userInput("Product type: ")
	area := userInput("area: ")
	return customerName,state,producttype, area
}

func createInput() (string,string,string,string,string,string) {
	month := userInput("Month: ")
	day := userInput("Day: ")
	customerName,state,producttype,area := commonInput()
	return month,day,customerName,state,producttype,area
}

func lineToBeUpdated(lines string) (string, error) {
	line := strings.Split(lines,",")
	customerName,state,producttype,area := commonInput()
	if customerName == "" {
		customerName = string(line[1])
	}

	if state == "" {
		state = string(line[2])
	}
	if producttype == "" {
		producttype = string(line[4])
	}
	//area
	if area == "" {
		area = string(line[5])
	}
	strc, err := validateOrderWithoutDate(customerName, state, producttype, area)
	if err != true{
		return "", errors.New("It's wrong")
	}else{
		orderNumber, _ := strconv.Atoi(string(lines[0]))
		strc.OrderNumber = orderNumber
		return strc.toString(), nil
	}
}

func userInput(input string) string{
	var response string
	fmt.Print(input)
	fmt.Scanln(&response)
	return response
}

func lineToBeDeleted(lines []string,date string, orderNumber string) {

	file1, _ := os.Create(filePath + "/Orders_"+date)
	defer file1.Close()

	datewriter := bufio.NewWriter(file1)

	for _ , line := range lines {
		updateLine, err := regexp.Compile("^" + orderNumber)
		if err != nil {
			fmt.Print(err)
		}
		var stringToWrite string = line
		if updateLine.MatchString(line) {
			continue
		}
		datewriter.WriteString(stringToWrite + "\n")
	}
	datewriter.Flush()

}

func openFileToBeDeleted(date string, orderNumber string) ([]string, bool){
	file1, err := os.Open(filePath + "/Orders_"+date)
	if err != nil {
		panic("File Not Found")
	}
	scanner := bufio.NewScanner(file1)
	scanner.Split(bufio.ScanLines)
	var lines []string
	var exists bool = false
	// var indexInt int = -1
	// var counter int
	for scanner.Scan() {
		updateLine, err := regexp.Compile("^" + orderNumber)
		if err != nil {
			fmt.Print(err)
		}
		line := scanner.Text()
		if updateLine.MatchString(line) {
			exists = true
			// indexInt = counter
		}
		// counter++
		lines = append(lines, scanner.Text())

	}
	os.Remove(filePath + "/Orders_"+date)
	// return lines, indexInt, exists
	return lines, exists
}

func main() {
	for {
		input := userInput("Action: ")
		if input == "create" {
			order1,err := validateOrder(createInput())
			if err == true {
				order1.writeOrder()
			} else {
				fmt.Println("Incorrect Order")
			}
		} else if input == "update" {
			date := userInput("Enter date to be read [MMDDYYYY]")
			order := userInput("Enter order to be read ")
			arr,val := openFileToBeDeleted(date,order)
			if val {
				stroke := userInput("Delete Record?: ")
				if stroke == "yes" {
					updateOrder(arr,date,order)
				} else {
					fmt.Println("Invalid Record To Update")
				}
			}
		} else if input == "delete" {
			date := userInput("Enter date to be read [MMDDYYYY]")
			order := userInput("Enter order to be read ")
			arr,val := openFileToBeDeleted(date,order)
			if val {
				stroke := userInput("Delete Record?: ")
				if stroke == "yes" {
					lineToBeDeleted(arr,date,order)
				} else {
					fmt.Println("Invalid Record To Delete")
				}
			}
		} else if input == "read" {
			fmt.Println(display(userInput("Enter date to be read [MMDDYYYY]")))
		} else if input == "exit" {
			break
		} else {
			fmt.Println("Invalid Input")
		}
	}
}