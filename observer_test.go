package observer

import (
	"context"
	"sync"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestObserver_Subscribe(t *testing.T) {
	type fields struct {
		observer      Observable[string]
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
				observer:      new(Observer[string]),
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
				observer:      new(Observer[string]),
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
				observer:      new(Observer[string]),
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
				observer:      new(Observer[string]),
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
				observer:      new(Observer[string]),
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
				observer:      new(Observer[string]),
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
				observer:      new(Observer[string]),
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
				observer:      new(Observer[string]),
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
				observer:      new(Observer[string]),
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
				observer:      new(Observer[string]),
				numberClients: 3,
			},
			expectClients: 3,
		},
		{
			name: "send message to all clients 3",
			fields: fields{
				observer:      new(Observer[string]),
				numberClients: 3,
			},
			expectClients: 3,
		},
		{
			name: "send message to all clients 4",
			fields: fields{
				observer:      new(Observer[string]),
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

func TestObserver_Handle(t *testing.T) {
	type fields struct {
		observer      *Observer[string]
		numberClients int
	}
	type args struct {
		message string
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		expectClients     int
		expectCancelError error
	}{
		{
			name: "send message to all clients 1",
			args: args{
				message: "client1",
			},
			fields: fields{
				observer:      new(Observer[string]),
				numberClients: 3,
			},
			expectClients: 3,
		},
		{
			name: "send message to all clients 2",
			args: args{
				message: "client2",
			},
			fields: fields{
				observer:      new(Observer[string]),
				numberClients: 3,
			},
			expectClients: 3,
		},
		{
			name: "send message to all clients 3",
			fields: fields{
				observer:      new(Observer[string]),
				numberClients: 3,
			},
			expectClients: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()
			cancelfuncs := []CancelFunc{}
			var wg sync.WaitGroup
			g, _ := errgroup.WithContext(ctx)
			for i := 0; i < tt.fields.numberClients; i++ {
				wg.Add(1)
				f, cf := Handle[string](ctx, tt.fields.observer, func(_ context.Context, s string) error {
					defer wg.Done()
					if s != tt.args.message {
						t.Errorf("expected %s got %s", tt.args.message, s)
					}
					return nil
				})
				g.Go(f)
				cancelfuncs = append(cancelfuncs, cf)
			}

			got := tt.fields.observer.Clients()
			if got != int64(tt.expectClients) {
				t.Errorf("TestObserver_Clients got = %v, want %v", got, tt.expectClients)
			}

			tt.fields.observer.NotifyAll(tt.args.message)
			wg.Wait()

			for _, cancel := range cancelfuncs {
				err := cancel()
				if err != tt.expectCancelError {
					t.Errorf("TestObserver_Clients.Cancel() = %v, want %v", err, tt.expectCancelError)
				}
			}
			if err := g.Wait(); err != nil {
				t.Fatal(err)
			}
			gotAfterCancel := tt.fields.observer.Clients()
			if gotAfterCancel != 0 {
				t.Errorf("TestObserver_Clients() after cancel = %v, want 0", got)
			}
		})
	}
}
