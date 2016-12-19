package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
)

// date required to make xml signature
type Signature struct {
	AuthorCertificate      []byte
	DistributorCertificate []byte
	references             []reference
}

type reference struct {
	Uri    string
	Digest string
}

func NewSignature(files []PackageFile) (*Signature, error) {
	for _, file := range files {
	}
	return nil, nil
}

func (this *Signature) Path() string {
	return "author-signature.xml"
}

func (this *Signature) WriteContent(writer io.Writer) error {
	return nil
}

func createReference(uri string, reader io.Reader) (*reference, error) {
	var ref reference
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("Unable to append reference: %v", err)
	}
	checksum := sha256.Sum256(bytes)
	ref.Uri = uri
	ref.Digest = base64.StdEncoding.EncodeToString([]byte(checksum[:]))

	return &ref, nil
}

var authorSignatureTmpl = template.Must(template.New("author-signature").Parse(`
<Signature xmlns="http://www.w3.org/2000/09/xmldsig#" Id="AuthorSignature">
<SignedInfo>
<CanonicalizationMethod Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"></CanonicalizationMethod>
<SignatureMethod Algorithm="http://www.w3.org/2001/04/xmldsig-more#rsa-sha256"></SignatureMethod>

<Reference URI="{{.Uri}}">
<DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"></DigestMethod>
<DigestValue>{{.Digest}}</DigestValue>

<Reference URI="#prop">
<Transforms>
<Transform Algorithm="http://www.w3.org/2006/12/xml-c14n11"></Transform>
</Transforms>
<DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"></DigestMethod>
<DigestValue></DigestValue>
</Reference>
</SignedInfo>

<SignatureValue>
</SignatureValue>
<KeyInfo>
<X509Data>
<X509Certificate>
</X509Certificate>
</X509Data>
</KeyInfo>
<Object Id="prop"><SignatureProperties xmlns:dsp="http://www.w3.org/2009/xmldsig-properties"><SignatureProperty Id="profile" Target="#AuthorSignature"><dsp:Profile URI="http://www.w3.org/ns/widgets-digsig#profile"></dsp:Profile></SignatureProperty><SignatureProperty Id="role" Target="#AuthorSignature"><dsp:Role URI="http://www.w3.org/ns/widgets-digsig#role-author"></dsp:Role></SignatureProperty><SignatureProperty Id="identifier" Target="#AuthorSignature"><dsp:Identifier></dsp:Identifier></SignatureProperty></SignatureProperties></Object>
</Signature>
`))
