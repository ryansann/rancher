package cluster

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/rancher/norman/api/access"
	"github.com/rancher/norman/types"
	mgmtv3 "github.com/rancher/rancher/pkg/generated/norman/management.cattle.io/v3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (a ActionHandler) RotateEncryptionKey(actionName string, action *types.Action, apiContext *types.APIContext) error {
	rtn := map[string]interface{}{
		"type":    "rotateEncryptionKeyOutput",
		"message": "rotating key for encrypted data",
	}

	var mgmtCluster mgmtv3.Cluster
	if err := access.ByID(apiContext, apiContext.Version, apiContext.Type, apiContext.ID, &mgmtCluster); err != nil {
		rtn["message"] = "cluster does not exist"
		apiContext.WriteResponse(http.StatusBadRequest, rtn)
		return errors.Wrapf(err, "failed to get cluster by ID %s", apiContext.ID)
	}

	cluster, err := a.ClusterClient.Get(apiContext.ID, v1.GetOptions{})
	if err != nil {
		rtn["message"] = "cluster does not exist"
		apiContext.WriteResponse(http.StatusBadRequest, rtn)
		return errors.Wrapf(err, "failed to get cluster by ID %s", apiContext.ID)
	}

	cluster.Spec.RancherKubernetesEngineConfig.RotateEncryptionKey = true
	if _, err := a.ClusterClient.Update(cluster); err != nil {
		rtn["message"] = "failed to update cluster object"
		apiContext.WriteResponse(http.StatusInternalServerError, rtn)
		return errors.Wrapf(err, "unable to update cluster %s", cluster.Name)
	}

	apiContext.WriteResponse(http.StatusOK, rtn)
	return nil
}
