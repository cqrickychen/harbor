package replication

import (
	"net/http"

	common_http "github.com/vmware/harbor/src/common/http"
	"github.com/vmware/harbor/src/common/utils/registry/auth"
	"github.com/vmware/harbor/src/jobservice_v2/env"
	"github.com/vmware/harbor/src/jobservice_v2/logger"
)

// Deleter deletes repository or images on the destination registry
type Deleter struct {
	ctx         env.JobContext
	repository  *repository
	dstRegistry *registry
	logger      logger.Interface
	retry       bool
}

// ShouldRetry : retry if the error is network error
func (d *Deleter) ShouldRetry() bool {
	return d.retry
}

// MaxFails ...
func (d *Deleter) MaxFails() uint {
	return 3
}

// Validate ....
func (d *Deleter) Validate(params map[string]interface{}) error {
	return nil
}

// Run ...
func (d *Deleter) Run(ctx env.JobContext, params map[string]interface{}) error {
	err := d.run(ctx, params)
	d.retry = retry(err)
	return err
}

func (d *Deleter) run(ctx env.JobContext, params map[string]interface{}) error {
	if err := d.init(ctx, params); err != nil {
		return err
	}

	return d.delete()
}

func (d *Deleter) init(ctx env.JobContext, params map[string]interface{}) error {
	d.logger = ctx.GetLogger()
	d.ctx = ctx

	if canceled(d.ctx) {
		d.logger.Warning(errCanceled.Error())
		return errCanceled
	}

	d.repository = &repository{
		name: params["repository"].(string),
	}
	if tags, ok := params["tags"]; ok {
		tgs := tags.([]interface{})
		for _, tg := range tgs {
			d.repository.tags = append(d.repository.tags, tg.(string))
		}
	}

	url := params["dst_registry_url"].(string)
	insecure := params["dst_registry_insecure"].(bool)
	cred := auth.NewBasicAuthCredential(
		params["dst_registry_username"].(string),
		params["dst_registry_password"].(string))

	var err error
	d.dstRegistry, err = initRegistry(url, insecure, cred, d.repository.name)
	if err != nil {
		d.logger.Errorf("failed to create client for destination registry: %v", err)
		return err
	}

	d.logger.Infof("initialization completed: repository: %s, tags: %v, destination URL: %s, insecure: %v",
		d.repository.name, d.repository.tags, d.dstRegistry.url, d.dstRegistry.insecure)

	return nil
}

func (d *Deleter) delete() error {
	repository := d.repository.name
	tags := d.repository.tags
	if len(tags) == 0 {
		if canceled(d.ctx) {
			d.logger.Warning(errCanceled.Error())
			return errCanceled
		}
		if err := d.dstRegistry.DeleteRepository(repository); err != nil {
			if e, ok := err.(*common_http.Error); ok && e.Code == http.StatusNotFound {
				d.logger.Warningf("repository %s not found", repository)
				return nil
			}
			d.logger.Errorf("failed to delete repository %s: %v", repository, err)
			return err
		}
		d.logger.Infof("repository %s has been deleted", repository)
		return nil
	}

	for _, tag := range tags {
		if canceled(d.ctx) {
			d.logger.Warning(errCanceled.Error())
			return errCanceled
		}
		if err := d.dstRegistry.DeleteImage(repository, tag); err != nil {
			if e, ok := err.(*common_http.Error); ok && e.Code == http.StatusNotFound {
				d.logger.Warningf("image %s:%s not found", repository, tag)
				return nil
			}
			d.logger.Errorf("failed to delete image %s:%s: %v", repository, tag, err)
			return err
		}
		d.logger.Infof("image %s:%s has been deleted", repository, tag)
	}
	return nil
}