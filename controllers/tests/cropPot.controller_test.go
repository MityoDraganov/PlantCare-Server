package controllers_test

import (
    "bytes"
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "PlantCare/controllers"
    "PlantCare/dtos"
    "PlantCare/initPackage"
    "PlantCare/models"

    "github.com/clerk/clerk-sdk-go/v2"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/gorilla/mux"
    "gorm.io/gorm"
)

type MockDB struct {
    mock.Mock
}

func (db *MockDB) Save(value interface{}) *gorm.DB {
    args := db.Called(value)
    return args.Get(0).(*gorm.DB)
}

func (db *MockDB) First(out interface{}, where ...interface{}) *gorm.DB {
    args := db.Called(append([]interface{}{out}, where...)...)
    return args.Get(0).(*gorm.DB)
}

func (db *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
    dbArgs := db.Called(append([]interface{}{query}, args...)...)
    return dbArgs.Get(0).(*gorm.DB)
}

func (db *MockDB) Find(out interface{}, where ...interface{}) *gorm.DB {
    args := db.Called(append([]interface{}{out}, where...)...)
    return args.Get(0).(*gorm.DB)
}

func TestAssignCropPotToUser(t *testing.T) {
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/assign-pot/valid-token", nil)

    claims := &clerk.SessionClaims{
        clerk.RegisteredClaims{Subject: "user"},
    }
	
    // Use a proper key for context (same as in actual handler)
    ctx := context.WithValue(req.Context(), clerk.Key, claims)
    req = req.WithContext(ctx)

    mockDB := new(MockDB)
    initPackage.Db = mockDB // Ensure to use the mock DB in the test

    cropPot := models.CropPot{
        Token: "valid-token",
    }

    mockDB.On("Where", "token = ?", "valid-token").Return(mockDB).Once()
    mockDB.On("First", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
        arg := args.Get(0).(*models.CropPot)
        *arg = cropPot
    }).Return(mockDB).Once()
    mockDB.On("Save", mock.Anything).Return(mockDB).Once()

    controllers.AssignCropPotToUser(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
}


func TestUpdateCropPot(t *testing.T) {
    w := httptest.NewRecorder()
    updateData := dtos.ControlDto{
        WateringInterval: 7,
    }
    jsonValue, _ := json.Marshal(updateData)
    req, _ := http.NewRequest("PUT", "/crop-pot/update/{id}", bytes.NewBuffer(jsonValue))

    params := map[string]string{
        "id": "pot-1",
    }
    req = mux.SetURLVars(req, params)

    mockDB := new(MockDB)
    db := &gorm.DB{}
    initPackage.Db = db

    controlSettings := models.Control{
        WateringInterval: 5,
    }
    cropPot := models.CropPot{
        Token: "pot-1",
    }

    mockDB.On("First", &cropPot, "id = ?", "pot-1").Return(mockDB).Once()
    mockDB.On("First", &controlSettings, "id = ?", "cs-1").Return(mockDB).Once()
    mockDB.On("Save", mock.Anything).Return(mockDB).Once()

    controllers.UpdateCropPot(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
}
