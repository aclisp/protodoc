package main

import "fmt"

func (pf protoFile) output() {
	pf.outputServices()
	pf.outputEnums()
	pf.outputObjects()
}

func (pf protoFile) outputServices() {
	for _, s := range pf.services {

		// service line
		fmt.Printf("SERVICE %s\n\n", s.packageName)

		// method group comment
		fmt.Printf("%s\n\n", s.comment)

		// list of methods
		for _, inf := range s.infs {

			// method line
			fmt.Printf("METHOD %s", inf.serviceName+"."+inf.methodName)
			if inf.typ != unary {
				fmt.Printf(" (%s)\n", inf.typ)
			} else {
				fmt.Printf("\n")
			}

			// rest api line
			fmt.Printf("%s %s\n", inf.httpMethod, inf.urlPath)

			// method comments
			fmt.Printf("%s\n\n", inf.comment)

			// method request
			fmt.Printf("REQUEST PARAMETERS (%s)\n", inf.req.typ)
			for _, f := range inf.req.params {
				fmt.Printf("    %s %s %s\n", f.Type(pf), f.name, f.comment)
			}

			// method response
			fmt.Printf("RESPONSE PARAMETERS (%s)\n", inf.res.typ)
			for _, f := range inf.res.params {
				fmt.Printf("    %s %s %s\n", f.Type(pf), f.name, f.comment)
			}

			// end of method
			fmt.Printf("\n")
		}
	}
}

func (pf protoFile) outputEnums() {
	for _, enum := range pf.enums {
		fmt.Printf("ENUM %s\n", enum.name)
		fmt.Printf("%s\n\n", enum.comment)
		fmt.Printf("CONSTANTS\n")
		for _, c := range enum.constants {
			fmt.Printf("    %s %s %s\n", c.name, c.val, c.comment)
		}
		fmt.Printf("\n")
	}
}

func (pf protoFile) outputObjects() {
	for _, obj := range pf.objects {
		fmt.Printf("OBJECT %s\n", obj.name)
		fmt.Printf("%s\n\n", obj.comment)
		fmt.Printf("ATTRIBUTES\n")
		for _, a := range obj.attrs {
			fmt.Printf("    %s %s %s\n", a.Type(pf), a.name, a.comment)
		}
		fmt.Printf("\n")
	}
}
