package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

// date required to make xml signature
type Signature struct {
	AuthorCertificate string
	AuthorPass        string
	references        []*reference
	name              SignatureType
}

type SignatureType string

var (
	AuthorSignature      SignatureType = "author-signature.xml"
	DistributorSignature SignatureType = "signature1.xml"
)

type reference struct {
	URI    string
	Digest string
}

func NewSignature(name SignatureType, files []PackageFile) (sig *Signature, err error) {
	sig = &Signature{}
	for _, file := range files {
		rc, err := file.GetReadCloser()
		if err != nil {
			return nil, err
		}
		ref, err := createReference(file.PackagePath(), rc)
		rc.Close()
		if err != nil {
			return nil, err
		}
		sig.references = append(sig.references, ref)
	}
	sig.name = name
	return sig, nil
}

func (sign *Signature) PackagePath() string {
	return string(sign.name)
}

// simple wrapper aroung bytes.Buffer to support io.ReadCloser
type signatureBuffer struct {
	*bytes.Buffer
}

func (sign *signatureBuffer) Close() error {
	return nil
}

func (sign *Signature) GetReadCloser() (io.ReadCloser, error) {
	buf := &bytes.Buffer{}

	if sign.name == AuthorSignature {
		if err := authorSignatureTmpl.Execute(buf, sign.references); err != nil {
			return nil, fmt.Errorf("Template execute failed: %v", err)
		}
	} else if sign.name == DistributorSignature {
		if err := distributorSignatureTmpl.Execute(buf, sign.references); err != nil {
			return nil, fmt.Errorf("Template execute failed: %v", err)
		}
	}
	file, err := ioutil.TempFile("/tmp/", "signtmp")
	if err != nil {
		return nil, fmt.Errorf("TempFile failed: %v", err)
	}
	_, err = io.Copy(file, buf)
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("TempFile failed: %v", err)
	}
	file.Sync()
	log.Print("Opening file: ", file.Name())
	cmd := exec.Command("xmlsec1", "--sign", "--output", sign.PackagePath(), "--pkcs12", sign.AuthorCertificate, "--pwd", sign.AuthorPass, file.Name())
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	file.Close()
	if err != nil {
		return nil, fmt.Errorf("Unable to run xmlsec1: %s", stderr.String())
	}
	file, err = os.Open(sign.PackagePath())
	if err != nil {
		return nil, fmt.Errorf("Open failed: %v", err)
	}
	return file, nil
}

func createReference(uri string, reader io.Reader) (*reference, error) {
	var ref reference
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("Unable to append reference: %v", err)
	}
	checksum := sha256.Sum256(bytes)
	ref.URI = uri
	ref.Digest = base64.StdEncoding.EncodeToString([]byte(checksum[:]))

	return &ref, nil
}

var authorSignatureTmpl = template.Must(template.New("author-signature").Parse(`
<Signature xmlns="http://www.w3.org/2000/09/xmldsig#" Id="AuthorSignature">
<SignedInfo>
<CanonicalizationMethod Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"></CanonicalizationMethod>
<SignatureMethod Algorithm="http://www.w3.org/2001/04/xmldsig-more#rsa-sha256"></SignatureMethod>
{{range .}}
<Reference URI="{{.Uri}}">
<DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"></DigestMethod>
<DigestValue>{{.Digest}}</DigestValue>
</Reference>
{{end}}
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

var distributorSignatureTmpl = template.Must(template.New("signature1").Parse(`
<Signature xmlns="http://www.w3.org/2000/09/xmldsig#" Id="DistributorSignature">
<SignedInfo>
<CanonicalizationMethod Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"></CanonicalizationMethod>
<SignatureMethod Algorithm="http://www.w3.org/2001/04/xmldsig-more#rsa-sha256"></SignatureMethod>
{{range .}}
<Reference URI="{{.Uri}}">
<DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"></DigestMethod>
<DigestValue>{{.Digest}}</DigestValue>
</Reference>
{{end}}
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
<Object Id="prop"><SignatureProperties xmlns:dsp="http://www.w3.org/2009/xmldsig-properties"><SignatureProperty Id="profile" Target="#DistributorSignature"><dsp:Profile URI="http://www.w3.org/ns/widgets-digsig#profile"></dsp:Profile></SignatureProperty><SignatureProperty Id="role" Target="#DistributorSignature"><dsp:Role URI="http://www.w3.org/ns/widgets-digsig#role-distributor"></dsp:Role></SignatureProperty><SignatureProperty Id="identifier" Target="#DistributorSignature"><dsp:Identifier></dsp:Identifier></SignatureProperty></SignatureProperties></Object>
</Signature>
`))
