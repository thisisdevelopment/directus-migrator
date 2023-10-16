package directus

import (
	"bytes"
	"context"
	"directus-migrator/pkg/shared"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

func (d *directus) GetSchema(ctx context.Context) (string, error) {

	var res DirectusSchema
	code, err := d.cli.Do(ctx, "GET", d.cfg.DirectusConfig.SnapshotURL, nil, &res)
	if err != nil {
		return "", err
	}
	if code != http.StatusOK {
		return "", fmt.Errorf("unexpected responsecode %d, expected 200", code)
	}

	return res.Convert(), nil
}

func (d *directus) Diff(ctx context.Context, schema string) (string, error) {

	// Create a buffer to store the request body
	var buf bytes.Buffer

	// Create a new multipart writer with the buffer
	w := multipart.NewWriter(&buf)

	fw, err := w.CreateFormFile("file", "schema.json")
	if err != nil {
		return "", err
	}

	diffIO := bytes.NewBufferString(schema)

	// Copy the contents of the file to the form field
	if _, err := io.Copy(fw, diffIO); err != nil {
		return "", err
	}

	if err := w.Close(); err != nil {
		return "", err
	}

	var diffUrl = d.cfg.EnvConfig[d.env].DirectusURL + "/" + d.cfg.DirectusConfig.DiffURL
	if ctx.Value(shared.VersionForcer("forceVersion")).(bool) {
		diffUrl += "?force=true"
	}

	// Send the request
	req, err := http.NewRequest(
		"POST",
		diffUrl,
		&buf)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+d.bearer)
	req.Header.Set("Content-Length", strconv.Itoa(buf.Len()))

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// if there's no diff, directus returns empty response.
	if len(body) == 0 {
		ret, err := json.Marshal(map[string]string{"result": "identical"})
		return string(ret), err
	}

	var res DirectusSchema
	err = json.Unmarshal(body, &res)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		d.log.Errorln(string(body))
		return "", err
	}

	return res.Convert(), nil

}

func (d *directus) Apply(ctx context.Context, diff string) (string, error) {

	buf := strings.NewReader(diff)

	var applyUrl = d.cfg.EnvConfig[d.env].DirectusURL + "/" + d.cfg.DirectusConfig.ApplyURL
	if ctx.Value(shared.VersionForcer("forceVersion")).(bool) {
		applyUrl += "?force=true"
	}

	// Send the request
	req, err := http.NewRequest(
		"POST",
		applyUrl,
		buf)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+d.bearer)
	req.Header.Set("Content-Length", strconv.Itoa(buf.Len()))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusNoContent {
		d.log.Errorln(string(body), "statuscode:", resp.StatusCode)
		return "", err
	}

	return string(body), nil

}
