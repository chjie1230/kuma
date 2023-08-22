package store

import (
	"fmt"
	"time"

	core_model "github.com/kumahq/kuma/pkg/core/resources/model"
)

type CreateOptions struct {
	Name         string
	Mesh         string
	CreationTime time.Time
	Owner        core_model.Resource
}

type CreateOptionsFunc func(*CreateOptions)

func NewCreateOptions(fs ...CreateOptionsFunc) *CreateOptions {
	opts := &CreateOptions{}
	for _, f := range fs {
		f(opts)
	}
	return opts
}

func CreateBy(key core_model.ResourceKey) CreateOptionsFunc {
	return CreateByKey(key.Name, key.Mesh)
}

func CreateByKey(name, mesh string) CreateOptionsFunc {
	return func(opts *CreateOptions) {
		opts.Name = name
		opts.Mesh = mesh
	}
}

func CreatedAt(creationTime time.Time) CreateOptionsFunc {
	return func(opts *CreateOptions) {
		opts.CreationTime = creationTime
	}
}

func CreateWithOwner(owner core_model.Resource) CreateOptionsFunc {
	return func(opts *CreateOptions) {
		opts.Owner = owner
	}
}

type UpdateOptions struct {
	ModificationTime time.Time
}

func ModifiedAt(modificationTime time.Time) UpdateOptionsFunc {
	return func(opts *UpdateOptions) {
		opts.ModificationTime = modificationTime
	}
}

type UpdateOptionsFunc func(*UpdateOptions)

func NewUpdateOptions(fs ...UpdateOptionsFunc) *UpdateOptions {
	opts := &UpdateOptions{}
	for _, f := range fs {
		f(opts)
	}
	return opts
}

type DeleteOptions struct {
	Name string
	Mesh string
}

type DeleteOptionsFunc func(*DeleteOptions)

func NewDeleteOptions(fs ...DeleteOptionsFunc) *DeleteOptions {
	opts := &DeleteOptions{}
	for _, f := range fs {
		f(opts)
	}
	return opts
}

func DeleteBy(key core_model.ResourceKey) DeleteOptionsFunc {
	return DeleteByKey(key.Name, key.Mesh)
}

func DeleteByKey(name, mesh string) DeleteOptionsFunc {
	return func(opts *DeleteOptions) {
		opts.Name = name
		opts.Mesh = mesh
	}
}

type DeleteAllOptions struct {
	Mesh string
}

type DeleteAllOptionsFunc func(*DeleteAllOptions)

func DeleteAllByMesh(mesh string) DeleteAllOptionsFunc {
	return func(opts *DeleteAllOptions) {
		opts.Mesh = mesh
	}
}

func NewDeleteAllOptions(fs ...DeleteAllOptionsFunc) *DeleteAllOptions {
	opts := &DeleteAllOptions{}
	for _, f := range fs {
		f(opts)
	}
	return opts
}

type GetOptions struct {
	Name    string
	Mesh    string
	Version string
}

type GetOptionsFunc func(*GetOptions)

func NewGetOptions(fs ...GetOptionsFunc) *GetOptions {
	opts := &GetOptions{}
	for _, f := range fs {
		f(opts)
	}
	return opts
}

func GetBy(key core_model.ResourceKey) GetOptionsFunc {
	return GetByKey(key.Name, key.Mesh)
}

func GetByKey(name, mesh string) GetOptionsFunc {
	return func(opts *GetOptions) {
		opts.Name = name
		opts.Mesh = mesh
	}
}

func GetByVersion(version string) GetOptionsFunc {
	return func(opts *GetOptions) {
		opts.Version = version
	}
}

func (g *GetOptions) HashCode() string {
	return fmt.Sprintf("%s:%s", g.Name, g.Mesh)
}

type (
	ListFilterFunc func(rs core_model.Resource) bool
)

type ListOptions struct {
	Mesh         string
	PageSize     int
	PageOffset   string
	FilterFunc   ListFilterFunc
	NameContains string
	Ordered      bool
}

type ListOptionsFunc func(*ListOptions)

func NewListOptions(fs ...ListOptionsFunc) *ListOptions {
	opts := &ListOptions{}
	for _, f := range fs {
		f(opts)
	}
	return opts
}

// Filter returns true if the item passes the filtering criteria
func (l *ListOptions) Filter(rs core_model.Resource) bool {
	if l.FilterFunc == nil {
		return true
	}

	return l.FilterFunc(rs)
}

func ListByNameContains(name string) ListOptionsFunc {
	return func(opts *ListOptions) {
		opts.NameContains = name
	}
}

func ListByMesh(mesh string) ListOptionsFunc {
	return func(opts *ListOptions) {
		opts.Mesh = mesh
	}
}

func ListByPage(size int, offset string) ListOptionsFunc {
	return func(opts *ListOptions) {
		opts.PageSize = size
		opts.PageOffset = offset
	}
}

func ListByFilterFunc(filterFunc ListFilterFunc) ListOptionsFunc {
	return func(opts *ListOptions) {
		opts.FilterFunc = filterFunc
	}
}

func ListOrdered() ListOptionsFunc {
	return func(opts *ListOptions) {
		opts.Ordered = true
	}
}

