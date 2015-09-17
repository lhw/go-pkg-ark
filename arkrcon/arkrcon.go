// Package ark provides the basic RCON commands for an ARK Surival Server
package arkrcon

import (
	"errors"
	"fmt"
	"github.com/james4k/rcon"
	"log"
	"regexp"
	"strings"
)

type ARKRcon struct {
	rc      *rcon.RemoteConsole
	address string
}

type ARKPlayer struct {
	Username string
	Steam64  string
}

type ARKChatMsg struct {
	Username      string
	Playername    string
	Message       string
	ServerMessage bool
}

var (
	EmptyResponse = errors.New("No Server Response")
	FailResponse  = errors.New("Server failed at request")
)

/*
  All command information based on:
  http://steamcommunity.com/sharedfiles/filedetails/?id=454529617
*/

// ListPlayers returns a list of online players or an empty list
func (a *ARKRcon) ListPlayers() ([]ARKPlayer, error) {
	/* CMD: listplayers
	    Success:
	   - No Players Connected
	   - 0. CyFreeze, 76561198025588951
	     ...
	*/
	resp, err := a.Query("listplayers")
	if err != nil {
		return nil, err
	}
	rex := regexp.MustCompile(`\d+\. ([^,]+), (\d+)`)
	list := make([]ARKPlayer, 0)
	all := rex.FindAllStringSubmatch(resp, -1)
	for _, m := range all {
		list = append(list, ARKPlayer{m[1], m[2]})
	}
	return list, nil
}

func (a *ARKRcon) SaveWorld() error {
	/* CMD: saveworld
	   Success: World Saved
	*/
	return a.simpleResponse("saveworld", "World Saved")
}

func (a *ARKRcon) DoExit() error {
	/* CMD: doexit
	   Success: Exiting...
	*/
	return a.simpleResponse("doexit", "Exiting")
}

func (a *ARKRcon) SendChatToPlayer(player, message string) error {
	/* CMD: serverchattoplayer "player" "msg"
	   Success: /
	*/
	return a.emptyResponse(fmt.Sprintf(`serverchattoplayer "%s" "%s"`, player, message))
}

func (a *ARKRcon) SendChatToID(steam64, message string) error {
	/* CMD: serverchatto "steam64" "msg"
	   Success: /
	*/
	return a.emptyResponse(fmt.Sprintf(`serverchatto "%s" "%s"`, steam64, message))
}

// GetChat returns a list of chat messages since the last call to getchat or
// an empty list in case there were none
func (a *ARKRcon) GetChat() ([]ARKChatMsg, error) {
	/* CMD: getchat
	   Success: - SERVER: foo
	              CyFreeze (Bob The Builder): foobar
	              Valki(Valki): wup wup
	*/
	resp, err := a.Query("getchat")
	if err != nil {
		return nil, err
	}
	rex := regexp.MustCompile(`(\w+)\s*(?:\(([\w\s]+)\))?:\s*(.*?)$`)
	list := make([]ARKChatMsg, 0)
	all := rex.FindAllStringSubmatch(resp, -1)
	for _, m := range all {
		list = append(list, ARKChatMsg{m[1], m[2], m[3], strings.HasPrefix(m[1], "SERVER")})
	}
	return list, nil
}

// SetTimeOfDay expects the time format to be hh:mm
func (a *ARKRcon) SetTimeOfDay(time string) error {
	/* CMD: settimeofday
	   Success: /
	*/
	return a.emptyResponse(fmt.Sprintf(`settimeofday %s`, time))
}

func (a *ARKRcon) WhitelistPlayer(steam64 string) error {
	/* CMD: allowplayertojoinnocheck steam64
	   Success: <steam64> Allow Player To Join No Check
	*/
	return a.simpleResponse(fmt.Sprintf(`allowplayertojoinnocheck %s`, steam64), fmt.Sprintf(`%s Allow`, steam64))
}

func (a *ARKRcon) RemoveWhitelist(steam64 string) error {
	/* CMD: disallowplayertojoinnocheck steam64
	   Success: <steam64> Disallowed Player To Join No Checknned
	*/
	return a.simpleResponse(fmt.Sprintf(`disallowplayertojoinnocheck %s`, steam64), fmt.Sprintf(`%s Disallowed`, steam64))
}

