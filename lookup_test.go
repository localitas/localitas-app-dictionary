package dictionary

import (
	"testing"
)

// Lookups ask through the SDK's upstream client, which is pooled and pings
// an idle HTTP/2 connection a NAT may have dropped; a plain client of their
// own would keep a dead connection and stall every try on it.
func TestLookupsUseTheSDKClient(t *testing.T) {
	if upstream.Client != nil {
		t.Fatal("lookups bring their own client, without the SDK's stale-connection check")
	}
}
