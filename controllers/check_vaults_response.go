package controllers

import "github.com/gofiber/fiber/v2"

func checkVaultsResponse(agent *fiber.Agent) (statusCode int, body, errString string) {
	statusCode, body, errs := agent.String()

	if len(errs) > 0 {
		for _, err := range errs {
			errString += err.Error() + ";;"
		}

		if body != "" {
			errString += body + ";;"
		}
	}

	if statusCode != 200 && statusCode != 201 && statusCode != 204 {
		if body != "" {
			errString += body + ";;"
		}
	}

	return statusCode, body, errString
}