func (l *ListOptions) IsCacheable() bool {
	return l.FilterFunc == nil
}

func (l *ListOptions) HashCode() string {
	return fmt.Sprintf("%s:%t:%s:%d:%s", l.Mesh, l.Ordered, l.NameContains, l.PageSize, l.PageOffset)
}

type resourceMetaWrap struct {
	resourceMeta     core_model.ResourceMeta
	meshName         *string
	name             *string
	version          *string
	nameExtensions   core_model.ResourceNameExtensions
	createTime       *time.Time
	modificationTime *time.Time
}

func (r resourceMetaWrap) GetName() string {
	if r.name == nil {
		return r.resourceMeta.GetName()
	} else {
		return *r.name
	}
}

func (r resourceMetaWrap) GetNameExtensions() core_model.ResourceNameExtensions {
	if r.nameExtensions == nil {
		return r.resourceMeta.GetNameExtensions()
	} else {
		return r.nameExtensions
	}
}

func (r resourceMetaWrap) GetVersion() string {
	if r.version == nil {
		return r.resourceMeta.GetVersion()
	} else {
		return *r.version
	}
}

func (r resourceMetaWrap) GetMesh() string {
	if r.meshName == nil {
		return r.resourceMeta.GetMesh()
	} else {
		return *r.meshName
	}
}

func (r resourceMetaWrap) GetCreationTime() time.Time {
	if r.createTime == nil {
		return r.resourceMeta.GetCreationTime()
	} else {
		return *r.createTime
	}
}

func (r resourceMetaWrap) GetModificationTime() time.Time {
	if r.modificationTime == nil {
		return r.resourceMeta.GetModificationTime()
	} else {
		return *r.modificationTime
	}
}

type ResourceMetaWrapFunc func(wrap *resourceMetaWrap)

func WithMeshName(meshName string) ResourceMetaWrapFunc {
	return func(wrap *resourceMetaWrap) {
		wrap.meshName = &meshName
	}
}

func WithResourceMeta(resourceMeta core_model.ResourceMeta) ResourceMetaWrapFunc {
	return func(wrap *resourceMetaWrap) {
		wrap.resourceMeta = resourceMeta
	}
}

func WithName(name string) ResourceMetaWrapFunc {
	return func(wrap *resourceMetaWrap) {
		wrap.name = &name
	}
}

func WithVersion(version string) ResourceMetaWrapFunc {
	return func(wrap *resourceMetaWrap) {
		wrap.version = &version
	}
}

func WithResourceNameExtensions(resourceNameExtensions core_model.ResourceNameExtensions) ResourceMetaWrapFunc {
	return func(wrap *resourceMetaWrap) {
		wrap.nameExtensions = resourceNameExtensions
	}
}

func WithCreationTime(creationTime time.Time) ResourceMetaWrapFunc {
	return func(wrap *resourceMetaWrap) {
		wrap.createTime = &creationTime
	}
}

func WithModificationTime(modificationTime time.Time) ResourceMetaWrapFunc {
	return func(wrap *resourceMetaWrap) {
		wrap.modificationTime = &modificationTime
	}
}

func ResouceMetaWrap(resourceMetaWrapFunc ...ResourceMetaWrapFunc) core_model.ResourceMeta {
	wrap := &resourceMetaWrap{}
	for _, fn := range resourceMetaWrapFunc {
		fn(wrap)
	}
	return wrap
}

type ListIndexOptions struct {
	Name string
	Mesh string
	Type core_model.ResourceType
}

type ListIndexOptionsFunc func(*ListIndexOptions)

func WithListIndexOptions(typ core_model.ResourceType, mesh, name string) ListIndexOptionsFunc {
	return func(opts *ListIndexOptions) {
		opts.Type = typ
		opts.Mesh = mesh
		opts.Name = name
	}
}

func NewListIndexOptions(fs ...ListIndexOptionsFunc) *ListIndexOptions {
	opts := &ListIndexOptions{}
	for _, f := range fs {
		f(opts)
	}
	return opts
}

type DeleteIndexOptions struct {
	Name string
	Mesh string
	Type core_model.ResourceType

	Owner struct {
		Name string
		Mesh string
		Type core_model.ResourceType
	}
}

type DeleteIndexOptionsFunc func(*DeleteIndexOptions)

func WithDeleteIndexOptions(typ core_model.ResourceType, mesh, name string) DeleteIndexOptionsFunc {
	return func(opts *DeleteIndexOptions) {
		opts.Type = typ
		opts.Mesh = mesh
		opts.Name = name
	}
}

func WithDeleteIndexOwnerOptions(typ core_model.ResourceType, mesh, name string) DeleteIndexOptionsFunc {
	return func(opts *DeleteIndexOptions) {
		opts.Owner.Type = typ
		opts.Owner.Mesh = mesh
		opts.Owner.Name = name
	}
}

func NewDeleteIndexOptions(fs ...DeleteIndexOptionsFunc) *DeleteIndexOptions {
	opts := &DeleteIndexOptions{}
	for _, f := range fs {
		f(opts)
	}
	return opts
}
