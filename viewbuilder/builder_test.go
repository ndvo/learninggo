package viewbuilder

import(
	"testing"
)

func TestFindViews(t *testing.T){
	println("Testing started")
	flist, err := FindViews("../api")
	if len(flist)==0 || err!= nil{
		t.Error("Nothing found")
	}
}
