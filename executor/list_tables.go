/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package executor

import (
	"admin-cli/tabular"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/XiaoMi/pegasus-go-client/idl/admin"
)

// ListTables command.
func ListTables(client *Client, useJSON bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	resp, err := client.Meta.ListApps(ctx, &admin.ListAppsRequest{
		Status: admin.AppStatus_AS_AVAILABLE,
	})
	if err != nil {
		return err
	}

	type tableStruct struct {
		AppID          int32             `json:"app_id"`
		Name           string            `json:"name"`
		PartitionCount int32             `json:"partition_count"`
		CreateTime     string            `json:"create_time"`
		Envs           map[string]string `json:"envs"`
	}
	var tbList []interface{}
	for _, tb := range resp.Infos {
		tbList = append(tbList, tableStruct{
			AppID:          tb.AppID,
			Name:           tb.AppName,
			PartitionCount: tb.PartitionCount,
			CreateTime:     time.Unix(tb.CreateSecond, 0).Format("2006-01-02"),
			Envs:           tb.Envs,
		})
	}

	if useJSON {
		// formats into JSON
		outputBytes, _ := json.MarshalIndent(tbList, "", "  ")
		fmt.Fprintln(client, string(outputBytes))
		return nil
	}

	// formats into tabular
	tabular.Print(client, tbList)
	return nil
}
