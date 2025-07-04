// Copyright 2012 The GoSNMP Authors. All rights reserved.  Use of this
// source code is governed by a BSD-style license that can be found in the
// LICENSE file.

package main

import (
	"fmt"
	g "github.com/Rakortx/gosnmp"
	"log"
	"strconv"
	"strings"
)

func main() {
	// Default is a pointer to a GoSNMP struct that contains sensible defaults
	// eg port 161, community public, etc
	g.Default.Target = "192.168.1.10"
	err := g.Default.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer g.Default.Conn.Close()

	oids := []string{"1.3.6.1.2.1.1.4.0", "1.3.6.1.2.1.1.7.0"}
	result, err2 := g.Default.Get(oids) // Get() accepts up to g.MAX_OIDS
	if err2 != nil {
		log.Fatalf("Get() err: %v", err2)
	}

	for i, variable := range result.Variables {
		fmt.Printf("%d: oid: %s ", i, variable.Name)
		switch variable.Type {
		case g.OctetString:
			value := variable.Value.([]byte)
			if strings.Contains(strconv.Quote(string(value)), "\\x") {
				tmp := ""
				for i := 0; i < len(value); i++ {
					tmp += fmt.Sprintf("%v", value[i])
					if i != (len(value) - 1) {
						tmp += " "
					}
				}
				fmt.Printf("Hex-String: %s\n", tmp)
			} else {
				fmt.Printf("string: %s\n", string(variable.Value.([]byte)))
			}
		default:
			// ... or often you're just interested in numeric values.
			// ToBigInt() will return the Value as a BigInt, for plugging
			// into your calculations.
			fmt.Printf("number: %d\n", g.ToBigInt(variable.Value))
		}
	}
}
