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
package apiDefinition

type Resource struct {
	path             string
	pathtype         string
	description      string
	consumes         []string
	schemes          []string
	tags             []string
	summary          string
	iD               string
	productionUrls   []Endpoint
	sandboxUrls      []Endpoint
	security         []map[string][]string
	vendorExtensible map[string]interface{}
}

func (resource *Resource) GetProdEndpoints() []Endpoint {
	return resource.productionUrls
}

func (resource *Resource) GetSandEndpoints() []Endpoint {
	return resource.sandboxUrls
}

func (resource *Resource) GetPath() string {
	return resource.path
}

func (resource *Resource) GetId() string {
	return resource.iD
}
