package main

import (
	"bytes"
	"compress/zlib"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

func main() {
	var (
		isReq bool
		isRsp bool
	)
	flag.BoolVar(&isReq, "q", false, "is request")
	flag.BoolVar(&isRsp, "p", false, "is response")

	flag.Parse()

	if (isReq && isRsp) || (!isReq && !isRsp) {
		flag.Usage()
		return
	}

	if flag.NArg() != 1 {
		flag.Usage()
		return
	}

	data, err := hex.DecodeString(flag.Arg(0))
	if err != nil {
		fmt.Printf("%s, decode string err:%s\n", flag.Arg(0), err)
		return
	}

	method := "unknown"
	t := "req"
	if isRsp {
		t = "rsp"
	}

	if isReq {
		h, err := proto.ParseReqHeader(data)
		if err != nil {
			fmt.Printf("%s, parse header err:%s\n", flag.Arg(0), err)
			return
		}
		method = fmt.Sprintf("%04X", h.Method)
	} else {
		h, err := tdx.ParseRespHeader(data)
		if err != nil {
			fmt.Printf("%s, parse header err:%s\n", flag.Arg(0), err)
			return
		}
		method = fmt.Sprintf("%04X", h.Method)
		if h.RawDataSize != h.PkgDataSize {
			r, _ := zlib.NewReader(bytes.NewReader(data[16:]))
			out := bytes.Buffer{}
			io.Copy(&out, r)
			data = data[:16]
			data = append(data, out.Bytes()...)
		}
	}

	{
		fmt.Printf("header: %s\n", method)
	}

	// save to file
	f, err := os.OpenFile(fmt.Sprintf("%d_%s_%s.bin", time.Now().Unix(), method, t), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		fmt.Println(err)
		return
	}
}
