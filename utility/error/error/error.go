package utility

import (
	"encoding/json"
	"log"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HandleServiceErr(err error) error {
	if err == nil {
		return nil
	}

	if vErrs, ok := err.(validation.Errors); ok {
		nvErrs := make(map[string]string)
		for k, v := range vErrs {
			// key := cases.Title(language.English)
			nvErrs[k] = v.Error()
		}
		log.Println(nvErrs)
		vErr, _ := json.Marshal(nvErrs)
		return status.Error(codes.InvalidArgument, string(vErr))
	}

	if strings.Contains(err.Error(), "pq:") || strings.Contains(err.Error(), "json:") {
		return status.Error(codes.Internal, MsgInternalServerError)
	}

	// Note: Add more error type conditions when needed
	return status.Convert(err).Err()
}
