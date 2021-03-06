package main

import (
	"log"

	"github.com/flynn/flynn/Godeps/_workspace/src/github.com/flynn/go-docopt"
	"github.com/flynn/flynn/controller/client"
	ct "github.com/flynn/flynn/controller/types"
)

func init() {
	register("resource", runResource, `
usage: flynn resource
       flynn resource add <provider>
       flynn resource remove <provider> <resource>

Manage resources for the app.

Commands:
       With no arguments, shows a list of resources.

       add     provisions a new resource for the app using <provider>.
       remove  removes the existing <resource> provided by <provider>.
`)
}

func runResource(args *docopt.Args, client controller.Client) error {
	if args.Bool["add"] {
		return runResourceAdd(args, client)
	}
	if args.Bool["remove"] {
		return runResourceRemove(args, client)
	}

	resources, err := client.AppResourceList(mustApp())
	if err != nil {
		return err
	}

	w := tabWriter()
	defer w.Flush()

	var provider *ct.Provider

	listRec(w, "ID", "Provider ID", "Provider Name")
	for _, j := range resources {
		provider, err = client.GetProvider(j.ProviderID)
		if err != nil {
			return err
		}
		listRec(w, j.ID, j.ProviderID, provider.Name)
	}

	return err
}

func runResourceAdd(args *docopt.Args, client controller.Client) error {
	provider := args.String["<provider>"]

	res, err := client.ProvisionResource(&ct.ResourceReq{ProviderID: provider, Apps: []string{mustApp()}})
	if err != nil {
		return err
	}

	env := make(map[string]*string)
	for k, v := range res.Env {
		s := v
		env[k] = &s
	}

	releaseID, err := setEnv(client, "", env)
	if err != nil {
		return err
	}

	log.Printf("Created resource %s and release %s.", res.ID, releaseID)

	return nil
}

func runResourceRemove(args *docopt.Args, client controller.Client) error {
	provider := args.String["<provider>"]
	resource := args.String["<resource>"]

	res, err := client.DeleteResource(provider, resource)
	if err != nil {
		return err
	}

	release, err := client.GetAppRelease(mustApp())
	if err != nil {
		return err
	}

	// Unset all the keys associated with this resource
	env := make(map[string]*string)
	for k := range res.Env {
		// Only unset the key if it hasn't been modified
		if release.Env[k] == res.Env[k] {
			env[k] = nil
		}
	}

	releaseID, err := setEnv(client, "", env)
	if err != nil {
		return err
	}

	log.Printf("Deleted resource %s, created release %s.", res.ID, releaseID)

	return nil
}
