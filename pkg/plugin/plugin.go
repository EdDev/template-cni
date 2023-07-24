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
	"github.com/containernetworking/plugins/pkg/ns"
)

func CmdAdd(args *skel.CmdArgs) error {
	netConf, cniVersion, err := loadConf(args.StdinData)
	if err != nil {
		return err
	}
	fmt.Printf("\nCNI Version: %s\nNetConf: %+v\nstdin: %s", cniVersion, netConf, string(args.StdinData))

	envArgs, err := getEnvArgs(args.Args)
	if err != nil {
		return err
	}
	fmt.Printf("\nEnvArgs: %+v\nargs: %s", envArgs, args.Args)

	if args.IfName == "" {
		// Nothing to do (probably not a fit in production).
		return nil
	}

	netns, err := ns.GetNS(args.Netns)
	if err != nil {
		return fmt.Errorf("failed to open netns %q: %v", netns, err)
	}
	defer netns.Close()

	err = netns.Do(func(_ ns.NetNS) error {
		dummy := netlink.NewDummy(args.IfName)
		if lerr := netlink.CreateLink(dummy); lerr != nil {
			return lerr
		}
		fmt.Printf("dummy link %s created", dummy.Name)
		return nil
	})

	return err
}

func CmdDel(args *skel.CmdArgs) error {
	return nil
}

func CmdCheck(args *skel.CmdArgs) error {
	return nil
}
