package tguard

import (
	"context"
	"testing"
	"time"
)

type TestData struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func TestTimeGuard(t *testing.T) {

	identityChecker := func(id string, data TestData) bool {
		return id == data.Id
	}
	var tests = []struct {
		name     string
		ttl      time.Duration
		interval time.Duration
		sleep    time.Duration
		cancel   bool
		data     []TestData
	}{
		{
			name:     "test timeout",
			ttl:      time.Millisecond * 200,
			interval: time.Millisecond * 100,
			sleep:    time.Millisecond * 700,
			cancel:   false,
			data: []TestData{
				{
					Id:   "1",
					Name: "test",
				},
			},
		},

		{
			name:     "test cancel before timeout",
			ttl:      time.Millisecond * 200,
			interval: time.Millisecond * 100,
			cancel:   true,
			data: []TestData{
				{
					Id:   "2",
					Name: "test",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := New(Config[TestData]{
				Fallback:        func(data TestData) {},
				IdentityChecker: identityChecker,
				DefaultTTL:      tt.ttl,
				Interval:        tt.interval,
			})
			ctx := context.Background()
			go g.Connect(ctx)
			for _, v := range tt.data {
				_ = g.Start(ctx, v)
			}
			if tt.sleep > 0 {
				time.Sleep(tt.sleep)
			}
			if tt.cancel {
				for _, v := range tt.data {
					err := g.Cancel(ctx, v.Id)
					if err != nil {
						t.Error(err)
					}
				}
			}

		})
	}

}
