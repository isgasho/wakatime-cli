package heartbeat_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/wakatime/wakatime-cli/pkg/heartbeat"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeartbeat_ID(t *testing.T) {
	h := heartbeat.Heartbeat{
		Branch:     heartbeat.String("heartbeat"),
		Category:   heartbeat.CodingCategory,
		Entity:     "/tmp/main.go",
		EntityType: heartbeat.FileType,
		IsWrite:    heartbeat.Bool(true),
		Project:    heartbeat.String("wakatime"),
		Time:       1592868313.541149,
	}
	assert.Equal(t, "1592868313.541149-file-coding-wakatime-heartbeat-/tmp/main.go-true", h.ID())
}

func TestHeartbeat_ID_NilFields(t *testing.T) {
	h := heartbeat.Heartbeat{
		Category:   heartbeat.CodingCategory,
		Entity:     "/tmp/main.go",
		EntityType: heartbeat.FileType,
		Time:       1592868313.541149,
	}
	assert.Equal(t, "1592868313.541149-file-coding---/tmp/main.go-false", h.ID())
}

func TestHeartbeat_JSON(t *testing.T) {
	h := heartbeat.Heartbeat{
		Branch:         heartbeat.String("heartbeat"),
		Category:       heartbeat.CodingCategory,
		CursorPosition: heartbeat.Int(12),
		Dependencies:   []string{"dep1", "dep2"},
		Entity:         "/tmp/main.go",
		EntityType:     heartbeat.FileType,
		IsWrite:        heartbeat.Bool(true),
		Language:       heartbeat.String("golang"),
		LineNumber:     heartbeat.Int(42),
		Lines:          heartbeat.Int(100),
		Project:        heartbeat.String("wakatime"),
		Time:           1585598060.1,
		UserAgent:      "wakatime/13.0.7",
	}

	jsonEncoded, err := json.Marshal(h)
	require.NoError(t, err)

	f, err := os.Open("./testdata/heartbeat.json")
	require.NoError(t, err)

	defer f.Close()

	expected, err := ioutil.ReadAll(f)
	require.NoError(t, err)

	assert.JSONEq(t, string(expected), string(jsonEncoded))
}

func TestHeartbeat_JSON_NilFields(t *testing.T) {
	h := heartbeat.Heartbeat{
		Category:   heartbeat.CodingCategory,
		Entity:     "/tmp/main.go",
		EntityType: heartbeat.FileType,
		Time:       1585598060,
		UserAgent:  "wakatime/13.0.7",
	}

	jsonEncoded, err := json.Marshal(h)
	require.NoError(t, err)

	f, err := os.Open("./testdata/heartbeat_null_fields.json")
	require.NoError(t, err)

	defer f.Close()

	expected, err := ioutil.ReadAll(f)
	require.NoError(t, err)

	assert.JSONEq(t, string(expected), string(jsonEncoded))
}

func TestNewHandle(t *testing.T) {
	sender := mockSender{
		SendFn: func(hh []heartbeat.Heartbeat) ([]heartbeat.Result, error) {
			assert.Equal(t, []heartbeat.Heartbeat{
				{
					Branch:     heartbeat.String("test"),
					Category:   heartbeat.CodingCategory,
					Entity:     "/tmp/main.go",
					EntityType: heartbeat.FileType,
					Time:       1585598060,
					UserAgent:  "wakatime/13.0.7",
				},
			}, hh)
			return []heartbeat.Result{
				{
					Status:    201,
					Heartbeat: heartbeat.Heartbeat{},
				},
			}, nil
		},
	}

	opts := []heartbeat.HandleOption{
		func(next heartbeat.Handle) heartbeat.Handle {
			return func(hh []heartbeat.Heartbeat) ([]heartbeat.Result, error) {
				for i := range hh {
					hh[i].Branch = heartbeat.String("test")
				}

				return next(hh)
			}
		},
	}

	handle := heartbeat.NewHandle(&sender, opts...)
	_, err := handle([]heartbeat.Heartbeat{
		{
			Category:   heartbeat.CodingCategory,
			Entity:     "/tmp/main.go",
			EntityType: heartbeat.FileType,
			Time:       1585598060,
			UserAgent:  "wakatime/13.0.7",
		},
	})
	require.NoError(t, err)
}

type mockSender struct {
	SendFn        func(hh []heartbeat.Heartbeat) ([]heartbeat.Result, error)
	SendFnInvoked bool
}

func (m *mockSender) Send(hh []heartbeat.Heartbeat) ([]heartbeat.Result, error) {
	m.SendFnInvoked = true
	return m.SendFn(hh)
}
