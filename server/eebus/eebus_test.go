package eebus

import (
	"testing"
	"time"

	eebusapi "github.com/enbility/eebus-go/api"
	eebusmocks "github.com/enbility/eebus-go/mocks"
	shipapi "github.com/enbility/ship-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/evcc-io/evcc/util"
	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v4"
)

func TestConfig(t *testing.T) {
	conf := `
certificate:
  private: |
    -----BEGIN EC PRIVATE KEY-----
    MHcCfoo==
    -----END EC PRIVATE KEY-----
  public: |
    -----BEGIN CERTIFICATE-----
    MIIBbar=
    -----END CERTIFICATE-----
`

	var res Config
	require.NoError(t, yaml.Unmarshal([]byte(conf), &res))
}

func TestServicePairingDetailUpdate(t *testing.T) {
	identity := shipapi.NewServiceIdentity("aabbcc", "fingerprint", "paired")
	dev := &mockDevice{}

	for _, tc := range []struct {
		name    string
		paired  []shipapi.ServiceIdentity
		clients map[string][]Device
		cancel  bool
	}{
		{
			name:    "configured ski",
			clients: map[string][]Device{identity.SKI: {dev}},
		},
		{
			name:    "paired fingerprint before ski discovery",
			paired:  []shipapi.ServiceIdentity{shipapi.NewServiceIdentity("", identity.Fingerprint, "paired")},
			clients: map[string][]Device{"": {dev}},
		},
		{
			name:   "paired ski after restart",
			paired: []shipapi.ServiceIdentity{shipapi.NewServiceIdentity(identity.SKI, "", "")},
		},
		{
			name:   "unknown ski",
			cancel: true,
		},
		{
			name:    "unknown ski with paired consumer",
			clients: map[string][]Device{"": {dev}},
			cancel:  true,
		},
		{
			name:    "different paired device",
			paired:  []shipapi.ServiceIdentity{shipapi.NewServiceIdentity("ddeeff", "other", "other")},
			clients: map[string][]Device{"": {dev}},
			cancel:  true,
		},
		{
			name:    "same ship id with different certificate",
			paired:  []shipapi.ServiceIdentity{shipapi.NewServiceIdentity("ddeeff", "other", identity.ShipID)},
			clients: map[string][]Device{"": {dev}},
			cancel:  true,
		},
		{
			name:    "same ski with different fingerprint",
			paired:  []shipapi.ServiceIdentity{shipapi.NewServiceIdentity(identity.SKI, "other", identity.ShipID)},
			clients: map[string][]Device{"": {dev}},
			cancel:  true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			service := eebusmocks.NewServiceInterface(t)
			c := &EEBus{service: service, paired: tc.paired, clients: tc.clients}
			detail := shipapi.NewConnectionStateDetail(shipapi.ConnectionStateReceivedPairingRequest, nil)
			if tc.cancel {
				service.EXPECT().CancelPairing(identity).Run(func(shipapi.ServiceIdentity) {
					require.True(t, c.mux.TryLock(), "CancelPairing called with server mutex held")
					c.mux.Unlock()
					detail.SetState(shipapi.ConnectionStateNone)
					c.ServicePairingDetailUpdate(identity, detail)
				}).Once()
			}
			c.ServicePairingDetailUpdate(identity, detail)
		})
	}
}

// mockDevice implements Device for testing
type mockDevice struct{}

func (d *mockDevice) Connect(connected bool) {}
func (d *mockDevice) UseCaseEvent(_ spineapi.DeviceRemoteInterface, entity spineapi.EntityRemoteInterface, event eebusapi.EventType) {
}

var _ Device = (*mockDevice)(nil)

// TestUnregisterDevice_MutexNotHeldDuringShipCall is the regression guard
// for issue #28942. It asserts that c.mux is NOT held at the point
// UnregisterRemoteService is called. The pre-fix code held c.mux across that
// cross-layer call, and ship-go's synchronous HandleConnectionClosed
// callback chain re-entered connect(ski, false) on the same goroutine,
// which then deadlocked on c.mux.Lock() (Go mutexes are non-reentrant).
//
// The assertion uses a goroutine that tries to briefly acquire c.mux from
// inside the mock's UnregisterRemoteService implementation; if the lock is
// held, the acquisition times out and the test fails.
func TestUnregisterDevice_MutexNotHeldDuringShipCall(t *testing.T) {
	dev := &mockDevice{}
	c := &EEBus{
		log:     util.NewLogger("test"),
		clients: map[string][]Device{"aabbcc": {dev}},
	}

	service := eebusmocks.NewServiceInterface(t)
	service.EXPECT().UnregisterRemoteService(shipapi.NewServiceIdentity("aabbcc", "", "")).Run(func(shipapi.ServiceIdentity) {
		acquired := make(chan struct{})
		go func() {
			c.mux.Lock()
			defer c.mux.Unlock()
			close(acquired)
		}()
		select {
		case <-acquired:
			// good — mutex was free
		case <-time.After(100 * time.Millisecond):
			t.Errorf("c.mux was held while UnregisterRemoteSKI was called — " +
				"regression to the cross-layer lock hold that caused #28942")
		}
	}).Once()
	c.service = service

	c.UnregisterDevice("aabbcc", dev)
}
