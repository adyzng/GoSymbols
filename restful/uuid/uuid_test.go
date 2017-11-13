package uuid

import (
	"testing"
)

const iterations = 1000

func TestGenerate(t *testing.T) {
	for i := 0; i < iterations; i++ {
		u := Generate()

		if u.Version() != 0x20 {
			t.Fatalf("generated uuid is not valid <%s>", u)
		}

		sid := u.String()
		if len(sid) != len("974AFFD3-1BCC-4475-8910-A967AFAE51FE") {
			t.Fatalf("uuid to string failed. <%s>", sid)
		}

		rid := u.RawString()
		if len(rid) != len("974AFFD31BCC44758910A967AFAE51FE") {
			t.Fatalf("uuid to raw string failed. <%s>", rid)
		}
	}
}

func TestParse(t *testing.T) {
	for i := 0; i < iterations; i++ {
		u := Generate()
		parsed, err := Parse(u.String())

		if err != nil {
			t.Fatalf("parse uuid failed %v: %v", u, err)
		}

		if parsed != u {
			t.Fatalf("parsed uuid not equal origin %v: %v", u, parsed)
		}
	}

	for _, c := range []string{
		"abcded",
		"{0AD4A3E1-AA6F-4986-97BC-26BCA7E113E5}", // invalid format
		"  0C79A2E2-D363-4E0F-8A61-091027200D6B", // leading space
		"59AEBE38-7B8A-4185-9D15-685B85154EE2  ", // trailing space
		"00000000-0000-0000-0000-x00000000000",   // correct length, invalid character
	} {
		if _, err := Parse(c); err == nil {
			t.Fatalf("parsing %q should fail", c)
		} else {
			t.Logf("parsing expected error : %v", err)
		}
	}
}

func BenchmarkGenerate(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Generate()
			//b.Logf("uuid: %v", u)
		}
	})
}

func BenchmarkParse(b *testing.B) {
	idx, cnt := 0, 4
	idList := []string{
		"D1A6710587134CE3AA6FF50465B2A37A",
		"B156EF0B-83A2-4922-8308-A24631E285A6",
		"D03A9D910A8F4ADBA41DCDD52EDF121C",
		"E45830C4-1E08-411D-99CE-6140102BD683",
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c := idx % cnt
			idx++

			if _, err := Parse(idList[c]); err != nil {
				b.Fatalf("parse uuid failed <%v : %v>", idList[c], err)
			}
		}
	})
}
