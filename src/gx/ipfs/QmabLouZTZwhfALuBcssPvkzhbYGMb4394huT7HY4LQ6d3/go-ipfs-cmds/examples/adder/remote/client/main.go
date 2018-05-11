package main

import (
	"context"
	"os"

	"gx/ipfs/QmabLouZTZwhfALuBcssPvkzhbYGMb4394huT7HY4LQ6d3/go-ipfs-cmds/examples/adder"

	cmds "gx/ipfs/QmabLouZTZwhfALuBcssPvkzhbYGMb4394huT7HY4LQ6d3/go-ipfs-cmds"
	cli "gx/ipfs/QmabLouZTZwhfALuBcssPvkzhbYGMb4394huT7HY4LQ6d3/go-ipfs-cmds/cli"
	http "gx/ipfs/QmabLouZTZwhfALuBcssPvkzhbYGMb4394huT7HY4LQ6d3/go-ipfs-cmds/http"
	cmdkit "gx/ipfs/QmceUdzxkimdYsgtX733uNgzf1DLHyBKN6ehGSp85ayppM/go-ipfs-cmdkit"
)

func main() {
	// parse the command path, arguments and options from the command line
	req, err := cli.Parse(context.TODO(), os.Args[1:], os.Stdin, adder.RootCmd)
	if err != nil {
		panic(err)
	}

	// create http rpc client
	client := http.NewClient(":6798")

	// send request to server
	res, err := client.Send(req)
	if err != nil {
		panic(err)
	}

	req.Options["encoding"] = cmds.Text

	// create an emitter
	re, retCh := cli.NewResponseEmitter(os.Stdout, os.Stderr, req.Command.Encoders["Text"], req)

	if pr, ok := req.Command.PostRun[cmds.CLI]; ok {
		re = pr(req, re)
	}

	wait := make(chan struct{})
	// copy received result into cli emitter
	go func() {
		err = cmds.Copy(re, res)
		if err != nil {
			re.SetError(err, cmdkit.ErrNormal|cmdkit.ErrFatal)
		}
		close(wait)
	}()

	// wait until command has returned and exit
	ret := <-retCh
	<-wait
	os.Exit(ret)
}
