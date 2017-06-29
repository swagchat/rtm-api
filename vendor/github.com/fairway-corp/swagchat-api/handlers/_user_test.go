package handlers

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/fairway-corp/swagchat-api/utils"
)

var createUserIds []string

type userStruct struct {
	UserId string `json:"userId,omitempty"`
}

/*
func TestPostUser(t *testing.T) {
	useDbForMain()

	ts := httptest.NewServer(Mux)
	defer ts.Close()

	testTable := []testRecord{
		{
			testNo: 1,
			in: `
				{
					"name": "dennis"
				}
			`,
			out:            `(?m)^{"id":[0-9]+,"userId":"[a-z0-9-]+","name":"dennis","unreadCount":0,"customData":{},"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 201,
		},
		{
			testNo: 2,
			in: `
				{
					"name": "dennis",
					"pictureUrl": "http://localhost/images/dennis.png"
				}
			`,
			out:            `(?m)^{"id":[0-9]+,"userId":"[a-z0-9-]+","name":"dennis","pictureUrl":"http://localhost/images/dennis.png","unreadCount":0,"customData":{},"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 201,
		},
		{
			testNo: 3,
			in: `
				{
					"name": "dennis",
					"pictureUrl": "http://localhost/images/dennis.png",
					"informationUrl": "http://localhost/dennis"
				}
			`,
			out:            `(?m)^{"id":[0-9]+,"userId":"[a-z0-9-]+","name":"dennis","pictureUrl":"http://localhost/images/dennis.png","informationUrl":"http://localhost/dennis","unreadCount":0,"customData":{},"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 201,
		},
		{
			testNo: 4,
			in: `
				{
					"name": "dennis",
					"pictureUrl": "http://localhost/images/dennis.png",
					"informationUrl": "http://localhost/dennis",
					"customData": {"key": "value"}
				}
			`,
			out:            `(?m)^{"id":[0-9]+,"userId":"[a-z0-9-]+","name":"dennis","pictureUrl":"http://localhost/images/dennis.png","informationUrl":"http://localhost/dennis","unreadCount":0,"customData":{"key":"value"},"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 201,
		},
		{
			testNo: 5,
			in: `
				{
					"userId": "custom-id",
					"name": "dennis"
				}
			`,
			out:            `(?m)^{"id":[0-9]+,"userId":"custom-id","name":"dennis","unreadCount":0,"customData":{},"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 201,
		},
		{
			testNo: 6,
			in: `
				{
					"userId": "custom-id",
					"name": "dennis"
				}
			`,
			out:            `(?m)^{"title":"An error occurred while creating user item.","status":500,"detail":".*","errorName":"database-error"}$`,
			httpStatusCode: 500,
		},
		{
			testNo: 7,
			in: `
			`,
			out:            `(?m)^{"title":"Json for user item creation is parse error.","status":400,"errorName":"invalid-json"}$`,
			httpStatusCode: 400,
		},
		{
			testNo: 8,
			in: `
				json
			`,
			out:            `(?m)^{"title":"Json for user item creation is parse error.","status":400,"errorName":"invalid-json"}$`,
			httpStatusCode: 400,
		},
		{
			testNo: 9,
			in: `
				{
					"pictureUrl": "http://example.com/picture.png"
				}
			`,
			out:            `(?m)^{"title":"Request parameter for user item creation is invalid.","status":400,"errorName":"invalid-param","invalidParams":\[{"name":"name","reason":"name is required, but it's empty."}\]}$`,
			httpStatusCode: 400,
		},
	}

	for _, testRecord := range testTable {
		reader := strings.NewReader(testRecord.in)
		res, err := http.Post(ts.URL+"/"+utils.API_VERSION+"/users", "application/json", reader)

		if err != nil {
			t.Fatalf("TestNo %d\nhttp request failed: %v", testRecord.testNo, err)
		}

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("TestNo %d\nError by ioutil.ReadAll(): %v", testRecord.testNo, err)
		}

		if res.StatusCode != testRecord.httpStatusCode {
			t.Fatalf("TestNo %d\nHTTP Status Code Failure\n[expected]%d\n[result  ]%d", testRecord.testNo, testRecord.httpStatusCode, res.StatusCode)
		}

		r := regexp.MustCompile(testRecord.out)
		if !r.MatchString(string(data)) {
			t.Fatalf("TestNo %d\nResponse Body Failure\n[expected]%s\n[result  ]%s", testRecord.testNo, testRecord.out, string(data))
		}

		if testRecord.httpStatusCode == 201 {
			user := &userStruct{}
			_ = json.Unmarshal(data, user)
			createUserIds = append(createUserIds, user.UserId)
		}
	}
}

func TestGetUser(t *testing.T) {
	useDbForMain()

	ts := httptest.NewServer(Mux)
	defer ts.Close()

	if len(createUserIds) != 5 {
		t.Fatalf("createUserIds length error \n[expected]%d\n[result  ]%d", 5, len(createUserIds))
		t.Failed()
	}

	testTable := []testRecord{
		{
			testNo:         1,
			in:             createUserIds[0],
			out:            `(?m)^{"id":[0-9]+,"userId":"[a-z0-9-]+","name":"dennis","unreadCount":0,"customData":{},"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         2,
			in:             createUserIds[1],
			out:            `(?m)^{"id":[0-9]+,"userId":"[a-z0-9-]+","name":"dennis","pictureUrl":"http://localhost/images/dennis.png","unreadCount":0,"customData":{},"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         3,
			in:             createUserIds[2],
			out:            `(?m)^{"id":[0-9]+,"userId":"[a-z0-9-]+","name":"dennis","pictureUrl":"http://localhost/images/dennis.png","informationUrl":"http://localhost/dennis","unreadCount":0,"customData":{},"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         4,
			in:             createUserIds[3],
			out:            `(?m)^{"id":[0-9]+,"userId":"[a-z0-9-]+","name":"dennis","pictureUrl":"http://localhost/images/dennis.png","informationUrl":"http://localhost/dennis","unreadCount":0,"customData":{"key":"value"},"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         5,
			in:             createUserIds[4],
			out:            `(?m)^{"id":[0-9]+,"userId":"custom-id","name":"dennis","unreadCount":0,"customData":{},"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 200,
		},
		{
			testNo:         6,
			in:             "not-exist-user-id",
			out:            `(?m)^{"title":"An error occurred while getting user item.","status":500,"detail":".*","errorName":"database-error"}$`,
			httpStatusCode: 500,
		},
	}

	for _, testRecord := range testTable {
		res, err := http.Get(ts.URL + "/" + utils.API_VERSION + "/users/" + testRecord.in)

		if err != nil {
			t.Fatalf("TestNo %d\nhttp request failed: %v", testRecord.testNo, err)
		}

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("TestNo %d\nError by ioutil.ReadAll(): %v", testRecord.testNo, err)
		}

		if res.StatusCode != testRecord.httpStatusCode {
			t.Fatalf("TestNo %d\nHTTP Status Code Failure\n[expected]%d\n[result  ]%d", testRecord.testNo, testRecord.httpStatusCode, res.StatusCode)
		}

		r := regexp.MustCompile(testRecord.out)
		if !r.MatchString(string(data)) {
			t.Fatalf("TestNo %d\nResponse Body Failure\n[expected]%s\n[result  ]%s", testRecord.testNo, testRecord.out, string(data))
		}
	}
}

func TestPutUser(t *testing.T) {
	useDbForMain()

	ts := httptest.NewServer(Mux)
	defer ts.Close()

	testTable := []testRecord{
		{
			testNo: 1,
			in: `
				{
					"name": "Jeremy"
				}
			`,
			out:            `(?m)^{"id":[0-9]+,"userId":"custom-id","name":"Jeremy","unreadCount":0,"customData":{},"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 204,
		},
		{
			testNo: 2,
			in: `
				{
					"pictureUrl": "http://localhost/images/jeremy.png"
				}
			`,
			out:            `(?m)^{"id":[0-9]+,"userId":"custom-id","name":"Jeremy","pictureUrl":"http://localhost/images/jeremy.png","unreadCount":0,"customData":{},"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 204,
		},
		{
			testNo: 3,
			in: `
				{
					"informationUrl": "http://localhost/jeremy"
				}
			`,
			out:            `(?m)^{"id":[0-9]+,"userId":"custom-id","name":"Jeremy","pictureUrl":"http://localhost/images/jeremy.png","informationUrl":"http://localhost/jeremy","unreadCount":0,"customData":{},"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 204,
		},
		{
			testNo: 4,
			in: `
				{
					"customData": {"key": "value"}
				}
			`,
			out:            `(?m)^{"id":[0-9]+,"userId":"custom-id","name":"Jeremy","pictureUrl":"http://localhost/images/jeremy.png","informationUrl":"http://localhost/jeremy","unreadCount":0,"customData":{"key":"value"},"created":[0-9]+,"modified":[0-9]+}$`,
			httpStatusCode: 204,
		},
		{
			testNo: 5,
			in: `
			`,
			out:            `(?m)^{"title":"Json for user item updating is parse error.","status":400,"errorName":"invalid-json"}$`,
			httpStatusCode: 400,
		},
		{
			testNo: 6,
			in: `
				json
			`,
			out:            `(?m)^{"title":"Json for user item updating is parse error.","status":400,"errorName":"invalid-json"}$`,
			httpStatusCode: 400,
		},
	}

	for _, testRecord := range testTable {
		reader := strings.NewReader(testRecord.in)
		req, _ := http.NewRequest("PUT", ts.URL+"/"+utils.API_VERSION+"/users/custom-id", reader)
		req.Header.Set("Content-Type", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("TestNo %d\nhttp request failed: %v", testRecord.testNo, err)
		}

		if res.StatusCode != testRecord.httpStatusCode {
			t.Fatalf("TestNo %d\nHTTP Status Code Failure\n[expected]%d\n[result  ]%d", testRecord.testNo, testRecord.httpStatusCode, res.StatusCode)
		}

		if testRecord.httpStatusCode == 204 {
			res, err = http.Get(ts.URL + "/" + utils.API_VERSION + "/users/custom-id")
		}

		data, err := ioutil.ReadAll(res.Body)
		r := regexp.MustCompile(testRecord.out)
		if !r.MatchString(string(data)) {
			t.Fatalf("TestNo %d\nResponse Body Failure\n[expected]%s\n[result  ]%s", testRecord.testNo, testRecord.out, string(data))
		}
	}
}

func TestDeleteUser(t *testing.T) {
	useDbForMain()

	ts := httptest.NewServer(Mux)
	defer ts.Close()

	testTable := []testRecord{
		{
			testNo:         1,
			in:             "custom-id",
			out:            `(?m)^$`,
			httpStatusCode: 204,
		},
	}

	for _, testRecord := range testTable {
		req, _ := http.NewRequest("DELETE", ts.URL+"/"+utils.API_VERSION+"/users/"+testRecord.in, nil)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("TestNo %d\nhttp request failed: %v", testRecord.testNo, err)
		}

		if res.StatusCode != testRecord.httpStatusCode {
			t.Fatalf("TestNo %d\nHTTP Status Code Failure\n[expected]%d\n[result  ]%d", testRecord.testNo, testRecord.httpStatusCode, res.StatusCode)
		}

		if testRecord.httpStatusCode == 204 {
			res, err = http.Get(ts.URL + "/" + utils.API_VERSION + "/users/" + testRecord.in)
		}

		data, err := ioutil.ReadAll(res.Body)
		r := regexp.MustCompile(testRecord.out)
		if !r.MatchString(string(data)) {
			t.Fatalf("TestNo %d\nResponse Body Failure\n[expected]%s\n[result  ]%s", testRecord.testNo, testRecord.out, string(data))
		}
	}
}

func TestGetUsers(t *testing.T) {
	useDbForMain()

	ts := httptest.NewServer(Mux)
	defer ts.Close()

	testTable := []testRecord{
		{
			testNo:         1,
			out:            `(?m)^{"users":[{"id":[0-9]+,"userId":"[a-z0-9-]+","name":"dennis","unreadCount":0,"customData":{},"created":[0-9]+,"modified":[0-9]+},{"id":[0-9]+,"userId":"[a-z0-9-]+","name":"dennis","pictureUrl":"http://localhost/images/dennis.png","unreadCount":0,"customData":{},"created":[0-9]+,"modified":[0-9]+},{"id":[0-9]+,"userId":"[a-z0-9-]+","name":"dennis","pictureUrl":"http://localhost/images/dennis.png","informationUrl":"http://localhost/dennis","unreadCount":0,"customData":{},"created":[0-9]+,"modified":[0-9]+},{"id":[0-9]+,"userId":"[a-z0-9-]+","name":"dennis","pictureUrl":"http://localhost/images/dennis.png","informationUrl":"http://localhost/dennis","unreadCount":0,"customData":{"key":"value"},"created":[0-9]+,"modified":[0-9]+}]}$`,
			httpStatusCode: 200,
		},
	}

	for _, testRecord := range testTable {
		res, err := http.Get(ts.URL + "/" + utils.API_VERSION + "/users")
		if err != nil {
			t.Fatalf("TestNo %d\nhttp request failed: %v", testRecord.testNo, err)
		}

		if res.StatusCode != testRecord.httpStatusCode {
			t.Fatalf("TestNo %d\nHTTP Status Code Failure\n[expected]%d\n[result  ]%d", testRecord.testNo, testRecord.httpStatusCode, res.StatusCode)
		}

		data, err := ioutil.ReadAll(res.Body)
		r := regexp.MustCompile(testRecord.out)
		if !r.MatchString(string(data)) {
			t.Fatalf("TestNo %d\nResponse Body Failure\n[expected]%s\n[result  ]%s", testRecord.testNo, testRecord.out, string(data))
		}
	}
}
*/
func TestGetUserRooms(t *testing.T) {
	useDbForTestGetUserRooms()
	log.Println(utils.Cfg.Sqlite.DatabasePath)

	ts := httptest.NewServer(Mux)
	defer ts.Close()

	testTable := []testRecord{
		{
			testNo:         1,
			in:             "custom-user-id-0001",
			out:            `(?m)^$`,
			httpStatusCode: 200,
		},
	}

	for _, testRecord := range testTable {
		res, err := http.Get(ts.URL + "/" + utils.API_VERSION + "/users/" + testRecord.in + "/rooms")
		if err != nil {
			t.Fatalf("TestNo %d\nhttp request failed: %v", testRecord.testNo, err)
		}

		if res.StatusCode != testRecord.httpStatusCode {
			t.Fatalf("TestNo %d\nHTTP Status Code Failure\n[expected]%d\n[result  ]%d", testRecord.testNo, testRecord.httpStatusCode, res.StatusCode)
		}

		data, err := ioutil.ReadAll(res.Body)
		r := regexp.MustCompile(testRecord.out)
		if !r.MatchString(string(data)) {
			t.Fatalf("TestNo %d\nResponse Body Failure\n[expected]%s\n[result  ]%s", testRecord.testNo, testRecord.out, string(data))
		}
	}
}
