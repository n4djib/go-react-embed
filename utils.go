package main

import (
	"context"
	"database/sql"
	"fmt"

	_ "embed"
	"errors"
	"flag"
	"go-react-embed/models"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema/schema.sql
var ddl string

func initAndLoadEnv() error {
	_, err := os.Stat(".env")
	if errors.Is(err, os.ErrNotExist) {
		err := createENV(".env")
		if err != nil {
			return err
		}
	}

	// load it
	return godotenv.Load(".env")
}

func createENV(out string) error {
	CONTENT := `MODE="PRODUCTION"
APP_URL="https://localhost"
APP_PORT="8080"
DATABASE="./go-react-embed.db"
SERVER_CRT="server.crt"
SERVER_KEY="server.key"`

	fmt.Println("creating ", out)
	return writeFile(out, CONTENT)
}

func initDatabaseModels() {
	// connect to database
	databaseFile := os.Getenv("DATABASE")
	db, err := sql.Open("sqlite3", databaseFile)
	if err != nil {
		log.Fatal("Connection to DB error\n", err)
	}
	// createTables
	ctx := context.Background()
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		log.Fatal("Table Cretation error\n", err)
	}
	queries := models.New(db)
	// assign to global variables in models package
	models.DB, models.CTX, models.QUERIES = db, ctx, queries
}

func openBrowser() error {
	url := os.Getenv("APP_URL") + ":" + os.Getenv("APP_PORT")
	// grab flag
	air_flag := flag.Bool("air", false, "detect if run by AIR")
	flag.Parse()
	air := bool(*air_flag)

	// open app url
	if !air {
		if err := openURL(url); err != nil {
			return err
		}
	}
	return nil
}

func openURL(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}

	return cmd.Start()
}

func checkSSLFilesExist(SERVER_CRT string, SERVER_KEY string) error {
	if _, err := os.Stat(SERVER_CRT); errors.Is(err, os.ErrNotExist) {
		fmt.Println("file ", SERVER_CRT, "not found")
		err := createServerCRT(SERVER_CRT)
		if err != nil {
			return err
		}
	}
	if _, err := os.Stat(SERVER_KEY); errors.Is(err, os.ErrNotExist) {
		fmt.Println("file ", SERVER_CRT, "not found")
		err := createServerKEY(SERVER_KEY)
		if err != nil {
			return err
		}
	}
	return nil
}

func createServerCRT(out string) error {
	CONTENT := `-----BEGIN CERTIFICATE-----
MIID7TCCAtWgAwIBAgIUBsUq4HmLJ9vRsVgdVTL5MjMwL04wDQYJKoZIhvcNAQEL
BQAwgYUxCzAJBgNVBAYTAmR6MQwwCgYDVQQIDANvZ3gxEDAOBgNVBAcMB291YXJn
bGExITAfBgNVBAoMGEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDESMBAGA1UEAwwJ
bG9jYWxob3N0MR8wHQYJKoZIhvcNAQkBFhBuNGRqaWJAZ21haWwuY29tMB4XDTI0
MDUyMTIwNTU1MloXDTI1MDUyMTIwNTU1MlowgYUxCzAJBgNVBAYTAmR6MQwwCgYD
VQQIDANvZ3gxEDAOBgNVBAcMB291YXJnbGExITAfBgNVBAoMGEludGVybmV0IFdp
ZGdpdHMgUHR5IEx0ZDESMBAGA1UEAwwJbG9jYWxob3N0MR8wHQYJKoZIhvcNAQkB
FhBuNGRqaWJAZ21haWwuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKC
AQEApiTaV4zdQo8AQOWLXsv8XKzZTGBpB8ePvfdg+dYMObu+cbxH6KSwN5NVj4yJ
uTrq0kWccqu87YKs2glbDAkfKGfx233mdpvw+nkSXqXfrx8NSSG4YRsHCTIGltWS
S1aoYG0KTaP/mG/AXVf6w89nxcBvV/NCERKG/ujYn6PfE3I5Fjx5W5I79M75nsyp
QfaY4J6aa3pnEKBNhqp9vOpKoZIIazvw3R8BQEgFKopJlDC2HGu/kQupcxXxzVzh
eFKxfleN0wKhEPIHraUpli2TKCyoFfOz7WxFuTVy0nhPol2tetthhnfc2rzsp4de
vZOure/J2lmgJ6oF+3aK+Vu6OQIDAQABo1MwUTAdBgNVHQ4EFgQUiIJIiXsigN6d
xhQQLBfPL0Lbqm4wHwYDVR0jBBgwFoAUiIJIiXsigN6dxhQQLBfPL0Lbqm4wDwYD
VR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAF/75sxWjsDjgeYnK6u29
XG1pXuz0U/cAel3VX80eQ2wtX6W2gsZbfdHwC/MD5FOtMeCT9rDyxCabbTK+myRg
gXwW+K2Wfyu+NZOg+gnpr4RV4FXSuZCZeb8ds+GOFCigqp4Vro8v26f0GeumR99c
hvDoMlUYXd04oLqMd/wpipxPQs7WwWopnig/e8S8BKJvQ2ELHG7aTvXwCgtZrE+r
yWrV61gTPRsAOQRQ4tTjA6v43+covmrNM36lYRW+rC0PpEmqVPRTxhAHnJ06o1uy
Gj4F4YoR78aZCS2ihk/I3yGH1LpfSlXnAEo4/JOPpsCOYImO2GwKgwbBVoGg+rQd
xw==
-----END CERTIFICATE-----`

	fmt.Println("creating ", out)
	return writeFile(out, CONTENT)
}

