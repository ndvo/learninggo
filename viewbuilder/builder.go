package viewbuilder

import( 
"os"
"io/ioutil"
"path/filepath"
)

var viewdir  string = "api"


type Explorer struct{
	journey_start int
}

// Explore walks through all directories searching untracked content
func (e Explorer) explore(d string){
	// checks everything in directory d


	//dircontents := os.Readdir
	//dirnames := os.Readdirnames
}

// 
func FindViews(dir string) (flist []os.FileInfo, err error){
	println("FindViews executed")
	println(dir)
	flist, err = ioutil.ReadDir(dir)
	println(err)
	for _, d := range flist{
		if !d.IsDir(){
			if filepath.Ext(d.Name()) == ".json" {
				println(d.Name())
			}
		}
		println(d.Name(), d.IsDir())
	}
	return
}

