package etcd

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	clientv3 "go.etcd.io/etcd/client/v3"

	core_model "github.com/kumahq/kuma/pkg/core/resources/model"
	"github.com/kumahq/kuma/pkg/core/resources/store"
	core_metrics "github.com/kumahq/kuma/pkg/metrics"
	"github.com/kumahq/kuma/pkg/plugins/common/etcd"
)

type ResoucesEtcdStore struct {
	prefix  string
	client  *clientv3.Client
	metrics core_metrics.Metrics
}

func (e ResoucesEtcdStore) Create(ctx context.Context, resource core_model.Resource, optionsFunc ...store.CreateOptionsFunc) error {
	opts := store.NewCreateOptions(optionsFunc...)
	specBytes, err := core_model.ToJSON(resource.GetSpec())
	if err != nil {
		return errors.Wrap(err, "failed to convert spec to json")
	}
	var owner *Owner
	if opts.Owner != nil {
		owner = &Owner{
			Type: opts.Owner.Descriptor().Name,
			Name: opts.Owner.GetMeta().GetName(),
			Mesh: opts.Owner.GetMeta().GetMesh(),
		}
	}

	resourceMetaWrap := store.ResouceMetaWrap(store.WithMeshName(opts.Mesh), store.WithName(opts.Name), store.WithVersion("0"),
		store.WithCreationTime(opts.CreationTime), store.WithModificationTime(opts.CreationTime))
	etcdResourceMeta := ToEtcdResourceMetaObject(resourceMetaWrap, owner)
	metaBytes, err := core_model.ToJSON(etcdResourceMeta)
	if err != nil {
		return errors.Wrap(err, "failed to convert spec to json")
	}
	value := resourceObject{
		MetaData: metaBytes,
		SpecData: specBytes,
	}

	bytes, err := json.Marshal(value)
	if err != nil {
		return errors.Wrap(err, "json.Marshal error")
	}
	key := etcd.NewEtcdResourcedKey(e.prefix, resource.Descriptor().Name, opts.Mesh, opts.Name).String()
	putResponse, err := e.client.Put(ctx, key, string(bytes), clientv3.WithPrevKV())
	if err != nil {
		return errors.Wrap(err, "etcd put error")
	}
	if putResponse.PrevKv != nil {
		return store.ErrorResourceAlreadyExists(resource.Descriptor().Name, opts.Name, opts.Mesh)
	}
	// update resource's meta with new version
	resource.SetMeta(resourceMetaWrap)

	return nil
}

func (e ResoucesEtcdStore) Update(ctx context.Context, resource core_model.Resource, optionsFunc ...store.UpdateOptionsFunc) error {
	opts := store.NewUpdateOptions(optionsFunc...)

	specBytes, err := core_model.ToJSON(resource.GetSpec())
	if err != nil {
		return errors.Wrap(err, "failed to convert spec to json")
	}

	version, err := strconv.Atoi(resource.GetMeta().GetVersion())
	if err != nil {
		return errors.Wrap(err, "failed to convert meta version to int")
	}
	newVersion := version + 1
	resourceMeta := resource.GetMeta()
	resourceMetaWrap := store.ResouceMetaWrap(store.WithResourceMeta(resourceMeta), store.WithVersion(strconv.Itoa(newVersion)), store.WithModificationTime(opts.ModificationTime))
	key := etcd.NewEtcdResourcedKey(e.prefix, resource.Descriptor().Name, resource.GetMeta().GetMesh(), resource.GetMeta().GetName()).String()
	etcdResourceMeta := ToEtcdResourceMetaObject(resourceMetaWrap, nil)
	metaBytes, err := core_model.ToJSON(etcdResourceMeta)
	if err != nil {
		return errors.Wrap(err, "failed to convert spec to json")
	}
	value := resourceObject{
		MetaData: metaBytes,
		SpecData: specBytes,
	}
	bytes, err := json.Marshal(value)
	if err != nil {
		return errors.Wrap(err, "json.Marshal error")
	}
	putResponse, err := e.client.Put(ctx, key, string(bytes), clientv3.WithPrevKV())
	if err != nil {
		return errors.Wrap(err, "etcd put error")
	}
	if putResponse.PrevKv == nil {
		return store.ErrorResourceConflict(resource.Descriptor().Name, resource.GetMeta().GetName(), resource.GetMeta().GetMesh())
	}
	// update resource's meta with new version
	resource.SetMeta(resourceMetaWrap)

	return nil
}

