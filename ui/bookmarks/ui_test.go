package bookmarks

import (
	"github.com/blackmann/go-gurl/lib"
	"github.com/blackmann/go-gurl/mock_lib"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func getTestBookmarkModal(controller *gomock.Controller) (Model, *mock_lib.MockPersistence) {
	persistence := mock_lib.NewMockPersistence(controller)
	return Model{persistence: persistence}, persistence
}

func TestModel_Init(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	bookmarkModal, persistence := getTestBookmarkModal(controller)
	persistence.EXPECT().GetBookmarks().Return([]lib.Bookmark{})

	// Intentionally sending nil so that it initializes
	bookmarkModal, _ = bookmarkModal.Update(nil)

	assert.Contains(t, bookmarkModal.View(), "No items")
}

func TestModel_Update_Filter(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	bookmarkModal, persistence := getTestBookmarkModal(controller)
	persistence.EXPECT().GetBookmarks().Return([]lib.Bookmark{
		{ID: 1, Name: "example", Url: "https://example"},
		{ID: 2, Name: "sample", Url: "https://sample"},
	})

	bookmarkModal, _ = bookmarkModal.Update(Filter("exa"))

	assert.Contains(t, bookmarkModal.View(), "example")
	assert.NotContains(t, bookmarkModal.View(), "sample")
}

func TestModel_Update_UpdateBookmarks(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	bookmarkModal, persistence := getTestBookmarkModal(controller)

	persistence.EXPECT().GetBookmarks().
		Return([]lib.Bookmark{})

	bookmarkModal, _ = bookmarkModal.Update(nil)
	// already nothing
	assert.Len(t, bookmarkModal.bookmarks, 0)

	persistence.EXPECT().GetBookmarks().
		Return([]lib.Bookmark{
			{ID: 1, Name: "example", Url: "https://example"},
			{ID: 2, Name: "sample", Url: "https://sample"},
		})
	bookmarkModal, _ = bookmarkModal.Update(lib.UpdateBookmarks)

	// couldn't find ways to better test this than access
	// the private field
	assert.Len(t, bookmarkModal.bookmarks, 2)
}

func TestModel_Update_KeyDown(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	bookmarkModal, persistence := getTestBookmarkModal(controller)

	persistence.EXPECT().GetBookmarks().
		Return([]lib.Bookmark{
			{ID: 1, Name: "example", Url: "https://example"},
			{ID: 2, Name: "sample", Url: "https://sample"},
		})

	// we need to send a filter so the list.Model is populated
	bookmarkModal, _ = bookmarkModal.Update(Filter(""))

	bookmarkModal, _ = bookmarkModal.Update(tea.KeyMsg{Type: tea.KeyDown})

	// couldn't find ways to better test this than access
	// the private field
	bookmark, _ := bookmarkModal.GetSelected()
	assert.Equal(t, bookmark.Name, "sample")
}
