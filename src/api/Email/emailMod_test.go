package email_test

import (
	"fmt"
	email "server/src/api/Email"
	"testing"
)

func TestEmailTest(t *testing.T) {
	t.Log("Test")
	e := email.GeneratEmail()
	r, er := e.SendEmail("testmogus", "sus", "julian.hareng@gmail.com")
	fmt.Println(er)
	fmt.Println(r)
}
