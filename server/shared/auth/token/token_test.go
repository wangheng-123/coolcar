package token

import (
	"github.com/dgrijalva/jwt-go"
	"testing"
	"time"
)

const publicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAxWwqOAWfdcJNQgBYnW8C
MX30vb3p5BhZM/B9uqKbuRqepTkOm9j+NlVmGHQrB8bi03y0JOO5pU6sAocCLNDr
S3o4sCCSTEoMpI2OxPmhbC/cetxha/AuJ6vwhSmjbdBTZouLMwj282ONeSWxUYdM
FWpJMbmgyd7LfPGJo/tZAbJaRw7MnA0+WQha+PRbc/gt/Qm77URVHmYvgtVeQLLs
7Sr555/dAwX+p3B6s/ve25VNF/Z76WJSYRZpDXi5FpYMb2Dm9BkaAd+tDTgPsWps
hXcWEuohsT04+WSKbAjwnptGe+cVK50X/bPY9EL4YcIMUGgXsx6T+UQKVp+R+g2n
5wIDAQAB
-----END PUBLIC KEY-----`

func TestVerify(t *testing.T) {
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if err != nil {
		t.Fatalf("cannot parse public key:%v", err)
	}
	v := &JWTTokenVerifier{
		PublicKey: pubKey,
	}

	cases := []struct {
		name    string
		tkn     string
		tm      time.Time
		want    string
		wantErr bool
	}{
		{
			name:    "valid_token",
			tkn:     "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.n4hPBm8kIcQ4ej5MR_E6Q7gl-11r3v9dt8CX1wL4IGij65o_9WojoYNn8Kw2ghH8Wdsan2c5b6XRw2Lzs5ihWGXlb8ckdWhACggtn49tfC4hPwG6dObYDxqfqi9BYkrw-vRCtNkeZMk6mzKWNwtaQM9BPunoJL7hFu5HBGTirR6jtYYYIZcGA1-yhRnb-Vrk7txp5qClEMY7u435De6qtro7iwv1E7ja_qkuwXJG75f45cBsGo6d5l7vv9RtdCck0cLx0iVo0KHvzal3YwsKvG7gomH-HrVdTWF_G5HquJcrYGLugrNzjBc8byBhcZB6rV4UQkCq4m3a7--tbfHw4A",
			tm:      time.Unix(1516239122, 0),
			want:    "1234567890",
			wantErr: false,
		},
		{
			name:    "token_expired",
			tkn:     "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.n4hPBm8kIcQ4ej5MR_E6Q7gl-11r3v9dt8CX1wL4IGij65o_9WojoYNn8Kw2ghH8Wdsan2c5b6XRw2Lzs5ihWGXlb8ckdWhACggtn49tfC4hPwG6dObYDxqfqi9BYkrw-vRCtNkeZMk6mzKWNwtaQM9BPunoJL7hFu5HBGTirR6jtYYYIZcGA1-yhRnb-Vrk7txp5qClEMY7u435De6qtro7iwv1E7ja_qkuwXJG75f45cBsGo6d5l7vv9RtdCck0cLx0iVo0KHvzal3YwsKvG7gomH-HrVdTWF_G5HquJcrYGLugrNzjBc8byBhcZB6rV4UQkCq4m3a7--tbfHw4A",
			tm:      time.Unix(1517239122, 0),
			want:    "",
			wantErr: true,
		},
		{
			name:    "bad_token",
			tkn:     "bad_token",
			tm:      time.Unix(1517239122, 0),
			want:    "",
			wantErr: true,
		},
		{
			name:    "wrong_signature",
			tkn:     "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyNH0.BA1GRYp2qOk3NHz5AN_k50UeJRfDZxaU8wQTlv_zqV-qu1tcVJEw15g9dXg6E-HLj9mwc52lJOINwdp-GaI7n1XsscMUwtqFNMM4uwrnzNGpcpyFrQwZjtkCQehA1gOyTGsCS0tdumzZ1jVWcxB-8UuhZKhFukzXCfFsjMY8fzEbj7XYK_q2hcMj-2DWRJz-njFjqxidd--bTZDxkIa0jf6_scp0rKzLtpEyMtyOpOM3WypYOilPTOx5BFcX4YmFzn6YGcLQtrfMUvHk8yFqc_nsRPqF8Gb2fqs2lxGJ3ezQwefgDzkl2BMvzQyM4Ig7qb2tKx2Q0X7Lu_RtTnFQNw",
			tm:      time.Unix(1517239122, 0),
			want:    "",
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			jwt.TimeFunc = func() time.Time {
				return c.tm
			}
			accountID, err := v.Verify(c.tkn)

			if !c.wantErr && err != nil {
				t.Errorf("verification failed:%v", err)
			}

			if c.wantErr && err != nil {
				t.Errorf("want err;got no error")
			}

			if accountID != c.want {
				t.Errorf("wrong account id. want:%q,got:%q", c.want, accountID)
			}
		})
	}

	//jwt.TimeFunc = func() time.Time {
	//	return time.Unix(1517239122, 0)
	//}
	//
	//accountID, err := v.Verify("eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.n4hPBm8kIcQ4ej5MR_E6Q7gl-11r3v9dt8CX1wL4IGij65o_9WojoYNn8Kw2ghH8Wdsan2c5b6XRw2Lzs5ihWGXlb8ckdWhACggtn49tfC4hPwG6dObYDxqfqi9BYkrw-vRCtNkeZMk6mzKWNwtaQM9BPunoJL7hFu5HBGTirR6jtYYYIZcGA1-yhRnb-Vrk7txp5qClEMY7u435De6qtro7iwv1E7ja_qkuwXJG75f45cBsGo6d5l7vv9RtdCck0cLx0iVo0KHvzal3YwsKvG7gomH-HrVdTWF_G5HquJcrYGLugrNzjBc8byBhcZB6rV4UQkCq4m3a7--tbfHw4A")
	//if err != nil {
	//	t.Errorf("verification failed:%v", err)
	//}
	//want := "1234567890"
	//if accountID != want {
	//	t.Errorf("wrong account id. want:%q,got:%q", want, accountID)
	//}
}