func (a *ARKRcon) SetMessageOfTheDay(motd string) error {
	/* CMD: setmessageoftheday motd
	   Success: Message of set to <motd>
	*/
	return a.simpleResponse(fmt.Sprintf(`setmessageoftheday %s`, motd), "Message of set to")
}

func (a *ARKRcon) Broadcast(message string) error {
	/* CMD: broadcast
	   Success: /
	*/
	return a.emptyResponse(fmt.Sprintf(`broadcast %s`, message))
}

func (a *ARKRcon) KickPlayer(steam64 string) error {
	/* CMD: kickplayer steam64
	   Success: <steam64> Kicked
	*/
	return a.simpleResponse(fmt.Sprintf(`kickplayer %s`, steam64), fmt.Sprintf(`%s Kicked`, steam64))
}

func (a *ARKRcon) BanPlayer(steam64 string) error {
	/* CMD: banplayer steam64
	   Success: <steam64> Banned
	*/
	return a.simpleResponse(fmt.Sprintf(`banplayer %s`, steam64), fmt.Sprintf(`%s Banned`, steam64))
}

func (a *ARKRcon) UnbanPlayer(steam64 string) error {
	/* CMD: unbanplayer steam64
	   Success: <steam64> Unbanned
	*/
	return a.simpleResponse(fmt.Sprintf(`unbanplayer %s`, steam64), fmt.Sprintf(`%s Unbanned`, steam64))
}

// Slomo modifier. Set to 1 to return to normal
func (a *ARKRcon) Slomo(multiplier int) error {
	/* CMD: slomo multiplier
	   Success: /
	*/
	return a.emptyResponse(fmt.Sprintf(`slomo %d`, multiplier))
}

// Dinos eventually respawn
func (a *ARKRcon) DestroyWildDinos() error {
	/* CMD: destroywilddinos
	   Success: /
	*/
	return a.emptyResponse("destroywilddinos")
}

// Same as DestroyWildDinos but also kills tamed Dinos
func (a *ARKRcon) DestroyAllEnemies() error {
	/* CMD: destroyallenemies
	   Success: /
	*/
	return a.emptyResponse("destroyallenemies")
}

/*
  No idea how to get ark player id yet
  Just keep them in mind for now
*/
func (a *ARKRcon) giveItemToPlayer(playerID, itemID, quantity, quality int, blueprint bool) {
	//giveitemnumtoplayer
}

func (a *ARKRcon) clearPlayerInventory(playerID int, clrInv, clrSlot, clrEquip bool) {
	//clearplayerinventory
}

func (a *ARKRcon) killPlayer(playerID int) {
	//killplayer
}

func (a *ARKRcon) giveExpToPlayer(playerID, exp int, fromtribe, preventshare bool) {
	//giveexptoplayer
}

func (a *ARKRcon) forcePlayerToJoinTribe(playerID, tribeID int) {
	//forceplayertojointribe
}

func (a *ARKRcon) emptyResponse(cmd string) error {
	_, err := a.Query(cmd)
	if err == EmptyResponse {
		return nil
	} else {
		return err
	}
}

func (a *ARKRcon) simpleResponse(cmd, exp string) error {
	resp, err := a.Query(cmd)
	if err != nil {
		return err
	}
	if !strings.Contains(resp, exp) {
		return FailResponse
	}
	return nil
}

func (a *ARKRcon) Query(cmd string) (string, error) {
	reqID, reqErr := a.rc.Write(cmd)
	if reqErr != nil {
		log.Println(reqID, reqErr)
		return "", reqErr
	}

	resp, respID, respErr := a.rc.Read()
	if respErr != nil {
		log.Println(resp, respID, respErr)
		return "", respErr
	}

	if strings.Contains(resp, "no response!!") {
		return "", EmptyResponse
	}
	return resp, nil
}

func NewARKRconConnection(address, password string) (*ARKRcon, error) {
	rc, err := rcon.Dial(address, password)
	if err != nil {
		return nil, err
	}
	return &ARKRcon{rc, address}, nil
}
