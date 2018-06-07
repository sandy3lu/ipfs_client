/*
Package corerepo provides pinning and garbage collection for local
IPFS block services.

IPFS nodes will keep local copies of any object that have either been
added or requested locally.  Not all of these objects are worth
preserving forever though, so the node adminstrator can pin objects
they want to keep and unpin objects that they don't care about.

Garbage collection sweeps iterate through the local block store
removing objects that aren't pinned, which frees storage space for new
objects.
*/
package corerepo

import (
	"context"
	"fmt"

	"github.com/ipfs/go-ipfs/core"
	path "github.com/ipfs/go-ipfs/path"
	resolver "github.com/ipfs/go-ipfs/path/resolver"
	uio "github.com/ipfs/go-ipfs/unixfs/io"

	cid "gx/ipfs/QmcZfnkapfECQGcLZaf9B79NRg7cRa9EnZh4LSbkCzwNvY/go-cid"
	"github.com/ipfs/go-ipfs/pin"
	b58 "gx/ipfs/QmWFAMPqsEyUX7gDUsRVmMWz59FxSpJ1b2v6bJ1yYzo7jY/go-base58-fast/base58"
)

func Pin(n *core.IpfsNode, ctx context.Context, paths []string, recursive bool) ([]*cid.Cid, error) {
	out := make([]*cid.Cid, len(paths))

	r := &resolver.Resolver{
		DAG:         n.DAG,
		ResolveOnce: uio.ResolveUnixfsOnce,
	}

	for i, fpath := range paths {
		p, err := path.ParsePath(fpath)
		if err != nil {
			return nil, err
		}

		dagnode, err := core.Resolve(ctx, n.Namesys, r, p)
		if err != nil {
			return nil, fmt.Errorf("pin: %s", err)
		}
		err = n.Pinning.Pin(ctx, dagnode, recursive)
		if err != nil {
			return nil, fmt.Errorf("pin: %s", err)
		}
		out[i] = dagnode.Cid()
	}

	err := n.Pinning.Flush()
	if err != nil {
		return nil, err
	}

	return out, nil
}

func Unpin(n *core.IpfsNode, ctx context.Context, paths []string, recursive bool) ([]*cid.Cid, error) {
	unpinned := make([]*cid.Cid, len(paths))

	r := &resolver.Resolver{
		DAG:         n.DAG,
		ResolveOnce: uio.ResolveUnixfsOnce,
	}

	for i, p := range paths {
		p, err := path.ParsePath(p)
		if err != nil {
			return nil, err
		}

		k, err := core.ResolveToCid(ctx, n.Namesys, r, p)
		if err != nil {
			return nil, err
		}

		err = n.Pinning.Unpin(ctx, k, recursive)
		if err != nil {
			return nil, err
		}
		unpinned[i] = k
	}

	err := n.Pinning.Flush()
	if err != nil {
		return nil, err
	}
	return unpinned, nil
}


func CheckForTask(ctx context.Context, n *core.IpfsNode){//TODO: sandy modified
	fmt.Println("[!!!!]receiving Task start......")
	for {
		select{
		case l:=<-pin.GetTask():
			ll:=b58.Encode([]byte(l))
			key := []string{ll}
			_, err:= Pin(n, ctx, key, true)
			n.Blockstore.PinLock().Unlock()
			if (err!=nil){
				fmt.Println("[!!!!]receive AddTask and pinned Err ", err)
				pin.SetTaskResult("Err")
			}else {
				fmt.Println("[!!!!]receive AddTask and pinned Success : ", ll)
				// broadcast message
				pin.SetTaskResult("OK" )
			}

		case l:=<-pin.GetRemoveTask():
			ll:=b58.Encode([]byte(l))
			key := []string{ll}
			_, err:= Unpin(n, ctx, key, true)
			n.Blockstore.PinLock().Unlock()
			if (err!=nil){
				fmt.Println("[!!!!]receive removeTask and unpinned Err ", err)

			}else {
				fmt.Println("[!!!!]receive removeTask and unpinned Success : ", ll)

			}

		case <-ctx.Done():
			fmt.Println("[!!!!]receiving Task  end ......")
			return
		}

	}
}