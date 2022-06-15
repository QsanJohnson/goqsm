package goqsm

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
)

type testConfig struct {
	ip       string
	user     string
	passwd   string
	scId     string
	systemOp *SystemOp
	volumeOp *VolumeOp
}

var testConf *testConfig

func TestMain(m *testing.M) {
	fmt.Println("------------Start of TestMain--------------")
	flag.Parse()

	logLevelStr := os.Getenv("GOQSM_LOG_LEVEL")
	logLevel, _ := strconv.Atoi(logLevelStr)
	if logLevel > 0 {
		flag.Set("alsologtostderr", "true")
		flag.Set("v", logLevelStr)
	}

	testProp, err := readTestConf("test.conf")
	if err != nil {
		panic("The system cannot find the file: test.conf")
	}

	ctx := context.Background()

	testConf = &testConfig{}
	testConf.ip = testProp["QSM_IP"]
	testConf.user = testProp["QSM_USERNAME"]
	testConf.passwd = testProp["QSM_PASSWORD"]
	testConf.scId = testProp["TEST_SC_ID"]
	fmt.Printf("TestConf: %s %s/%s\n", testConf.ip, testConf.user, testConf.passwd)

	testClient := getTestClient(testConf.ip)
	testAuthClient, err := testClient.GetAuthClient(ctx, testConf.user, testConf.passwd)
	if err != nil {
		panic(fmt.Sprintf("GetAuthClient failed: %v \n", err))
	}

	testConf.systemOp = NewSystem(testClient)
	testConf.volumeOp = NewVolume(testAuthClient)

	code := m.Run()
	fmt.Println("------------End of TestMain--------------")
	os.Exit(code)
}

func getTestClient(ip string) *Client {
	return NewClient(ip)
}

func readTestConf(filename string) (map[string]string, error) {
	configPropertiesMap := map[string]string{}
	if len(filename) == 0 {
		return nil, errors.New("Error reading conf file " + filename)
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')

		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				configPropertiesMap[key] = value
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	return configPropertiesMap, nil
}
