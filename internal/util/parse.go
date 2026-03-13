package util

import (
	"encoding/json"
	"golang-grpc/internal/log"
	"net/http"
)

func ParseBody(request *http.Request, result any) error {
	decoder := json.NewDecoder(request.Body)
	decodeError := decoder.Decode(result)
	if decodeError != nil {
		return decodeError
	}

	return nil
}

func WriteError(writer http.ResponseWriter, statusCode int, err error) {
	writer.WriteHeader(statusCode)
	writer.Header().Set("Content-Type", "application/json")
	_, writingError := writer.Write(json.RawMessage(`{"message":"` + err.Error() + `"}`))
	if writingError != nil {
		log.PrintError("Failed to marshal the error message json", writingError)
	}
}

type Metadata struct {
	Total  uint64  `json:"total"`
	Limit  *uint64 `json:"limit,omitempty"`
	Offset *uint64 `json:"offset,omitempty"`
}

type ResponseWrapper struct {
	Data     any       `json:"data"`
	Metadata *Metadata `json:"metadata,omitempty"`
}

func WriteResponse(writer http.ResponseWriter, statusCode int, result interface{}) error {
	writer.WriteHeader(statusCode)
	writer.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(result)

	if err != nil {
		log.Errorln("Failed to marshal the response json")
		return err
	}
	_, writingError := writer.Write(response)
	if writingError != nil {
		log.PrintError("Failed to marshal the response json", writingError)
	}

	return nil
}
