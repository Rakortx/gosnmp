// Copyright 2012 The GoSNMP Authors. All rights reserved.  Use of this
// source code is governed by a BSD-style license that can be found in the
// LICENSE file.

// This program demonstrates BulkWalk.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Rakortx/gosnmp"
)

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage:\n")
		fmt.Printf("   %s [-community=<community>] host [oid]\n", filepath.Base(os.Args[0]))
		fmt.Printf("     host      - the host to walk/scan\n")
		fmt.Printf("     oid       - the MIB/Oid defining a subtree of values\n\n")
		flag.PrintDefaults()
	}

	var community string
	flag.StringVar(&community, "community", "public", "the community string for device")

	flag.Parse()

	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(1)
	}
	target := flag.Args()[0]
	var oid string
	if len(flag.Args()) > 1 {
		oid = flag.Args()[1]
	}

	gosnmp.Default.Target = target
	gosnmp.Default.Transport = "tcp"
	gosnmp.Default.Community = community
	gosnmp.Default.Timeout = 10 * time.Second // Timeout better suited to walking
	gosnmp.Default.Logger = gosnmp.NewLogger(log.New(os.Stdout, "", 0))
	err := gosnmp.Default.Connect()
	if err != nil {
		fmt.Printf("Connect err: %v\n", err)
		os.Exit(1)
	}
	defer gosnmp.Default.Conn.Close()

	err = gosnmp.Default.BulkWalk(oid, printValue)
	if err != nil {
		fmt.Printf("Walk Error: %v\n", err)
		os.Exit(1)
	}

	// This may lead to the remote server closing the TCP connection
	time.Sleep(15 * time.Second)
	err = gosnmp.Default.BulkWalk(oid, printValue)
	if err != nil {
		fmt.Printf("Walk Error: %v\n", err)
		os.Exit(1)
	}
}

func printValue(pdu gosnmp.SnmpPDU) error {
	fmt.Printf("%s = ", pdu.Name)

	switch pdu.Type {
	case gosnmp.OctetString:
		b := pdu.Value.([]byte)
		fmt.Printf("STRING: %s\n", string(b))
	default:
		fmt.Printf("TYPE %d: %d\n", pdu.Type, gosnmp.ToBigInt(pdu.Value))
	}
	return nil
}
