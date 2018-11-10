package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net"

	quic "github.com/lucas-clemente/quic-go"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("File Transfers", func() {
	const samples = 6

	setupQUIC := func() (quic.Session /* client */, quic.Session /* server */) {
		ln, err := quic.ListenAddr("localhost:0", getTLSConfig(), nil)
		Expect(err).ToNot(HaveOccurred())

		serverSessionChan := make(chan quic.Session)
		go func() {
			defer GinkgoRecover()
			sess, err := ln.Accept()
			Expect(err).ToNot(HaveOccurred())
			serverSessionChan <- sess
		}()
		clientSession, err := quic.DialAddr(ln.Addr().String(), &tls.Config{InsecureSkipVerify: true}, nil)
		Expect(err).ToNot(HaveOccurred())
		return clientSession, <-serverSessionChan
	}

	setupTCP := func() (net.Conn /* client */, net.Conn /* server */) {
		ln, err := tls.Listen("tcp4", "localhost:0", getTLSConfig())
		Expect(err).ToNot(HaveOccurred())

		serverConnChan := make(chan net.Conn)
		go func() {
			defer GinkgoRecover()
			conn, err := ln.Accept()
			Expect(err).ToNot(HaveOccurred())
			// need to send something to make the client accept the conn
			_, err = conn.Write([]byte{'a'})
			Expect(err).ToNot(HaveOccurred())
			serverConnChan <- conn
		}()
		clientConn, err := tls.Dial("tcp4", ln.Addr().String(), &tls.Config{InsecureSkipVerify: true})
		Expect(err).ToNot(HaveOccurred())
		// read the first byte sent by the server
		_, err = clientConn.Read([]byte{0})
		Expect(err).ToNot(HaveOccurred())
		return clientConn, <-serverConnChan
	}

	const MB = 1 << 20
	sizesMB := []int{1, 5, 25, 100}
	for i := range sizesMB {
		sizeMB := sizesMB[i]
		size := int(sizeMB * MB)

		Measure(fmt.Sprintf("QUIC: transfering %d MB of data", sizeMB), func(b Benchmarker) {
			client, server := setupQUIC()
			defer client.Close()
			defer server.Close()
			go func() {
				defer GinkgoRecover()
				buf := make([]byte, size)
				str, err := server.AcceptStream()
				if err != nil {
					return
				}
				_, err = io.ReadFull(str, buf)
				Expect(err).ToNot(HaveOccurred())
				str.Close()
			}()

			data := bytes.Repeat([]byte{0x42}, size)
			b.Time("runtime", func() {
				str, err := client.OpenStream()
				Expect(err).ToNot(HaveOccurred())
				go func() {
					defer GinkgoRecover()
					_, err := str.Write(data)
					Expect(err).ToNot(HaveOccurred())
				}()
				_, err = str.Read([]byte{0})
				Expect(err).To(MatchError(io.EOF))
			})
		}, samples)

		Measure(fmt.Sprintf("TCP: transfering %d MB of data", sizeMB), func(b Benchmarker) {
			client, server := setupTCP()
			go func() {
				defer GinkgoRecover()
				buf := make([]byte, size)
				_, err := io.ReadFull(server, buf)
				Expect(err).ToNot(HaveOccurred())
				server.Close()
			}()

			data := bytes.Repeat([]byte{0x42}, size)
			b.Time("runtime", func() {
				go func() {
					defer GinkgoRecover()
					_, err := client.Write(data)
					Expect(err).ToNot(HaveOccurred())
				}()
				_, err := client.Read([]byte{0})
				Expect(err).To(MatchError(io.EOF))
			})
			client.Close()
		}, samples)
	}
})
