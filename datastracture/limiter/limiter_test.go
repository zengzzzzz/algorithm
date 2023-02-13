package limiter

import (
	"testing"
	"time"
)

func TestOfficalLimitFlowAllow(t *testing.T) {
	LimitFlowAllow()
}

func TestOfficalLimitFlowWait(t *testing.T) {
	LimitFlowWait()
}

func TestOfficalLimitFlowReserve(t *testing.T) {
	LimitFlowReserve()
}

func TestTokenBucketLimit(t *testing.T) {
	type args struct {
		capacity int
		rate     int
	}
	tests := []struct {
		name string
		args args
		want *TokenBucketLimiter
	}{
		{
			name: "60",
			args: args{
				capacity: 60,
				rate:     10,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewTokenBucketLimiter(tt.args.capacity, tt.args.rate)
			time.Sleep(time.Second)
			successCount := 0
			for i := 0; i < tt.args.rate; i++ {
				if l.TryAcquire() {
					successCount++
				}
			}

			if successCount != tt.args.rate {
				t.Errorf("NewTokenBucketLimiter() = %v, want %v", successCount, tt.args.rate)
				return
			}
			successCount = 0
			for i := 0; i < tt.args.capacity; i++ {
				if l.TryAcquire() {
					successCount++
				}
				time.Sleep(time.Second / 10)
			}
			if successCount != tt.args.capacity-tt.args.rate {
				t.Errorf("NewTokenBucketLimiter() = %v, want %v", successCount, tt.args.capacity-tt.args.rate)
				return
			}
		})
	}
}
