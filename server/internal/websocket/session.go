package websocket

import (
	"n.eko.moe/neko/internal/event"
	"n.eko.moe/neko/internal/message"
	"n.eko.moe/neko/internal/session"
)

func (h *MessageHandler) SessionCreated(id string, session *session.Session) error {
	if err := session.Send(message.Identity{
		Message: message.Message{Event: event.IDENTITY_PROVIDE},
		ID:      id,
	}); err != nil {
		return err
	}

	return nil
}

func (h *MessageHandler) SessionConnected(id string, session *session.Session) error {
	// send list of members to session
	if err := session.Send(message.Members{
		Message:  message.Message{Event: event.MEMBER_LIST},
		Memebers: h.sessions.GetConnected(),
	}); err != nil {
		h.logger.Warn().Str("id", id).Err(err).Msgf("sending event %s has failed", event.MEMBER_LIST)
		return err
	}

	// tell session there is a host
	host, ok := h.sessions.GetHost()
	if ok {
		if err := session.Send(message.Control{
			Message: message.Message{Event: event.CONTROL_LOCKED},
			ID:      host.ID,
		}); err != nil {
			h.logger.Warn().Str("id", id).Err(err).Msgf("sending event %s has failed", event.CONTROL_LOCKED)
			return err
		}
	}

	// let everyone know there is a new session
	if err := h.sessions.Brodcast(
		message.Member{
			Message: message.Message{Event: event.MEMBER_CONNECTED},
			Session: session,
		}, nil); err != nil {
		h.logger.Warn().Err(err).Msgf("brodcasting event %s has failed", event.CONTROL_RELEASE)
		return err
	}

	return nil
}

func (h *MessageHandler) SessionDestroyed(id string) error {
	// clear host if exists
	if h.sessions.IsHost(id) {
		h.sessions.ClearHost()
		if err := h.sessions.Brodcast(message.Control{
			Message: message.Message{Event: event.CONTROL_RELEASE},
			ID:      id,
		}, nil); err != nil {
			h.logger.Warn().Err(err).Msgf("brodcasting event %s has failed", event.CONTROL_RELEASE)
		}
	}

	// let everyone know session disconnected
	if err := h.sessions.Brodcast(
		message.MemberDisconnected{
			Message: message.Message{Event: event.MEMBER_DISCONNECTED},
			ID:      id,
		}, nil); err != nil {
		h.logger.Warn().Err(err).Msgf("brodcasting event %s has failed", event.MEMBER_DISCONNECTED)
		return err
	}

	return nil
}