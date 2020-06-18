package chatpreference

const defaultTimeZone = "Europe/London"

type Service struct {
	store Storer
}

func NewService(store Storer) *Service {
	return &Service{store: store}
}

func (s *Service) CreateDefaultChatPreferences(chats []int) {
	for _, chatID := range chats {
		_, err := s.store.GetChatPreference(chatID)
		if err != nil && err == ErrNotFound {
			_ = s.store.UpsertChatPreference(&ChatPreference{ChatID: chatID, TimeZone: defaultTimeZone})
		}
	}
}
