package main

import (
	"fmt"
	"strconv"

	"github.com/mit-dci/lit/litrpc"
)

// Shellparse parses user input and hands it to command functions if matching
func (lc *litAfClient) Shellparse(cmdslice []string) error {
	var err error
	var args []string
	cmd := cmdslice[0]
	if len(cmdslice) > 1 {
		args = cmdslice[1:]
	}
	if cmd == "exit" || cmd == "quit" {
		return fmt.Errorf("User exit")
	}
	// help gives you really terse help.  Just a list of commands.
	if cmd == "help" {
		err = Help(args)
		if err != nil {
			fmt.Printf("help error: %s\n", err)
		}
		return nil
	}
	// adr generates a new address and displays it
	if cmd == "adr" {
		err = lc.Adr(args)
		if err != nil {
			fmt.Printf("adr error: %s\n", err)
		}
		return nil
	}

	// bal shows the current set of utxos, addresses and score
	if cmd == "ls" {
		err = lc.Ls(args)
		if err != nil {
			fmt.Printf("ls error: %s\n", err)
		}
		return nil
	}

	// send sends coins to the address specified
	if cmd == "send" {
		err = lc.Send(args)
		if err != nil {
			fmt.Printf("send error: %s\n", err)
		}
		return nil
	}

	if cmd == "lis" { // listen for lnd peers
		err = lc.Lis(args)
		if err != nil {
			fmt.Printf("lis error: %s\n", err)
		}
		return nil
	}

	if cmd == "stop" { // stop remote node
		// actually returns an error
		return lc.Stop(args)
	}

	if cmd == "sweep" { // make lots of 1-in 1-out txs
		err = lc.Sweep(args)
		if err != nil {
			fmt.Printf("sweep error: %s\n", err)
		}
		return nil
	}

	// push money in a channel away from you
	if cmd == "push" {
		err = lc.Push(args)
		if err != nil {
			fmt.Printf("push error: %s\n", err)
		}
		return nil
	}

	if cmd == "con" { // connect to lnd host
		err = lc.Connect(args)
		if err != nil {
			fmt.Printf("con error: %s\n", err)
		}
		return nil
	}

	/*
		if cmd == "msend" {
			err = MSend(args)
			if err != nil {
				fmt.Printf("Msend error: %s\n", err)
			}
			return nil
		}
		if cmd == "rsend" {
			err = RSend(args)
			if err != nil {
				fmt.Printf("Rsend error: %s\n", err)
			}
			return nil
		}
		if cmd == "nsend" {
			err = NSend(args)
			if err != nil {
				fmt.Printf("Nsend error: %s\n", err)
			}
			return nil
		}

		if cmd == "fan" { // fan-out tx
			err = Fan(args)
			if err != nil {
				fmt.Printf("fan error: %s\n", err)
			}
			return nil
		}

		if cmd == "txs" { // show all txs
			err = Txs(args)
			if err != nil {
				fmt.Printf("txs error: %s\n", err)
			}
			return nil
		}

		if cmd == "wcon" { // connect to watch tower
			err = WCon(args)
			if err != nil {
				fmt.Printf("wcon error: %s\n", err)
			}
			return nil
		}

		if cmd == "watch" { // connect to watch tower
			err = Watch(args)
			if err != nil {
				fmt.Printf("watch error: %s\n", err)
			}
			return nil
		}



		// Peer to peer actions
		// send text message
		if cmd == "say" {
			err = Say(args)
			if err != nil {
				fmt.Printf("say error: %s\n", err)
			}
			return nil
		}
		// fund and create a new channel
		if cmd == "fund" {
			err = FundChannel(args)
			if err != nil {
				fmt.Printf("fund error: %s\n", err)
			}
			return nil
		}

		// cooperateive close of a channel
		if cmd == "cclose" {
			err = CloseChannel(args)
			if err != nil {
				fmt.Printf("cclose error: %s\n", err)
			}
			return nil
		}
		if cmd == "break" {
			err = BreakChannel(args)
			if err != nil {
				fmt.Printf("break error: %s\n", err)
			}
			return nil
		}

		if cmd == "fix" {
			err = Resume(args)
			if err != nil {
				fmt.Printf("fix error: %s\n", err)
			}
			return nil
		}
	*/
	fmt.Printf("Command not recognized. type help for command list.\n")
	return nil
}

func (lc *litAfClient) Ls(textArgs []string) error {

	aReply := new(litrpc.AdrReply)
	tReply := new(litrpc.TxoListReply)
	bReply := new(litrpc.BalReply)

	err := lc.rpccon.Call("LitRPC.TxoList", nil, tReply)
	if err != nil {
		return err
	}
	fmt.Printf("\tTxos:\n")
	for i, t := range tReply.Txos {
		fmt.Printf("%d %s h:%d amt:%d %s\n",
			i, t.OutPoint, t.Height, t.Amt, t.KeyPath)
	}

	err = lc.rpccon.Call("LitRPC.Address", nil, aReply)
	if err != nil {
		return err
	}
	fmt.Printf("\tAddresses:\n")
	for i, a := range aReply.Addresses {
		fmt.Printf("%d %s \n", i, a)
	}
	err = lc.rpccon.Call("LitRPC.Bal", nil, bReply)
	if err != nil {
		return err
	}
	fmt.Printf("\tTxo: %d Conf:%d Channel: %d\n",
		bReply.TxoTotal, bReply.Mature, bReply.ChanTotal)

	return nil
}

