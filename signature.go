package main

import (
	_ "github.com/lestrrat/go-xmlsec"
	_ "github.com/lestrrat/go-xmlsec/dsig"
	"io"
)

// date required to make xml signature
type Signature struct {
	AuthorCertificate      []byte
	DistributorCertificate []byte
	Files                  []string
}

func (this *Signature) MarshalXMLSignature(writer io.Writer) error {
	/*
		if err := xmlsec.Init(); err != nil {
			return err
		}

		defer xmlsec.Shutdown()
	*/
	return nil
}

/*
func (this *XMLSignature) AppendReference(uri string, reader io.Reader) error {
	var ref reference
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("Unable to append reference: %v", err)
	}
	checksum := sha256.Sum256(bytes)
	ref.Uri = uri
	ref.Value = base64.StdEncoding.EncodeToString([]byte(checksum[:]))
	ref.Method.Algorithm = sha256algorithm

	this.References = append(this.References, ref)
	return nil
}
*/