func (e ResoucesEtcdStore) Delete(ctx context.Context, resource core_model.Resource, optionsFunc ...store.DeleteOptionsFunc) error {
	opts := store.NewDeleteOptions(optionsFunc...)
	key := etcd.NewEtcdResourcedKey(e.prefix, resource.Descriptor().Name, opts.Mesh, opts.Name).String()
	response, err := e.client.Delete(ctx, key)
	if err != nil {
		return errors.Wrap(err, "etcd delete error")
	}
	if response.Deleted == 0 {
		return store.ErrorResourceNotFound(resource.Descriptor().Name, opts.Name, opts.Mesh)
	}

	return nil
}

func (e ResoucesEtcdStore) Get(ctx context.Context, resource core_model.Resource, optionsFunc ...store.GetOptionsFunc) error {
	opts := store.NewGetOptions(optionsFunc...)
	key := etcd.NewEtcdResourcedKey(e.prefix, resource.Descriptor().Name, opts.Mesh, opts.Name).String()

	response, err := e.client.Get(ctx, key)
	if err != nil {
		return errors.Wrap(err, "etcd get error")
	}
	if response.Count == 0 {
		return store.ErrorResourceNotFound(resource.Descriptor().Name, opts.Name, opts.Mesh)
	}

	if response.Count != 1 {
		return errors.New("get count is not one")
	}

	value := response.Kvs[0].Value
	var resourceObject resourceObject
	err = json.Unmarshal(value, &resourceObject)
	if err != nil {
		return errors.Wrap(err, "json.Unmarshal error")
	}

	if resource.GetSpec() != nil {
		if err := core_model.FromJSON(resourceObject.SpecData, resource.GetSpec()); err != nil {
			return errors.Wrap(err, "failed to convert json to spec")
		}
	}

	etcdResourceMeta := ToEtcdResourceMetaObject(resource.GetMeta(), nil)
	if err := core_model.FromJSON(resourceObject.MetaData, etcdResourceMeta); err != nil {
		return errors.Wrap(err, "failed to convert json to meta")
	}

	resource.SetMeta(etcdResourceMeta)

	if opts.Version != "" && resource.GetMeta().GetVersion() != opts.Version {
		return store.ErrorResourcePreconditionFailed(resource.Descriptor().Name, opts.Name, opts.Mesh)
	}

	return nil
}

func (e ResoucesEtcdStore) List(ctx context.Context, list core_model.ResourceList, optionsFunc ...store.ListOptionsFunc) error {
	opts := store.NewListOptions(optionsFunc...)
	key := etcd.NewEtcdResourcedKey(e.prefix, list.GetItemType(), opts.Mesh, "").Prefix()

	response, err := e.client.Get(ctx, key, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
	if err != nil {
		return errors.Wrap(err, "etcd get error")
	}
	length := 0
	for _, kv := range response.Kvs {
		value := kv.Value

		var resourceObject resourceObject
		err = json.Unmarshal(value, &resourceObject)
		if err != nil {
			return errors.Wrap(err, "json.Unmarshal error")
		}
		item := list.NewItem()

		if item.GetSpec() != nil {
			if err := core_model.FromJSON(resourceObject.SpecData, item.GetSpec()); err != nil {
				return errors.Wrap(err, "failed to convert json to spec")
			}
		}

		etcdResourceMeta := ToEtcdResourceMetaObject(item.GetMeta(), nil)
		if err := core_model.FromJSON(resourceObject.MetaData, etcdResourceMeta); err != nil {
			return errors.Wrap(err, "failed to convert json to meta")
		}

		if !strings.Contains(etcdResourceMeta.GetName(), opts.NameContains) {
			continue
		}

		item.SetMeta(etcdResourceMeta)

		length++
		if err := list.AddItem(item); err != nil {
			return errors.Wrap(err, "list.AddItem error")
		}
	}
	list.GetPagination().SetTotal(uint32(length))
	return nil
}

func NewResoucesEtcdStore(prefix string, metrics core_metrics.Metrics, client *clientv3.Client) (*ResoucesEtcdStore, error) {
	return &ResoucesEtcdStore{
		metrics: metrics,
		client:  client,
		prefix:  prefix,
	}, nil
}
