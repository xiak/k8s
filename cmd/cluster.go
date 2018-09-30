package main

import (
	"fmt"
	"github.com/xiak/k8s/pkg/ssh"
)

func main() {
	//command := cmd.NewDefaultClusterCommand()
	//if err := command.Execute(); err != nil {
	//	fmt.Fprintf(os.Stderr, "%v\n", err)
	//	os.Exit(1)
	//}
	s, err := ssh.NewSSHTunnel("10.xx.xx.xx", "root", "password")
	if err != nil {
		fmt.Errorf("Err: %s", err.Error())
		return
	}
	err = s.Dial()
	if err != nil {
		fmt.Errorf("Err: %s", err.Error())
		return
	}
	outs, _, code, err := s.RunCommond("mkdir -p /xiak && pwd")
	if err != nil {
		fmt.Printf("Err Msg: (%d) %s", code, err.Error())
	}
	fmt.Printf("%s", outs)
	outs, _, code, err = s.RunCommond("cd /xiak && pwd")
	if err != nil {
		fmt.Printf("Err Msg: (%d) %s", code, err.Error())
	}
	fmt.Printf("%s", outs)
	outs, _, code, err = s.RunCommond("pwd")
	if err != nil {
		fmt.Printf("Err Msg: (%d) %s", code, err.Error())
	}
	fmt.Printf("%s", outs)
	s.Close()
}