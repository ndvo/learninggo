package viewbuilder

import(
	"testing"
)

func TestFindViews(t *testing.T){
	flist, err := FindViews("../api")
	if len(flist)==0 || err!= nil{
		t.Error("Nothing found")
	}
}

func  TestBuildIndex(t *testing.T){
	viewslist, err := FindViews("../api")
	println(err.Error())
	for _, v := range viewslist{
		err:=IndexView(v)
		println(err.Error())
	}
}
