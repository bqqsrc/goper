TestExample1val1: TestExample1
TestExample1val2: 999823
TestExample2val1: 896359
KeyWords: KeyWords1
TestExample1val3: {
	TestExampleVal1: TestExampleVal1	
	KeyWords: KeyWords2
	TestExampleVal2: 8563.333
}
TestExample2val2: -96354
TestExample1val4: 986.33 36.66 589.33 4256.33
TestExample2val3: TestExample2val3_1 TestExample2val3_2 TestExample2val3_3 
TestExample2val4: false
KeyWords: KeyWords3
TestExample1val5: true
TestExample3: {
	Val1: 110
	Val2: 33.8563
	Val3: true false false true false true
	Val4: ConfigTestExample3
}

TestExample4val1: {
	
	KeyWords: KeyWords4
	TestExample4val1: {
		TestValue1: {
			Value1: 333
			Value2: value2 			
			KeyWords: KeyWords5
			Value3: true
		}		
		KeyWords: KeyWords6
		TestValue2: 99.88
	}
	TestExample1val2: {
		TestExampleVal1: TestExampleVal2
		TestExampleVal2: 85.333
	}	
	KeyWords: KeyWords7
	TestExample1val3: {
		TestExampleVal1: TestExampleVal3
		TestExampleVal2: 8563.333
	}
}
KeyWords: KeyWords8

include: ./include1.gs 
TestInclude1: {
	TestValue1: {include: ./include2.gs}
	TestValue2: {
		include: ./include3.gs
	}
	TestValue3: {
		Value1: 952
		Value2: TestValue3Val2
		include: ./include.gs
		Value3: true
	}
	include: ./include4.gs
}