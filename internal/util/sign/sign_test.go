package sign

import "testing"

func TestHash(t *testing.T) {
	type args struct {
		s   string
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "normal test",
			args:    args{s: "some metric", key: "some key"},
			want:    "2ac2b6ca2632df66613b05bdc6cfddf6d4e7a2ff389a0078b46a58d5a4b1c8d3",
			wantErr: false,
		},
		{
			name:    "empty strings",
			args:    args{s: "", key: ""},
			want:    "b613679a0814d9ec772f95d778c35fc5ff1697c493715653c6c712144292c5ad",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Hash(tt.args.s, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}
