package api

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/viniosilva/where-are-my-fruits/mocks"
)

func TestServer_ConfigServer(t *testing.T) {
	type args struct {
		host string
		port string
	}
	tests := map[string]struct {
		args args
	}{
		"should be successful": {
			args: args{
				host: "0.0.0.0",
				port: "3000",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			healthControllerMock := mocks.NewMockHealthController(ctrl)

			// when
			got := ConfigServer(tt.args.host, tt.args.port, nil, healthControllerMock)

			// then
			assert.NotNil(t, got)
		})
	}
}
