package endpoint

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/demostack/auth/pkg"

	"github.com/akrylysov/algnhsa"
	"github.com/aws/aws-sdk-go/aws"
)

// Core .
type Core struct {
	DB         *pkg.DynamoDB
	region     string
	phone      string
	isSAMLocal bool
}

// Hander .
func (c *Core) Hander(w http.ResponseWriter, r *http.Request) {
	switch true {
	case r.URL.Path == "/healthcheck":
		c.Healthcheck(w, r)
	case r.URL.Path == "/auth":
		c.Auth(w, r)
	case strings.HasPrefix(r.URL.Path, "/verify"):
		c.Verify(w, r)
	}
}

// Healthcheck .
func (c *Core) Healthcheck(w http.ResponseWriter, r *http.Request) {
	pkg.Render(w, http.StatusOK, "ok")
}

// Auth .
func (c *Core) Auth(w http.ResponseWriter, r *http.Request) {
	// Only accept POST.
	if r.Method != http.MethodPost {
		pkg.Render(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	req := Request{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		pkg.Render(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.Username == "" || req.Password == "" {
		pkg.Render(w, http.StatusBadRequest, "one or more of the request fields are missing")
		return
	}

	fmt.Println("auth request for username:", req.Username)
	//fmt.Println("auth request for password:", req.Password)

	if req.Username != "username" || req.Password != "password" {
		fmt.Println("auth invalid")
		pkg.Render(w, http.StatusBadRequest, err.Error())
		return
	}

	id, _ := pkg.UUID()
	c.DB.SetVerify(id, false)

	// Set the callback URL.
	u := "http://localhost:8080/verify/" + id
	// If running in lambda, rebuild the URL from the context.
	p, ok := algnhsa.ProxyRequestFromContext(r.Context())
	if ok && !c.isSAMLocal {
		u = fmt.Sprintf("https://%v/%v/verify/%v", p.Headers["Host"], p.RequestContext.Stage, id)

		fmt.Println("request header information:", p.Headers)
	}

	// Number of seconds to wait.
	wait := 20
	fmt.Printf("waiting %v second(s) for user to verify: %v\n", wait, u)

	err = pkg.SendMessage(c.region, c.phone, fmt.Sprintf(`MFA request: %v. Approve: %v`, req.Username, u))
	if err != nil {
		fmt.Println("Send error:", err)
	}

	for i := 1; i <= wait; i++ {
		time.Sleep(1 * time.Second)

		verified, ok := c.DB.Verified(id)
		if !ok {
			fmt.Println("internal error: auth code is missing")
			continue
		}

		if verified {
			fmt.Println("user verified")
			c.DB.DeleteVerify(id)

			creds, err := pkg.GetSessionToken(c.region, 900)
			if err != nil {
				fmt.Println("error getting sts credentials:", err)
			}

			pkg.RenderAccessKeys(w, http.StatusOK,
				aws.StringValue(creds.AccessKeyId),
				aws.StringValue(creds.SecretAccessKey),
				aws.StringValue(creds.SessionToken))
			return
		}
	}

	// Delete the value to clean up.
	c.DB.DeleteVerify(id)
	pkg.SendMessage(c.region, c.phone, "Request not approved.")

	fmt.Println("user did not respond in time")
	pkg.Render(w, http.StatusBadRequest, "user did not respond in time")

}

// Verify .
func (c *Core) Verify(w http.ResponseWriter, r *http.Request) {
	id := strings.Replace(r.URL.Path, "/verify/", "", -1)

	_, found := c.DB.Verified(id)
	if !found {
		pkg.Render(w, http.StatusBadRequest, "invalid ID")
		return
	}

	err := c.DB.SetVerify(id, true)
	if err != nil {
		fmt.Println("Could not set the token to valid:", err)
	}
	pkg.SendMessage(c.region, c.phone, "Request approved.")

	pkg.Render(w, http.StatusOK, "marked as approved")
}
