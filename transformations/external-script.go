package transformations

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gabriel-vasile/mimetype"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgtype"

	"external-api-service/globals"
	"external-api-service/types"
)

type ExternalScript struct{}

const ExternalScriptTransformationKey = "external-script"
const ExternalScriptUUIDKey = "script"

var ErrScriptNameNotSet = errors.New("the script name has not been set")

func (t ExternalScript) ApplyAfter(r *http.Response, options map[string]interface{}) error {
	transformationID := make([]byte, 32)
	_, err := rand.Read(transformationID)
	if err != nil {
		return err
	}

	responseFile, err := os.CreateTemp("", fmt.Sprintf("%s-*.response", transformationID))
	if err != nil {
		return err
	}
	_, err = io.Copy(responseFile, r.Body)
	if err != nil {
		return err
	}
	err = responseFile.Sync()
	if err != nil {
		return err
	}

	err = responseFile.Close()
	if err != nil {
		return err
	}

	scriptIDRaw, scriptNameSet := options[ExternalScriptUUIDKey]
	if !scriptNameSet {
		return ErrScriptNameNotSet
	}

	scriptID := pgtype.UUID{}
	err = scriptID.Scan(scriptIDRaw)
	if err != nil {
		return err
	}

	query, err := globals.SqlQueries.Raw("get-transformation-script")
	if err != nil {
		return err
	}

	scripts, err := globals.Db.Query(context.Background(), query, scriptID)
	var script types.ExternalTransformationScript
	err = pgxscan.ScanOne(&script, scripts)
	if err != nil {
		return err
	}

	scriptFile, err := os.CreateTemp("", fmt.Sprintf("%s-*.py", transformationID))
	if err != nil {
		return err
	}
	_, err = scriptFile.Write([]byte(script.ScriptContents.String))
	if err != nil {
		return err
	}
	err = scriptFile.Sync()
	if err != nil {
		return err
	}
	err = scriptFile.Close()
	if err != nil {
		return err
	}

	var errorOutput strings.Builder

	for _, requirement := range script.RequiredPackages {
		requirementInstallation := exec.Command("python", "-m", "pip", "install", requirement)
		requirementInstallation.Stdout = os.Stdout
		requirementInstallation.Stderr = &errorOutput
		err = requirementInstallation.Run()
		if err != nil {
			return fmt.Errorf("failed to install requirement %s: %s", requirement, errorOutput.String())
		}
		errorOutput.Reset()
	}
	var responseOutput bytes.Buffer
	scriptExecution := exec.Command("python", scriptFile.Name(), responseFile.Name())
	scriptExecution.Stdout = &responseOutput
	scriptExecution.Stderr = &errorOutput
	err = scriptExecution.Run()
	if err != nil {
		return err
	}
	contentType, err := mimetype.DetectReader(&responseOutput)
	if err != nil {
		return err
	}

	r.Header.Set("Content-Type", contentType.String())
	r.Body = io.NopCloser(&responseOutput)
	return nil
}
