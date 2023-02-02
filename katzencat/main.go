package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/katzenpost/katzenpost/client"
	"github.com/katzenpost/katzenpost/client/config"
	"github.com/katzenpost/katzenpost/core/crypto/rand"
	"github.com/katzenpost/katzenpost/core/epochtime"
	"github.com/katzenpost/katzenpost/core/pki"
	mClient "github.com/katzenpost/katzenpost/map/client"
	"github.com/katzenpost/katzenpost/stream"
	"golang.org/x/crypto/hkdf"
	"io"
	"os"
	"time"
)

const (
	keySize = 32
)

var cConf = flag.String("cfg", "namenlos.toml", "config file")
var outFile = flag.String("o", "", "file to write to")
var inFile = flag.String("i", "", "file to write to")

func getSession() (*client.Session, error) {
	cfg, err := config.LoadFile(*cConf)
	cfg.Logging.Level = "NOTICE"
	if err != nil {
		return nil, err
	}
	cc, err := client.New(cfg)
	if err != nil {
		return nil, err
	}
	for {
		session, err := cc.NewTOFUSession()
		switch err {
		case nil:
		case pki.ErrNoDocument:
			_, _, till := epochtime.Now()
			<-time.After(till)
			continue
		default:
			return nil, err
		}
		session.WaitForDocument()
		return session, nil
	}
}

func main() {
	secret := flag.String("s", "", "Secret given by initiator, or empty if initiating")
	flag.Parse()

	if *secret == "" && *inFile == "" {
		fmt.Println("Must specify input file")
		os.Exit(-1)
	}
	if *secret != "" && *outFile == "" {
		fmt.Println("Must specify output file")
		os.Exit(-1)
	}

	salt := []byte("katzencat_initiator_receiver_secret")
	isinitiator := *secret == ""
	if *secret == "" {
		newsecret := &[keySize]byte{}
		io.ReadFull(rand.Reader, newsecret[:])
		*secret = base64.StdEncoding.EncodeToString(newsecret[:])
		fmt.Println("KatzenCat secret: ", *secret)
	}
	hash := sha256.New

	keymaterial := hkdf.New(hash, []byte(*secret), salt, nil)
	asecret := &[keySize]byte{}
	bsecret := &[keySize]byte{}

	if isinitiator {
		io.ReadFull(keymaterial, asecret[:])
		io.ReadFull(keymaterial, bsecret[:])
	} else {
		io.ReadFull(keymaterial, bsecret[:])
		io.ReadFull(keymaterial, asecret[:])
	}
	s, err := getSession()
	if err != nil {
		panic(err)
	}
	c, err := mClient.NewClient(s)
	if err != nil {
		panic(err)
	}

	st := stream.NewStream(c, asecret[:], bsecret[:])

	if isinitiator {
		fi, err := os.Stat(*inFile)
		if err != nil {
			panic(err)
		}
		f, err := os.Open(*inFile)
		if err != nil {
			panic(err)
		}
		// XXX: length prefix payload so client can close connection upon complete
		lengthprefix := make([]byte, 8)
		binary.BigEndian.PutUint64(lengthprefix, uint64(fi.Size()))
		st.Write(lengthprefix)
		total := int64(0)
		for {
			n, err := io.Copy(st, f)
			total += n
			switch err {
			case io.EOF, nil:
				// XXX: unsure how we panic() with io.EOF here!
				// theory: peer Close received in same Frame that Ack'd a blocked Write?
				break
			case os.ErrDeadlineExceeded:
				// this happen when Writes block until timeout
				continue
			default:
				panic(err)
			}
		}
		fmt.Println("Wrote ", total, "bytes")
		// try to read a response from the client until defaultTimeout, and log status
		for {
			b, _ := io.ReadAll(st)
			if len(b) > 0 {
				fmt.Println("peer has completed transfer")
				break
			}
			fmt.Println("peer did not finish reading")
		}
		st.Close()
	} else {
		f, err := os.Create(*outFile)
		if err != nil {
			panic(err)
		}
		lengthprefix := make([]byte, 8)
		n, err := st.Read(lengthprefix)
		if err != nil {
			panic(err)
		}
		if n != 8 {
			panic("failed to read lenghtprefix")
		}
		payloadlen := binary.BigEndian.Uint64(lengthprefix)
		limited := io.LimitReader(st, int64(payloadlen))
		total := int64(0)
		for {
			nn, err := io.Copy(f, limited)
			total += nn
			switch err {
			case os.ErrDeadlineExceeded:
				continue
			case io.EOF:
				panic("failed with short Read, wrong lenght prefix or bug")
			default:
				panic(err)
			case nil:
				break
			}
		}
		fmt.Println("Read ", total, "bytes")
		for {
			_, err = st.Write([]byte{0x42})
			switch err {
			case io.EOF:
				panic("server closed connection prematurely")
			case os.ErrDeadlineExceeded:
				continue
			default:
				panic(err)
			case nil:
			}
			// Hangup reader
		}
		// Hangup reader
		st.Close()
		f.Close()
	}
	// it seems that messages may get lost in the send queue if exit happens immediately after Close()
	<-time.After(10 * time.Second)
	s.Shutdown()
}