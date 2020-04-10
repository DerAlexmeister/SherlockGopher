package sherlockanalyser

import "testing"

func Test1(t *testing.T) {
	atask := NewTask()

	atask.setHTMLCode("<html><head><title>MMIX</title></head><body><p>Private Website von Johannes Walter. <br>E-Mail: someone@walterj.de <br><a href='/impressum.html'>Impressum</a><a href='/impressum.html'>Datenschutz</a></p><br><br><h1>MMIX</h1><br><p>Tutoriumsunterlagen: <a href='https://github.com/jwalteri/MMIX'>GitHub von Johannes Walter</a></p></body></html>")
	atask.setAddr("https://walterj.de/mmix.html")
	atask.setTaskID(1)

	atask.Execute()

	expected := make([]string, 0)
	expected = append(expected, "https://walterj.de/impressum.html")
	expected = append(expected, "https://walterj.de/impressum.html")
	expected = append(expected, "https://github.com/jwalteri/MMIX")

	if len(expected) != len(atask.getFoundLinks()) {
		t.Errorf("got %d elements, want %d elements", len(atask.foundLinks), len(expected))
	}

	for i, ele := range atask.getFoundLinks() {
		if ele != expected[i] {
			t.Errorf("got %q, want %q", ele, expected[i])
		}
	}
}
