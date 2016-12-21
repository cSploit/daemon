package rogue_ap

import (
	"testing"
	"github.com/stretchr/testify/require"
)

func TestRogueAP_Start(t *testing.T) {
	rogue_bad_bssid := RogueAP{
		BSSID: "IAmNotAMacAddress",
	}

	err := rogue_bad_bssid.Start()
	require.NotNil(t, err)

	rogue_bad_mac_deny := RogueAP{
		DenyMac:[]string{"11:22:33:44:55:66", "NotAMacAddr"},
	}

	err = rogue_bad_mac_deny.Start()
	require.NotNil(t, err)

	rogue_bad_mac_allow := RogueAP{
		AllowMac:[]string{"11:22:33:44:55:66", "NotAMacAddr"},
	}

	err = rogue_bad_mac_allow.Start()
	require.NotNil(t, err)

	rogue_no_iface := RogueAP{}

	err = rogue_no_iface.Start()
	require.NotNil(t, err)
}
