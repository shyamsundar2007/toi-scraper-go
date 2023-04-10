package notify

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"reflect"
	"testing"
)

type mockTgClient struct{}

var message string

func (t mockTgClient) Send(msg tgbotapi.Chattable) (tgbotapi.Message, error) {
	message = msg.(tgbotapi.MessageConfig).Text
	return tgbotapi.Message{}, nil
}

func TestTelegramNotifier_Notify(t *testing.T) {
	type fields struct {
		tgClient TgAPI
		chatId   int64
	}
	type args struct {
		movies []Movie
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "testNotifyMoviesNil",
			fields: fields{
				tgClient: mockTgClient{},
				chatId:   0,
			},
			args:    args{movies: nil},
			want:    "",
			wantErr: true,
		},
		{
			name: "testNotifyValidMovie",
			fields: fields{
				tgClient: mockTgClient{},
				chatId:   0,
			},
			args: args{movies: []Movie{
				{
					MovieName:    "testMovie",
					UserRating:   Rating("4.0"),
					CriticRating: Rating("3.5"),
					Language:     "Tamil",
				},
			}},
			want:    "testMovie (Tamil)\nCritic: 3.5\nUser: 4.0\n--------------------\n",
			wantErr: false,
		},
		{
			name: "testNotifyValidMovies",
			fields: fields{
				tgClient: mockTgClient{},
				chatId:   0,
			},
			args: args{movies: []Movie{
				{
					MovieName:    "testMovie",
					UserRating:   Rating("4.0"),
					CriticRating: Rating("3.5"),
					Language:     "Tamil",
				},
				{
					MovieName:    "testMovie2",
					UserRating:   Rating("2.0"),
					CriticRating: Rating("1.5"),
					Language:     "Hindi",
				},
			}},
			want: "testMovie (Tamil)\nCritic: 3.5\nUser: 4.0\n--------------------\n" +
				"testMovie2 (Hindi)\nCritic: 1.5\nUser: 2.0\n--------------------\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tgNotifier := &TelegramNotifier{
				tgClient: tt.fields.tgClient,
				chatId:   tt.fields.chatId,
			}
			err := tgNotifier.Notify(tt.args.movies)
			if (err != nil) != tt.wantErr {
				t.Errorf("Notify() error = %v, wantErr %v", err, tt.wantErr)
			} else if !reflect.DeepEqual(message, tt.want) {
				t.Errorf("Notify() = %v, want %v", message, tt.want)
			}
		})
	}
}
