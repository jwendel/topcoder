package auth

import (
	"testing"
)

func TestDomainExists(t *testing.T) {
	s := datastore{}
	s.Init("test_data.json")

	if ok := s.DomainExists("topcoder.com"); !ok {
		t.Errorf("DomainExists('topcoder.com') = %v, wanted true", ok)
	}

	if ok := s.DomainExists("appirio.com"); !ok {
		t.Errorf("DomainExists('appirio.com') = %v, wanted true", ok)
	}

	if ok := s.DomainExists("google.com"); ok {
		t.Errorf("DomainExists('google.com') = %v, wanted false", ok)
	}

	if ok := s.DomainExists(""); ok {
		t.Errorf("DomainExists('') = %v, wanted false", ok)
	}
}

func TestUserPasswordValid(t *testing.T) {
	s := datastore{}
	s.Init("test_data.json")

	if ok := s.UserPasswordValid("topcoder.com", "teru", EncryptPassword("ilovejava")); !ok {
		t.Errorf("UserPasswordValid(topcoder.com, teru, ilovejava) = %v, wanted true", ok)
	}

	if ok := s.UserPasswordValid("appirio.com", "chris", EncryptPassword("ilovesushi")); !ok {
		t.Errorf("UserPasswordValid(appirio.com, chris, ilovesushi) = %v, wanted true", ok)
	}

	if ok := s.UserPasswordValid("appirio.com", "narinder", EncryptPassword("ilovesamurai")); !ok {
		t.Errorf("UserPasswordValid(appirio.com, narinder, ilovesamurai) = %v, wanted true", ok)
	}

	if ok := s.UserPasswordValid("topcoder.com", "narinder", EncryptPassword("ilovesamurai")); ok {
		t.Errorf("UserPasswordValid(topcoder.com, narinder, ilovesamurai) = %v, wanted false", ok)
	}

	if ok := s.UserPasswordValid("appirio.com", "jun", EncryptPassword("ilovesamurai")); ok {
		t.Errorf("UserPasswordValid(appirio.com, jun, ilovesamurai) = %v, wanted false", ok)
	}

	if ok := s.UserPasswordValid("appirio.com", "narinder", EncryptPassword("ihatesamurai")); ok {
		t.Errorf("UserPasswordValid(appirio.com, narinder, ihatesamurai) = %v, wanted false", ok)
	}

	if ok := s.UserPasswordValid("", "", EncryptPassword("")); ok {
		t.Errorf("UserPasswordValid('','','') = %v, wanted false", ok)
	}

	if ok := s.UserPasswordValid("", "", ""); ok {
		t.Errorf("UserPasswordValid('','','') = %v, wanted false", ok)
	}

	if ok := s.UserPasswordValid("topcoder.com", "teru", ""); ok {
		t.Errorf("UserPasswordValid(topcoder.com, teru, '') = %v, wanted false", ok)
	}
}
