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

	t.Run("removeSoftItem test", func(t *testing.T) {
		slice := []Data[TestData]{
			{
				ExpireTime: time.Now().Add(time.Millisecond * 50).Unix(),
				Original: TestData{
					Id:   "1",
					Name: "test",
				},
			},
			{
				ExpireTime: time.Now().Add(time.Millisecond * 30).Unix(),
				Original: TestData{
					Id:   "2",
					Name: "test 2",
				},
			},
			{
				ExpireTime: time.Now().Add(time.Millisecond * 30).Unix(),
				Original: TestData{
					Id:   "3",
					Name: "test 3",
				},
			},
		}
		slice, maxIdx := removeSoftItem(slice, 0, len(slice)-1)
		if len(slice) != 2 {
			t.Errorf("invalid slice length")
		}
		if maxIdx != 1 {
			t.Errorf("invalid maxIdx")
		}
	})

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
		{
			name:     "test timeout with multiple data",
			ttl:      time.Millisecond * 200,
			interval: time.Millisecond * 100,
			sleep:    time.Millisecond * 6000,
			cancel:   false,
			data: []TestData{
				{
					Id:   "2",
					Name: "test",
				},
				{
					Id:   "3",
					Name: "test 2",
				},
				{
					Id:   "4",
					Name: "test 3",
				},
				{
					Id:   "5",
					Name: "test 4",
				},
				{
					Id:   "6",
					Name: "test 5",
				},
				{
					Id:   "7",
					Name: "test 6",
				},
				{
					Id:   "8",
					Name: "test 7",
				},
				{
					Id:   "9",
					Name: "test 8",
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
