package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_mainHandler_PostGet(t *testing.T) {
	testLink := `practicum.yandex.ru`
	// POST
	bodyReader := strings.NewReader(testLink)
	postRequest := httptest.NewRequest(http.MethodPost, "/", bodyReader)
	postRequest.Header.Add("Content-Type", "text/plain")
	// создаём новый Recorder
	postRecorder := httptest.NewRecorder()
	mainHandler(postRecorder, postRequest)

	resPost := postRecorder.Result()
	// проверяем код ответа
	assert.Equal(t, http.StatusCreated, resPost.StatusCode)
	// получаем и проверяем тело запроса
	defer resPost.Body.Close()
	resBody, err := io.ReadAll(resPost.Body)
	resStr := string(resBody)

	require.NoError(t, err)
	require.True(t, strings.Contains(resStr, "localhost:8080/"))

	// GET
	// получаем хэш сокращённого адреса
	cacheId := resStr[strings.LastIndex(resStr, "/")+1:]
	assert.Equal(t, len(cacheId), 8)

	// отправляем новый запрос
	getRequest := httptest.NewRequest(http.MethodGet, "/"+cacheId, nil)

	getRecorder := httptest.NewRecorder()
	mainHandler(getRecorder, getRequest)

	// получаем ответ
	resGet := getRecorder.Result()
	defer resGet.Body.Close()
	assert.Equal(t, http.StatusTemporaryRedirect, resGet.StatusCode)
	resGetStr := resGet.Header.Get("Location")

	require.NotEmpty(t, resGetStr)
	assert.Equal(t, resGetStr, testLink)
}

func Test_mainHandler_BadMethod(t *testing.T) {
	putRequest := httptest.NewRequest(http.MethodPut, "/", nil)
	// создаём новый Recorder
	recorder := httptest.NewRecorder()
	mainHandler(recorder, putRequest)

	res := recorder.Result()
	defer res.Body.Close()
	// проверяем код ответа
	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
}
