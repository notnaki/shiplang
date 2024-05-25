package runtime

import (
	"fmt"
	"time"
)

func showFN(args []RuntimeVal) RuntimeVal {
	for i, arg := range args {
		if i > 0 {
			fmt.Print(", ")
		}
		switch val := arg.(type) {
		case String:
			fmt.Print(val.Value)
		case Number:
			fmt.Print(val.Value)
		case Bool:
			fmt.Print(val.Value)
		case Array:
			fmt.Print("[")
			for j, element := range val.Elements {
				if j > 0 {
					fmt.Print(", ")
				}
				fmt.Print(element.Inspect())
			}
			fmt.Print("]")
		case Struct:
			fmt.Println("{ ")
			for propName, propVal := range val.Properties {

				fmt.Printf("  %s: ", propName)
				var mslice []RuntimeVal
				mslice = append(mslice, propVal)

				showFN(mslice)

			}

			fmt.Print("}")
		default:
			fmt.Print(val.Inspect())
		}
	}
	fmt.Println()
	return MKNULL()
}

func timeFN(_ []RuntimeVal) RuntimeVal {

	return MKNUM(float64(time.Now().UnixMilli()))
}

func dateFN(_ []RuntimeVal) RuntimeVal {
	currentTime := time.Now()
	formattedDateTime := currentTime.Format("15:04:05.000 02-01-2006")

	return MKSTR(formattedDateTime)
}

func rangeFN(args []RuntimeVal) RuntimeVal {
	if len(args) < 1 || len(args) > 2 {
		panic("range function expects one or two arguments")
	}

	// Ensure the first argument is a number
	if args[0].Type() != NumberType {
		panic("range function expects a number as the first argument")
	}

	start := int(args[0].(Number).Value)
	var end int

	// If two arguments are provided, ensure the second argument is also a number
	if len(args) == 2 {
		if args[1].Type() != NumberType {
			panic("range function expects a number as the second argument")
		}
		end = int(args[1].(Number).Value)
	} else {
		// If only one argument is provided, default end value to start
		end = start
		start = 0
	}

	// Create a range of numbers between start and end (inclusive)
	result := make([]RuntimeVal, end-start+1)
	for i := 0; i <= end-start; i++ {
		result[i] = MKNUM(float64(start + i))
	}

	return Array{Elements: result}
}
