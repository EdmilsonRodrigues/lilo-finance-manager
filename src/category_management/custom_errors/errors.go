package customerrors

import "github.com/EdmilsonRodrigues/lilo-finance-manager/src/common_utils/go/serialization/http_serialization"

var BadRequestResponse = httpserialization.ErrorResponse{
	Details: httpserialization.ErrorDetails{
		Status:  400,
		Message: "Bad Request",
	},
}

var ForbiddenResponse = httpserialization.ErrorResponse{
	Details: httpserialization.ErrorDetails{
		Status:  403,
		Message: "Forbidden: User does not have access to this resource",
	},
}

var UnprocessableEntityResponse = httpserialization.ErrorResponse{
	Details: httpserialization.ErrorDetails{
		Status:  422,
		Message: "Unprocessable Entity, check the request body and path parameters",
	},
}

var InternalServerErrorResponse = httpserialization.ErrorResponse{
	Details: httpserialization.ErrorDetails{
		Status:  500,
		Message: "Internal Server Error",
	},
}

func InternalServerError(message string) httpserialization.ErrorResponse {
	err := InternalServerErrorResponse
	err.Details.Message += ": " + message
	return err
}
