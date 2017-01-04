package main

const MAX = 5

var Public, private string

type Opaque struct {
	Public int
	private int
}

func main() {
	semicolin()
	control()
	assignmentAndRedeclaration()
	returnTypes()
	deferprint()
	newdata()
}

func semicolin() {
	foo := make([]int,MAX)
	for i:=0; i<MAX;i++ {
		foo[i] += i
	}
	if a := MAX; a > foo[MAX-1] {
		println("semicolin")
	}
}

func control() {
	a, b, c := 0, 1, 2
	if a > c {
		print(a)
	} else if b > c {
		print(b)
	} else {
		print(c)
	}
	foo := make([]int,MAX)
	for i := range foo {
		foo[i] = i
		switch foo[i] {
		case a:
			print(a)
		case b:
			print(b)
		case c:
			print(c)
		default:
			print()
		}
	}
	println()
}

func assignmentAndRedeclaration() {
	var a int //new variable
	b := 5    //new variable
	if b == 5 {
		b = 4	//b assigned 0
		a := b 	//new variable a
		print(a)//prints 4
	}
	println(a)//prints 0
}

func returnTypes() {
	print(single())
	a, b := double()
	println(a,b)
	a, b = doubleName()
	println(a,b)
}

func single() int {
	return 5
}
func double() (int, string) {
	return 5, "five"
}
func doubleName() (number int, text string) {
	number = 5
	text = "five"
	return
}

func deferprint() {
	//prints abcxyz
	defer println()
	defer print("z")
	print("a")
	defer print("y")
	defer print("x")
	print("b")
	print("c")
}


func newdata() {
	a := new([]int) //points to unallocated array
	b := make([]int,MAX)
	c := []int{0, 0, 0, 0, 0}
	//a = nil, b = c
	//prints addr [5/5]addr [5/5]addr
	println(a,b,c)
}
	
