package observer

import (
	"testing"
)

func TestObserver_RegisterClient(t *testing.T) {
	type fields struct {
		observer *Observer[string, string]
	}
	type args struct {
		id      string
		channel chan string
	}
	tests := []struct {
		name   string
		args   args
		fields fields
	}{
		{
			name: "Register client",
			args: args{
				id:      "client1",
				channel: make(chan string),
			},
			fields: fields{
				observer: NewObserver[string, string](),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.observer.RegisterClient(tt.args.id, tt.args.channel)
		})
	}
}

func TestObserver_DeRegisterClient(t *testing.T) {
	type fields struct {
		observer *Observer[string, string]
		clients  map[string]chan string
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		fields  fields
		wantErr bool
	}{
		{
			name: "DeRegister client1",
			args: args{
				id: "client1",
			},
			fields: fields{
				observer: NewObserver[string, string](),
				clients: map[string]chan string{
					"client1": make(chan string),
				},
			},
			wantErr: false,
		},
		{
			name: "DeRegister client1",
			args: args{
				id: "client1",
			},
			fields: fields{
				observer: NewObserver[string, string](),
				clients: map[string]chan string{
					"client1": make(chan string),
					"client2": make(chan string),
				},
			},
			wantErr: false,
		},
		{
			name: "DeRegister client2",
			args: args{
				id: "client2",
			},
			fields: fields{
				observer: NewObserver[string, string](),
				clients: map[string]chan string{
					"client1": make(chan string),
					"client2": make(chan string),
					"client3": make(chan string),
				},
			},
			wantErr: false,
		},
		{
			name: "DeRegister client5",
			args: args{
				id: "client5",
			},
			fields: fields{
				observer: NewObserver[string, string](),
				clients: map[string]chan string{
					"client1": make(chan string),
					"client2": make(chan string),
					"client3": make(chan string),
				},
			},
			wantErr: true,
		},
		{
			name: "DeRegister client19",
			args: args{
				id: "client5",
			},
			fields: fields{
				observer: NewObserver[string, string](),
				clients: map[string]chan string{
					"client1": make(chan string),
					"client2": make(chan string),
					"client3": make(chan string),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for id, client := range tt.fields.clients {
				// add clients to observer
				tt.fields.observer.RegisterClient(id, client)
			}
			// deregister client
			err := tt.fields.observer.DeRegisterClient(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Observer[string,string].DeRegisterClient() error = %v wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestObserver_NotifyAll(t *testing.T) {
	type fields struct {
		observer *Observer[string, string]
		clients  map[string]chan string
	}
	type args struct {
		message string
	}
	tests := []struct {
		name   string
		args   args
		fields fields
	}{
		{
			name: "send message to all clients 1",
			args: args{
				message: "client1",
			},
			fields: fields{
				observer: NewObserver[string, string](),
				clients: map[string]chan string{
					"client1": make(chan string, 1),
					"client2": make(chan string, 1),
					"client3": make(chan string, 1),
				},
			},
		},
		{
			name: "send message to all clients 3",
			args: args{
				message: "hello world",
			},
			fields: fields{
				observer: NewObserver[string, string](),
				clients: map[string]chan string{
					"client1": make(chan string, 1),
					"client2": make(chan string, 1),
					"client3": make(chan string, 1),
				},
			},
		},
		{
			name: "send message to all clients 4",
			args: args{
				message: "foo bar",
			},
			fields: fields{
				observer: NewObserver[string, string](),
				clients: map[string]chan string{
					"client1": make(chan string, 1),
					"client2": make(chan string, 1),
					"client3": make(chan string, 1),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for id, client := range tt.fields.clients {
				// add clients to observer
				tt.fields.observer.RegisterClient(id, client)
			}
			// notify all clients
			tt.fields.observer.NotifyAll(tt.args.message)
		})
	}
}

func TestObserver_NotifyClient(t *testing.T) {
	type fields struct {
		observer *Observer[string, string]
		clients  map[string]chan string
	}
	type args struct {
		clientID string
		message  string
	}
	tests := []struct {
		name    string
		args    args
		fields  fields
		wantErr bool
	}{
		{
			name: "notify client100 (not registered)",
			args: args{
				clientID: "client100",
				message:  "client1",
			},
			fields: fields{
				observer: NewObserver[string, string](),
				clients: map[string]chan string{
					"client1": make(chan string, 1),
					"client2": make(chan string, 1),
					"client3": make(chan string, 1),
				},
			},
			wantErr: true,
		},
		{
			name: "notify client1",
			args: args{
				clientID: "client1",
				message:  "client1",
			},
			fields: fields{
				observer: NewObserver[string, string](),
				clients: map[string]chan string{
					"client1": make(chan string, 1),
					"client2": make(chan string, 1),
					"client3": make(chan string, 1),
				},
			},
			wantErr: false,
		},
		{
			name: "notify client2",
			args: args{
				clientID: "client2",
				message:  "hello world",
			},
			fields: fields{
				observer: NewObserver[string, string](),
				clients: map[string]chan string{
					"client1": make(chan string, 1),
					"client2": make(chan string, 1),
					"client3": make(chan string, 1),
				},
			},
			wantErr: false,
		},
		{
			name: "notify client3",
			args: args{
				clientID: "client3",
				message:  "foo bar",
			},
			fields: fields{
				observer: NewObserver[string, string](),
				clients: map[string]chan string{
					"client1": make(chan string, 1),
					"client2": make(chan string, 1),
					"client3": make(chan string, 1),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for id, client := range tt.fields.clients {
				// add clients to observer
				tt.fields.observer.RegisterClient(id, client)
			}
			// notify client
			err := tt.fields.observer.NotifyClient(tt.args.clientID, tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("Observer[string,string].NotifyClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestObserver_Clients(t *testing.T) {
	type fields struct {
		observer *Observer[string, string]
		clients  map[string]chan string
	}
	type args struct {
		message string
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{
			name: "1 client",
			fields: fields{
				observer: NewObserver[string, string](),
				clients: map[string]chan string{
					"client1": make(chan string),
				},
			},
			want: 1,
		},
		{
			name: "2 clients",
			fields: fields{
				observer: NewObserver[string, string](),
				clients: map[string]chan string{
					"client1": make(chan string),
					"client2": make(chan string),
				},
			},
			want: 2,
		},
		{
			name: "3 clients",
			fields: fields{
				observer: NewObserver[string, string](),
				clients: map[string]chan string{
					"client1": make(chan string),
					"client2": make(chan string),
					"client3": make(chan string),
				},
			},
			want: 3,
		},
		{
			name: "10 client",
			fields: fields{
				observer: NewObserver[string, string](),
				clients: map[string]chan string{
					"client1":  make(chan string),
					"client2":  make(chan string),
					"client3":  make(chan string),
					"client4":  make(chan string),
					"client5":  make(chan string),
					"client6":  make(chan string),
					"client7":  make(chan string),
					"client8":  make(chan string),
					"client9":  make(chan string),
					"client10": make(chan string),
				},
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for id, client := range tt.fields.clients {
				// add clients to observer
				tt.fields.observer.RegisterClient(id, client)
			}

			got := tt.fields.observer.Clients()
			if got != tt.want {
				t.Errorf("Observer[string,string].Clients() = %v, want %v", got, tt.want)
			}
		})
	}
}
