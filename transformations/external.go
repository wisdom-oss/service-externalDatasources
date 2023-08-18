package transformations

import (
	"encoding/json"
	"errors"
	"external-api-service/globals"
	databaseTypes "external-api-service/types/database"
	"fmt"
	"github.com/blockloop/scan/v2"
	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	wisdomUtils "github.com/wisdom-oss/microservice-utils"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var supportedExecutors = []string{
	"python",
	"rscript",
}

type ExternalTransformation struct {
	Enabled       bool        `json:"enabled"`
	Script        uuid.UUID   `json:"script"`
	Configuration interface{} `json:"configuration"`
}

func (et ExternalTransformation) ApplyBefore(r *http.Request, data []byte) error {
	return errors.New("transformation is post-proxy only")
}

func (et ExternalTransformation) ApplyAfter(r *http.Response, data []byte) error {
	// try to parse the configuration for the transformation
	if err := json.Unmarshal(data, &et); err != nil {
		return err
	}
	if !et.Enabled {
		return nil
	}
	// now create a new uuid to identify the request/responses
	responseId := uuid.NewString()

	// now pull the script from the script database
	scriptRows, err := globals.SqlQueries.Query(globals.Db, "get-transformation-script", et.Script.String())
	if err != nil {
		return fmt.Errorf("unable to pull transformation script: %w", err)
	}
	var script databaseTypes.TransformationScript
	err = scan.Row(&script, scriptRows)
	if err != nil {
		return fmt.Errorf("unable to parse script definition: %w", err)
	}

	// now check if the executor is supported by the container
	if !wisdomUtils.ArrayContains(supportedExecutors, script.Engine) {
		return errors.New("unsupported engine configured for script. supported engines are: python, rscript")
	}

	// now store the script contents into a temporary file
	scriptFile, err := os.CreateTemp("", fmt.Sprintf("%s-*.script", responseId))
	if err != nil {
		return fmt.Errorf("unable to create temporary file for script: %w", err)
	}
	_, err = scriptFile.WriteString(script.Contents)
	if err != nil {
		return fmt.Errorf("unable to store script in temporary file: %w", err)
	}

	// now store the response in a temporary file
	responseFile, err := os.CreateTemp("", fmt.Sprintf("%s-*.initialResponse", responseId))
	_, err = io.Copy(responseFile, r.Body)
	if err != nil {
		return fmt.Errorf("unable to store response in temporary file: %w", err)
	}

	// now check for possibly required packages and install them
	var errorOutput strings.Builder
	if script.RequiredPackages != nil && len(script.RequiredPackages) > 0 {
		switch script.Engine {
		case "python":
			for _, requiredPackage := range script.RequiredPackages {
				requirementsInstall := exec.Command(
					"python", "-m", "pip", "install", requiredPackage,
				)
				requirementsInstall.Stdout = os.Stdout
				requirementsInstall.Stderr = &errorOutput
				err = requirementsInstall.Run()
				if err != nil {
					return fmt.Errorf("unable to install configured requirements: %s", errorOutput.String())
				}
				errorOutput.Reset()
			}

		}
	}

	// now execute the temporary script with the configured executor
	externalTransformation := exec.Command(script.Engine, scriptFile.Name(), responseFile.Name())
	externalTransformation.Stderr = &errorOutput
	resultFilePathBytes, err := externalTransformation.Output()
	if err != nil {
		return fmt.Errorf("error while executing external script: %s", errorOutput.String())
	}

	// now clean up the path
	resultFilePath := string(resultFilePathBytes)
	resultFilePath = strings.TrimSpace(resultFilePath)

	// now try to guess the content type of the file
	contentType, err := mimetype.DetectFile(resultFilePath)
	if err != nil {
		return fmt.Errorf("unable to guess the content type of the changed response: %w", err)
	}

	r.Header.Set("Content-Type", contentType.String())

	// now open the result file
	resultFile, err := os.Open(resultFilePath)
	if err != nil {
		return fmt.Errorf("error while opening result of external script: %w", err)
	}
	// and set the contents as the new readcloser to the response
	resultReadCloser := io.NopCloser(resultFile)
	r.Body = resultReadCloser
	return nil
}
