package v1

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/adyzng/goPDB/restful"
	"github.com/adyzng/goPDB/symbol"
	"github.com/gorilla/mux"

	log "gopkg.in/clog.v1"
)

// RestBranchList response to restful API
//	[:]/api/v1/branch
//
//	@ return {
//		Total: 		int
//		Branchs: 	[]*symbol.Branch
//	}
//
func RestBranchList(w http.ResponseWriter, r *http.Request) {
	bs := restful.BranchList{}
	symbol.GetServer().WalkBuilders(func(bu symbol.Builder) error {
		if b, ok := bu.(*symbol.BrBuilder); ok {
			bs.Total++
			bs.Branchs = append(bs.Branchs, &b.Branch)
		}
		return nil
	})
	resp := restful.RestResponse{
		Data: &bs,
	}
	resp.WriteJSON(w)
}

// RestBuildList response to restful API
//	[:]/api/v1/branch/:name
//
//	@:name {branch name}
//
//	@return {
//		Total: 		int
//		Builds: 	[]*symbol.Build
//	}
//
func RestBuildList(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	resp := restful.RestResponse{
		ErrCodeMsg: restful.ErrInvalidParam,
	}

	if sname, ok := vars["name"]; ok {
		builder := symbol.GetServer().GetBuilder(sname)
		if builder != nil {
			blst := restful.BuildList{
				Branch: sname,
			}
			_, err := builder.ParseBuilds(func(build *symbol.Build) error {
				blst.Total++
				blst.Builds = append(blst.Builds, build)
				return nil
			})
			if err != nil {
				log.Error(2, "[Restful] Parse builds for %s failed: %v.", sname, err)
			}
			resp.Data = blst
			resp.ErrCodeMsg = restful.ErrSucceed
		} else {
			resp.ErrCodeMsg.Message = "no such branch"
		}
	}
	resp.WriteJSON(w)
}

// RestSymbolList response to restful API
//	[:]/api/v1/branch/:name/:bid
//
//	@:name {branch name}
//	@:bid  {build id}
//
//	@ return {
//		Total: 		int
//		Builds: 	[]*symbol.Build
//	}
//
func RestSymbolList(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	resp := restful.RestResponse{
		ErrCodeMsg: restful.ErrInvalidParam,
	}

	sname, bid := vars["name"], vars["bid"]
	if sname != "" && bid != "" {
		buider := symbol.GetServer().GetBuilder(sname)
		if buider != nil {
			symLst := restful.SymbolList{
				Branch: sname,
				Build:  bid,
			}
			_, err := buider.ParseSymbols(bid, func(sym *symbol.Symbol) error {
				symLst.Total++
				symLst.Symbols = append(symLst.Symbols, sym)
				return nil
			})
			if err != nil {
				log.Error(2, "[Restful] Parse symbols for %s:%s failed: %v.",
					sname, bid, err)
			}
			resp.Data = symLst
			resp.ErrCodeMsg = restful.ErrSucceed
		} else {
			resp.ErrCodeMsg.Message = "no such build"
		}
	}
	resp.WriteJSON(w)
}

// DownloadSymbol response download symbol file api
//	[:]/api/v1/branch/:name/:hash?name=name
//
//	@:name	{branch name}
//	@:hash	{file hash}
//	@?name	{file name}
//
//	@ return file
//
func DownloadSymbol(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sname, hash := vars["name"], vars["hash"]

	r.ParseForm()
	fname := r.FormValue("name")

	if sname == "" || hash == "" || fname == "" {
		log.Warn("[Restful] Download symbol invalid param: [%s, %s, %s]",
			sname, hash, fname)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	buider := symbol.GetServer().GetBuilder(sname)
	if buider == nil {
		log.Warn("[Restful] Download symbol branch not exist: [%s, %s, %s]",
			sname, hash, fname)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fpath := buider.GetSymbolPath(hash, fname)
	fd, err := os.OpenFile(fpath, os.O_RDONLY, 666)
	if err != nil {
		log.Warn("[Restful] Open symbol file %s failed: %v.", fpath, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer fd.Close()

	// set response header
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fname))

	// send fil content
	var size int64
	if size, err = io.Copy(w, fd); err != nil {
		log.Error(2, "[Restful] Send file failed: %v.", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//w.WriteHeader(http.StatusOK)
	log.Trace("[Restful] Send file complete. [%d: %s]", size, fpath)
}
