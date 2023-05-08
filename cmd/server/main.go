package main

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

type WrapWriter struct {
	Writer io.Writer
	Name   string
}

func (w *WrapWriter) Write(p []byte) (n int, err error) {
	w.Writer.Write([]byte(w.Name + ": "))
	fmt.Println()
	return w.Writer.Write(p)
}

func main() {
	portForwards := map[string]string{
		"material": "-n qa3-elan port-forward svc/material-pkg-common-service 9002:9000",
		"archive":  "-n qa3-elan port-forward svc/archive-common-service 9001:9000",
		"scm":      "-n qa3-elan port-forward svc/scm-backend-service 9003:9000",
		"order":    "-n qa3-ssmrt port-forward svc/ssc-order-common-service 9004:9000",
		//"material": "-n dev1-elan port-forward svc/material-pkg-common-service 9001:9000",
		//"archive":  "-n dev1-elan port-forward svc/archive-common-service 9002:9000",
		//"scm":      "-n dev1-elan port-forward svc/scm-backend-service 9003:9000",
		//"order":    "-n dev1-ssmrt port-forward svc/ssc-order-common-service 9004:9000",
	}
	g := errgroup.Group{}
	for k, v := range portForwards {
		name, cmd := k, v
		g.Go(func() error {
			for {
				fmt.Printf("%s start port forwarding \n", name)
				p := exec.Command("/usr/bin/kubectl", strings.Split(cmd, " ")...)
				p.Stdout = &WrapWriter{Writer: os.Stdout, Name: name}
				p.Stderr = &WrapWriter{Writer: os.Stderr, Name: name}
				if err := p.Start(); err != nil {
					fmt.Printf("%s start prot forwarding failed error:%v \n", name, err)
					return err
				}
				err := p.Wait()
				fmt.Printf("%s prot forwarding exited error:%v \n", name, err)
				time.Sleep(time.Second * 30)
			}
		})
	}
	if err := g.Wait(); err != nil {
		panic(err)
	}
}
