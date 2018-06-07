package worker

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/choerodon/choerodon-agent/pkg/common"
	"github.com/choerodon/choerodon-agent/pkg/controls"
	"github.com/choerodon/choerodon-agent/pkg/model"
	model_kubernetes "github.com/choerodon/choerodon-agent/pkg/model/kubernetes"
)

func init() {
	registerCmdFunc(model.KubernetesGetLogs, func(w *workerManager, cmd *model.Command) ([]*model.Command, *model.Response, bool) {
		return w.GetLogsByKubernetes(cmd)
	})
	registerCmdFunc(model.KubernetesExec, func(w *workerManager, cmd *model.Command) ([]*model.Command, *model.Response, bool) {
		return w.ExecByKubernetes(cmd)
	})
}

func (w *workerManager) GetLogsByKubernetes(cmd *model.Command) ([]*model.Command, *model.Response, bool) {
	var req *model_kubernetes.GetLogsByKubernetesRequest
	err := json.Unmarshal([]byte(cmd.Payload), &req)
	if err != nil {
		return nil, NewResponseError(cmd.Key, model.KubernetesGetLogsFailed, err), false
	}
	readCloser, err := w.kubeClient.GetLogs(w.namespace, req.PodName, req.ContainerName)
	if err != nil {
		return nil, NewResponseError(cmd.Key, model.KubernetesGetLogsFailed, err), false
	}
	readWriter := struct {
		io.Reader
		io.Writer
	}{
		readCloser,
		ioutil.Discard,
	}
	pipe, err := controls.NewPipeFromEnds(nil, readWriter, w.appClient, req.PipeID, common.Log)
	if err != nil {
		return nil, NewResponseError(cmd.Key, model.KubernetesGetLogsFailed, err), false
	}
	pipe.OnClose(func() {
		readCloser.Close()
	})
	return nil, nil, true
}

func (w *workerManager) ExecByKubernetes(cmd *model.Command) ([]*model.Command, *model.Response, bool) {
	var req *model_kubernetes.ExecByKubernetesRequest
	err := json.Unmarshal([]byte(cmd.Payload), &req)
	if err != nil {
		return nil, NewResponseError(cmd.Key, model.KubernetesExecFailed, err), false
	}
	pipe, err := controls.NewPipe(w.appClient, req.PipeID, common.Exec)
	if err != nil {
		return nil, NewResponseError(cmd.Key, model.KubernetesExecFailed, err), false
	}
	local, _ := pipe.Ends()
	w.kubeClient.Exec(w.namespace, req.PodName, req.ContainerName, local)
	return nil, nil, true
}
