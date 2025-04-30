package customerrors

import "github.com/EdmilsonRodrigues/lilo-finance-manager/src/common_utils/go/serialization"

var BadRequestResponse = serialization.ErrorResponse{Details: serialization.ErrorDetails{Status: 400, Message: "Bad Request"}}
var ForbiddenResponse = serialization.ErrorResponse{Details: serialization.ErrorDetails{Status: 403, Message: "Forbidden: User does not have access to this resource"}}
var UnprocessableEntityResponse = serialization.ErrorResponse{Details: serialization.ErrorDetails{Status: 422, Message: "Unprocessable Entity, check the request body and path parameters"}}
var InternalServerErrorResponse = serialization.ErrorResponse{Details: serialization.ErrorDetails{Status: 500, Message: "Internal Server Error"}}

func InternalServerError(message string) serialization.ErrorResponse {
	err := InternalServerErrorResponse
	err.Details.Message += ": " + message
	return err
}
