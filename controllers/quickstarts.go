/*
Copyright 2021 Red Hat OpenShift Data Foundation.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"github.com/ghodss/yaml"
	"github.com/go-logr/logr"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	consolev1 "github.com/openshift/api/console/v1"
)

func (r *StorageSystemReconciler) ensureQuickStarts(logger logr.Logger) error {
	qscrd := extv1.CustomResourceDefinition{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{Name: "consolequickstarts.console.openshift.io", Namespace: ""}, &qscrd)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.V(2).Info("No custom resource definition found for consolequickstart. Skipping quickstart initialization")
			return nil
		}
		return err
	}
	if len(AllQuickStarts) == 0 {
		logger.Info("No quickstarts found")
		return nil
	}
	for _, qs := range AllQuickStarts {
		cqs := consolev1.ConsoleQuickStart{}
		err := yaml.Unmarshal(qs, &cqs)
		if err != nil {
			logger.Error(err, "Failed to unmarshal ConsoleQuickStart", "ConsoleQuickStartString", string(qs))
			continue
		}
		found := consolev1.ConsoleQuickStart{}
		err = r.Client.Get(context.TODO(), types.NamespacedName{Name: cqs.Name, Namespace: cqs.Namespace}, &found)
		if err != nil {
			if errors.IsNotFound(err) {
				err = r.Client.Create(context.TODO(), &cqs)
				if err != nil {
					logger.Error(err, "Failed to create quickstart", "Name", cqs.Name, "Namespace", cqs.Namespace)
					return nil
				}
				logger.Info("Creating quickstarts", "Name", cqs.Name, "Namespace", cqs.Namespace)
				continue
			}
			logger.Error(err, "Error has occurred when fetching quickstarts")
			return nil
		}
		found.Spec = cqs.Spec
		err = r.Client.Update(context.TODO(), &found)
		if err != nil {
			logger.Error(err, "Failed to update quickstart", "Name", cqs.Name, "Namespace", cqs.Namespace)
			return nil
		}
		logger.Info("Updating quickstarts", "Name", cqs.Name, "Namespace", cqs.Namespace)
	}
	return nil
}
