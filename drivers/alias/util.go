package alias

import (
	"context"
	stdpath "path"
	"strings"

	"github.com/alist-org/alist/v3/internal/fs"
	"github.com/alist-org/alist/v3/internal/model"
	"github.com/alist-org/alist/v3/pkg/utils"
)

func (d *Alias) listRoot() []model.Obj {
	var objs []model.Obj
	for k, _ := range d.pathMap {
		obj := model.Object{
			Name:     k,
			IsFolder: true,
			Modified: d.Modified,
		}
		objs = append(objs, &obj)
	}
	return objs
}

// do others that not defined in Driver interface
func getPair(path string) (string, string) {
	//path = strings.TrimSpace(path)
	if strings.Contains(path, ":") {
		pair := strings.SplitN(path, ":", 2)
		if !strings.Contains(pair[0], "/") {
			return pair[0], pair[1]
		}
	}
	return stdpath.Base(path), path
}

func getRootAndPath(path string) (string, string) {
	path = strings.TrimPrefix(path, "/")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], parts[1]
}

func (d *Alias) get(ctx context.Context, path string, dst, sub string) (model.Obj, error) {
	obj, err := fs.Get(ctx, stdpath.Join(dst, sub))
	if err != nil {
		return nil, err
	}
	return &model.Object{
		Path:     path,
		Name:     obj.GetName(),
		Size:     obj.GetSize(),
		Modified: obj.ModTime(),
		IsFolder: obj.IsDir(),
	}, nil
}

func (d *Alias) list(ctx context.Context, dst, sub string) ([]model.Obj, error) {
	objs, err := fs.List(ctx, stdpath.Join(dst, sub))
	// the obj must implement the model.SetPath interface
	// return objs, err
	if err != nil {
		return nil, err
	}
	return utils.SliceConvert(objs, func(obj model.Obj) (model.Obj, error) {
		return &model.Object{
			Name:     obj.GetName(),
			Size:     obj.GetSize(),
			Modified: obj.ModTime(),
			IsFolder: obj.IsDir(),
		}, nil
	})
}

func (d *Alias) link(ctx context.Context, dst, sub string, args model.LinkArgs) (*model.Link, error) {
	link, _, err := fs.Link(ctx, stdpath.Join(dst, sub), args)
	return link, err
}
