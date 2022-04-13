package observer

import (
	"sync"
	"testing"
)

func TestObserver_Subscribe(t *testing.T) {
	type fields struct {
		observer      *Observer[string]
		numberClients int
	}
	type args struct {
		message string
	}
	tests := []struct {
		name              string
		args              args
		fields            fields
		expectMessage     string
		expectClients     int
		expectCancelError error
	}{
		{
			name: "send message to all clients 1",
			args: args{
				message: "client1",
			},
			fields: fields{
				observer:      NewObserver[string](),
				numberClients: 3,
			},
			expectMessage: "client1",
			expectClients: 3,
		},
		{
			name: "send message to all clients 3",
			args: args{
				message: "hello world",
			},
			fields: fields{
				observer:      NewObserver[string](),
				numberClients: 3,
			},
			expectMessage: "hello world",
			expectClients: 3,
		},
		{
			name: "send message to all clients 4",
			args: args{
				message: "foo bar",
			},
			fields: fields{
				observer:      NewObserver[string](),
				numberClients: 3,
			},
			expectMessage: "foo bar",
			expectClients: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cancelfuncs := []CancelFunc{}
			for i := 0; i < tt.fields.numberClients; i++ {
				_, cf := tt.fields.observer.Subscribe()
				cancelfuncs = append(cancelfuncs, cf)
			}
			if len(cancelfuncs) != tt.expectClients {
				t.Errorf("Observer.Subscribe() expected %v clients, got %v", tt.expectClients, len(cancelfuncs))
			}
		})
	}
}

func TestObserver_deleteClient(t *testing.T) {
	type fields struct {
		observer      *Observer[string]
		numberClients int
	}
	type args struct {
		message string
	}
	tests := []struct {
		name              string
		args              args
		fields            fields
		expectClients     int
		expectCancelError error
	}{
		{
			name: "send message to all clients 1",
			args: args{
				message: "client1",
			},
			fields: fields{
				observer:      NewObserver[string](),
				numberClients: 3,
			},
			expectClients: 3,
		},
		{
			name: "send message to all clients 3",
			args: args{
				message: "hello world",
			},
			fields: fields{
				observer:      NewObserver[string](),
				numberClients: 3,
			},
			expectClients: 3,
		},
		{
			name: "send message to all clients 4",
			args: args{
				message: "foo bar",
			},
			fields: fields{
				observer:      NewObserver[string](),
				numberClients: 3,
			},
			expectClients: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cancelfuncs := []CancelFunc{}
			for i := 0; i < tt.fields.numberClients; i++ {
				_, cf := tt.fields.observer.Subscribe()
				cancelfuncs = append(cancelfuncs, cf)
			}

			got := tt.fields.observer.Clients()
			if got != int64(tt.expectClients) {
				t.Errorf("TestObserver.Observer[string].Clients() = %v, want %v", got, tt.expectClients)
			}

			for _, cancel := range cancelfuncs {
				err := cancel()
				if err != tt.expectCancelError {
					t.Errorf("TestObserver.Observer[string].Cancel() = %v, want %v", err, tt.expectCancelError)
				}
			}

			gotAfterCancel := tt.fields.observer.Clients()
			if gotAfterCancel != 0 {
				t.Errorf("TestObserver.Observer[string].Clients() after cancel = %v, want 0", got)
			}
		})
	}
}

func TestObserver_NotifyAll(t *testing.T) {
	type fields struct {
		observer      *Observer[string]
		numberClients int
	}
	type args struct {
		message string
	}
	tests := []struct {
		name          string
		args          args
		fields        fields
		expectMessage string
	}{
		{
			name: "send message to all clients 1",
			args: args{
				message: "client1",
			},
			fields: fields{
				observer:      NewObserver[string](),
				numberClients: 3,
			},
			expectMessage: "client1",
		},
		{
			name: "send message to all clients 3",
			args: args{
				message: "hello world",
			},
			fields: fields{
				observer:      NewObserver[string](),
				numberClients: 3,
			},
			expectMessage: "hello world",
		},
		{
			name: "send message to all clients 4",
			args: args{
				message: "foo bar",
			},
			fields: fields{
				observer:      NewObserver[string](),
				numberClients: 3,
			},
			expectMessage: "foo bar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wg := &sync.WaitGroup{}
			wg.Add(tt.fields.numberClients)
			for i := 0; i < tt.fields.numberClients; i++ {
				ch, _ := tt.fields.observer.Subscribe()
				go func(ch <-chan string) {
					for {
						message := <-ch
						if tt.expectMessage == message {
							wg.Done()
							return
						}
					}
				}(ch)
			}
			// notify all clients
			tt.fields.observer.NotifyAll(tt.args.message)
			wg.Wait()
		})
	}
}

func TestObserver_Clients(t *testing.T) {
	type fields struct {
		observer      *Observer[string]
		numberClients int
	}
	tests := []struct {
		name              string
		fields            fields
		expectClients     int
		expectCancelError error
	}{
		{
			name: "send message to all clients 1",
			fields: fields{
				observer:      NewObserver[string](),
				numberClients: 3,
			},
			expectClients: 3,
		},
		{
			name: "send message to all clients 3",
			fields: fields{
				observer:      NewObserver[string](),
				numberClients: 3,
			},
			expectClients: 3,
		},
		{
			name: "send message to all clients 4",
			fields: fields{
				observer:      NewObserver[string](),
				numberClients: 3,
			},
			expectClients: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cancelfuncs := []CancelFunc{}

			for i := 0; i < tt.fields.numberClients; i++ {
				_, cf := tt.fields.observer.Subscribe()
				cancelfuncs = append(cancelfuncs, cf)
			}

			got := tt.fields.observer.Clients()
			if got != int64(tt.expectClients) {
				t.Errorf("TestObserver_Clients got = %v, want %v", got, tt.expectClients)
			}

			for _, cancel := range cancelfuncs {
				err := cancel()
				if err != tt.expectCancelError {
					t.Errorf("TestObserver_Clients.Cancel() = %v, want %v", err, tt.expectCancelError)
				}
			}

			gotAfterCancel := tt.fields.observer.Clients()
			if gotAfterCancel != 0 {
				t.Errorf("TestObserver_Clients() after cancel = %v, want 0", got)
			}
		})
	}
}
