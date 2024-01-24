package factories

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/viniosilva/where-are-my-fruits/mocks"
)

func TestFactory_Build(t *testing.T) {
	tests := map[string]struct {
		mock    func(db *mocks.MockFactoryDB)
		wantErr string
	}{
		"should be successful": {
			mock: func(db *mocks.MockFactoryDB) {
				db.EXPECT().DB().Return(nil, nil)
			},
		},
		"should throw error on get DB": {
			mock: func(db *mocks.MockFactoryDB) {
				db.EXPECT().DB().Return(nil, fmt.Errorf("error"))
			},
			wantErr: "error",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			dbMock := mocks.NewMockFactoryDB(ctrl)

			tt.mock(dbMock)

			// when
			got, err := Build(dbMock, nil)

			// then
			assert.NotNil(t, got)
			if err != nil || tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}
