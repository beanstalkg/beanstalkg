package testintegration

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"os"
	"os/exec"
	"time"
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
				println("id is " + id)
				Ω(id).Should(MatchRegexp("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}"))
			})
		})
	})

	AfterEach(func() {
		gexec.Kill()
		conn.Close()
	})
})
