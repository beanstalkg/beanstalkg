package testintegration

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"os"
	"os/exec"
	"time"
	// "strconv"
)

var _ = Describe("Put", func() {

	var bs_session *gexec.Session
	var conn *Conn

	BeforeEach(func() {
		var err error
		command := exec.Command(os.Getenv("GOPATH") + "/bin/beanstalkg")
		bs_session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Ω(err).ShouldNot(HaveOccurred())
		// wait for the server to become ready
		time.Sleep(500 * time.Millisecond)
		conn, err = dial("tcp", "127.0.0.1:11300")
		Ω(err).ShouldNot(HaveOccurred())
	})

	Describe("Put command", func() {
		Context("Used with default tube", func() {
			It("should correctly put first job", func() {
				r, err := conn.cmd([]byte("hello"), "put", 0, dur(1), dur(1))
				Ω(err).ShouldNot(HaveOccurred())
				var id string
				_, err = conn.readResp(r, false, "INSERTED %s", &id)
				Ω(err).ShouldNot(HaveOccurred())
				// println("successfully put job with id " + id)
				Ω(id).Should(MatchRegexp("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}"))
			})
		})
	})

	Describe("Produce and Consume job", func() {
		Context("Used with default tube", func() {
			It("should correctly put job and then reserve it", func() {
				put_payload := "hello"
				// produce job
				r1, err := conn.cmd([]byte(put_payload), "put", 0, dur(1), dur(1))
				Ω(err).ShouldNot(HaveOccurred())
				var put_id string
				_, err = conn.readResp(r1, false, "INSERTED %s", &put_id)
				Ω(err).ShouldNot(HaveOccurred())
				// consume job
				r2, err := conn.cmd(nil, "reserve")
				Ω(err).ShouldNot(HaveOccurred())
				// var info int
				var reserved_id string
				var reserved_payload []byte
				reserved_payload, err = conn.readResp(r2, true, "RESERVED %s", &reserved_id)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(reserved_id).Should(BeIdenticalTo(put_id))
				Ω(string(reserved_payload)).Should(BeIdenticalTo(put_payload))
			})
		})
	})

	AfterEach(func() {
		gexec.Kill()
		conn.Close()
	})
})
