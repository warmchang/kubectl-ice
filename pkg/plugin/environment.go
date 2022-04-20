package plugin

import (
	"fmt"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var environmentShort = "List the env name and value for each container"

var environmentDescription = ` Print the the environment variables used in running containers in a pod, single pods
and containers can be selected by name. If no name is specified the environment details of all pods in
the current namespace are shown.

The T column in the table output denotes S for Standard and I for init containers`

var environmentExample = `  # List containers env info from pods
  %[1]s env

  # List container env info from pods output in JSON format
  %[1]s env -o json

  # List container env info from a single pod
  %[1]s env my-pod-4jh36

  # List env info for all containers named web-container searching all 
  # pods in the current namespace
  %[1]s env -c web-container

  # List env info for all containers called web-container searching all pods in current
  # namespace sorted by container name in descending order (notice the ! charator)
  %[1]s env -c web-container --sort '!CONTAINER'

  # List env info for all containers called web-container searching all pods in current
  # namespace sorted by pod name in ascending order
  %[1]s env -c web-container --sort PODNAME

  # List container env info from all pods where label app matches web
  %[1]s env -l app=web

  # List container env info from all pods where the pod label app is either web or mail
  %[1]s env -l "app in (web,mail)"`

func environment(cmd *cobra.Command, kubeFlags *genericclioptions.ConfigFlags, args []string) error {
	var podname []string
	var showPodName bool = true
	var translateConfigMap bool

	connect := Connector{}
	if err := connect.LoadConfig(kubeFlags); err != nil {
		return err
	}

	// if a single pod is selected we dont need to show its name
	if len(args) >= 1 {
		podname = args
		if len(podname[0]) >= 1 {
			showPodName = false
		}
	}
	commonFlagList, err := processCommonFlags(cmd)
	if err != nil {
		return err
	}
	connect.Flags = commonFlagList

	podList, err := connect.GetPods(podname)
	if err != nil {
		return err
	}

	if cmd.Flag("translate").Value.String() == "true" {
		translateConfigMap = true
	}

	table := Table{}
	table.SetHeader(
		"T", "PODNAME", "CONTAINER", "NAME", "VALUE",
	)

	if len(commonFlagList.filterList) >= 1 {
		err = table.SetFilter(commonFlagList.filterList)
		if err != nil {
			return err
		}
	}

	if !showPodName {
		// we need to hide the pod name in the table
		table.HideColumn(1)
	}

	for _, pod := range podList {
		connect.SetNamespace(pod.Namespace)
		for _, container := range pod.Spec.Containers {
			// should the container be processed
			if skipContainerName(commonFlagList, container.Name) {
				continue
			}
			allRows := buildEnvFromContainer(container)
			for _, envRow := range allRows {
				tblOut := envBuildRow(container, pod.Name, "S", envRow, connect, translateConfigMap)
				table.AddRow(tblOut...)
			}
		}
		for _, container := range pod.Spec.InitContainers {
			// should the container be processed
			if skipContainerName(commonFlagList, container.Name) {
				continue
			}
			allRows := buildEnvFromContainer(container)
			for _, envRow := range allRows {
				tblOut := envBuildRow(container, pod.Name, "I", envRow, connect, translateConfigMap)
				table.AddRow(tblOut...)
			}
		}
	}

	if err := table.SortByNames(commonFlagList.sortList...); err != nil {
		return err
	}

	outputTableAs(table, commonFlagList.outputAs)
	return nil

}

func envBuildRow(container v1.Container, podName string, containerType string, env v1.EnvVar, connect Connector, translate bool) []Cell {
	var envKey, envValue string
	var configMap string
	var key string

	envKey = env.Name
	if len(env.Value) == 0 {
		if env.ValueFrom.ConfigMapKeyRef != nil {
			configMap = env.ValueFrom.ConfigMapKeyRef.LocalObjectReference.Name
			key = env.ValueFrom.ConfigMapKeyRef.Key
			fmt.Println("^", configMap, key)
			envValue = "CONFIGMAP:" + configMap + " KEY:" + key
		}

		if env.ValueFrom.SecretKeyRef != nil {
			configMap = env.ValueFrom.SecretKeyRef.LocalObjectReference.Name
			key = env.ValueFrom.SecretKeyRef.Key
			fmt.Println("$", configMap, key)
			envValue = "SECRETMAP:" + configMap + " KEY:" + key
			translate = false //never translate secrets
		}

		if translate {
			envValue = connect.GetConfigMapValue(configMap, key)
		}

	} else {
		envValue = env.Value
	}

	return []Cell{
		NewCellText(containerType),
		NewCellText(podName),
		NewCellText(container.Name),
		NewCellText(envKey),
		NewCellText(envValue),
	}
}

func buildEnvFromContainer(container v1.Container) []v1.EnvVar {
	if len(container.Env) == 0 {
		return []v1.EnvVar{}
	}
	return container.Env
}