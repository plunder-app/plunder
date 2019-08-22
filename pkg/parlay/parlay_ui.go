package parlay

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey"
	"github.com/plunder-app/plunder/pkg/parlay/parlaytypes"
)

func contains(v string, a []string) bool {
	for _, i := range a {
		if strings.Contains(v, i) {
			return true
		}
	}
	return false
}

// StartUI will enable parlay to provide an easier way of selecting which operations will be performed
func StartUI(m *parlaytypes.TreasureMap) (*parlaytypes.TreasureMap, error) {

	deployments := []string{}
	for i := range m.Deployments {
		deployments = append(deployments, m.Deployments[i].Name)
	}
	if len(deployments) == 0 {
		return nil, fmt.Errorf("No Deployments were found")
	}

	var multiQs = []*survey.Question{
		{
			Name: "letter",
			Prompt: &survey.MultiSelect{
				Message: "Select deployment(s)",
				Options: deployments,
			},
		},
	}
	deploymentAnswers := []string{}

	// ask the question
	err := survey.Ask(multiQs, &deploymentAnswers)

	if err != nil {
		return nil, err
	}

	// Create a new TreasureMap from the answered questions
	newMap, err := m.FindDeployments(deploymentAnswers)
	if err != nil {
		return nil, err
	}

	for i := range newMap.Deployments {

		// Ask for Hosts
		multiQs[0].Prompt = &survey.MultiSelect{
			Message: fmt.Sprintf("Select Hosts(s) for [%s]", newMap.Deployments[i].Name),
			Options: newMap.Deployments[i].Hosts,
		}

		hostAnswers := []string{}
		err := survey.Ask(multiQs, &hostAnswers)
		if err != nil {
			return nil, err
		}

		// Ask for Actions
		actions := []string{}
		for y := range newMap.Deployments[i].Actions {
			actions = append(actions, m.Deployments[i].Actions[y].Name)
		}

		if len(actions) == 0 {
			return nil, fmt.Errorf("No Deployments were found")
		}
		multiQs[0].Prompt = &survey.MultiSelect{
			Message: fmt.Sprintf("Select Actions(s) for [%s]", newMap.Deployments[i].Name),
			Options: actions,
		}

		deploymentAnswers := []string{}
		err = survey.Ask(multiQs, &deploymentAnswers)
		if err != nil {
			return nil, err
		}

		newMap.Deployments[i].Hosts = hostAnswers
		foundActions, err := newMap.Deployments[i].FindActions(deploymentAnswers)
		if err != nil {
			return nil, err
		}
		newMap.Deployments[i].Actions = foundActions
	}

	return newMap, nil
}
