package bulletinimpl

type Config struct {
	Xmlns                     string `json:"xmlns"`
	XmlnsCvrf                 string `json:"xmlns_cvrf"`
	ContactDetails            string `json:"contact_details"`
	IssuingAuthority          string `json:"issuing_authority"`
	SecurityBulletinUrlPrefix string `json:"security_bulletin_url_prefix"`
	DefectUrlPrefix           string `json:"defect_url_prefix"`
}

func (c *Config) SetDefault() {
	if c.Xmlns == "" {
		c.Xmlns = "http://www.icasi.org/CVRF/schema/cvrf/1.1"
	}

	if c.XmlnsCvrf == "" {
		c.XmlnsCvrf = "http://www.icasi.org/CVRF/schema/cvrf/1.1"
	}

	if c.ContactDetails == "" {
		c.ContactDetails = "openeuler-release@openeuler.org"
	}

	if c.IssuingAuthority == "" {
		c.IssuingAuthority = "openEuler release SIG"
	}

	if c.SecurityBulletinUrlPrefix == "" {
		c.SecurityBulletinUrlPrefix = "https://www.openeuler.org/en/security/safety-bulletin/detail.html?id="
	}

	if c.DefectUrlPrefix == "" {
		c.DefectUrlPrefix = "https://www.openeuler.org/en/security/cve/detail.html?id="
	}
}
