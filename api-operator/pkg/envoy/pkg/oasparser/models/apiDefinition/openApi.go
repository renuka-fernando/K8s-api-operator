/*
 *  Copyright (c) 2020, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

//OpenApi version 3
package apiDefinition

import (
	"encoding/json"
	"github.com/getkin/kin-openapi/openapi3"
	logger "github.com/wso2/k8s-api-operator/api-operator/pkg/envoy/loggers"
	"github.com/wso2/k8s-api-operator/api-operator/pkg/envoy/pkg/configs"
	"net/url"
	"strconv"
	"strings"
)

/**
 * Set openApi3 data to mgwSwagger  Instance.
 *
 * @param swagger3  OpenApi3 unmarshalled data
 */
func (swagger *MgwSwagger) SetInfoOpenApi(swagger3 openapi3.Swagger) {
	swagger.swaggerVersion = swagger3.OpenAPI
	if swagger3.Info != nil {
		swagger.description = swagger3.Info.Description
		swagger.title = swagger3.Info.Title
		swagger.version = swagger3.Info.Version
	}

	swagger.vendorExtensible = convertExtensibletoReadableFormat(swagger3.ExtensionProps)
	swagger.resources = SetResourcesOpenApi(swagger3)

	if IsServerUrlIsAvailable(swagger3) {
		for i, _ := range swagger3.Servers {
			endpoint := getHostandBasepathandPort(swagger3.Servers[i].URL)
			swagger.productionUrls = append(swagger.productionUrls, endpoint)
		}
	}
}

/**
 * Set swagger3 resource path details to mgwSwagger  Instance.
 *
 * @param path  Resource path
 * @param pathtype  Path type(Get, Post ... )
 * @param operation  Operation type
 * @return Resource  MgwSwagger resource instance
 */
func setOperationOpenApi(path string, pathtype string, operation *openapi3.Operation) Resource {
	var resource Resource
	if operation != nil {
		resource = Resource{
			path:        path,
			pathtype:    pathtype,
			iD:          operation.OperationID,
			summary:     operation.Summary,
			description: operation.Description,
			//Schemes: operation.,
			//tags: operation.Tags,
			//Security: operation.Security.,
			vendorExtensible: convertExtensibletoReadableFormat(operation.ExtensionProps)}
	}
	return resource
}

/**
 * Set swagger3 all resource to mgwSwagger resources.
 *
 * @param openApi  Swagger3 unmarshalled data
 * @return []Resource  MgwSwagger resource array
 */
func SetResourcesOpenApi(openApi openapi3.Swagger) []Resource {
	var resources []Resource
	if openApi.Paths != nil {
		for path, pathItem := range openApi.Paths {
			var resource Resource
			if pathItem.Get != nil {
				resource = setOperationOpenApi(path, "get", pathItem.Get)
			} else if pathItem.Post != nil {
				resource = setOperationOpenApi(path, "post", pathItem.Post)
			} else if pathItem.Put != nil {
				resource = setOperationOpenApi(path, "put", pathItem.Put)
			} else if pathItem.Delete != nil {
				resource = setOperationOpenApi(path, "delete", pathItem.Delete)
			} else if pathItem.Head != nil {
				resource = setOperationOpenApi(path, "head", pathItem.Head)
			} else if pathItem.Patch != nil {
				resource = setOperationOpenApi(path, "patch", pathItem.Patch)
			} else {
				//resource = setOperation(contxt,"get",pathItem.Get)
			}
			resources = append(resources, resource)
		}
	}

	return resources
}

/**
 * Retrieve host, basepath and port from the endpoint defintion from of the swaggers.
 *
 * @param rawUrl  RawUrl defintion
 * @return Endpoint  Endpoint instance
 */
func getHostandBasepathandPort(rawUrl string) Endpoint {
	var (
		basepath string
		host     string
		port     uint32
	)
	if !strings.Contains(rawUrl, "://") {
		rawUrl = "http://" + rawUrl
	}
	u, err := url.Parse(rawUrl)
	if err != nil {
		logger.LoggerOasparser.Fatal(err)
	}

	host = u.Hostname()
	basepath = u.Path
	if u.Port() != "" {
		u32, err := strconv.ParseUint(u.Port(), 10, 32)
		if err != nil {
			logger.LoggerOasparser.Error("Error passing port value to mgwSwagger", err)
		}
		port = uint32(u32)
	} else {
		//read default port from configs
		conf, errReadConfig := configs.ReadConfigs()
		if errReadConfig != nil {
			logger.LoggerOasparser.Fatal("Error loading configuration. ", errReadConfig)
		}
		port = conf.Envoy.ApiDefaultPort
	}
	return Endpoint{Host: host, Basepath: basepath, Port: port}
}

/**
 * Check the availability od server url in openApi3
 *
 * @param swagger3  Swagger3 unmarshalled data
 * @return bool  Bool value of availability
 */
func IsServerUrlIsAvailable(swagger3 openapi3.Swagger) bool {
	if swagger3.Servers != nil {
		if len(swagger3.Servers) > 0 && (swagger3.Servers[0].URL != "") {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

/**
 * Unmarshall the vendo extensible in open api3.
 *
 * @param vendorExtensible  VendorExtensible data of open api3
 * @return map[string]interface{}  Map of the vendorExtensible
 */
func convertExtensibletoReadableFormat(vendorExtensible openapi3.ExtensionProps) map[string]interface{} {
	jsnRawExtensible := vendorExtensible.Extensions
	b, err := json.Marshal(jsnRawExtensible)
	if err != nil {
		logger.LoggerOasparser.Error("Error marsheling vendor extenstions: ", err)
	}

	var extensible map[string]interface{}
	err = json.Unmarshal(b, &extensible)
	if err != nil {
		logger.LoggerOasparser.Error("Error unmarsheling vendor extenstions:", err)
	}
	return extensible
}
