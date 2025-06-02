package interaction

import (
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"
)

// InteractionResponder 用於管理互動回應的狀態
type InteractionResponder struct {
	session      *discordgo.Session
	interaction  *discordgo.Interaction
	hasResponded bool
	mu           sync.Mutex
}

// NewInteractionResponder 創建新的回應管理器
func NewInteractionResponder(s *discordgo.Session, i *discordgo.Interaction) *InteractionResponder {
	return &InteractionResponder{
		session:      s,
		interaction:  i,
		hasResponded: false,
	}
}

// Respond 智能判斷使用 Response 或 Edit
func (r *InteractionResponder) Respond(content string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.hasResponded {
		// 第一次回應
		err := r.session.InteractionRespond(r.interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: content,
			},
		})
		if err == nil {
			r.hasResponded = true
		}
		return err
	} else {
		// 後續編輯
		_, err := r.session.InteractionResponseEdit(r.interaction, &discordgo.WebhookEdit{
			Content: &content,
		})
		return err
	}
}

func (r *InteractionResponder) RespondWithError(msg string, err error) error {
	if err != nil {
		return r.Respond(fmt.Sprintf("%s: %v", msg, err))
	}
	return r.Respond(msg)
}
