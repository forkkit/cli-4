/*
 * Copyright 2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package commands

import (
	"context"
	"fmt"

	"github.com/projectriff/riff/pkg/cli"
	"github.com/projectriff/riff/pkg/parsers"
	"github.com/projectriff/riff/pkg/validation"
	requestv1alpha1 "github.com/projectriff/system/pkg/apis/request/v1alpha1"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type RequestProcessorCreateOptions struct {
	cli.ResourceOptions

	ItemName       string
	Image          string
	ApplicationRef string
	FunctionRef    string

	Env []string
	// TODO implement
	// EnvFrom []string
}

func (opts *RequestProcessorCreateOptions) Validate(ctx context.Context) *cli.FieldError {
	errs := &cli.FieldError{}

	errs = errs.Also(opts.ResourceOptions.Validate((ctx)))

	if opts.ItemName == "" {
		errs = errs.Also(cli.ErrMissingField(cli.ItemFlagName))
	} else {
		errs = errs.Also(validation.K8sName(opts.ItemName, cli.ItemFlagName))
	}

	// application-ref, build-ref and image are mutually exclusive
	used := []string{}
	unused := []string{}

	if opts.ApplicationRef != "" {
		used = append(used, cli.ApplicationRefFlagName)
	} else {
		unused = append(unused, cli.ApplicationRefFlagName)
	}

	if opts.FunctionRef != "" {
		used = append(used, cli.FunctionRefFlagName)
	} else {
		unused = append(unused, cli.FunctionRefFlagName)
	}

	if opts.Image != "" {
		used = append(used, cli.ImageFlagName)
	} else {
		unused = append(unused, cli.ImageFlagName)
	}

	if len(used) == 0 {
		errs = errs.Also(cli.ErrMissingOneOf(unused...))
	} else if len(used) > 1 {
		errs = errs.Also(cli.ErrMultipleOneOf(used...))
	}

	errs = errs.Also(validation.EnvVars(opts.Env, cli.EnvFlagName))

	return errs
}

func NewRequestProcessorCreateCommand(c *cli.Config) *cobra.Command {
	opts := &RequestProcessorCreateOptions{}

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "<todo>",
		Example: "<todo>",
		Args: cli.Args(
			cli.NameArg(&opts.Name),
		),
		PreRunE: cli.ValidateOptions(opts),
		RunE: func(cmd *cobra.Command, args []string) error {
			processor := &requestv1alpha1.RequestProcessor{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: opts.Namespace,
					Name:      opts.Name,
				},
				Spec: requestv1alpha1.RequestProcessorSpec{
					{
						Name: opts.ItemName,
						Template: &corev1.PodSpec{
							Containers: []corev1.Container{{}},
						},
					},
				},
			}

			if opts.ApplicationRef != "" {
				processor.Spec[0].Build = &requestv1alpha1.Build{
					ApplicationRef: opts.ApplicationRef,
				}
			}
			if opts.FunctionRef != "" {
				processor.Spec[0].Build = &requestv1alpha1.Build{
					FunctionRef: opts.FunctionRef,
				}
			}
			if opts.Image != "" {
				processor.Spec[0].Template.Containers[0].Image = opts.Image
			}

			for _, env := range opts.Env {
				if processor.Spec[0].Template.Containers[0].Env == nil {
					processor.Spec[0].Template.Containers[0].Env = []corev1.EnvVar{}
				}
				processor.Spec[0].Template.Containers[0].Env = append(processor.Spec[0].Template.Containers[0].Env, parsers.EnvVar(env))
			}

			processor, err := c.Request().RequestProcessors(opts.Namespace).Create(processor)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Created request processor %q\n", processor.Name)
			return nil
		},
	}

	cli.NamespaceFlag(cmd, c, &opts.Namespace)
	cmd.Flags().StringVar(&opts.ItemName, cli.StripDash(cli.ItemFlagName), "", "<todo>")
	cmd.Flags().StringVar(&opts.Image, cli.StripDash(cli.ImageFlagName), "", "<todo>")
	cmd.Flags().StringVar(&opts.ApplicationRef, cli.StripDash(cli.ApplicationRefFlagName), "", "<todo>")
	cmd.Flags().StringVar(&opts.FunctionRef, cli.StripDash(cli.FunctionRefFlagName), "", "<todo>")
	cmd.Flags().StringArrayVar(&opts.Env, cli.StripDash(cli.EnvFlagName), []string{}, "<todo>")

	return cmd
}
