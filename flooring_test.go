package main
import (
	"reflect"
	"testing"
	"fmt"
)
func TestValName(t *testing.T) {
	name := "Andrew"
	if validateName(name) != true {
		t.Errorf("Should return true")
	}
}
func TestValName2(t *testing.T) {
	name := "Andrew*"
	if validateName(name) != false {
		t.Errorf("Should return true")
	}
}
func TestValName3(t *testing.T) {
	name := "Andrew "
	if validateName(name) != true {
		t.Errorf("Should return true")
	}
}
func TestValDate(t *testing.T) {
	month := "12"
	day := "01"
	time1,val := validateDate(month,day)
	if val != true && time1.Day() != 1 && time1.Month() != 12{
		t.Errorf("Should return true")
	}
}
func TestValDate2(t *testing.T) {
	month := "13"
	day := "01"
	_,val := validateDate(month,day)
	if val != false {
		t.Errorf("Should return true")
	}
}
func TestValDate3(t *testing.T) {
	month := "05"
	day := "01"
	_,val := validateDate(month,day)
	if val != false {
		t.Errorf("Should return true")
	}
}
func TestValDate4(t *testing.T) {
	month := "07"
	day := "01"
	_,val := validateDate(month,day)
	if val != false {
		t.Errorf("Should return true")
	}
}
func TestValDate5(t *testing.T) {
	month := "07"
	day := "02"
	t1,val := validateDate(month,day)
	if val != true && t1.Month() == 7 && t1.Day() == 2 {
		t.Errorf("Should return true")
	}
}
func TestState1(t *testing.T) {
	line := "OH,Ohio,6.25"
	testObject ,_ := parseState(line,"Ohio")
	if testObject.TaxRate != 6.25 {
		t.Errorf("Invalid Struct")
	}
}
func TestState2(t *testing.T) {
	testObject ,_ := validateState("Ohio")
	if testObject.TaxRate != 6.25 {
		t.Errorf("Invalid Struct")
	}
}
func TestState3(t *testing.T) {
	_ ,val := validateState("MD")
	if val != false {
		t.Errorf("Invalid Struct")
	}
}
func TestProduct1(t *testing.T) {
	line := "Tile,3.50,4.15"
	testObject ,_ := parseProduct(line,"Tile")
	if testObject.LaborCostPerSquareFoot != 4.15 && testObject.CostPerSquareFoot != 3.50 {
		t.Errorf("Invalid Struct")
	}
}
func TestProduct2(t *testing.T) {
	testObject ,_ := validateProduct("Tile")
	if testObject.LaborCostPerSquareFoot != 4.15 && testObject.CostPerSquareFoot != 3.50 {
		t.Errorf("Invalid Struct")
	}
}
func TestValArea(t *testing.T) {
	name := "3.50"
	_,val := validateArea(name)
	if val != false {
		t.Errorf("Should return true")
	}
}
func TestValArea2(t *testing.T) {
	name := "100.50"
	_,val := validateArea(name)
	if val != true {
		t.Errorf("Should return true")
	}
}
func TestValOrder(t *testing.T) {
	order1,_ := validateOrder("09","01","Wise","OH","Wood","100.00")
	date ,_ := validateDate("09","01")
	expected := &order{OrderDate: date, CustomerName: "Wise",State: "OH",taxRate: 6.25,ProductType: "Wood",Area: 100.00,CostPerSquareFoot: 5.15,LaborCostPerSquareFoot: 4.75,MaterialCost: 515.00,LaborCost: 475,Tax: 61.88,Total: 1051.88}
	if !reflect.DeepEqual(*order1,*expected){
		t.Errorf("Structs do not match")
	}
}
func TestValOrder1(t *testing.T) {
	_,val1 := validateOrder("09","01","Wise","OH","Wood","10.00")
	if val1 != false{
		t.Errorf("Structs do not match")
	}
}
func TestValOrder2(t *testing.T) {
	order1,_ := validateOrder("09","01","Wise","OH","Wood","100.00")
	date ,_ := validateDate("09","01")
	expected := &order{OrderNumber: 1, OrderDate: date, CustomerName: "Wise",State: "OH",taxRate: 6.5,ProductType: "Wood",Area: 100.00,CostPerSquareFoot: 5.15,LaborCostPerSquareFoot: 4.75,MaterialCost: 515.00,LaborCost: 475,Tax: 61.88,Total: 1051.88}
	order1.writeOrder()
	expected.writeOrder()
}
func TestLineToBeDeleted(t *testing.T){
	line,_  := openFileToBeDeleted("912020", "0")
	lineToBeDeleted(line, "912020", "0")
	line,_  = openFileToBeDeleted("912020", "1")
	lineToBeDeleted(line, "912020", "1")
	fmt.Println(display("912020"))
	if display("912020") != "********************************"{
		t.Errorf("String should be empty")
	}
}