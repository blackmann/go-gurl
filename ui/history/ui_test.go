package history

import (
	"github.com/blackmann/gurl/lib"
	"github.com/blackmann/gurl/mock_lib"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func getTestHistoryModal(controller *gomock.Controller) (Model, *mock_lib.MockPersistence) {
	persistence := mock_lib.NewMockPersistence(controller)
	return Model{persistence: persistence}, persistence
}

func TestModel_Init(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	instance, persistence := getTestHistoryModal(controller)
	persistence.EXPECT().GetHistory().Return([]lib.History{})

	instance, _ = instance.Update(nil)

	assert.Contains(t, instance.View(), "No items")
}

func TestModel_Update_Filter_ID(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	instance, persistence := getTestHistoryModal(controller)
	persistence.EXPECT().GetHistory().Return([]lib.History{
		{ID: 1, Date: time.Now(), Url: "endpoint-1", Method: "PUT", Annotation: "def"},
		{ID: 2, Date: time.Now(), Url: "endpoint-2", Method: "GET"},
	})

	// initialize
	instance, _ = instance.Update(nil)

	instance, _ = instance.Update(Filter("2"))

	assert.Equal(t, 1, strings.Count(instance.View(), "$"))
	assert.Contains(t, instance.View(), "endpoint-2")
	assert.NotContains(t, instance.View(), "endpoint-1")
}

func TestModel_Update_Filter_Annotation(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	instance, persistence := getTestHistoryModal(controller)
	persistence.EXPECT().GetHistory().Return([]lib.History{
		{ID: 1, Date: time.Now(), Url: "endpoint-1", Method: "PUT", Annotation: "def"},
		{ID: 2, Date: time.Now(), Url: "endpoint-2", Method: "GET"},
	})

	// initialize
	instance, _ = instance.Update(nil)

	instance, _ = instance.Update(Filter("def"))

	assert.Contains(t, instance.View(), "endpoint-1")
	assert.NotContains(t, instance.View(), "endpoint-2")
}

func TestModel_Update_UpdateHistory(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	instance, persistence := getTestHistoryModal(controller)
	persistence.EXPECT().GetHistory().Return([]lib.History{})

	// initialize
	instance, _ = instance.Update(nil)

	persistence.EXPECT().GetHistory().Return([]lib.History{
		{ID: 2, Date: time.Now(), Url: "endpoint-2", Method: "GET"},
	})

	instance, _ = instance.Update(lib.UpdateHistory)

	// the list model is only updated when there's a filter submitted
	instance, _ = instance.Update(Filter(""))

	assert.Equal(t, 1, strings.Count(instance.View(), "$"))
}
