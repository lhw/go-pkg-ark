package arkrcon

import "testing"

func TestPlayerList(t *testing.T) {
	ark, err := newARKRconConnectionEnv()
	if err != nil {
		t.Skip("No connection available: ", err)
	}

	if playerList, err := ark.ListPlayers(); err != nil {
		t.Error(err)
	} else {
		t.Log("Players online: ", len(playerList))
	}
}
