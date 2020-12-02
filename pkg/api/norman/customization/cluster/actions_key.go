package cluster

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"

	v3client "github.com/rancher/rancher/pkg/client/generated/management/v3"

	"github.com/rancher/rancher/pkg/controllers/management/etcdbackup"

	"github.com/pkg/errors"
	"github.com/rancher/norman/api/access"
	"github.com/rancher/norman/types"
	mgmtv3 "github.com/rancher/rancher/pkg/generated/norman/management.cattle.io/v3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (a ActionHandler) RotateEncryptionKey(actionName string, action *types.Action, apiContext *types.APIContext) error {
	response := map[string]interface{}{
		"type": v3client.RotateEncryptionKeyOutputType,
		v3client.RotateEncryptionKeyOutputFieldMessage: "etcd backup created, rotating encryption key",
	}

	var mgmtCluster mgmtv3.Cluster
	if err := access.ByID(apiContext, apiContext.Version, apiContext.Type, apiContext.ID, &mgmtCluster); err != nil {
		response[v3client.RotateEncryptionKeyOutputFieldMessage] = "cluster does not exist"
		apiContext.WriteResponse(http.StatusBadRequest, response)
		return errors.Wrapf(err, "failed to get cluster by ID %s", apiContext.ID)
	}

	cluster, err := a.ClusterClient.Get(apiContext.ID, v1.GetOptions{})
	if err != nil {
		response[v3client.RotateEncryptionKeyOutputFieldMessage] = "cluster does not exist"
		apiContext.WriteResponse(http.StatusBadRequest, response)
		return errors.Wrapf(err, "failed to get cluster by ID %s", apiContext.ID)
	}

	// create etcd backup so that if the rotate encryption action fails, the user can restore from this snapshot
	newBackup, err := etcdbackup.NewBackupObject(cluster, true)
	if err != nil {
		response[v3client.RotateEncryptionKeyOutputFieldMessage] = "failed to initialize etcdbackup object"
		apiContext.WriteResponse(http.StatusInternalServerError, response)
		return errors.Wrapf(err, "failed to initialize etcdbackup object")
	}

	backup, err := a.BackupClient.Create(newBackup)
	if err != nil {
		response[v3client.RotateEncryptionKeyOutputFieldMessage] = "failed to create etcdbackup object"
		apiContext.WriteResponse(http.StatusInternalServerError, response)
		return errors.Wrapf(err, "failed to cteate etcdbackup object")
	}

	response[v3client.RotateEncryptionKeyOutputFieldBackup] = backup

	cluster.Spec.RancherKubernetesEngineConfig.RotateEncryptionKey = true
	if _, err := a.ClusterClient.Update(cluster); err != nil {
		response[v3client.RotateEncryptionKeyOutputFieldMessage] = "failed to update cluster object"
		apiContext.WriteResponse(http.StatusInternalServerError, response)
		return errors.Wrapf(err, "unable to update cluster %s", cluster.Name)
	}

	res, err := json.Marshal(response)
	if err != nil {
		return err
	}

	apiContext.Response.Header().Set("Content-Type", "application/json")
	http.ServeContent(apiContext.Response, apiContext.Request, v3.ClusterActionRotateEncryptionKey, time.Now(), bytes.NewReader(res))
	return nil
}
