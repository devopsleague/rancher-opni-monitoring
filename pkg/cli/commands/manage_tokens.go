package commands

import (
	"fmt"
	"time"

	cliutil "github.com/rancher/opni-monitoring/pkg/cli/util"
	"github.com/rancher/opni-monitoring/pkg/core"
	"github.com/rancher/opni-monitoring/pkg/management"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
)

func BuildTokensCmd() *cobra.Command {
	tokensCmd := &cobra.Command{
		Use:     "tokens",
		Aliases: []string{"token"},
		Short:   "Manage bootstrap tokens",
	}
	tokensCmd.AddCommand(BuildTokensCreateCmd())
	tokensCmd.AddCommand(BuildTokensRevokeCmd())
	tokensCmd.AddCommand(BuildTokensListCmd())
	return tokensCmd
}

func BuildTokensCreateCmd() *cobra.Command {
	var ttl string
	var labels []string
	tokensCreateCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a bootstrap token",
		Run: func(cmd *cobra.Command, args []string) {
			duration, err := time.ParseDuration(ttl)
			if err != nil {
				lg.Fatal(err)
			}
			labelMap, err := cliutil.ParseKeyValuePairs(labels)
			if err != nil {
				lg.Fatal(err)
			}
			t, err := client.CreateBootstrapToken(cmd.Context(),
				&management.CreateBootstrapTokenRequest{
					Ttl:    durationpb.New(duration),
					Labels: labelMap,
				})
			if err != nil {
				lg.Fatal(err)
			}
			fmt.Println(cliutil.RenderBootstrapToken(t))
		},
	}
	tokensCreateCmd.Flags().StringVar(&ttl, "ttl", "300s", "Time to live")
	tokensCreateCmd.Flags().StringSliceVar(&labels, "labels", []string{}, "Labels which will be auto-applied to any clusters created with this token")
	return tokensCreateCmd
}

func BuildTokensRevokeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "revoke <token>",
		Short: "Revoke a bootstrap token",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			for _, token := range args {
				_, err := client.RevokeBootstrapToken(cmd.Context(),
					&core.Reference{
						Id: token,
					})
				if err != nil {
					lg.Fatal(err)
				}
				lg.Info("Revoked token %s", token)
			}
		},
	}
}

func BuildTokensListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List bootstrap tokens",
		Run: func(cmd *cobra.Command, args []string) {
			t, err := client.ListBootstrapTokens(cmd.Context(), &emptypb.Empty{})
			if err != nil {
				lg.Fatal(err)
			}
			cliutil.RenderBootstrapTokenList(t)
		},
	}
}

func BuildTokensGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <token-id>...",
		Short: "get bootstrap tokens",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			tokenList := []*core.BootstrapToken{}
			for _, id := range args {
				t, err := client.GetBootstrapToken(cmd.Context(), &core.Reference{
					Id: id,
				})
				if err != nil {
					lg.Fatal(err)
				}
				tokenList = append(tokenList, t)
			}
			cliutil.RenderBootstrapTokenList(&core.BootstrapTokenList{
				Items: tokenList,
			})
		},
	}
}
