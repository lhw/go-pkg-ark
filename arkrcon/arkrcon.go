// Package ark provides the basic RCON commands for an ARK Surival Server
package arkrcon

import (
	"errors"
	"fmt"
	"github.com/james4k/rcon"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type ARKRcon struct {
	rc       *rcon.RemoteConsole
	address  string
	password string
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
	NoConnection  = errors.New("No connection to RCON")
	EmptyResponse = errors.New("No Server Response")
	FailResponse  = errors.New("Server failed at executing request")
	EnvMissing    = errors.New("One or more environment Variables not set")

	//RegEx
	playerRegex = regexp.MustCompile(`\d+\. ([^,]+), (\d+)`)
	chatRegex   = regexp.MustCompile(`(\w+)\s*(?:\(([\w\s]+)\))?:\s*(.*?)$`)
)

/*
  All command information based on:
  http://steamcommunity.com/sharedfiles/filedetails/?id=454529617
*/

// ListPlayers returns a list of online players or an empty list
func (a *ARKRcon) ListPlayers() (list []ARKPlayer, err error) {
	/* CMD: listplayers
	    Success:
	   - No Players Connected
	   - 0. CyFreeze, 76561198025588951
	     ...
	*/
	var resp string
	if resp, err = a.Query("listplayers"); err != nil {
		return
	}

	list = make([]ARKPlayer, 0)
	all := playerRegex.FindAllStringSubmatch(resp, -1)
	for _, m := range all {
		list = append(list, ARKPlayer{m[1], m[2]})
	}
	return
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
func (a *ARKRcon) GetChat() (list []ARKChatMsg, err error) {
	/* CMD: getchat
	   Success: - SERVER: foo
	              CyFreeze (Bob The Builder): foobar
	              Valki(Valki): wup wup
	*/
	var resp string
	if resp, err = a.Query("getchat"); err != nil {
		return
	}

	list = make([]ARKChatMsg, 0)
	all := chatRegex.FindAllStringSubmatch(resp, -1)
	for _, m := range all {
		list = append(list, ARKChatMsg{m[1], m[2], m[3], strings.HasPrefix(m[1], "SERVER")})
	}
	return
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

func (a *ARKRcon) SpawnDino(blueprint string, x_off, y_off, z_off, level int) error {
	/* CMD: spawndino blueprintPath x_offset y_offset z_offset level
	   Success: TBD
	*/
	return a.emptyResponse(fmt.Sprintf(`spawndino %s %d %d %d %d`, blueprint, x_off, y_off, z_off, level))
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

func (a *ARKRcon) simpleResponse(cmd, exp string) (err error) {
	var resp string
	if resp, err = a.Query(cmd); err != nil {
		return err
	}
	if !strings.Contains(resp, exp) {
		return FailResponse
	}
	return
}

func (a *ARKRcon) Query(cmd string) (resp string, err error) {
	if a == nil {
		return "", NoConnection
	}

	if _, err = a.rc.Write(cmd); err != nil {
		return
	}

	if resp, _, err = a.rc.Read(); err != nil {
		return
	}

	if strings.Contains(resp, "no response!!") {
		return "", EmptyResponse
	}
	return
}

func NewARKRconConnection(address, password string) (*ARKRcon, error) {
	var err error
	var rc *rcon.RemoteConsole

	if rc, err = rcon.Dial(address, password); err != nil {
		return nil, err
	}
	return &ARKRcon{rc, address, password}, nil
}

func newARKRconConnectionEnv() (*ARKRcon, error) {
	addr := os.Getenv("ADDRESS")
	port, err := strconv.Atoi(os.Getenv("PORT"))
	pass := os.Getenv("ADMIN_PASSWORD")
	if addr == "" || pass == "" || err != nil {
		return nil, EnvMissing
	}
	return NewARKRconConnection(fmt.Sprintf("%s:%d", addr, port), pass)
}
