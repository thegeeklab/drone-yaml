// Copyright (c) 2019, Drone IO Inc.
// Copyright (c) 2021, Robert Kaussow <mail@thegeeklab.de>

package pretty

import "github.com/drone/drone-yaml/yaml"

// helper function pretty prints the cron resource.
func printCron(w writer, v *yaml.Cron) {
	_, _ = w.WriteString("---")
	w.WriteTagValue("version", v.Version)
	w.WriteTagValue("kind", v.Kind)
	w.WriteTagValue("name", v.Name)
	printSpec(w, v)
	_ = w.WriteByte('\n')
	_ = w.WriteByte('\n')
}

// helper function pretty prints the spec block.
func printSpec(w writer, v *yaml.Cron) {
	w.WriteTag("spec")

	w.IndentIncrease()
	w.WriteTagValue("schedule", v.Spec.Schedule)
	w.WriteTagValue("branch", v.Spec.Branch)

	if hasDeployment(v) {
		printDeploy(w, v)
	}

	w.IndentDecrease()
}

// helper function pretty prints the deploy block.
func printDeploy(w writer, v *yaml.Cron) {
	w.WriteTag("deployment")
	w.IndentIncrease()
	w.WriteTagValue("target", v.Spec.Deployment.Target)
	w.IndentDecrease()
}

// helper function returns true if the deployment
// object is empty.
func hasDeployment(v *yaml.Cron) bool {
	return v.Spec.Deployment.Target != ""
}
