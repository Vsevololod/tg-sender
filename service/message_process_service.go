package service

import (
	"log/slog"
	"tg-sender/domain"

	messages "github.com/Vsevololod/tg-api-contracts-lib/gen/go/messages"
)

type TgClientI interface {
	SendMessageAll(userId uint64, text string) bool
	SendMessage(userID int, title, filePath string) bool
	SendPhoto(userID uint64, title, filePath, photoURL string) bool
}

// MessageProcessService — сервис обработки сообщений
type MessageProcessService struct {
	inputMessageChannel chan domain.MessageWithContext
	tgClient            TgClientI
	log                 *slog.Logger
}

// NewMessageProcessService создает новый сервис и принимает канал для сообщений
func NewMessageProcessService(inputMessageChannel chan domain.MessageWithContext,
	tgClient TgClientI, log *slog.Logger) *MessageProcessService {
	return &MessageProcessService{
		inputMessageChannel: inputMessageChannel,
		tgClient:            tgClient,
		log:                 log,
	}
}

// StartProcessing запускает обработку сообщений в отдельной горутине
func (s *MessageProcessService) StartProcessing(workerCount int) {
	for i := 0; i < workerCount; i++ {
		go func(workerID int) {
			for msg := range s.inputMessageChannel {
				s.ProcessMessage(workerID, msg)
			}
		}(i)
	}
}

// ProcessMessage выполняет обработку сообщения
func (s *MessageProcessService) ProcessMessage(workerID int, msg domain.MessageWithContext) {
	s.log.Info("Processing Message with id", slog.String("uuid", msg.UUID), slog.Int64("worker", int64(workerID)))

	switch msg.Message.Type {
	case messages.MessageType_TEXT:
		s.tgClient.SendMessageAll(msg.Message.UserId, msg.Message.Text)
	case messages.MessageType_IMAGE:
		s.tgClient.SendPhoto(
			msg.Message.UserId,
			msg.Message.Text,
			msg.Message.Params[messages.MessageParams_FILE_URL.String()],
			msg.Message.Params[messages.MessageParams_PHOTO_URL.String()])

	}

}

func (s *MessageProcessService) StopProcessing() {

}