func (lc *litAfClient) Adr(textArgs []string) error {
	args := new(litrpc.AdrArgs)
	args.NumToMake = 1
	reply := new(litrpc.AdrReply)
	err := lc.rpccon.Call("LitRPC.Address", args, reply)
	if err != nil {
		return err
	}
	fmt.Printf("new adr(s): %s\n", reply.Addresses)
	return nil
}

func (lc *litAfClient) Send(textArgs []string) error {
	args := new(litrpc.SendArgs)
	reply := new(litrpc.TxidsReply)

	// need args, fail
	if len(textArgs) < 2 {
		return fmt.Errorf("need args: ssend address amount(satoshis) wit?")
	}
	/*
		adr, err := btcutil.DecodeAddress(args[0], lc.Param)
		if err != nil {
			fmt.Printf("error parsing %s as address\t", args[0])
			return err
		}
	*/
	amt, err := strconv.ParseInt(textArgs[1], 10, 64)
	if err != nil {
		return err
	}

	fmt.Printf("send %d to address: %s \n", amt, textArgs[0])

	args.DestAddrs = []string{textArgs[0]}
	args.Amts = []int64{amt}

	err = lc.rpccon.Call("LitRPC.Send", args, reply)
	if err != nil {
		return err
	}
	fmt.Printf("sent txid(s):\n")
	for i, t := range reply.Txids {
		fmt.Printf("\t%d %s\n", i, t)
	}
	return nil
}

// Lis starts listening.  Takes args of port to listen on.
func (lc *litAfClient) Lis(textArgs []string) error {
	args := new(litrpc.ListenArgs)
	reply := new(litrpc.StatusReply)

	args.Port = ":2448"
	if len(textArgs) > 0 {
		args.Port = ":" + textArgs[0]
	}

	err := lc.rpccon.Call("LitRPC.Listen", args, reply)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", reply.Status)
	return nil
}

func (lc *litAfClient) Connect(textArgs []string) error {
	args := new(litrpc.ConnectArgs)
	reply := new(litrpc.StatusReply)

	if len(textArgs) == 0 {
		return fmt.Errorf("need: con pubkeyhash@hostname:port")
	}

	args.LNAddr = textArgs[0]

	err := lc.rpccon.Call("LitRPC.Connect", args, reply)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", reply.Status)
	return nil
}

func (lc *litAfClient) Stop(textArgs []string) error {
	reply := new(litrpc.StatusReply)

	err := lc.rpccon.Call("LitRPC.Stop", nil, reply)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", reply.Status)

	lc.rpccon.Close()
	return fmt.Errorf("stopped remote lit node")
}

func (lc *litAfClient) Sweep(textArgs []string) error {
	args := new(litrpc.SweepArgs)
	reply := new(litrpc.TxidsReply)

	var err error

	if len(textArgs) < 2 {
		return fmt.Errorf("sweep syntax: sweep adr howmany (drop)")
	}

	args.DestAdr = textArgs[0]
	args.NumTx, err = strconv.Atoi(textArgs[1])
	if err != nil {
		return err
	}

	if len(textArgs) > 2 {
		args.Drop = true
	}

	err = lc.rpccon.Call("LitRPC.Sweep", args, reply)
	if err != nil {
		return err
	}
	fmt.Printf("Swept\n")
	for i, t := range reply.Txids {
		fmt.Printf("%d %s\n", i, t)
	}

	return nil
}

// Push is the shell command which calls PushChannel
func (lc *litAfClient) Push(textArgs []string) error {
	args := new(litrpc.PushArgs)
	reply := new(litrpc.PushReply)

	if len(textArgs) < 3 {
		return fmt.Errorf("need args: push peerIdx chanIdx amt (times)")
	}

	// this stuff is all the same as in cclose, should put into a function...
	peerIdx64, err := strconv.ParseInt(textArgs[0], 10, 32)
	if err != nil {
		return err
	}
	cIdx64, err := strconv.ParseInt(textArgs[1], 10, 32)
	if err != nil {
		return err
	}
	amt, err := strconv.ParseInt(textArgs[2], 10, 32)
	if err != nil {
		return err
	}

	times := int64(1)
	if len(textArgs) > 3 {
		times, err = strconv.ParseInt(textArgs[3], 10, 32)
		if err != nil {
			return err
		}
	}

	args.PeerIdx = uint32(peerIdx64)
	args.QChanIdx = uint32(cIdx64)
	args.Amt = amt

	for times > 0 {
		err := lc.rpccon.Call("LitRPC.Push", args, reply)
		if err != nil {
			return err
		}
		fmt.Printf("Pushed %d new balance %d at state %d\n",
			amt, reply.MyAmt, reply.StateIndex)
	}

	return nil
}

func Help(args []string) error {
	fmt.Printf("commands:\n")
	fmt.Printf("help adr bal send fake fan sweep lis con fund push cclose break exit\n")
	return nil
}