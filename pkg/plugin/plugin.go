/*
 * This file is part of the KubeVirt project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2023 Red Hat, Inc.
 *
 */

package plugin

import (
	"fmt"

	"github.com/eddev/template-cni/pkg/plugin/netlink"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	type100 "github.com/containernetworking/cni/pkg/types/100"
	"github.com/containernetworking/plugins/pkg/ns"
)

func CmdAdd(args *skel.CmdArgs) error {
	result, err := CmdAddResult(args)
	if err != nil {
		return err
	}
	return result.Print()
}

func CmdAddResult(args *skel.CmdArgs) (types.Result, error) {
	netConf, cniVersion, err := loadConf(args.StdinData)
	if err != nil {
		return nil, err
	}
	fmt.Printf("\nCNI Version: %s\nNetConf: %+v\nstdin: %s", cniVersion, netConf, string(args.StdinData))

	envArgs, err := getEnvArgs(args.Args)
	if err != nil {
		return nil, err
	}
	fmt.Printf("\nEnvArgs: %+v\nargs: %s", envArgs, args.Args)

	netns, err := ns.GetNS(args.Netns)
	if err != nil {
		return nil, fmt.Errorf("failed to open netns %q: %v", netns, err)
	}
	defer netns.Close()

	result := type100.Result{CNIVersion: cniVersion}

	err = netns.Do(func(_ ns.NetNS) error {
		dummy := netlink.NewDummy(args.IfName)
		if lerr := netlink.CreateLink(dummy); lerr != nil {
			return lerr
		}
		fmt.Printf("dummy link %s created", dummy.Name)

		dummyLink, lerr := netlink.ReadLink(dummy.Name)
		if lerr != nil {
			return lerr
		}

		result.Interfaces = append(result.Interfaces, &type100.Interface{
			Name:    dummyLink.Attrs().Name,
			Mac:     "00:00:00:00:00:01",
			Sandbox: netns.Path(),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func CmdDel(args *skel.CmdArgs) error {
	return nil
}

func CmdCheck(args *skel.CmdArgs) error {
	return nil
}
