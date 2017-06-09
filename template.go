package saml

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
)

// TemplateProvider is an interface used by IdpAuthnRequest to render
// any needed templates.
type TemplateProvider interface {
	MakeHTTPPostTemplate(w http.ResponseWriter, url string, samlResponse string, relayState string) error
}

// DefaultTemplateProvider produces a set of default templates
type DefaultTemplateProvider struct {
}

// MakeHTTPPostTemplate will write out the template for ACS endpoints
// that use HTTP Post
func (dtp *DefaultTemplateProvider) MakeHTTPPostTemplate(w http.ResponseWriter, url string, samlResponse string, relayState string) error {
	tmpl := template.Must(template.New("saml-post-form").Parse(`<html>` +
		`<form method="post" action="{{.URL}}" id="SAMLResponseForm">` +
		`<input type="hidden" name="SAMLResponse" value="{{.SAMLResponse}}" />` +
		`<input type="hidden" name="RelayState" value="{{.RelayState}}" />` +
		`<input type="submit" value="Continue" />` +
		`</form>` +
		`<script>document.getElementById('SAMLResponseForm').submit();</script>` +
		`</html>`))
	data := struct {
		URL          string
		SAMLResponse string
		RelayState   string
	}{
		URL:          url,
		SAMLResponse: samlResponse,
		RelayState:   relayState,
	}

	buf := bytes.NewBuffer(nil)
	if err := tmpl.Execute(buf, data); err != nil {
		return err
	}
	if _, err := io.Copy(w, buf); err != nil {
		return err
	}
	return nil
}
