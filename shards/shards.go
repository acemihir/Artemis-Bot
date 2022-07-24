// Edited version of "shards" by "servusdei2018":
// Repository: https://github.com/servusdei2018/shards
// License: https://github.com/servusdei2018/shards/blob/master/LICENSE.md
package shards

import (
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	TIMELIMIT = time.Second * 5
	VERSION   = "1.2.2"
)

type Shard struct {
	sync.RWMutex
	Session    *discordgo.Session
	ID         int
	ShardCount int
	handlers   []interface{}
}

// AddHandler registers an event handler for a Shard.
// Shouldn't be called after Init or results in undefined behavior.
func (s *Shard) AddHandler(handler interface{}) {
	s.Lock()
	defer s.Unlock()

	s.handlers = append(s.handlers, handler)
}

// ApplicationCommandCreate registers an application command for a Shard.
// Shouldn't be called before Initialization.
func (s *Shard) ApplicationCommandCreate(guildID string, cmd *discordgo.ApplicationCommand) error {
	s.Lock()
	defer s.Unlock()

	// Referencing s.Session before Initialization will result in a nil pointer dereference panic.
	if s.Session == nil {
		return fmt.Errorf("error: shard.ApplicationCommandCreate must not be called before shard.Init")
	}

	_, err := s.Session.ApplicationCommandCreate(s.Session.State.User.ID, guildID, cmd)
	return err
}

func (s *Shard) ApplicationCommandDelete(guildID string, cmd *discordgo.ApplicationCommand) error {
	s.Lock()
	defer s.Unlock()

	// Referencing s.Session before Initialization will result in a nil pointer dereference panic.
	if s.Session == nil {
		return fmt.Errorf("error: shard.ApplicationCommandCreate must not be called before shard.Init")
	}

	err := s.Session.ApplicationCommandDelete(s.Session.State.User.ID, guildID, cmd.ID)
	return err
}

// GuildCount returns the amount of guilds that a Shard is handling.
func (s *Shard) GuildCount() (count int) {
	s.RLock()
	defer s.RUnlock()

	if s.Session != nil {
		s.Session.State.RLock()
		count += len(s.Session.State.Guilds)
		s.Session.State.RUnlock()
	}

	return
}

// Init initializes a shard with a bot token, its Shard ID, the total amount of shards, and a Discord intent.
func (s *Shard) Init(token string, ID, ShardCount int, intent discordgo.Intent) (err error) {
	s.Lock()
	defer s.Unlock()

	s.ID = ID
	s.ShardCount = ShardCount

	// Create the session.
	s.Session, err = discordgo.New(token)
	if err != nil {
		return
	}

	s.Session.ShardCount = s.ShardCount
	s.Session.ShardID = s.ID
	s.Session.Identify.Intents = intent

	// Add handlers to the session.
	for _, handler := range s.handlers {
		s.Session.AddHandler(handler)
	}

	err = s.Session.Open()

	return
}

// Stop stops a shard.
func (s *Shard) Stop() (err error) {
	s.Lock()
	defer s.Unlock()

	err = s.Session.Close()

	return
}
