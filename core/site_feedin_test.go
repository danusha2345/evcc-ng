package core

import (
	"testing"
	"time"

	"github.com/evcc-io/evcc/api"
	"github.com/evcc-io/evcc/tariff"
	"github.com/evcc-io/evcc/util"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func mockFeedInTariff(t *testing.T, price float64) api.Tariff {
	t.Helper()
	ctrl := gomock.NewController(t)
	tf := api.NewMockTariff(ctrl)
	now := time.Now()
	tf.EXPECT().Rates().Return(api.Rates{{
		Start: now.Add(-time.Hour),
		End:   now.Add(time.Hour),
		Value: price,
	}}, nil).AnyTimes()
	tf.EXPECT().Type().Return(api.TariffTypePriceDynamic).AnyTimes()
	return tf
}

func TestShouldFeedInCurtail(t *testing.T) {
	for _, tc := range []struct {
		name      string
		enabled   bool
		threshold float64
		tariff    api.Tariff
		want      bool
	}{
		{name: "disabled", enabled: false, want: false},
		{name: "enabled no tariff", enabled: true, want: false},
		{name: "price above threshold", enabled: true, threshold: 0, tariff: mockFeedInTariff(t, 0.05), want: false},
		{name: "price at threshold (curtail)", enabled: true, threshold: 0, tariff: mockFeedInTariff(t, 0), want: true},
		{name: "price below zero threshold", enabled: true, threshold: 0, tariff: mockFeedInTariff(t, -0.02), want: true},
		{name: "negative threshold not yet hit", enabled: true, threshold: -0.05, tariff: mockFeedInTariff(t, -0.02), want: false},
		{name: "negative threshold hit", enabled: true, threshold: -0.05, tariff: mockFeedInTariff(t, -0.10), want: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := &Site{
				log: util.NewLogger("test"),
				tariffs: &tariff.Tariffs{
					FeedIn: tc.tariff,
				},
				feedInControl:          tc.enabled,
				feedInControlThreshold: tc.threshold,
			}
			assert.Equal(t, tc.want, s.shouldFeedInCurtail())
		})
	}
}
