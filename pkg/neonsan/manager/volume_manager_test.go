/*
Copyright 2018 Yunify, Inc.

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

package manager

import (
	"github.com/yunify/qingstor-csi/pkg/neonsan/util"
	"strings"
	"testing"
	"time"
)

const (
	TestPoolName           = "csi"
	TestNormalVolumeName   = "foo"
	TestNotFoundVolumeName = "nofound"
)

func TestCreateVolume(t *testing.T) {
	Pools = []string{TestPoolName}
	tests := []struct {
		name      string
		volName   string
		volPool   string
		volSize64 int64
		replicas  int
		infoExist bool
		errStr    string
	}{
		{
			name:      "create succeed",
			volName:   TestNormalVolumeName,
			volPool:   TestPoolName,
			volSize64: 2 * util.Gib,
			replicas:  1,
			infoExist: true,
			errStr:    "",
		},
		{
			name:      "create failed",
			volName:   TestNormalVolumeName,
			volPool:   TestPoolName,
			volSize64: 2 * util.Gib,
			replicas:  1,
			infoExist: false,
			errStr:    "Volume already existed",
		},
	}
	for _, v := range tests {
		volInfo, err := CreateVolume(v.volName, v.volPool, v.volSize64, v.replicas)

		// check volume info
		if (v.infoExist == false && volInfo != nil) || (v.infoExist == true && volInfo == nil) {
			t.Errorf("name [%s]:  volume info expect [%t], but actually [%t]", v.name, v.infoExist, volInfo == nil)
		}

		// check error
		if v.errStr != "" && err != nil {
			if !strings.Contains(err.Error(), v.errStr) {
				t.Errorf("name [%s]: error expect [%s], but actually [%s]", v.name, v.errStr, err.Error())
			}
		} else if v.errStr == "" && err == nil {
			continue
		} else {
			t.Errorf("name [%s]: error expect [%s], but actually [%v]", v.name, v.errStr, err)
		}
	}
}

func TestFindVolume(t *testing.T) {
	Pools = []string{TestPoolName}
	tests := []struct {
		name    string
		volName string
		volPool string
		info    *VolumeInfo
	}{
		{
			name:    "found volume",
			volName: TestNormalVolumeName,
			volPool: TestPoolName,
			info: &VolumeInfo{
				Name:     "foo",
				Pool:     TestPoolName,
				SizeByte: 2 * util.Gib,
				Status:   VolumeStatusOk,
				Replicas: 1,
			},
		},
		{
			name:    "not found volume",
			volName: TestNotFoundVolumeName,
			volPool: TestPoolName,
			info:    nil,
		},
	}
	for _, v := range tests {
		volInfo, err := FindVolume(v.volName, v.volPool)
		if err != nil {
			t.Errorf("name [%s]: volume error [%s]", v.name, err.Error())
		}

		// check volume info
		if v.info != nil && volInfo != nil {
			if v.info.Name != volInfo.Name || v.info.Pool != volInfo.Pool {
				t.Errorf("name [%s]: volume info expect [%v], but actually [%v]", v.name, v.info, volInfo)
			}
		} else if v.info == nil && volInfo == nil {
			continue
		} else {
			t.Errorf("name [%s]: volume info expect [%v], but actually [%v]", v.name, v.info, volInfo)
		}
	}
}

func TestFindVolumeWithoutPool(t *testing.T) {
	Pools = []string{TestPoolName}
	tests := []struct {
		name    string
		volName string
		volPool string
	}{
		{
			name:    "found volume in pool",
			volName: TestNormalVolumeName,
			volPool: TestPoolName,
		},
		{
			name:    "not found volume in pool",
			volName: TestNotFoundVolumeName,
			volPool: "",
		},
	}
	for _, v := range tests {
		ret, err := FindVolumeWithoutPool(v.volName)
		if err != nil {
			t.Errorf("name [%s]: volume error [%s]", v.name, err.Error())
		}
		if v.volPool != "" && ret != nil {
			if v.volPool != ret.Pool {
				t.Errorf("name [%s]: volume pool expect [%s], but actually [%s]", v.name, v.volPool, ret.Pool)
			}
		} else if v.volPool == "" && ret == nil {
			continue
		} else {
			t.Errorf("name [%s]: volume pool expect [%s], but actually [%v]", v.name, v.volPool, ret)
		}
	}
}

func TestListVolumeByPool(t *testing.T) {
	Pools = []string{TestPoolName}
	tests := []struct {
		name    string
		volName string
		volPool string
		info    []*VolumeInfo
	}{
		{
			name:    "found volume",
			volPool: TestPoolName,
			info: []*VolumeInfo{
				{
					Name:     TestNormalVolumeName,
					Pool:     TestPoolName,
					SizeByte: 2 * util.Gib,
				},
			},
		},
		{
			name:    "not found volume",
			volName: TestNotFoundVolumeName,
			volPool: TestPoolName,
			info:    nil,
		},
	}

	for _, v := range tests {
		volList, err := ListVolumeByPool(v.volPool)
		if err != nil {
			t.Errorf("name [%s]: volume error [%s]", v.name, err.Error())
		}
		// verify array
		if len(v.info) != len(volList) {
			t.Errorf("name [%s]: expect [%d], but actually [%d]", v.name, len(v.info), len(volList))
		}
		// check each array element
		for i := range v.info {
			if v.info[i].Name != volList[i].Name || v.info[i].Pool != volList[i].Pool {
				t.Errorf("name [%s]: index [%d] expect [%v], but actually [%v]", v.name, i, v.info[i], volList[i])
			}
		}
	}
}

func TestAttachVolume(t *testing.T) {
	Pools = []string{TestPoolName}
	tests := []struct {
		name   string
		volume string
		pool   string
		errStr string
	}{
		{
			name:   "attach foo image",
			volume: TestNormalVolumeName,
			pool:   TestPoolName,
			errStr: "",
		},
		{
			name:   "reattach foo image",
			volume: TestNormalVolumeName,
			pool:   TestPoolName,
			errStr: "exit status 17",
		},
		{
			name:   "attach not exists image",
			volume: TestNotFoundVolumeName,
			pool:   TestPoolName,
			errStr: "exit status 154",
		},
	}
	for _, v := range tests {
		err := AttachVolume(v.volume, v.pool)
		if err == nil && v.errStr == "" {
			continue
		} else if err != nil && v.errStr != "" {
			if !strings.Contains(err.Error(), v.errStr) {
				t.Errorf("name [%s]: expect contains [%s], but actually [%s]", v.name, v.errStr, err.Error())
			}
		} else {
			t.Errorf("name [%s]: expect [%v], but actually [%v]", v.name, v.errStr, err)
		}
	}
}

func TestFindAttachedVolumeWithoutPool(t *testing.T) {
	Pools = []string{TestPoolName}
	tests := []struct {
		name   string
		volume string
		pool   string
		errStr string
	}{
		{
			name:   "attach info",
			volume: TestNormalVolumeName,
			pool:   TestPoolName,
			errStr: "",
		},
		{
			name:   "nil attach info",
			volume: TestNotFoundVolumeName,
			pool:   "",
			errStr: "",
		},
	}
	for _, v := range tests {
		info, err := FindAttachedVolumeWithoutPool(v.volume)
		if err != nil && v.errStr != "" {
			if !strings.Contains(err.Error(), v.errStr) {
				t.Errorf("name [%s]: expect [%v], but actually [%v]", v.name, v.errStr, err)
			}
		} else if !(err == nil && v.errStr == "") {
			t.Errorf("name [%s]: expect [%v], but actually [%v]", v.name, v.errStr, err)
		}
		if v.pool == "" && info == nil {
			continue
		} else if v.pool != "" && info != nil {
			if v.pool != info.Pool {
				t.Errorf("name [%s]: expect [%v], but actually [%v]", v.name, v.pool, info)
			}
		} else {
			t.Errorf("name [%s]: expect [%v], but actually [%v]", v.name, v.pool, info)
		}
	}
}

func TestDetachVolume(t *testing.T) {
	Pools = []string{TestPoolName}
	time.Sleep(3 * time.Second)
	tests := []struct {
		name   string
		volume string
		pool   string
		errStr string
	}{
		{
			name:   "detach foo image",
			volume: "foo",
			pool:   TestPoolName,
			errStr: "",
		},
		{
			name:   "re-detach foo image",
			volume: "foo",
			pool:   TestPoolName,
			errStr: "exit status 25",
		},
		{
			name:   "detach not exists image",
			volume: "nofound",
			pool:   TestPoolName,
			errStr: "exit status 25",
		},
	}
	for _, v := range tests {
		err := DetachVolume(v.volume, v.pool)
		if err == nil && v.errStr == "" {
			continue
		} else if err != nil && v.errStr != "" {
			if !strings.Contains(err.Error(), v.errStr) {
				t.Errorf("name [%s]: expect contains [%s], but actually [%s]", v.name, v.errStr, err.Error())
			}
		} else {
			t.Errorf("name [%s]: expect [%v], but actually [%v]", v.name, v.errStr, err)
		}
	}
}

func TestDeleteVolume(t *testing.T) {
	Pools = []string{TestPoolName}
	tests := []struct {
		name    string
		volName string
		volPool string
		errStr  string
	}{
		{
			name:    "delete success",
			volName: "foo",
			volPool: TestPoolName,
			errStr:  "",
		},
		{
			name:    "delete failed",
			volName: "nofound",
			volPool: TestPoolName,
			errStr:  "Volume not exists",
		},
	}
	for _, v := range tests {
		err := DeleteVolume(v.volName, v.volPool)
		if v.errStr == "" && err == nil {
			continue
		} else if v.errStr != "" && err != nil {
			if !strings.Contains(err.Error(), v.errStr) {
				t.Errorf("name [%s]: error expect [%s], but actually [%s]", v.name, v.errStr, err.Error())
			}
		} else {
			t.Errorf("name [%s]: error expect [%s], but actually [%v]", v.name, v.errStr, err)
		}
	}
}

func TestProbeNeonsanCommand(t *testing.T) {
	Pools = []string{TestPoolName}
	tests := []struct {
		name   string
		nilErr bool
	}{
		{
			name:   "Probe Neonsan",
			nilErr: true,
		},
	}
	for _, v := range tests {
		err := ProbeNeonsanCommand()
		if (err == nil) != v.nilErr {
			t.Errorf("name [%s]: expect [%t], but actually [%t], error [%v].", v.name, v.nilErr, err == nil, err)
		}
	}
}

func TestProbeQbdCommand(t *testing.T) {
	Pools = []string{TestPoolName}
	tests := []struct {
		name   string
		nilErr bool
	}{
		{
			name:   "Probe Qbd",
			nilErr: true,
		},
	}
	for _, v := range tests {
		err := ProbeQbdCommand()
		if (err == nil) != v.nilErr {
			t.Errorf("name [%s]: expect [%t], but actually [%t], error [%v].", v.name, v.nilErr, err == nil, err)
		}
	}
}
