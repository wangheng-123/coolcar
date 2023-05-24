package token

import (
	"github.com/dgrijalva/jwt-go"
	"testing"
	"time"
)

const privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAxWwqOAWfdcJNQgBYnW8CMX30vb3p5BhZM/B9uqKbuRqepTkO
m9j+NlVmGHQrB8bi03y0JOO5pU6sAocCLNDrS3o4sCCSTEoMpI2OxPmhbC/cetxh
a/AuJ6vwhSmjbdBTZouLMwj282ONeSWxUYdMFWpJMbmgyd7LfPGJo/tZAbJaRw7M
nA0+WQha+PRbc/gt/Qm77URVHmYvgtVeQLLs7Sr555/dAwX+p3B6s/ve25VNF/Z7
6WJSYRZpDXi5FpYMb2Dm9BkaAd+tDTgPsWpshXcWEuohsT04+WSKbAjwnptGe+cV
K50X/bPY9EL4YcIMUGgXsx6T+UQKVp+R+g2n5wIDAQABAoIBAD16zrb9gUgRxy9r
ni0A2eUBvYqtLr46muTi2rSZWPdPA+KeHx+HdrsC4zVwT8ovNIn5xrvKMxD4q/Zp
htLgCLZLjIXEzup1EPTZpIFQ3+UqDnYwVSJg/G6cS92cNunGu4IuqJ2vCACJmxAE
HfEuuhHdNIgIr7rN5/6z1VI8vt4Y2UoLp/st9RNfQkMA/+UGhI01Gf3+HHUmSYSi
fYccYARxSZG0v8cjTCVwIuaDccj2E15kmsvwD/JQNgGNf5OpakhcsoPCu2WSw2g8
N/tHMlDKSIjpjntKFHtzOnEHZjISIlevieX31hpGpslc5sptyZ9f5uibCBK8vmMy
ELs8WYkCgYEA8QtIfp/rnsiYPt6GUO9maP0NxGqpLPuphJuGxoEkTM4n0cqL6ub9
fZDFE+dGbY4LxHDKs0pSAzfAjQshhaLGl6Gwi18snOAKsJsBQ7UelTYpqUnoJ7NE
SxiutzRgDhtCUzKGegTpIIZugrnibqd1fT79ypLzBbjza3yKxaqFlIsCgYEA0awA
eyGGNNCk0nUHNfsQAzGpKi9HRdkSgb/VIdT7jS7JyFtQj/VJExXjZ2stx/K8mBz+
CQNGf5/x0jP+dEZtmvA/F3PLGJegN0Zt721caoUxvMXFNgCGMmI2fCZdPqoHBg3P
UfDoRakb9IBmFI4NCbqOvI49byVy9ogXYgnZ+ZUCgYEA652I3mV5zrrrvsCLcH+i
jkuVcoKEHbldyYaxJkZD8mOtnq8rN5FVBbFGQx9Vw28O55UNPlYOdqC/sd7IhLVJ
BB0D0ihVFn8VU+4gPUvEuju4W4cny+66eeGFnwUuQ3u3yFViB9HXA3kEevoycNF1
0diKAcLElLpmDpItn+wAKOMCgYEAv2QhoEQOzMEz4wR+i5Dcof2/7Ejx51lp4lRo
yhQvd0WhXam1FWOSy8AsL2gPhzgVXUkBvtpljPREeluJxzvOqyLohJDncFBgKHS6
v1Z1iKqCp01kYpIB7ZXnJFakwSHVfXo3qBWfAI+IfByEkfjE//9ycb3paD6n/VBm
/8/8UC0CgYEAjJrTwgFe9L5JNvp4X7fzAz4fmqYKcMwGkpOjT9grOg0t43X/13bn
DecvQvKO+j+uYvFkfadJc+nH2mq/0Oqxbd8pIuKBUFQFH/pr51UclKXwj+c8VWF6
9Lpwm9kJ+87vWDicHVzIu3Yn0noDQ5X06urCcZ35QKYVkjdfm4Ucp1A=
-----END RSA PRIVATE KEY-----`

func TestGenerateToken(t *testing.T) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		t.Fatalf("cannot parse private key:%v", err)
	}
	g := NewJWTTokenGen("server/auth", key)
	g.nowFunc = func() time.Time {
		return time.Unix(1516239022, 0)
	}
	token, err := g.GenerateToken("1234567890", 2*time.Hour)
	if err != nil {
		t.Errorf("cannot generate token :%v", err)
	}
	want := "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTYyNDYyMjIsImlhdCI6MTUxNjIzOTAyMiwiaXNzIjoic2VydmVyL2F1dGgiLCJzdWIiOiIxMjM0NTY3ODkwIn0.DuE3Pw8kelcrPl9EJ3xR_NrqU5ZOI4Zc7otAXAcLq2jynceJeBv18T89URX11ifatit3HDaYYYLXHl3B4HYbBkdleBhbXYqHFdTEhbsJ9ssIfTPwWHd1KDOV3CdjNUgFA_8VxKoMdZBjl9TKiMyBGckVwIZ-ClKFV2mi6WwWMAMisYTPWGAbGj13AHU4GyIgTqFUrc96Os_-4s4tYcHznUp6gZ6MHLwuCNiaQa9nMaOTm2uk2yqytBAF0sciyBdWVPesJDIL50PqgKghqprwDqkZE0ZEhqNinw4wN6gjyJP53ARz4-Ul4b-hbaySkVQYuR0je1We0nF_krczdFmm0A"
	if token != want {
		t.Errorf("wrong token generated. want:%q ;got:%q", want, token)
	}
}