func createServerKEY(out string) error {
	CONTENT := `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCmJNpXjN1CjwBA
5Ytey/xcrNlMYGkHx4+992D51gw5u75xvEfopLA3k1WPjIm5OurSRZxyq7ztgqza
CVsMCR8oZ/HbfeZ2m/D6eRJepd+vHw1JIbhhGwcJMgaW1ZJLVqhgbQpNo/+Yb8Bd
V/rDz2fFwG9X80IREob+6Nifo98TcjkWPHlbkjv0zvmezKlB9pjgnppremcQoE2G
qn286kqhkghrO/DdHwFASAUqikmUMLYca7+RC6lzFfHNXOF4UrF+V43TAqEQ8get
pSmWLZMoLKgV87PtbEW5NXLSeE+iXa1622GGd9zavOynh169k66t78naWaAnqgX7
dor5W7o5AgMBAAECggEACADkxllqK9BLaC6DaXuYBCRGeBjf2y8iAOB0TgZUte/e
0EgZF0wRrUexbc+zFb0QsU4g2t998Fuxm6cRbGTz5YrzIKPx8KF5kqzqgkYC0cL3
U9NmRQ7bs3XZUYK0M3r9jV5k8LfXZcUIxjMpO+DXcQ3Tv0lsxt/dDwjX+jUkCqHY
Y4upSkjl2jzurXyyvLv7kfPsD+s3R2lcz2NtiY7XZ/y82jJAE8ZO70Pi5hma4DJQ
euawkUehNItkydSUo6KZbBttU3pt+nNqj90D0oUb0Sl2biujJPdJJX9Dq/uJJHiW
z/KOrzJIER75gqusola/9lN1e/J2FKie8WuiiixGlQKBgQDh9w3LxqjSDUVG0Evz
QAG+lYLtvdYbGRZhDICWzArhmh3fIxgR50WadTxzBfGP8VQWfSZ+EcfPAhkEwwAD
XQPviMcmwL0YgOE2XXs35Kr850WVtIEpOqZQqGemS27bXqwIywUeFUyQ5OctdEJl
DnUK/xl/CUYywZ2Ls9cbcU8bbQKBgQC8OkH6bfAIEo52DLQcZCcgi6OLgoRHePC0
nYgF+X10UoixlJkeZdZkdLVvAiuS9gGtwjbFnvRhOoAxct84OAi8vJNCH38lbPu2
nc/XAzM28t5vv0X5c93VLt7e3KAFY0FH0DrsiACYMf/H0PinkqrFq0aM1p5fwP9Q
FQRwQazufQKBgQC31cztijPyoEVKNVB1GA/TQ8P/M0CrTx+72PYMuPfpTv8aeGyu
tB8WaGbDlYRPfSDSIwNb8Y9DRQuhqhuqoNQA3qBXUNsGwmN3XVpPwMOzeVxNTUr/
he2lFT0uN5R6+GyxwqnpLZ7bCr9hZYJWwQpL5fqSNbNcu9Q2whsxAmA/iQKBgQCN
EKw53wKxShbyafrh/D1GquBaweoZFo5vDlDPCXf4IZLIY7GNkozmpIEFPP8jGLOR
Yahi2woThCBm7sxT+cqyiDFksO49QjwzVHpbjc5oNAR4g0UR+sAZ8RKeu4JCB2z5
QRmoAxVO+snTGs3/6G+LzR0GmCIBaUbu4ZF9//p2kQKBgEoJiOlHedGzX0NSjAmq
Y7OSNmsdocAXBLHi2GWrdwFvlHxZIDad8/NtM7jdy0wS2WYGlmCik7qASGZWVhfn
iSU8UKNGgQXBDkDhpo60TlItnpOSjVWl8hT2TGMi4IilkmRI4gO+faKlqd4lC9/m
vRW5n/iormjNUUjBq2lv8uqt
-----END PRIVATE KEY-----`

	fmt.Println("creating ", out)
	return writeFile(out, CONTENT)
}

func writeFile(out string, content string) error {
	// create .env file
	f, err := os.Create(out)
	fmt.Println("creating SERVER_CRT 3 ", out)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	return err
}
